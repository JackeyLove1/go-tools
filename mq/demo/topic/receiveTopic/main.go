package main

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "go-tools/mq/demo/RabbitMQ"
)

const ExchangeName = "testTopic"
const RouteKeyPrefix = "test.topic.route"
const RouteKey1 = "test.topic.route1"
const RouteKey2 = "test.topic.route2"

var routes = []string{
    "#",
    "*.*.*",
    "#.*",
}

func ConsumeWorker(ctx context.Context, idx int) {
    r, err := RabbitMQ.NewRabbitMQ("", ExchangeName, routes[idx])
    RabbitMQ.FailOnErr(err, "")
    defer r.Destory()
    // try to get or declare exchange
    err = r.Channel.ExchangeDeclare(
        r.Exchange,
        "topic",
        true,
        false,
        false,
        false,
        nil,
    )
    RabbitMQ.FailOnErr(err, "")

    // declare a queue
    q, err := r.Channel.QueueDeclare(
        "",
        false,
        false,
        true,
        false,
        nil,
    )
    RabbitMQ.FailOnErr(err, "")
    r.QueueName = q.Name
    log.Printf("queue name:%s, route key:%s \n", r.QueueName, r.Key)

    // bind queue
    err = r.Channel.QueueBind(
        r.QueueName,
        r.Key,
        r.Exchange,
        false,
        nil,
    )
    RabbitMQ.FailOnErr(err, "")

    // consume
    // r.Channel.Qos(1, 0, false)
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
            log.Printf("exit\n")
            return
        case <-time.After(time.Second * 100):
            log.Printf("timeout\n")
            return
        case msg := <-msgs:
            var randomNum RabbitMQ.RandomNumber
            err := json.Unmarshal(msg.Body, &randomNum)
            RabbitMQ.FailOnErr(err, "")
            log.Printf("Consumer %d Received a message: %s", idx, randomNum.String())
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    for i := 0; i < 3; i++ {
        go ConsumeWorker(ctx, i)
    }
    time.Sleep(time.Second * 100)
}
