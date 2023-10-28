package RabbitMQ

import (
    "fmt"
    "log"

    "github.com/streadway/amqp"
)

const MQURL = "amqp://guest:guest@localhost:5672/"

type RabbitMQ struct {
    conn      *amqp.Connection
    channel   *amqp.Channel
    QueueName string
    Exchange  string
    Key       string
    MQurl     string
}

func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
    return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, MQurl: MQURL}
}

func (r *RabbitMQ) failOnErr(err error, msg string) {
    if err != nil {
        log.Fatalf("%s:%s", msg, err)
        panic(fmt.Sprintf("%s:%s", msg, err))
    }
}

func (r *RabbitMQ) Destory() {
    _ = r.channel.Close()
    _ = r.conn.Close()
}

// create simple mode mq
func NewRabbitMQSimple(queueName string) *RabbitMQ {
    r := NewRabbitMQ(queueName, "", "")
    var err error
    r.conn, err = amqp.Dial(r.MQurl)
    r.failOnErr(err, fmt.Sprintf("failed to connect rabbitmq, url:%s", r.MQurl))
    r.channel, err = r.conn.Channel()
    r.failOnErr(err, "failed to open a channel")
    return r
}

// Simple mode publisher
func (r *RabbitMQ) PublishSimple(message string) {
    // 1. create r.QueueName queue, if exist then skip
    _, err := r.channel.QueueDeclare(
        r.QueueName,
        // durable: flush to the disk
        false,
        // delete when unused
        false,
        // exclusive: only one consumer
        false,
        // non-blocking
        false,
        // other arguments
        nil,
    )
    r.failOnErr(err, "failed to declare a queue")

    // 2. publish message
    r.channel.Publish(
        r.Exchange,
        r.QueueName,
        // 强制退货，如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
        false,
        // 强制返还，如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
        false,
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(message),
        })
}

func (r *RabbitMQ) ConsumeSimple() {
    // 1. apply for the queue
    q, err := r.channel.QueueDeclare(
        r.QueueName,
        false,
        false,
        false,
        false,
        nil,
    )
    r.failOnErr(err, "failed to declare a queue")

    // 2. receive message
    messages, err := r.channel.Consume(
        q.Name,
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    r.failOnErr(err, "failed to consume")

    // 3. handle message
    forever := make(chan bool)
    go func() {
        for d := range messages {
            s := string(d.Body)
            log.Printf("Received a message: %s", s)
        }
    }()
    log.Println(" [*] Waiting for messages. To exit press CTRL+C")
    <-forever
}

// create pub/sub
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
    r := NewRabbitMQ("", exchangeName, "")
    var err error
    r.conn, err = amqp.Dial(r.MQurl)
    r.failOnErr(err, fmt.Sprintf("failed to connect rabbitmq, url:%s", r.MQurl))
    r.channel, err = r.conn.Channel()
    r.failOnErr(err, "failed to open a channel")
    return r
}

// publish
func (r *RabbitMQ) PublishPub(message string) {
    // 1. create exchange if not exist
    err := r.channel.ExchangeDeclare(
        r.Exchange,
        "fanout",
        true,
        false,
        false,
        false,
        nil,
    )
    r.failOnErr(err, "failed to declare an exchange")

    // 2. publish message
    err = r.channel.Publish(
        r.Exchange,
        "",
        false,
        false,
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(message),
        },
    )
    r.failOnErr(err, "failed to publish a message")
}

// subscribe
func (r *RabbitMQ) RecieveSub() {
    // 1. try to create exchange if not exist
    err := r.channel.ExchangeDeclare(
        r.Exchange,
        "fanout",
        true,
        false,
        false,
        false,
        nil,
    )
    r.failOnErr(err, "failed to declare an exchange")

    // 2. try to create queue if not exist
    q, err := r.channel.QueueDeclare(
        "",
        false,
        false,
        true,
        false,
        nil,
    )
    r.failOnErr(err, "failed to declare a queue")

    // 3. bind queue
    err = r.channel.QueueBind(
        q.Name,
        "",
        r.Exchange,
        false,
        nil,
    )

    // 4. receive message
    messages, err := r.channel.Consume(
        q.Name,
        "",
        true,
        false,
        false,
        false,
        nil,
    )

    forever := make(chan bool)
    go func() {
        for d := range messages {
            s := string(d.Body)
            log.Printf("Received a message: %s", s)
        }
    }()

    fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
    <-forever
}

func NewRabbitMQRouting(exchange string, routingKey string) *RabbitMQ {
    r := NewRabbitMQ("", exchange, routingKey)
    var err error
    r.conn, err = amqp.Dial(r.MQurl)
    r.failOnErr(err, fmt.Sprintf("failed to connect rabbitmq, url:%s", r.MQurl))
    r.channel, err = r.conn.Channel()
    r.failOnErr(err, "failed to open a channel")
    return r
}

func (r *RabbitMQ) PublishRouting(message string) {
    err := r.channel.Publish(
        r.Exchange,
        r.Key,
        false,
        false,
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(message),
        },
    )
    r.failOnErr(err, "failed to publish a message")

    err = r.channel.Publish(
        r.Exchange,
        r.Key,
        false,
        false,
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(message),
        })
    r.failOnErr(err, "failed to publish a message")
}

func (r *RabbitMQ) RecieveRouting() {
    err := r.channel.ExchangeDeclare(
        r.Exchange,
        "direct",
        true,
        false,
        false,
        false,
        nil,
    )
    r.failOnErr(err, "failed to declare an exchange")

    q, err := r.channel.QueueDeclare(
        "",
        false,
        false,
        true,
        false,
        nil,
    )
    r.failOnErr(err, "failed to declare a queue")

    // bind queue
    err = r.channel.QueueBind(
        q.Name,
        r.Key,
        r.Exchange,
        false,
        nil,
    )

    // receive message
    messages, err := r.channel.Consume(
        q.Name,
        "",
        true,
        false,
        false,
        false,
        nil,
    )

    forever := make(chan bool)
    go func() {
        for d := range messages {
            log.Printf("Received a message: %s", d.Body)
        }
    }()

    fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
    <-forever
}
