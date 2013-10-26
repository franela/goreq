package request

import (
    "testing"
    . "github.com/onsi/gomega"
    . "github.com/franela/goblin"
    "net/http/httptest"
    "net/http"
    "fmt"
    "io/ioutil"
    "strings"
)

func TestRequest(t *testing.T) {
    g := Goblin(t)

    RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

    g.Describe("Request", func() {

        g.Describe("General request methods", func() {
            var ts *httptest.Server

            g.Before(func() {
                ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                    if (r.Method == "GET" || r.Method == "OPTIONS" || r.Method == "TRACE" || r.Method == "PATCH" || r.Method == "FOOBAR") && r.URL.Path == "/foo" {
                        w.WriteHeader(200)
                        fmt.Fprint(w, "bar")
                    }
                    if r.Method == "POST" && r.URL.Path == "/" {
                        body, _ := ioutil.ReadAll(r.Body)
                        w.Header().Add("Location", ts.URL + "/123")
                        w.WriteHeader(201)
                        fmt.Fprint(w, string(body))
                    }
                    if r.Method == "PUT" && r.URL.Path == "/foo/123" {
                        body, _ := ioutil.ReadAll(r.Body)
                        w.WriteHeader(200)
                        fmt.Fprint(w, string(body))
                    }
                    if r.Method == "DELETE" && r.URL.Path == "/foo/123" {
                        w.WriteHeader(204)
                    }
                }))
            })

            g.After(func() {
                ts.Close()
            })

            g.It("Should do a GET", func() {
                res, err := Request{ Uri: ts.URL + "/foo" }.Do()

                Expect(err).Should(BeNil())
                Expect(res.Body.ToString()).Should(Equal("bar"))
                Expect(res.StatusCode).Should(Equal(200))
            })

            g.Describe("POST", func() {
                g.It("Should send a string", func() {
                    res, err := Request{ Method: "POST", Uri: ts.URL, Body: "foo" }.Do()

                    Expect(err).Should(BeNil())
                    Expect(res.Body.ToString()).Should(Equal("foo"))
                    Expect(res.StatusCode).Should(Equal(201))
                    Expect(res.Header.Get("Location")).Should(Equal(ts.URL + "/123"))
                })

                g.It("Should send a Reader", func() {
                    res, err := Request{ Method: "POST", Uri: ts.URL, Body: strings.NewReader("foo") }.Do()

                    Expect(err).Should(BeNil())
                    Expect(res.Body.ToString()).Should(Equal("foo"))
                    Expect(res.StatusCode).Should(Equal(201))
                    Expect(res.Header.Get("Location")).Should(Equal(ts.URL + "/123"))
                })

                g.It("Send any object that is json encodable", func() {
                    obj := map[string]string {"foo": "bar"}
                    res, err := Request{ Method: "POST", Uri: ts.URL, Body: obj}.Do()

                    Expect(err).Should(BeNil())
                    Expect(res.Body.ToString()).Should(Equal(`{"foo":"bar"}`))
                    Expect(res.StatusCode).Should(Equal(201))
                    Expect(res.Header.Get("Location")).Should(Equal(ts.URL + "/123"))
                })
            })

            g.It("Should do a PUT", func() {
                res, err := Request{ Method: "PUT", Uri: ts.URL + "/foo/123", Body: "foo" }.Do()

                Expect(err).Should(BeNil())
                Expect(res.Body.ToString()).Should(Equal("foo"))
                Expect(res.StatusCode).Should(Equal(200))
            })

            g.It("Should do a DELETE", func() {
                res, err := Request{ Method: "DELETE", Uri: ts.URL + "/foo/123" }.Do()

                Expect(err).Should(BeNil())
                Expect(res.StatusCode).Should(Equal(204))
            })

            g.It("Should do a OPTIONS", func() {
                res, err := Request{ Method: "OPTIONS", Uri: ts.URL + "/foo" }.Do()

                Expect(err).Should(BeNil())
                Expect(res.Body.ToString()).Should(Equal("bar"))
                Expect(res.StatusCode).Should(Equal(200))
            })

            g.It("Should do a PATCH", func() {
                res, err := Request{ Method: "PATCH", Uri: ts.URL + "/foo" }.Do()

                Expect(err).Should(BeNil())
                Expect(res.Body.ToString()).Should(Equal("bar"))
                Expect(res.StatusCode).Should(Equal(200))
            })

            g.It("Should do a TRACE", func() {
                res, err := Request{ Method: "TRACE", Uri: ts.URL + "/foo" }.Do()

                Expect(err).Should(BeNil())
                Expect(res.Body.ToString()).Should(Equal("bar"))
                Expect(res.StatusCode).Should(Equal(200))
            })

            g.It("Should do a custom method", func() {
                res, err := Request{ Method: "FOOBAR", Uri: ts.URL + "/foo" }.Do()

                Expect(err).Should(BeNil())
                Expect(res.Body.ToString()).Should(Equal("bar"))
                Expect(res.StatusCode).Should(Equal(200))
            })

            g.Describe("Responses", func() {
                g.It("Should handle strings", func() {
                    res, _ := Request{ Method: "POST", Uri: ts.URL, Body: "foo bar" }.Do()

                    Expect(res.Body.ToString()).Should(Equal("foo bar"))
                });

                g.It("Should handle io.ReaderCloser", func() {
                    res, _ := Request{ Method: "POST", Uri: ts.URL, Body: "foo bar" }.Do()

                    body, _ := ioutil.ReadAll(res.Body)
                    Expect(string(body)).Should(Equal("foo bar"))
                });

                g.It("Should handle parsing JSON", func() {
                    res, _ := Request{ Method: "POST", Uri: ts.URL, Body: `{"foo": "bar"}` }.Do()

                    var foobar map[string]string

                    res.Body.FromJsonTo(&foobar)

                    Expect(foobar).Should(Equal(map[string]string { "foo": "bar" }))
                });
            })
        })

        g.Describe("Timeouts", func() {
            g.It("Should timeout after a specified amount of ms")
            g.It("Should connect timeout after a specified amount of ms")
        })

        g.Describe("Misc", func() {
            g.It("Should offer to set request headers")
        })
    })
}
