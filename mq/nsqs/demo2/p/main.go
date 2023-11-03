package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/nsqio/go-nsq"
)

type MyHandle struct {
    Title string
}

func (m *MyHandle) HandleMessage(msg *nsq.Message) error {
    fmt.Printf("%s recv from %v, msg:%v\n", m.Title, msg.NSQDAddress, string(msg.Body))
    return nil
}

func initConsumer(topic string, channel string, address string) error {
    config := nsq.NewConfig()
    config.LookupdPollInterval = 15 * time.Second
    c, err := nsq.NewConsumer(topic, channel, config)
    if err != nil {
        fmt.Printf("create consumer failed, err:%v\n", err)
        return err
    }
    consumer := &MyHandle{
        Title: "topic",
    }
    c.AddHandler(consumer)
    if err = c.ConnectToNSQLookupd(address); err != nil {
        return err
    }
    return nil
}

func main() {
    err := initConsumer("topic_demo", "first", "127.0.0.1:4161")
    if err != nil {
        fmt.Printf("init consumer failed, err:%v\n", err)
        return
    }
    c := make(chan os.Signal)
    signal.Notify(c, syscall.SIGINT)
    <-c
}
