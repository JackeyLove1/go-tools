package controller

import (
    "time"

    "tiktok-backend/api"
)

var DemoVideos = []api.Video{
    {
        Id:            1,
        Author:        DemoUser,
        PlayUrl:       "https://www.w3schools.com/html/movie.mp4",
        CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
        FavoriteCount: 0,
        CommentCount:  0,
        IsFavorite:    false,
    },
}

var DemoComments = []api.Comment{
    {
        Id:         1,
        User:       DemoUser,
        Content:    "Test Comment",
        CreateTime: time.Now(),
    },
}

var DemoUser = api.User{
    Id:            1,
    Name:          "TestUser",
    FollowCount:   0,
    FollowerCount: 0,
    IsFollow:      false,
}
