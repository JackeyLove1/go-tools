package main

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/whuanle/lsm"
    "github.com/whuanle/lsm/config"
)

type TestValue struct {
    A int64
    B int64
    C int64
    D string
}

func failOnErr(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
    dir, err := os.Getwd()
    failOnErr(err)
    path := filepath.Join(dir, "temp")
    defer os.RemoveAll(path)
    os.MkdirAll(path, 0755)
    lsm.Start(config.Config{
        DataDir:    path,
        Level0Size: 1,
        PartSize:   4,
        Threshold:  500,
    })
    // 64 个字节
    testV := TestValue{
        A: 1,
        B: 1,
        C: 3,
        D: "00000000000000000000000000000000000000",
    }

    lsm.Set("aaa", testV)

    value, success := lsm.Get[TestValue]("aaa")
    if success {
        fmt.Println(value)
    }

    lsm.Delete[string]("aaa")
}
