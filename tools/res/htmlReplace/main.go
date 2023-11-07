package main

import (
    "log"
    "net/http"
    "strings"

    "github.com/PuerkitoBio/goquery"
)

var NasPorts = []string{"8500", "8600", "8700"}

func replaceLinks(href string) string {
    parts := strings.Split(href, "/")
    if len(parts) <= 3 {
        return href
    }
    address := parts[3]
    path := "/" + strings.Join(parts[4:], "/")

    println("address:", address, " path:", path)
    return href
}

func main() {
    url := "http://example.com" // Replace with the URL of the webpage you want to scrape

    // Make a HTTP GET request to the webpage
    response, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    // Create a goquery document from the response body
    doc, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        log.Fatal(err)
    }

    // Find all anchor tags (a) that contain links
    doc.Find("a").Each(func(i int, s *goquery.Selection) {
        // Get the href attribute of each link
        href, exists := s.Attr("href")
        if exists {
            replaceLinks(href)
            s.SetAttr("href", "jacky!")
        }
    })
    // println(doc.Text())
}
