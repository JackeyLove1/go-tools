package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"

    "github.com/gin-gonic/gin"
)

//TODO: add jwt token

const MaxFileSize = 1024 * 1024 * 1024

var FilePath = "server_data"

func failOnErr(err error, msg string) {
    if err != nil {
        log.Println(msg)
        panic(err)
    }
}

func init() {
    currDir, err := os.Getwd()
    if err != nil {
        failOnErr(err, "failed to get current dir")
    }
    FilePath = filepath.Join(currDir, FilePath)
    err = os.MkdirAll(FilePath, 0755)
    if err != nil {
        failOnErr(err, fmt.Sprintf("failed to create dir: %s, err: %w", FilePath, err))
    }
}

func FileLimiter() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxFileSize)
        c.Next()
    }
}

func main() {
    router := gin.Default()
    router.Use(gin.Recovery())
    router.Use(gin.Logger())
    router.Use(FileLimiter())
    router.StaticFS("/", gin.Dir("/", true))
    router.POST("/upload", func(c *gin.Context) {
        file, err := c.FormFile("file")
        if err != nil {
            c.String(http.StatusBadRequest, "form file err: %s", err.Error())
            return
        }
        filename := filepath.Base(file.Filename)
        savePath := filepath.Join(FilePath, filename)
        log.Println("savePath: ", savePath)
        if err = c.SaveUploadedFile(file, savePath); err != nil {
            c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
            return
        }
        c.String(http.StatusOK, "Uploaded successfully files with fields")
    })
    router.Run(":8080")
}
