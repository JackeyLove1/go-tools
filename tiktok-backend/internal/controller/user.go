package controller

import (
    "github.com/gin-gonic/gin"
    "tiktok-backend/api"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]api.User{
    "zhangleidouyin": {
        Id:            1,
        Name:          "zhanglei",
        FollowCount:   10,
        FollowerCount: 5,
        IsFollow:      true,
    },
}

// Register 处理用户登录请求的RPC远程调用
func Register(c *gin.Context) {

}

func UserInfo(c *gin.Context) {

}
