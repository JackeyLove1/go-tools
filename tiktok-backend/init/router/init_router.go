package router

import (
    "github.com/gin-gonic/gin"
    "tiktok-backend/internal/controller"
)

func InitRouter(r *gin.Engine) {
    // 用户注册与登录需要进行鉴权, Feed可授权可不授权
    r.POST("/douyin/user/register/", controller.Register)
    // r.POST("/douyin/user/login/", jwt.JwtMiddleware.LoginHandler)
    r.GET("/douyin/feed/", controller.Feed)

    // 鉴权authorization
    // auth := r.Group("/douyin", jwt.JwtMiddleware.MiddlewareFunc())
    auth := r.Group("/douyin")
    // basic apis
    auth.GET("/user/", controller.UserInfo)
    auth.POST("/publish/action/", controller.Publish)
    auth.GET("/publish/list/", controller.PublishList)
}
