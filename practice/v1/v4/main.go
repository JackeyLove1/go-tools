package main

import (
    "strconv"
    "time"

    "github.com/syndtr/goleveldb/leveldb"
    "github.com/syndtr/goleveldb/leveldb/util"
)

const (
    Number = 10000
    Key    = "key-"
    Target = "target-"
)

var db *leveldb.DB

func PutData() {
    start := time.Now()
    for i := 0; i < Number; i++ {
        dbKey := Key + strconv.Itoa(i)
        dbValue := "values"
        err := db.Put([]byte(dbKey), []byte(dbValue), nil)
        if err != nil {
            panic(err)
        }
    }
    println("Put Cost:", time.Since(start).Milliseconds(), "ms")
}

func Iter() {
    start := time.Now()
    iter := db.NewIterator(&util.Range{
        Start: []byte(Target),
        Limit: nil,
    }, nil)
    for iter.Next() {
        val := iter.Value()
        println("val:", string(val))
    }
    iter.Release()
    println("Iter Cost:", time.Since(start).Milliseconds(), "ms")
    err := iter.Error()
    println("Error:", err)
}

func main() {
    var err error
    db, err = leveldb.OpenFile("db", nil)
    if err != nil {
        panic(err)
    }
    defer db.Close()
    PutData()
    _ = db.Put([]byte(Target), []byte("answer"), nil)
    Iter()
}
