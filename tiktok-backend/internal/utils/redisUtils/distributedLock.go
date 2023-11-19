package redisUtils

import (
    "context"
    "errors"
    "time"

    "github.com/go-redis/redis/v8"
    "tiktok-backend/internal/utils/constants"
)

type DistributedLock interface {
    //Lock 尝试加锁
    Lock(context.Context) error
    //TryLock 在duration过后，自动放弃锁
    TryLock(context.Context) error
    //UnLock 解锁
    UnLock(ctx context.Context) error
}

type Lock struct {
    client          *redis.Client // redis客户端
    script          *redis.Script // 解锁脚本
    resource        string        // 锁资源
    randomValue     string        // 随机值
    watchDog        chan struct{} // 看门狗
    ttl             time.Duration // 过期时间
    tryLockInterval time.Duration // 重新获得锁间隔
}

func (l *Lock) TryLock(ctx context.Context) error {
    success, err := l.client.SetNX(ctx, l.resource, l.randomValue, l.ttl).Result()
    if err != nil {
        return nil
    }
    if !success {
        return constants.LockFailedErr
    }
    go l.startWatchDog()
    return nil
}

func (l *Lock) startWatchDog() {
    ticker := time.NewTicker(l.ttl / 3)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            // delay the lock time
            ctx, cancel := context.WithTimeout(context.Background(), l.ttl/3*2)
            ok, err := l.client.Expire(ctx, l.resource, l.ttl).Result()
            cancel()
            // 异常或锁已经不存在则不再续期
            if err != nil || !ok {
                return
            }
        case <-l.watchDog:
            // has already unlock
            return
        }
    }
}

func (l *Lock) Lock(ctx context.Context) error {
    // try to lock
    err := l.TryLock(ctx)
    if err == nil {
        return nil
    }
    if !errors.Is(constants.LockFailedErr, err) {
        return err
    }
    // failed to lock and always try
    ticker := time.NewTicker(l.tryLockInterval)
    defer ticker.Stop()
    for {
        select {
        case <-ctx.Done():
            // out of time
            return constants.TimeOutErr
        case <-ticker.C:
            err := l.TryLock(ctx)
            if err != nil {
                return nil
            }
            if errors.Is(constants.LockFailedErr, err) {
                return err
            }
        }
    }
}

func (l *Lock) UnLock(ctx context.Context) error {
    err := l.script.Run(ctx, l.client, []string{l.resource}, l.randomValue).Err()
    // close watchdog
    return err
}
