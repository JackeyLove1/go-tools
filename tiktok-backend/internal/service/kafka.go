package service

import (
    "sync"

    "github.com/IBM/sarama"
    initialization "ticktok/init"
)

var (
    kafkaServer sarama.SyncProducer
    kafkaOnce   sync.Once
)

func initKafka() {
    kafkaOnce.Do(func() {
        kafkaServer = initialization.GetKafkaServer()
    })
}
