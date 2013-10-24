package request

import (
    "io/ioutil"
    "net/http"
)

type Get struct {
    Uri string
}

type Response struct {
    StatusCode int
    Body string
}

type Error struct {

}

func (g Get) Do() (Response, *Error) {
    response := Response{}

    res, err := http.Get(g.Uri)

    if err != nil {
        // TODO: Generate the right error
    } else {
        body, _ := ioutil.ReadAll(res.Body)
        response.Body = string(body)
        response.StatusCode = res.StatusCode
    }

    return response, nil
}
