package main

import (
    "context"
    "fmt"
    "math/rand"
    "time"

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
    _, err = rdb.Del(ctx, "key2").Result()
    failOnErr(err, "DEL")
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

func SetNX(ctx context.Context) {
    const key = "nx_key"
    const value = "nx_value"
    set, err := rdb.SetNX(ctx, key, value, 10*time.Millisecond).Result()
    failOnErr(err, "SetNX")
    println("set:", set)
    time.Sleep(10 * time.Millisecond)
    _, err = rdb.Get(ctx, key).Result()
    if err != redis.Nil {
        println("failed to get nx_key")
    }
}

func Iter(ctx context.Context) {
    iter := rdb.Scan(ctx, 0, "k*", 0).Iterator()
    for iter.Next(ctx) {
        println("key:", iter.Val())
    }
    if err := iter.Err(); err != nil {
        panic(err)
    }
}

func DeleteWithoutTTL(ctx context.Context) {
    rdb.SetNX(ctx, "key3", "values", 10*time.Second)
    iter := rdb.Scan(ctx, 0, "k*", 0).Iterator()
    for iter.Next(ctx) {
        key := iter.Val()
        d, err := rdb.TTL(ctx, key).Result()
        if err != nil {
            panic(err)
        }
        if d == -1 {
            // means no ttl
            if err := rdb.Del(ctx, key).Err(); err != nil {
                panic(err)
            }
            println("Delete key:", key)
        }
    }
    if err := iter.Err(); err != nil {
        panic(err)
    }
}

func Pipeline(ctx context.Context) {
    if _, err := rdb.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
        pipeliner.HSet(ctx, "hash_key1", "str1", "value")
        pipeliner.HSet(ctx, "hash_key2", "str2", "value")
        pipeliner.HSet(ctx, "hash_key3", "str3", "value")
        val, err := rdb.HGet(ctx, "hash_key1", "str1").Result()
        if err != nil {
            failOnErr(err, "pipeliner HGet")
        }
        println("val:", val)
        return nil
    }); err != nil {
        panic(err)
    }
}

func bloomFilter(ctx context.Context) {
    inserted, err := rdb.Do(ctx, "BF.ADD", "bf_key", "item0").Bool()
    if err != nil {
        panic(err)
    }
    if inserted {
        fmt.Println("item0 was inserted")
    }

    for _, item := range []string{"item0", "item1"} {
        exists, err := rdb.Do(ctx, "BF.EXISTS", "bf_key", item).Bool()
        if err != nil {
            panic(err)
        }
        if exists {
            fmt.Printf("%s does exist\n", item)
        } else {
            fmt.Printf("%s does not exist\n", item)
        }
    }
    _, err = rdb.Do(ctx, "BF.DEL", "bf_key", "item0").Bool()
    failOnErr(err, "BF.DEL")
}

func cuckooFilter(ctx context.Context) {
    inserted, err := rdb.Do(ctx, "CF.ADDNX", "cf_key", "item0").Bool()
    if err != nil {
        panic(err)
    }
    if inserted {
        fmt.Println("item0 was inserted")
    } else {
        fmt.Println("item0 already exists")
    }

    for _, item := range []string{"item0", "item1"} {
        exists, err := rdb.Do(ctx, "CF.EXISTS", "cf_key", item).Bool()
        if err != nil {
            panic(err)
        }
        if exists {
            fmt.Printf("%s does exist\n", item)
        } else {
            fmt.Printf("%s does not exist\n", item)
        }
    }

    deleted, err := rdb.Do(ctx, "CF.DEL", "cf_key", "item0").Bool()
    if err != nil {
        panic(err)
    }
    if deleted {
        fmt.Println("item0 was deleted")
    }
}

func countMinSketch(ctx context.Context) {
    if err := rdb.Do(ctx, "CMS.INITBYPROB", "count_min", 0.001, 0.01).Err(); err != nil {
        panic(err)
    }

    items := []string{"item1", "item2", "item3", "item4", "item5"}
    counts := make(map[string]int, len(items))

    for i := 0; i < 10000; i++ {
        n := rand.Intn(len(items))
        item := items[n]

        if err := rdb.Do(ctx, "CMS.INCRBY", "count_min", item, 1).Err(); err != nil {
            panic(err)
        }
        counts[item]++
    }

    for item, count := range counts {
        ns, err := rdb.Do(ctx, "CMS.QUERY", "count_min", item).Int64Slice()
        if err != nil {
            panic(err)
        }
        fmt.Printf("%s: count-min=%d actual=%d\n", item, ns[0], count)
    }
}

func topK(ctx context.Context) {
    if err := rdb.Do(ctx, "TOPK.RESERVE", "top_items", 3).Err(); err != nil {
        panic(err)
    }

    counts := map[string]int{
        "item1": 1000,
        "item2": 2000,
        "item3": 3000,
        "item4": 4000,
        "item5": 5000,
        "item6": 6000,
    }

    for item, count := range counts {
        for i := 0; i < count; i++ {
            if err := rdb.Do(ctx, "TOPK.INCRBY", "top_items", item, 1).Err(); err != nil {
                panic(err)
            }
        }
    }

    items, err := rdb.Do(ctx, "TOPK.LIST", "top_items").StringSlice()
    if err != nil {
        panic(err)
    }

    for _, item := range items {
        ns, err := rdb.Do(ctx, "TOPK.COUNT", "top_items", item).Int64Slice()
        if err != nil {
            panic(err)
        }
        fmt.Printf("%s: top-k=%d actual=%d\n", item, ns[0], counts[item])
    }
}

func HyperLogLog(ctx context.Context) {
    for i := 0; i < 10; i++ {
        if err := rdb.PFAdd(ctx, "myset", fmt.Sprint(i)).Err(); err != nil {
            panic(err)
        }
    }

    card, err := rdb.PFCount(ctx, "myset").Result()
    if err != nil {
        panic(err)
    }
    println("set cardinality", card)
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
    SetNX(ctx)
    Iter(ctx)
    DeleteWithoutTTL(ctx)
    Pipeline(ctx)
    bloomFilter(ctx)
    cuckooFilter(ctx)
    countMinSketch(ctx)
    topK(ctx)
    HyperLogLog(ctx)
}
