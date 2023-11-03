package RabbitMQ

import (
    "fmt"
    "log"

    "github.com/streadway/amqp"
)

const MQURL = "amqp://guest:guest@localhost:5672/"

func FailOnErr(err error, msg string) {
    if err != nil {
        log.Fatalf("%s:%s", msg, err)
    }
}

type RabbitMQ struct {
    Conn      *amqp.Connection
    Channel   *amqp.Channel
    QueueName string
    Exchange  string
    Key       string
    Mqurl     string
}

func NewRabbitMQ(queueName string, exchange string, routeKey string) (*RabbitMQ, error) {
    rabbit := &RabbitMQ{
        QueueName: queueName,
        Exchange:  exchange,
        Key:       routeKey,
        Mqurl:     MQURL,
    }
    var err error
    rabbit.Conn, err = amqp.Dial(rabbit.Mqurl)
    if err != nil {
        return nil, err
    }
    rabbit.Channel, err = rabbit.Conn.Channel()
    if err != nil {
        return nil, err
    }
    return rabbit, nil
}

func (r *RabbitMQ) Destory() {
    r.Channel.Close()
    r.Conn.Close()
}

// return nil not panic
func (r *RabbitMQ) failOnErr(err error, message string) {
    if err != nil {
        log.Fatalf("%s:%s", message, err)
    }
}

func NewRabbitMQSimple(queueName string) (*RabbitMQ, error) {
    return NewRabbitMQ(queueName, "", "")
}

func (r *RabbitMQ) PublishSimple(message string) error {
    _, err := r.Channel.QueueDeclare(
        r.QueueName,
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }
    r.Channel.Publish(
        r.Exchange,
        r.QueueName,
        false,
        false,
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(message),
        })
    return nil
}

//simple 模式下消费者
func (r *RabbitMQ) ConsumeSimple() {
    //1.申请队列，如果队列不存在会自动创建，存在则跳过创建
    q, err := r.Channel.QueueDeclare(
        r.QueueName,
        //是否持久化
        false,
        //是否自动删除
        false,
        //是否具有排他性
        false,
        //是否阻塞处理
        false,
        //额外的属性
        nil,
    )
    if err != nil {
        fmt.Println(err)
    }

    //接收消息
    msgs, err := r.Channel.Consume(
        q.Name, // queue
        //用来区分多个消费者
        "", // consumer
        //是否自动应答
        true, // auto-ack
        //是否独有
        false, // exclusive
        //设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
        false, // no-local
        //列是否阻塞
        false, // no-wait
        nil,   // args
    )
    if err != nil {
        fmt.Println(err)
    }

    forever := make(chan bool)
    //启用协程处理消息
    go func() {
        for d := range msgs {
            //消息逻辑处理，可以自行设计逻辑
            log.Printf("Received a message: %s", d.Body)

        }
    }()

    log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
    <-forever

}

//订阅模式创建RabbitMQ实例
func NewRabbitMQPubSub(exchangeName string) (*RabbitMQ, error) {
    return NewRabbitMQ("", exchangeName, "")
}

//订阅模式生产
func (r *RabbitMQ) PublishPub(message string) {
    //1.尝试创建交换机
    err := r.Channel.ExchangeDeclare(
        r.Exchange,
        "fanout",
        true,
        false,
        //true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
        false,
        false,
        nil,
    )

    r.failOnErr(err, "Failed to declare an excha"+
        "nge")

    //2.发送消息
    err = r.Channel.Publish(
        r.Exchange,
        "",
        false,
        false,
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(message),
        })
}

//订阅模式消费端代码
func (r *RabbitMQ) RecieveSub() {
    //1.试探性创建交换机
    err := r.Channel.ExchangeDeclare(
        r.Exchange,
        //交换机类型
        "fanout",
        true,
        false,
        //YES表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
        false,
        false,
        nil,
    )
    r.failOnErr(err, "Failed to declare an exch"+
        "ange")
    //2.试探性创建队列，这里注意队列名称不要写
    q, err := r.Channel.QueueDeclare(
        "", //随机生产队列名称
        false,
        false,
        true,
        false,
        nil,
    )
    r.failOnErr(err, "Failed to declare a queue")

    //绑定队列到 exchange 中
    err = r.Channel.QueueBind(
        q.Name,
        //在pub/sub模式下，这里的key要为空
        "",
        r.Exchange,
        false,
        nil)

    //消费消息
    messges, err := r.Channel.Consume(
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
        for d := range messges {
            log.Printf("Received a message: %s", d.Body)
        }
    }()

    fmt.Println("退出请按 CTRL+C\n")
    <-forever
}

//路由模式
//创建RabbitMQ实例
func NewRabbitMQRouting(exchangeName string, routingKey string) (*RabbitMQ, error) {
    return NewRabbitMQ("", exchangeName, routingKey)
}

//路由模式发送消息
func (r *RabbitMQ) PublishRouting(message string) {
    //1.尝试创建交换机
    err := r.Channel.ExchangeDeclare(
        r.Exchange,
        //要改成direct
        "direct",
        true,
        false,
        false,
        false,
        nil,
    )

    r.failOnErr(err, "Failed to declare an excha"+
        "nge")

    //2.发送消息
    err = r.Channel.Publish(
        r.Exchange,
        //要设置
        r.Key,
        false,
        false,
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(message),
        })
}

//路由模式接受消息
func (r *RabbitMQ) RecieveRouting() {
    //1.试探性创建交换机
    err := r.Channel.ExchangeDeclare(
        r.Exchange,
        //交换机类型
        "direct",
        true,
        false,
        false,
        false,
        nil,
    )
    r.failOnErr(err, "Failed to declare an exch"+
        "ange")
    //2.试探性创建队列，这里注意队列名称不要写
    q, err := r.Channel.QueueDeclare(
        "", //随机生产队列名称
        false,
        false,
        true,
        false,
        nil,
    )
    r.failOnErr(err, "Failed to declare a queue")

    //绑定队列到 exchange 中
    err = r.Channel.QueueBind(
        q.Name,
        //需要绑定key
        r.Key,
        r.Exchange,
        false,
        nil)

    //消费消息
    messges, err := r.Channel.Consume(
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
        for d := range messges {
            log.Printf("Received a message: %s", d.Body)
        }
    }()

    fmt.Println("退出请按 CTRL+C\n")
    <-forever
}

//话题模式
//创建RabbitMQ实例
func NewRabbitMQTopic(exchangeName string, routingKey string) (*RabbitMQ, error) {
    return NewRabbitMQ("", exchangeName, routingKey)
}

//话题模式发送消息
func (r *RabbitMQ) PublishTopic(message string) {
    //1.尝试创建交换机
    err := r.Channel.ExchangeDeclare(
        r.Exchange,
        //要改成topic
        "topic",
        true,
        false,
        false,
        false,
        nil,
    )

    r.failOnErr(err, "Failed to declare an excha"+
        "nge")

    //2.发送消息
    err = r.Channel.Publish(
        r.Exchange,
        //要设置
        r.Key,
        false,
        false,
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(message),
        })
}

//话题模式接受消息
//要注意key,规则
//其中“*”用于匹配一个单词，“#”用于匹配多个单词（可以是零个）
//匹配 imooc.* 表示匹配 imooc.hello, 但是imooc.hello.one需要用imooc.#才能匹配到
func (r *RabbitMQ) RecieveTopic() {
    //1.试探性创建交换机
    err := r.Channel.ExchangeDeclare(
        r.Exchange,
        //交换机类型
        "topic",
        true,
        false,
        false,
        false,
        nil,
    )
    r.failOnErr(err, "Failed to declare an exch"+
        "ange")
    //2.试探性创建队列，这里注意队列名称不要写
    q, err := r.Channel.QueueDeclare(
        "", //随机生产队列名称
        false,
        false,
        true,
        false,
        nil,
    )
    r.failOnErr(err, "Failed to declare a queue")

    //绑定队列到 exchange 中
    err = r.Channel.QueueBind(
        q.Name,
        //在pub/sub模式下，这里的key要为空
        r.Key,
        r.Exchange,
        false,
        nil)

    //消费消息
    messges, err := r.Channel.Consume(
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
        for d := range messges {
            log.Printf("Received a message: %s", d.Body)
        }
    }()

    fmt.Println("退出请按 CTRL+C\n")
    <-forever
}
