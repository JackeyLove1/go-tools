package main

import (
    "errors"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"

    "github.com/gin-gonic/gin"
)

func ParseAddress(address string) (string, string, error) {
    if strings.HasPrefix(address, "n") {
        return parseIPV4(address)
    }
    return parseIPV6(address)
}

func parseIPV6(address string) (string, string, error) {
    res := strings.Split(address, ":")
    if len(res) != 2 {
        return "", "", errors.New("invalid address")
    }
    ip, port := res[0], res[1]
    ipParts := strings.Split(ip, "-")
    for idx, ipPart := range ipParts {
        if idx > 0 {
            ipParts[idx] = strings.TrimPrefix(ipPart[1:], "0")
        }
        if idx == len(ipParts)-1 {
            ipParts[idx] = ":" + ipParts[idx]
        }
    }
    ip = "[fdbd:" + strings.Join(ipParts, ":") + "]"
    log.Println("ip address:", ip)
    return ip, port, nil
}

func parseIPV4(address string) (string, string, error) {
    address = strings.TrimPrefix(address, "n")
    res := strings.Split(address, ":")
    if len(res) != 2 {
        return "", "", errors.New("invalid address")
    }
    ip, port := res[0], res[1]
    ipParts := strings.Split(ip, "-")
    for idx, ipPart := range ipParts {
        ipParts[idx] = strings.TrimPrefix(ipPart, "0")
    }
    ip = "10." + strings.Join(ipParts, ".")
    log.Println("ip address:", ip)
    return ip, port, nil
}

func main() {
    router := gin.Default()

    // This handler will match /url:port/*path
    router.Any("/proxy/http/:address/*path", func(c *gin.Context) {
        // Extract the parameters from the path
        targetAddress := c.Param("address")
        targetPath := c.Param("path")
        targetURL, targetPort, err := ParseAddress(targetAddress)
        if err != nil {
            log.Printf("Error parsing the target URL: %v", err)
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
        // Construct the destination URL
        destinationURL := "http://" + targetURL + ":" + targetPort + targetPath
        desURLWithoutPath := "http://" + targetURL + ":" + targetPort
        log.Println("des:", destinationURL)
        // Parse the destination URL
        url, err := url.Parse(destinationURL)
        log.Println("url:", url.String())
        if err != nil {
            log.Printf("Error parsing the target URL: %v", err)
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
        urlNoPath, err := url.Parse(desURLWithoutPath)
        log.Println("url no path:", urlNoPath.String())
        if err != nil {
            log.Printf("Error parsing the target URL: %v", err)
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
        // Create the reverse proxy
        proxy := httputil.NewSingleHostReverseProxy(urlNoPath)

        // Update the request to reflect the target's URL and scheme
        c.Request.URL = url
        c.Request.URL.Scheme = url.Scheme
        c.Request.URL.Host = url.Host
        c.Request.Header.Set("X-Forwarded-Host", c.Request.Host)
        c.Request.Host = url.Host

        // Use the proxy to forward the request
        proxy.ServeHTTP(c.Writer, c.Request)
    })

    // Run the server
    router.Run(":8080")
}
