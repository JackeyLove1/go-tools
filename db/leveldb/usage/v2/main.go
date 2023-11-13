package main

import (
    "os"
    "path/filepath"

    goleveldb "github.com/syndtr/goleveldb/leveldb"
)

func failOnErr(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
    dir, err := os.Getwd()
    failOnErr(err)
    path := filepath.Join(dir, "temp")
    db, err := goleveldb.OpenFile(path, nil)
    failOnErr(err)
    defer db.Close()
    err = db.Put([]byte("key"), []byte("value"), nil)
    failOnErr(err)
    data, err := db.Get([]byte("key"), nil)
    failOnErr(err)
    println("Value:", string(data))
    err = db.Delete([]byte("key"), nil)
    failOnErr(err)
}
