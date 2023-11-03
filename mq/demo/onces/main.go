package main

import (
    "sync"
    "testing"
)

var once sync.Once

var onceFunc sync.Once

func TestOnceFunc() {
    calls := 0
    f := sync.OnceFunc(func() {
        calls++
        println("calls:", calls)
    })
    avg := testing.AllocsPerRun(100, f)
    println("avg:", avg)
}

func main() {
    wg := &sync.WaitGroup{}
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            once.Do(func() {
                println("Hello, World")
            })
        }()
    }
    wg.Wait()
    TestOnceFunc()
}
