package main

import "go-tools/mq/RabbitMQ"

func main() {
    r := RabbitMQ.NewRabbitMQRouting("exec", "one")
    r.RecieveRouting()
}
