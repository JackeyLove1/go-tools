package main

import "go-tools/mq/RabbitMQ"

func main() {
    r := RabbitMQ.NewRabbitMQPubSub("newProduct")
    r.RecieveSub()
}
