package goreq

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
    headers []headerTuple
    Method string
    Uri string
    Body interface{}
    Timeout time.Duration
    ContentType string
    Accept string
    Host string
    UserAgent string
}

type Response struct {
    StatusCode int
    Body Body
    Header http.Header
}

type headerTuple struct {
    name string
    value string
}

type Body struct {
    io.ReadCloser
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

func (b *Body) FromJsonTo(o interface{}) error {
    if body, err := ioutil.ReadAll(b); err != nil {
        return err
    } else if err := json.Unmarshal(body, o); err != nil {
        return err
    }

    return nil
}

func (b *Body) ToString() (string) {
    body, _ := ioutil.ReadAll(b)
    // TODO: handle error
    return string(body)
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
    // TODO: handle error
    return &Response{ StatusCode: res.StatusCode, Header: res.Header, Body: Body{ res.Body } }
}

var dialer = &net.Dialer{ Timeout: 1000 * time.Millisecond }
var transport = &http.Transport{ Dial: dialer.Dial }

func SetConnectTimeout(duration time.Duration) {
    dialer.Timeout = duration
}

func (r *Request) AddHeader(name string, value string) {
    if r.headers == nil {
        r.headers = []headerTuple {}
    }
    r.headers = append(r.headers, headerTuple { name: name, value: value })
}

func (r Request) Do() (*Response, *Error) {
    client := &http.Client{ Transport: transport }
    b := prepareRequestBody(r.Body)
    req, _ := http.NewRequest(r.Method, r.Uri, b)

    // add headers to the request
    req.Host = r.Host
    req.Header.Add("User-Agent", r.UserAgent)
    req.Header.Add("Content-Type", r.ContentType)
    req.Header.Add("Accept", r.Accept)
    if r.headers != nil {
        for _, header := range(r.headers) {
            req.Header.Add(header.name, header.value)
        }
    }

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
