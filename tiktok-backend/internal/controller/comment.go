package controller

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "tiktok-backend/api"
)

type CommentListResponse struct {
    api.Response
    CommentList []api.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
    api.Response
    Comment api.Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
    token := c.Query("token")
    actionType := c.Query("action_type")

    if user, exist := usersLoginInfo[token]; exist {
        if actionType == "1" {
            text := c.Query("comment_text")
            c.JSON(http.StatusOK, CommentActionResponse{
                Response: api.Response{StatusCode: 0},
                Comment: api.Comment{
                    Id:         1,
                    User:       user,
                    Content:    text,
                    CreateTime: time.Now(),
                },
            })
            return
        }
        c.JSON(http.StatusOK, api.Response{
            StatusCode: 0,
        })
    } else {
        c.JSON(http.StatusOK, api.Response{
            StatusCode: 1,
            StatusMsg:  "User doesn't exist",
        })
    }
}

func CommentList(c *gin.Context) {
    c.JSON(http.StatusOK, CommentListResponse{
        Response:    api.Response{StatusCode: 0},
        CommentList: DemoComments,
    })

}
