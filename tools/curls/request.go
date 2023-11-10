package curls

import (
    "io"
    "net/http"
    "net/http/cookiejar"
    "time"
)

type Options struct {
    Headers        map[string]any
    BaseURL        string
    FormParams     map[string]any
    JSON           any
    XML            string
    Timeout        float32
    timeout        time.Duration
    Cookies        any
    Proxy          string
    SetRespCharset string
}

type Request struct {
    opts       Options
    cli        *http.Client
    req        *http.Request
    body       io.Reader
    cookiesJar *cookiejar.Jar
}

func (r *Request) Get(uri string, opts ...Options) () {
    panic("unimplemented error")
}
