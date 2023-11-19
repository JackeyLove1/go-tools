package oss

import (
    "io"
    "sync"

    "github.com/aliyun/aliyun-oss-go-sdk/oss"
    initialization "ticktok/init"
)

var (
    bucket     *oss.Bucket = nil
    bucketOnce sync.Once
)

func initBucket() {
    bucketOnce.Do(func() {
        bucket = initialization.GetBucket()
    })
}

func UploadFromFile(ossPath, localFilePath string) error {
    initBucket()
    return bucket.PutObjectFromFile(ossPath, localFilePath)
}

func UploadFromReader(ossPath string, srcReader io.Reader) error {
    initBucket()
    return bucket.PutObject(ossPath, srcReader)
}
