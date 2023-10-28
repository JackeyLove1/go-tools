package main

import "go-tools/mq/RabbitMQ"

func main() {
    r := RabbitMQ.NewRabbitMQSimple("test")
    r.ConsumeSimple()
}
