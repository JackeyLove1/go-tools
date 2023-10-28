package main

import (
    "fmt"
    "strconv"
    "time"

    "go-tools/mq/RabbitMQ"
)

func main() {
    r := RabbitMQ.NewRabbitMQPubSub("" + "newProduct")
    for i := 0; i < 1000; i++ {
        r.PublishPub("订阅模式生产第" + strconv.Itoa(i) + "条" + "数据")
        fmt.Println("订阅模式生产第" + strconv.Itoa(i) + "条" + "数据")
        time.Sleep(1 * time.Second)
    }
}
