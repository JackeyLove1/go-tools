package main

import (
    "fmt"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

const (
    DBName = "test2"
)

var db *gorm.DB
var err error

func failOnErr(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
    dsn := fmt.Sprintf("root:@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", DBName)
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    failOnErr(err)

}
