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
}


type Post struct {
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

func (r Get) Do() (Response, *Error) {
    response := Response{}

    res, err := http.Get(r.Uri)

    if err != nil {
        // TODO: Generate the right error
    } else {
        body, _ := ioutil.ReadAll(res.Body)
        response.Body = string(body)
        response.StatusCode = res.StatusCode
    }

    return response, nil
}

func (r Post) Do() (Response, *Error) {
    response := Response{}

    var body io.Reader

    if sb, ok := r.Body.(string); ok {
        // treat is as text
        body = strings.NewReader(sb)
    } else if rb, ok := r.Body.(io.Reader); ok {
        // treat is as text
        body = rb
    } else {
        // try to jsonify it
        if j, err := json.Marshal(r.Body); err == nil {
            body = bytes.NewReader(j)
        } else {
            // TODO: handle error. don't know what to do with this.
        }
    }
    res, err := http.Post(r.Uri, "text/plain", body)

    if err != nil {
        // TODO: Generate the right error
    } else {
        body, _ := ioutil.ReadAll(res.Body)
        response.Body = string(body)
        response.StatusCode = res.StatusCode
        response.Header = res.Header
    }

    return response, nil
}
