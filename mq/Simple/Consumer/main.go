package main

import (
    "fmt"

    "go-tools/mq/RabbitMQ"
)

func main() {
    r := RabbitMQ.NewRabbitMQSimple("test")
    r.PublishSimple("Hello,World!")
    fmt.Println("Send Message Success!")
}
