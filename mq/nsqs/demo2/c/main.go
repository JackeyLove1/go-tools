package main

import (
    "fmt"
    "strconv"

    "github.com/nsqio/go-nsq"
)

var producer *nsq.Producer

func initProducer(str string) (err error) {
    config := nsq.NewConfig()
    producer, err = nsq.NewProducer(str, config)
    if err != nil {
        panic(err)
    }
    return nil
}

func main() {
    nsqAddress := "127.0.0.1:4150"
    err := initProducer(nsqAddress)
    if err != nil {
        fmt.Printf("init producer failed, err:%v\n", err)
        return
    }

    for i := 0; i < 100; i++ {
        data := strconv.Itoa(i)
        err = producer.Publish("topic_demo", []byte(data))
        if err != nil {
            fmt.Printf("publish msg to nsq failed, err:%v\n", err)
            continue
        }
    }
}
