package request

import (
    "io/ioutil"
    "net/http"
    "strings"
    "io"
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
        body = strings.NewReader(sb)
    } else if rb, ok := r.Body.(io.Reader); ok {
        body = rb
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
