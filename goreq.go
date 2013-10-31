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
    timeout bool
    Err error
}

func (e *Error) Timeout() (bool) {
    return e.timeout
}

func (e *Error) Error() (string) {
    return e.Err.Error()
}

func (b *Body) FromJsonTo(o interface{}) {
    body, _ := ioutil.ReadAll(b)
    // TODO: handle error

    json.Unmarshal(body, o)
    // TODO: handle error
}

func (b *Body) ToString() (string) {
    body, _ := ioutil.ReadAll(b)
    // TODO: handle error
    return string(body)
}

func prepareRequestBody(b interface{}) (io.Reader, error) {
    var body io.Reader

    if sb, ok := b.(string); ok {
        // treat is as text
        body = strings.NewReader(sb)
    } else if rb, ok := b.(io.Reader); ok {
        // treat is as text
        body = rb
    } else {
        // try to jsonify it
        j, err := json.Marshal(b)
        if err == nil {
            body = bytes.NewReader(j)
        } else {
            return nil, err
        }
    }

   return body, nil
}

func newResponse(res *http.Response) (*Response) {
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
    b, e := prepareRequestBody(r.Body)

    if e != nil {
        // there was a problem marshaling the body
        return nil, &Error{ Err: e }
    }

    req, er := http.NewRequest(r.Method, r.Uri, b)

    if er != nil {
        // we couldn't parse the URL.
        return nil, &Error{ Err: er }
    }

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

    timeout := false
    var timer *time.Timer
    if r.Timeout > 0 {
        timer = time.AfterFunc(r.Timeout, func() {
            transport.CancelRequest(req)
            timeout = true
        })
    }

    res, err := client.Do(req)
    if timer != nil {
        timer.Stop()
    }

    if err != nil {
        if op, ok := err.(*net.OpError); !timeout && ok {
            timeout = op.Timeout()
        }
        return nil, &Error{ timeout: timeout, Err: err }
    }
    return newResponse(res), nil
}
