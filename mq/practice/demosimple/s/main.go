package main

import (
    "context"
    "log"
    "time"

    amqp "github.com/rabbitmq/amqp091-go"
)

func failOnErr(err error, msg string) {
    if err != nil {
        log.Fatalf("msg:%s, err:%w", msg, err)
    }
}

const QueueName = "test"

func main() {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnErr(err, "failed to connect rabbitmq")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnErr(err, "failed to create channel")
    defer ch.Close()

    q, err := ch.QueueDeclare(
        QueueName,
        false,
        false,
        false,
        false,
        nil,
    )
    failOnErr(err, "failed to declare a queue")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    body := "Hello, World!"
    err = ch.PublishWithContext(ctx,
        "",
        q.Name,
        false,
        false,
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(body),
        })
    failOnErr(err, "failed to pushlish msg")
    log.Printf(" [x] Sent %s\n", body)
}
