package main

import (
    "context"
    "encoding/json"
    "log"
    "strconv"
    "sync"
    "time"

    "github.com/streadway/amqp"
    "go-tools/mq/demo/RabbitMQ"
)

func failOnErr(err error, msg string) {
    if err != nil {
        log.Fatalf("%s:%s", msg, err)
    }
}

const ExchangeName = "testTopic"
const RouteKeyPrefix = "test.topic.route"
const RouteKey1 = "test.topic.route1"
const RouteKey2 = "test.topic.route2"

func SendWorker(cxt context.Context, idx int, wg *sync.WaitGroup) {
    defer wg.Done()
    r, err := RabbitMQ.NewRabbitMQ("", ExchangeName, RouteKeyPrefix+strconv.Itoa(idx))
    failOnErr(err, "")

    // declare exchange
    err = r.Channel.ExchangeDeclare(
        r.Exchange,
        "topic",
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

func main() {
    ctx := context.Background()
    wg := &sync.WaitGroup{}
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go SendWorker(ctx, i, wg)
    }
    wg.Wait()
}
