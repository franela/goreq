package request

import (
    "io/ioutil"
    "net/http"
    "strings"
    "io"
    "encoding/json"
    "bytes"
)

type Get struct {
    Uri string
    Body interface{}
}


type Post struct {
    Uri string
    Body interface{}
}

type Put struct {
    Uri string
    Body interface{}
}

type Delete struct {
    Uri string
    Body interface{}
}

type Response struct {
    StatusCode int
    Body string
    Header http.Header
}

type Error struct {

}

func makeRequest(method string, uri string, body interface{}) (Response) {
    client := &http.Client{}
    b := prepareRequestBody(body)
    req, _ := http.NewRequest(method, uri, b)
    // TODO: handler error
    res, _ := client.Do(req)
    // TODO: handler error
    return newResponse(res)
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
    body, _ := ioutil.ReadAll(res.Body)
    // TODO: handle error
    return Response{ StatusCode: res.StatusCode, Header: res.Header, Body: string(body) }
}

func (r Get) Do() (Response, *Error) {
    return makeRequest("GET", r.Uri, r.Body), nil
}


func (r Put) Do() (Response, *Error) {
    return makeRequest("PUT", r.Uri, r.Body), nil
}

func (r Post) Do() (Response, *Error) {
    return makeRequest("POST", r.Uri, r.Body), nil
}

func (r Delete) Do() (Response, *Error) {
    return makeRequest("DELETE", r.Uri, r.Body), nil
}
