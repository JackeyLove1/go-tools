package test

import (
    "context"
    "log"
    "testing"
    "time"

    "github.com/go-redis/redis/v7"
)

var ctx = context.Background()
var rdbClient *redis.Client

type RedisLock struct {
    client  *redis.Client
    key     string
    value   string
    timeout time.Duration
}

func NewRedisLock(client *redis.Client, key, value string, timeout time.Duration) *RedisLock {
    return &RedisLock{
        client:  client,
        key:     key,
        value:   value,
        timeout: timeout,
    }
}

func (lock *RedisLock) TryLock() bool {
    result, err := lock.client.SetNX(lock.key, lock.value, lock.timeout).Result()
    if err != nil {
        log.Printf("failed to acquire lock:%e\n", err)
        return false
    }
    return result
}

func (lock *RedisLock) Unlock() error {
    _, err := lock.client.Del(lock.key).Result()
    if err != nil {
        log.Printf("failed to release lock:%e\n", err)
    }
    return err
}

func init() {
    rdbClient = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    _, err := rdbClient.Ping().Result()
    if err != nil {
        panic(err)
    }
}

func TestRedis(t *testing.T) {
    lock := NewRedisLock(rdbClient, "lock-test", "value", 3*time.Second)
    if lock.TryLock() {
        time.Sleep(1 * time.Second)
        log.Printf("succeed to get the lock")
        lock.Unlock()
    } else {
        log.Printf("failed to get the lock")
    }
}
