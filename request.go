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

func (g Get) Do() (Response, bool, Error) {
    ok := true
    response := Response{}

    res, err := http.Get(g.Uri)

    if err != nil {
        ok = false
        // TODO: Generate the right error
    } else {
        body, e := ioutil.ReadAll(res.Body)

        if e != nil {
            ok = false

            // TODO: Generate the right error
        } else {
            response.Body = string(body)
            response.StatusCode = res.StatusCode
        }
    }

    return response, ok, Error{}
}
