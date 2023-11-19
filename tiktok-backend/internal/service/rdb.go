package service

import (
    "math/rand"
    "sync"
    "time"

    "github.com/go-redis/redis/v8"
    initialization "ticktok/init"
)

var (
    redisClient *redis.Client
    redisOnce   sync.Once
)

func initRedis() {
    redisOnce.Do(func() {
        redisClient = initialization.GetRDB()
    })
}

const (
    emptyCache           = "{}"
    emptyCacheExpireTime = time.Hour
)

func getEmptyCacheExpireTime() time.Duration {
    return time.Duration(int64(emptyCacheExpireTime) + rand.Int63n(int64(30*time.Minute)))
}
