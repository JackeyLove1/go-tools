package init

import (
    "context"
    "sync"

    "github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var rdbOnce sync.Once

func InitRDB() {
    rdbOnce.Do(func() {
        rdb = redis.NewClient(&redis.Options{
            Addr:     "localhost:6379",
            Password: "", // no password set
            DB:       0,  // use default DB
        })
    })
    _, err := rdb.Ping(context.Background()).Result()
    if err != nil {
        stdOutLogger.Panic().Caller().Str("Error occurs in InitRDB,", err.Error())
    }
}

func GetRDB() *redis.Client {
    return rdb
}
