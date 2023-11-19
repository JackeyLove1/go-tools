package service

import (
    "sync"
    "time"

    "ticktok/api"
    "ticktok/internal/dao"
    "ticktok/internal/utils/constants"
    "ticktok/internal/utils/logger"
)

type feedService struct{}

var (
    feedServiceInstance *feedService
    feedOnce            sync.Once
)

func GetFeedServiceInstance() *feedService {
    initRedis()
    initKafka()
    feedOnce.Do(func() {
        feedServiceInstance = &feedService{}
    })
    return feedServiceInstance
}

//Feed service层获取视频流
func (f *feedService) Feed(userId int64, latestTime time.Time) (int64, []api.Video, error) {
    videos, err := dao.GetVideoDaoInstance().GetFeedList(latestTime)
    logger.GlobalLogger.Printf("get Videos From FeedList")
    if err != nil {
        logger.GlobalLogger.Printf("dao.NewVideoDaoInstance().GetLatest error: %s", err)
        return -1, nil, err
    }
    if len(videos) == 0 {
        logger.GlobalLogger.Printf("没有早于当前时间的视频")
        return -1, nil, constants.NoVideoErr
    }
    videoList, err := getVideoListByModel(userId, videos)
    if err != nil {
        return -1, nil, err
    }
    return videos[len(videos)-1].CreatedAt.UnixMilli(), videoList, nil
}
