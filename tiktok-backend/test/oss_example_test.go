package test

import (
    "context"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "testing"

    "github.com/tencentyun/cos-go-sdk-v5"
)

func TestQCloud(t *testing.T) {
    // Set your bucket URL and credentials
    u, _ := url.Parse("https://<your-bucket-name>.cos.<region>.myqcloud.com")
    b := &cos.BaseURL{BucketURL: u}
    c := cos.NewClient(b, &http.Client{
        Transport: &cos.AuthorizationTransport{
            SecretID:  "<your-secret-id>",  // Replace with your SecretID
            SecretKey: "<your-secret-key>", // Replace with your SecretKey
        },
    })
    // Example: Upload a file to the bucket
    name := "Hello.txt"
    _, _, err := c.Object.Upload(
        context.Background(), name, name, nil,
    )
    if err != nil {
        log.Panicf("failed to upload file, err:%s", err.Error())
    }
    fmt.Println("File uploaded successfully")

    // Example: Download a file from the bucket
    resp, err := c.Object.Get(context.Background(), name, nil)
    if err != nil {
        log.Panicf("failed to get file, err:%s", err.Error())
    }
    defer resp.Body.Close()
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Panicf("failed to read resp body, err:%s", err.Error())
    }
    println("succeed to read data:", string(data))
    // More operations (listing, deleting, etc.) can be performed in a similar manner
}
