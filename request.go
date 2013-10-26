package request

import (
    "io/ioutil"
    "net/http"
    "strings"
    "io"
    "encoding/json"
    "bytes"
)

type Request struct {
    Method string
    Uri string
    Body interface{}
}

type Response struct {
    StatusCode int
    Body Body
    Header http.Header
}

type Body struct {
    io.ReadCloser
}

type Error struct {

}


func (b *Body) ToString() (string) {
    body, _ := ioutil.ReadAll(b)
    // TODO: handle error
    return string(body)
}

func (b *Body) FromJsonTo(o interface{}) {
    body, _ := ioutil.ReadAll(b)
    // TODO: handle error

    json.Unmarshal(body, o)
    // TODO: handle error
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

func newResponse(res *http.Response) (Response) {
    return Response{ StatusCode: res.StatusCode, Header: res.Header, Body: Body{ res.Body } }
}

func (r Request) Do() (Response, *Error) {
    client := &http.Client{}
    b := prepareRequestBody(r.Body)
    req, _ := http.NewRequest(r.Method, r.Uri, b)
    // TODO: handler error
    res, _ := client.Do(req)
    // TODO: handler error
    return newResponse(res), nil
}
