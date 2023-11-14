package main

import (
    "errors"
    "io"
    "net/http"
    "net/url"
    "os"
    "strings"
)

func parseFileNameFromUrl(remote string) (string, error) {
    decodedUrl, err := url.QueryUnescape(remote)
    if err != nil {
        return "", nil
    }
    fileUrl, err := url.Parse(decodedUrl)
    if err != nil {
        return "", nil
    }
    parts := strings.Split(fileUrl.Path, "/")
    if len(parts) == 0 {
        return "", errors.New("invalid url path")
    }
    fileName := strings.Replace(parts[len(parts)-1], " ", "-", -1)
    return fileName, nil
}

func main() {
    url := "http://localhost:8000/Desktop/%E5%9C%A8%20Apple%20Silicon%20Mac%20%E4%B8%8A%E5%85%A5%E9%97%A8%E6%B1%87%E7%BC%96%E8%AF%AD%E8%A8%80.pdf"
    resp, err := http.Get(url)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    fileName, err := parseFileNameFromUrl(url)
    if err != nil {
        panic(err)
    }
    output, err := os.Create(fileName)
    if err != nil {
        panic(err)
    }
    defer output.Close()
    _, err = io.Copy(output, resp.Body)
    if err != nil {
        panic(err)
    }
    println("Finish download:", fileName)
}
