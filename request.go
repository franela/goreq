package request

import (
    "io/ioutil"
    "net/http"
    "strings"
    "io"
    "encoding/json"
    "bytes"
    "time"
    "net"
)

type Request struct {
    Method string
    Uri string
    Body interface{}
    Timeout time.Duration
}

type Response struct {
    StatusCode int
    Body string
    Header http.Header
}

type Error struct {
    requestTimeout bool
}

func (e *Error) ConnectTimeout() (bool) {
    return !e.requestTimeout
}
func (e *Error) RequestTimeout() (bool) {
    return e.requestTimeout
}

func prepareRequestBody(b interface{}) (io.Reader) {
    var body io.Reader

    if sb, ok := b.(string); ok {
        // treat is as text
        body = strings.NewReader(sb)
    } else if rb, ok := b.(io.Reader); ok {
        // treat is as text
        body = rb
    } else {
        // try to jsonify it
        if j, err := json.Marshal(b); err == nil {
            body = bytes.NewReader(j)
        } else {
            // TODO: handle error. don't know what to do with this.
        }
    }

   return body
}

func newResponse(res *http.Response) (*Response) {
    body, _ := ioutil.ReadAll(res.Body)
    // TODO: handle error
    return &Response{ StatusCode: res.StatusCode, Header: res.Header, Body: string(body) }
}

var dialer = &net.Dialer{ Timeout: 1000 * time.Millisecond }
var transport = &http.Transport{ Dial: dialer.Dial }

func SetConnectTimeout(duration time.Duration) {
    dialer.Timeout = duration
}

func (r Request) Do() (*Response, *Error) {
    client := &http.Client{ Transport: transport }
    b := prepareRequestBody(r.Body)
    req, _ := http.NewRequest(r.Method, r.Uri, b)
    // TODO: handler error
    requestTimeout := false
    var timer *time.Timer
    if r.Timeout > 0 {
        timer = time.AfterFunc(r.Timeout, func() {
            transport.CancelRequest(req)
            requestTimeout = true
        })
    }
    res, err := client.Do(req)
    if timer != nil {
        timer.Stop()
    }
    // TODO: handler error
    if err != nil {
        return nil, &Error{ requestTimeout: requestTimeout }
    }
    return newResponse(res), nil
}
