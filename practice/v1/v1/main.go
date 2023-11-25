package main

import (
    "context"
    "fmt"

    "github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func failOnErr(err error, msg string) {
    if err != nil {
        panic(fmt.Sprintf("Msg:%s, Error:%s", msg, err.Error()))
    }
}

func Ping(ctx context.Context) error {
    pong, err := rdb.Ping(ctx).Result()
    if err != nil {
        return err
    }
    println("pong:", pong)
    return nil
}

func String(ctx context.Context) {
    println("Start String ... ")
    const key = "string_key"
    const value = "string_value"
    err := rdb.Set(ctx, key, value, 0).Err()
    failOnErr(err, "Set")

    val, err := rdb.Get(ctx, key).Result()
    failOnErr(err, "Get")
    println("val:", val)
    defer func() {
        _, err = rdb.Del(ctx, key).Result()
        if err != nil {
            println("Del err:", err)
        }
    }()

    val2, err := rdb.Get(ctx, "key2").Result()
    if err == redis.Nil {
        println("key2 does not exist")
    } else if err != nil {
        failOnErr(err, "Get")
    } else {
        println("val2:", val2)
    }
}

func List(ctx context.Context) {
    println("Start List ... ")
    const key = "list_key"
    value := []any{1, 2, 3, "list_value", 3.14}
    err := rdb.RPush(ctx, key, value...).Err()
    failOnErr(err, "RPush")
    defer func() {
        _, err = rdb.Del(ctx, key).Result()
        if err != nil {
            println("Del err:", err)
        }
    }()
    v2, err := rdb.LRange(ctx, key, 0, -1).Result()
    failOnErr(err, "LRange")
    for _, val := range v2 {
        println("val:", val)
    }
}

func Hashset(ctx context.Context) {
    const key = "hash_key"
    err := rdb.HSet(ctx, key, "Name", "name", "Age", 18, "Height", 3.14).Err()
    failOnErr(err, "HSet")
    defer rdb.Del(ctx, key)
    age, err := rdb.HGet(ctx, key, "Age").Result()
    failOnErr(err, "Hash Get Filed")
    println("age:", age)
    for filed, value := range rdb.HGetAll(ctx, key).Val() {
        println("filed:", filed, "value:", value)
    }
}

func main() {
    rdb = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })
    ctx := context.TODO()
    var err error
    err = Ping(ctx)
    failOnErr(err, "Ping")
    String(ctx)
    List(ctx)
    Hashset(ctx)
}
