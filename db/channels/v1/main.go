package main

func main() {
    ch1 := make(chan int, 10)
    ch2 := make(chan int, 10)
    go func() {
        defer close(ch1)
        for i := 0; i < 10; i++ {
            ch1 <- i
        }
    }()
    go func() {
        defer close(ch2)
        for {
            select {
            case v, ok := <-ch1:
                if !ok {
                    return
                }
                ch2 <- v * v
            }
        }
    }()

    for i := range ch2 {
        println(i)
    }
}
