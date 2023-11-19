package test

import (
    "fmt"
    "testing"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/schema"
)

type TestUser struct {
    gorm.Model
    Name  string `gorm:"type:varchar(255);not null"`
    Email string `gorm:"type:varchar(255);not null;uniqueIndex:idx_user_email"`
}

var db *gorm.DB

func init() {
    var err error
    dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        "root",
        "",
        "localhost",
        "3306",
        "test")
    db, err = gorm.Open(mysql.Open(dns), &gorm.Config{
        DisableForeignKeyConstraintWhenMigrating: true,
        NamingStrategy: schema.NamingStrategy{
            SingularTable: false,
            TablePrefix:   "test_ticktok_",
        },
        SkipDefaultTransaction: true,
    })
    FailOnErr(err)

    err = db.Migrator().DropTable(&TestUser{})
    FailOnErr(err)

    err = db.AutoMigrate(&TestUser{})
    FailOnErr(err)
}

func TestUniqueIndex(t *testing.T) {
    user1 := TestUser{
        Name:  "John Wick",
        Email: "john@test.com",
    }
    result := db.Create(&user1)
    if result.Error != nil {
        t.Errorf("failed to create db, err:%s", result.Error.Error())
        t.FailNow()
    }
    user2 := TestUser{
        Name:  "John Wick",
        Email: "john@test.com",
    }
    result = db.Create(&user2)
    if result.Error == nil {
        t.Errorf("shouldn't insert two same email entry")
        t.FailNow()
    }
}
