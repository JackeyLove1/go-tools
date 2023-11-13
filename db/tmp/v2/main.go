package main

import "sync/atomic"

func main() {
    var count atomic.Uint32
    println(count.Load())
}
