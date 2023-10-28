package main

import (
    "fmt"
    "log"
    "time"

    "github.com/streadway/amqp"
)

const (
    queueName = "test"
    exchange  = ""
    mqurl     = "amqp://guest:guest@localhost:5672/"
)

var (
    conn    *amqp.Connection
    channel *amqp.Channel
    count   int = 0
)

func failOnErr(err error, msg string) {
    if err != nil {
        log.Fatalf("%s:%s", msg, err)
        panic(fmt.Sprintf("%s:%s", msg, err))
    }
}

func Connect() {
    var err error
    conn, err = amqp.Dial(mqurl)
    failOnErr(err, "failed to connect rabbitmq")
    channel, err = conn.Channel()
    failOnErr(err, "failed to open a channel")
}

func Close() {
    channel.Close()
    conn.Close()
}

func push() {
    if channel == nil {
        Connect()
    }
    message := "Hello,World!"
    q, err := channel.QueueDeclare(
        queueName,
        false,
        false,
        false,
        false,
        nil)
    failOnErr(err, "failed to declare a queue")
    err = channel.Publish(
        exchange,
        q.Name,
        false,
        false,
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(message),
        },
    )
}

func receive() {
    if channel == nil {
        Connect()
    }

    q, err := channel.QueueDeclare(
        queueName,
        false,
        false,
        false,
        false,
        nil,
    )
    failOnErr(err, "failed to declare a queue")

    msg, err := channel.Consume(
        q.Name,
        "",
        false,
        false,
        false,
        false,
        nil)
    failOnErr(err, "failed to register a consumer")

    msgForever := make(chan bool)

    go func() {
        for d := range msg {
            s := string(d.Body)
            count++
            log.Printf("Received count:%d,  message: %s", count, s)
        }
    }()
    fmt.Println("退出请按 CTRL+C\n")
    <-msgForever
}

func main() {
    Connect()
    go func() {
        for {
            push()
            time.Sleep(5 * time.Second)
        }
    }()

    receive()

    fmt.Println("生产消费完成")
    Close()
}
