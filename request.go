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
}

type Error struct {

}

func (g Get) Do() (Response, string, *Error) {
    response := Response{}
    var body []byte

    res, err := http.Get(g.Uri)

    if err != nil {
        // TODO: Generate the right error
    } else {
        body, _ = ioutil.ReadAll(res.Body)
        response.StatusCode = res.StatusCode
    }

    return response, string(body), nil
}
