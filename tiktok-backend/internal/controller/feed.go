package controller

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "tiktok-backend/api"
)

type FeedResponse struct {
    api.Response
    Video    []api.Video `json:"video_list,omitempty"`
    NextTime int64       `json:"next_time,omitempty"`
}

func Feed(c *gin.Context) {
    token := c.Query("token")
    if token != "" {
        c.JSON(http.StatusOK, FeedResponse{
            Response: api.Response{StatusCode: 0},
            Video:    DemoVideos,
            NextTime: 0,
        })
        return
    }
}
