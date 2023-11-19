package dao

import (
    "sync"

    "gorm.io/gorm"
    "tiktok-backend/init"
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
