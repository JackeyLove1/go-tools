package test

import (
    "testing"
    "time"

    "github.com/robfig/cron/v3"
)

// https://github.com/robfig/cron/blob/master/doc.go
func TestCron(t *testing.T) {
    c := cron.New()
    const msgNum = 5
    msgQueue := make([]string, 0, 10)
    _, err := c.AddFunc("@every 1s", func() {
        msgQueue = append(msgQueue, "hello")
    })
    if err != nil {
        t.Errorf("failed to init cron task")
    }
    quit := make(chan struct{})
    c.Start()
    go func() {
        defer c.Stop()
        time.Sleep(msgNum * time.Second)
        quit <- struct{}{}
    }()
    <-quit
    if len(msgQueue) != msgNum {
        t.Errorf("msgQueue length is %d, want %d", len(msgQueue), msgNum)
    }
}
