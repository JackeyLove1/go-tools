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
    // url := "http://localhost:8000/.cache/torch/text/datasets/WikiText2/wikitext-2/wiki.test.tokens"
    // url := "http://localhost:6789/proxy/download/localhost:8000/.cache/torch/text/datasets/WikiText2/wikitext-2/wiki.test.tokens"
    url := "http://localhost:6789/proxy/download/localhost:8000/Desktop/tmp/randomfile"
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
