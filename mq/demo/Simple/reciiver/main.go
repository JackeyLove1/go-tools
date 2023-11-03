package main

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "go-tools/mq/demo/RabbitMQ"
)

func ConsumeWorker(ctx context.Context, idx int) {
    r, err := RabbitMQ.NewRabbitMQ("test", "", "")
    RabbitMQ.FailOnErr(err, "")
    // declare a queue
    _, err = r.Channel.QueueDeclare(
        r.QueueName,
        false,
        false,
        false,
        false,
        nil,
    )
    RabbitMQ.FailOnErr(err, "")
    // consume
    r.Channel.Qos(1, 0, false)
    msgs, err := r.Channel.Consume(
        r.QueueName,
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    RabbitMQ.FailOnErr(err, "")

    for {
        select {
        case <-ctx.Done():
            log.Printf("Consumer%d timeout\n", idx)
            return
        case msg := <-msgs:
            var randomNum RabbitMQ.RandomNumber
            err := json.Unmarshal(msg.Body, &randomNum)
            RabbitMQ.FailOnErr(err, "")
            log.Printf("%d consumer Received a message: %s", idx, randomNum.String())
        }
    }
}

func main() {
    ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
    for i := 0; i < 3; i++ {
        go ConsumeWorker(ctx, i)
    }
    for {
        select {
        case <-ctx.Done():
            log.Printf("Timeout")
            return
        }
    }
}
