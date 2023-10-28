package main

import (
    "fmt"
    "strconv"
    "time"

    "go-tools/mq/RabbitMQ"
)

func main() {
    r1 := RabbitMQ.NewRabbitMQRouting("exec", "one")
    r2 := RabbitMQ.NewRabbitMQRouting("exec", "1233444")
    for i := 0; i < 10; i++ {
        r1.PublishRouting("Hello imooc one!" + strconv.Itoa(i))
        r2.PublishRouting("Hello imooc two!" + strconv.Itoa(i))
        time.Sleep(1 * time.Second)
        fmt.Printf("Send :%d\n", i)
    }
}
