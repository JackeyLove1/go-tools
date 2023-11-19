package dao

import (
    "sync"

    "gorm.io/gorm"
    "ticktok/init"
)

var (
    db     *gorm.DB
    dbOnce sync.Once
)

func DaoInitialization() {
    dbOnce.Do(func() {
        db = init.GetDB()
        initKafkaClient()
    })
}
