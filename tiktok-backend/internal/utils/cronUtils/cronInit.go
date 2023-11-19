package cronUtils

import (
    "github.com/robfig/cron/v3"
)

// CronLab 分布式定时任务组件
var CronLab *cron.Cron

func InitCron() {
    CronLab = cron.New()
    CronLab.Start()
}
