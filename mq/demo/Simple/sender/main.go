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

func main() {
    r, err := RabbitMQ.NewRabbitMQ("test", "", "")
    failOnErr(err, "")
    // declare a queue if not existed
    _, err = r.Channel.QueueDeclare(
        r.QueueName,
        false,
        false,
        false,
        false,
        nil,
    )
    failOnErr(err, "Failed to declare a queue")
    for i := 0; i < 100; i++ {
        data := RabbitMQ.GenerateRandomNumber()
        dataBytes, err := json.Marshal(data)
        failOnErr(err, "")
        r.Channel.Publish(
            r.Exchange,
            r.QueueName,
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
