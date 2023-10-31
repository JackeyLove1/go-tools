package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type Body struct {
    Price uint `json:"price" binding:"required,gte=10,lte=1000"`
}

func main() {
    engine := gin.New()
    engine.POST("/test", func(context *gin.Context) {
        body := Body{}
        if err := context.ShouldBind(&body); err != nil {
            context.AbortWithStatusJSON(http.StatusBadRequest,
                gin.H{
                    "error":   "VALIDATEERR-1",
                    "message": "Invalid inputs. Please check your inputs"})
            return
        }
        context.JSON(http.StatusAccepted, &body)
    })
    engine.Run(":3000")
}
