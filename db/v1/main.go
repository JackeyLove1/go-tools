package main

import (
    "errors"
    "fmt"
    "log"
    "time"

    "github.com/brianvoe/gofakeit/v6"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Name     string    `gorm:"column:name;type:varchar(256);not null"`
    Birthday time.Time `gorm:"column:birthday;type:datetime;not null"`
    Age      uint16    `gorm:"column:age;type:smallint unsigned;not null"`
}

func (u *User) String() string {
    return fmt.Sprintf("UserInfos. Name: %s, Age: %d, birthday:%s", u.Name, u.Age, u.Birthday.String())
}

type Product struct {
    gorm.Model
    Code  string
    Price uint
}

func init() {
    gofakeit.Seed(time.Now().Unix())
}

var db *gorm.DB

func printResult(result *gorm.DB) {
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            println("ErrRecordNotFound")
            return
        }
        println("Error:", result.Error.Error())
    }
    println("RowsAffected:", result.RowsAffected)
}

func InitDB() {
    var err error
    db, err = gorm.Open(mysql.New(mysql.Config{
        DSN:                       "root:@tcp(127.0.0.1:3306)/test2?charset=utf8&parseTime=True&loc=Local", // data source name
        DefaultStringSize:         256,                                                                     // default size for string fields
        DisableDatetimePrecision:  true,                                                                    // disable datetime precision, which not supported before MySQL 5.6
        DontSupportRenameIndex:    true,                                                                    // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
        DontSupportRenameColumn:   true,                                                                    // `change` when rename column, rename column not supported before MySQL 8, MariaDB
        SkipInitializeWithVersion: false,                                                                   // auto configure based on currently MySQL version
    }), &gorm.Config{})
    if err != nil {
        panic(err)
    }

    db.AutoMigrate(&Product{})
    db.AutoMigrate(&User{})
}

func InsertUser() {
    user1 := User{
        Name:     gofakeit.Name(),
        Birthday: gofakeit.Date(),
        Age:      gofakeit.Uint16(),
    }
    result := db.Create(&user1)
    if result.Error != nil {
        log.Fatalln(result.Error)
    }
}

func BatchInsert() {
    batchSize := 5
    users := make([]User, 0, batchSize)
    for i := 0; i < batchSize; i++ {
        users = append(users, User{
            Name:     gofakeit.Name(),
            Birthday: gofakeit.Date(),
            Age:      gofakeit.Uint16(),
        })
    }
    result := db.Create(&users)
    printResult(result)
}

func TakeFirst() {
    // select * from users order by id limit 1;
    user := User{}
    result := db.First(&user)
    printResult(result)
    println(user.String())
}

func TakeAll() {
    user := make([]User, 0)
    result := db.Find(&user)
    printResult(result)
    for _, u := range user {
        println(u.String())
    }
}

func TakeOne() {
    var user User
    result := db.Where("id = ?", "1").Find(&user)
    printResult(result)
    println(user.String())
}

func RawTake() {
    user := make([]User, 0)
    result := db.Raw("select * from users").Scan(&user)
    printResult(result)
    for _, u := range user {
        println(u.String())
    }
}

func TakeOr() {
    user := make([]User, 0)
    result := db.Where("id = ?", "1").Or("age = ?", "30787").Find(&user)
    printResult(result)
    for _, u := range user {
        println(u.String())
    }
}

func SelectFields() {

    db.Select("name", "age")
}
func main() {
    InitDB()
    InsertUser()
    // BatchInsert()
    // TakeFirst()
    // TakeAll()
    TakeOne()
    RawTake()
    TakeOr()
}
