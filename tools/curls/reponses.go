package curls

import (
    "net/http"
    "net/http/cookiejar"
)

type Response struct {
    resp           *http.Response
    req            *http.Request
    cookiesJar     *cookiejar.Jar
    err            error
    setRespCharset string
}

func (r *Response) GetCookies() []*http.Cookie {
    return r.cookiesJar.Cookies(r.req.URL)
}

func (r *Response) GetCookie(cookieName string) *http.Cookie {
    cookies := r.cookiesJar.Cookies(r.req.URL)
    for _, cookie := range cookies {
        if cookie.Name == cookieName {
            return cookie
        }
    }
    return nil
}

func (r *Response) GetRequest() *http.Request {
    return r.req
}

func (r *Response) GetResponse() *http.Response {
    return r.resp
}

func (r *Response) GetStatusCode() string {
    return r.resp.Status
}

func (r *Response) GetHeaders() http.Header {
    return r.req.Header
}

func (r *Response) GetHeader(key string) []string {
    value, ok := r.req.Header[key]
    if !ok {
        return nil
    }
    return value
}
