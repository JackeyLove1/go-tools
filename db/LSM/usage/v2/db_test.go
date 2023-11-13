package main

import (
    "math/rand"
    "os"
    "path/filepath"
    "sync"
    "testing"

    goleveldb "github.com/syndtr/goleveldb/leveldb"
)

var once sync.Once

const randomSeed = 42
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const length = 10

func GenerateRandomStr(length int) string {
    once.Do(func() {
        rand.Seed(randomSeed)
    })
    randomBytes := make([]byte, length)
    charsetLength := len(charset)
    for i := 0; i < length; i++ {
        randomIndex := rand.Int() % charsetLength
        randomBytes[i] = charset[randomIndex]
    }
    return string(randomBytes)
}

func TestRandomString(t *testing.T) {
    str1 := GenerateRandomStr(length)
    str2 := GenerateRandomStr(length)
    if len(str1) != length || len(str2) != length || str1 == str2 {
        t.FailNow()
    }
}

func BenchmarkDB(b *testing.B) {
    dir, err := os.Getwd()
    failOnErr(err)
    path := filepath.Join(dir, "temp")
    defer os.RemoveAll(path)
    os.MkdirAll(path, 0755)
    db, err := goleveldb.OpenFile(path, nil)
    failOnErr(err)
    defer db.Close()
    keys := make([]string, b.N)
    values := make([]string, b.N)
    for i := 0; i < b.N; i++ {
        keys[i] = GenerateRandomStr(length)
        values[i] = GenerateRandomStr(length)
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        err = db.Put([]byte(keys[i]), []byte(values[i]), nil)
        if err != nil {
            b.FailNow()
        }
    }
    for i := 0; i < b.N; i++ {
        value, err := db.Get([]byte(keys[i]), nil)
        if err != nil {
            b.FailNow()
        }
        if string(value) != values[i] {
            b.FailNow()
        }
    }
}
