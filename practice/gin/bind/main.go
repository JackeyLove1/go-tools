package main

import (
    "fmt"

    "github.com/gin-gonic/gin"
)

type User struct {
    ID   string `json:"id" uri:"id" binding:"required"`
    Name string `json:"name" binding:"required"`
    Age  int    `json:"age" binding:"required"`
}

func (u *User) String() string {
    return fmt.Sprintf("ID:%s, Name:%s, Age:%d", u.ID, u.Name, u.Age)
}

type UserUri struct {
    ID string `uri:"id" binding:"required"`
}

func GetUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"err": "invalid parameter"})
        return
    }
    c.JSON(200, gin.H{"user": user.String()})
}

func GetUserUri(c *gin.Context) {
    var user UserUri
    if err := c.ShouldBindUri(&user); err != nil {
        c.JSON(400, gin.H{"err": "invalid parameter"})
        return
    }
    c.JSON(200, gin.H{"user": user.ID})

}

func main() {
    r := gin.Default()
    r.Any("/user/:id", GetUserUri)
    r.Any("/user", GetUser)
    r.Run()
}
