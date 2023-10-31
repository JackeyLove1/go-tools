package main

import (
    "context"
    "time"
)

func main() {
    ticker := time.NewTicker(1 * time.Second)
    idx := 0
    ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
    for {
        select {
        case <-ticker.C:
            idx++
            println("ticker:", idx)
        case <-ctx.Done():
            ticker.Stop()
            println("timeout")
            return
        }
    }
}
