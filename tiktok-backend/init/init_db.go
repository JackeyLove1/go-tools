package init

import (
    "fmt"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "gorm.io/gorm/schema"
)

var db *gorm.DB

func InitDB() {
    stdOutLogger.Print("In InitDataBase")
    dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        dbUser,
        dbPassWord,
        dbHost,
        dbPort,
        dbName)
    var err error
    logLevelMap := map[string]logger.LogLevel{
        "silent": logger.Silent,
        "error":  logger.Error,
        "warn":   logger.Warn,
        "info":   logger.Info,
    }

    db, err = gorm.Open(mysql.Open(dns), &gorm.Config{
        DisableForeignKeyConstraintWhenMigrating: true,
        NamingStrategy: schema.NamingStrategy{
            SingularTable: false,
            TablePrefix:   "t_ticktok_",
        },
        SkipDefaultTransaction: true,
        Logger:                 logger.Default.LogMode(logLevelMap[dbLogLevel]),
    })
    if err != nil {
        stdOutLogger.Panic().Caller().Str("failed to init db", err.Error())
    }

    err = db.AutoMigrate()
}
