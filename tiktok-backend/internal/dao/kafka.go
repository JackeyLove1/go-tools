package dao

import (
    "sync"

    "github.com/IBM/sarama"
    "ticktok/init"
)

var (
    kafkaClient sarama.Consumer
    kafkaOnce   sync.Once
)

func initKafkaClient() {
    kafkaOnce.Do(func() {
        kafkaClient = init.GetKafkaClient()
        go func() {
            for {
                err := GetFavoriteDaoInstance().getFromMessageQueue()
                if err == nil {
                    break
                }
            }
        }()
    })
}
