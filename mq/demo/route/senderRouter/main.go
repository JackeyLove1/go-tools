package main

import (
    "encoding/json"
    "log"
    "time"

    "github.com/streadway/amqp"
    "go-tools/mq/demo/RabbitMQ"
)

func failOnErr(err error, msg string) {
    if err != nil {
        log.Fatalf("%s:%s", msg, err)
    }
}

const ExchangeName = "testRoute"
const RouteKey1 = "route1"

func main() {
    r, err := RabbitMQ.NewRabbitMQ("", ExchangeName, RouteKey1)
    failOnErr(err, "")

    // declare exchange
    err = r.Channel.ExchangeDeclare(
        r.Exchange,
        "direct",
        true,
        false,
        false,
        false,
        nil,
    )
    failOnErr(err, "")

    for i := 0; i < 100; i++ {
        data := RabbitMQ.GenerateRandomNumber()
        dataBytes, err := json.Marshal(data)
        failOnErr(err, "")
        r.Channel.Publish(
            r.Exchange,
            r.Key,
            false,
            false,
            amqp.Publishing{
                ContentType: "text/plain",
                Body:        dataBytes,
            })
        log.Printf("[send] %s", data.String())
        time.Sleep(time.Millisecond * 50)
    }
}
