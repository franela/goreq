package request

import (
    "io/ioutil"
    "net/http"
    "strings"
)

type Get struct {
    Uri string
}


type Post struct {
    Uri string
    Body string
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

    res, err := http.Post(r.Uri, "text/plain", strings.NewReader(r.Body))

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
