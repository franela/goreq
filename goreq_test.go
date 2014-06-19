package goreq

import (
	"fmt"
	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type Query struct {
	Limit int
	Skip  int
}

func TestRequest(t *testing.T) {

	query := Query{
		Limit: 3,
		Skip:  5,
	}

        valuesQuery := url.Values{};
        valuesQuery.Set("name", "marcos")
        valuesQuery.Add("friend", "jonas")
        valuesQuery.Add("friend", "peter")

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
					if r.Method == "GET" && r.URL.Path == "/getquery" {
						w.WriteHeader(200)
						fmt.Fprint(w, fmt.Sprintf("%v", r.URL))
					}
					if r.Method == "GET" && r.URL.Path == "/getbody" {
						w.WriteHeader(200)
						io.Copy(w, r.Body)
					}
					if r.Method == "POST" && r.URL.Path == "/" {
						w.Header().Add("Location", ts.URL+"/123")
						w.WriteHeader(201)
						io.Copy(w, r.Body)
					}
					if r.Method == "POST" && r.URL.Path == "/getquery" {
						w.WriteHeader(200)
						fmt.Fprint(w, fmt.Sprintf("%v", r.URL))
					}
					if r.Method == "PUT" && r.URL.Path == "/foo/123" {
						w.WriteHeader(200)
						io.Copy(w, r.Body)
					}
					if r.Method == "DELETE" && r.URL.Path == "/foo/123" {
						w.WriteHeader(204)
					}
					if r.Method == "GET" && r.URL.Path == "/redirect_test/301" {
						http.Redirect(w, r, "/redirect_test/302", 301)
					}
					if r.Method == "GET" && r.URL.Path == "/redirect_test/302" {
						http.Redirect(w, r, "/redirect_test/303", 302)
					}
					if r.Method == "GET" && r.URL.Path == "/redirect_test/303" {
						http.Redirect(w, r, "/redirect_test/307", 303)
					}
					if r.Method == "GET" && r.URL.Path == "/redirect_test/307" {
						http.Redirect(w, r, "/getquery", 307)
					}
				}))
			})

			g.After(func() {
				ts.Close()
			})

			g.Describe("GET", func() {

				g.It("Should do a GET", func() {
					res, err := Request{Uri: ts.URL + "/foo"}.Do()

					Expect(err).Should(BeNil())
					str, _ := res.Body.ToString()
					Expect(str).Should(Equal("bar"))
					Expect(res.StatusCode).Should(Equal(200))
				})

				g.It("Should do a GET with querystring", func() {
					res, err := Request{
						Uri:         ts.URL + "/getquery",
						QueryString: query,
					}.Do()

					Expect(err).Should(BeNil())
					str, _ := res.Body.ToString()
					Expect(str).Should(Equal("/getquery?limit=3&skip=5"))
					Expect(res.StatusCode).Should(Equal(200))
				})


				g.It("Should support url.Values in querystring", func() {
					res, err := Request{
						Uri:         ts.URL + "/getquery",
						QueryString: valuesQuery,
					}.Do()

					Expect(err).Should(BeNil())
					str, _ := res.Body.ToString()
					Expect(str).Should(Equal("/getquery?friend=jonas&friend=peter&name=marcos"))
					Expect(res.StatusCode).Should(Equal(200))
				})

				g.It("Should support sending string body", func() {
					res, err := Request{Uri: ts.URL + "/getbody", Body: "foo"}.Do()

					Expect(err).Should(BeNil())
					str, _ := res.Body.ToString()
					Expect(str).Should(Equal("foo"))
					Expect(res.StatusCode).Should(Equal(200))
				})

				g.It("Shoulds support sending a Reader body", func() {
					res, err := Request{Uri: ts.URL + "/getbody", Body: strings.NewReader("foo")}.Do()

					Expect(err).Should(BeNil())
					str, _ := res.Body.ToString()
					Expect(str).Should(Equal("foo"))
					Expect(res.StatusCode).Should(Equal(200))
				})

				g.It("Support sending any object that is json encodable", func() {
					obj := map[string]string{"foo": "bar"}
					res, err := Request{Uri: ts.URL + "/getbody", Body: obj}.Do()

					Expect(err).Should(BeNil())
					str, _ := res.Body.ToString()
					Expect(str).Should(Equal(`{"foo":"bar"}`))
					Expect(res.StatusCode).Should(Equal(200))
				})

				g.It("Support sending an array of bytes body", func() {
					bdy := []byte{'f', 'o', 'o'}
					res, err := Request{Uri: ts.URL + "/getbody", Body: bdy}.Do()

					Expect(err).Should(BeNil())
					str, _ := res.Body.ToString()
					Expect(str).Should(Equal("foo"))
					Expect(res.StatusCode).Should(Equal(200))
				})

				g.It("Should return an error when body is not JSON encodable", func() {
					res, err := Request{Uri: ts.URL + "/getbody", Body: math.NaN()}.Do()

					Expect(res).Should(BeNil())
					Expect(err).ShouldNot(BeNil())
				})

			})

			g.Describe("POST", func() {
				g.It("Should send a string", func() {
					res, err := Request{Method: "POST", Uri: ts.URL, Body: "foo"}.Do()

					Expect(err).Should(BeNil())
					str, _ := res.Body.ToString()
					Expect(str).Should(Equal("foo"))
					Expect(res.StatusCode).Should(Equal(201))
					Expect(res.Header.Get("Location")).Should(Equal(ts.URL + "/123"))
				})

				g.It("Should send a Reader", func() {
					res, err := Request{Method: "POST", Uri: ts.URL, Body: strings.NewReader("foo")}.Do()

					Expect(err).Should(BeNil())
					str, _ := res.Body.ToString()
					Expect(str).Should(Equal("foo"))
					Expect(res.StatusCode).Should(Equal(201))
					Expect(res.Header.Get("Location")).Should(Equal(ts.URL + "/123"))
				})

				g.It("Send any object that is json encodable", func() {
					obj := map[string]string{"foo": "bar"}
					res, err := Request{Method: "POST", Uri: ts.URL, Body: obj}.Do()

					Expect(err).Should(BeNil())
					str, _ := res.Body.ToString()
					Expect(str).Should(Equal(`{"foo":"bar"}`))
					Expect(res.StatusCode).Should(Equal(201))
					Expect(res.Header.Get("Location")).Should(Equal(ts.URL + "/123"))
				})

				g.It("Send an array of bytes", func() {
					bdy := []byte{'f', 'o', 'o'}
					res, err := Request{Method: "POST", Uri: ts.URL, Body: bdy}.Do()

					Expect(err).Should(BeNil())
					str, _ := res.Body.ToString()
					Expect(str).Should(Equal("foo"))
					Expect(res.StatusCode).Should(Equal(201))
					Expect(res.Header.Get("Location")).Should(Equal(ts.URL + "/123"))
				})

				g.It("Should return an error when body is not JSON encodable", func() {
					res, err := Request{Method: "POST", Uri: ts.URL, Body: math.NaN()}.Do()

					Expect(res).Should(BeNil())
					Expect(err).ShouldNot(BeNil())
				})

				g.It("Should do a POST with querystring", func() {
					bdy := []byte{'f', 'o', 'o'}
					res, err := Request{
            Method: "POST",
						Uri:         ts.URL + "/getquery",
            Body: bdy,
						QueryString: query,
					}.Do()

					Expect(err).Should(BeNil())
					str, _ := res.Body.ToString()
					Expect(str).Should(Equal("/getquery?limit=3&skip=5"))
					Expect(res.StatusCode).Should(Equal(200))
				})
			})

			g.It("Should do a PUT", func() {
				res, err := Request{Method: "PUT", Uri: ts.URL + "/foo/123", Body: "foo"}.Do()

				Expect(err).Should(BeNil())
				str, _ := res.Body.ToString()
				Expect(str).Should(Equal("foo"))
				Expect(res.StatusCode).Should(Equal(200))
			})

			g.It("Should do a DELETE", func() {
				res, err := Request{Method: "DELETE", Uri: ts.URL + "/foo/123"}.Do()

				Expect(err).Should(BeNil())
				Expect(res.StatusCode).Should(Equal(204))
			})

			g.It("Should do a OPTIONS", func() {
				res, err := Request{Method: "OPTIONS", Uri: ts.URL + "/foo"}.Do()

				Expect(err).Should(BeNil())
				str, _ := res.Body.ToString()
				Expect(str).Should(Equal("bar"))
				Expect(res.StatusCode).Should(Equal(200))
			})

			g.It("Should do a PATCH", func() {
				res, err := Request{Method: "PATCH", Uri: ts.URL + "/foo"}.Do()

				Expect(err).Should(BeNil())
				str, _ := res.Body.ToString()
				Expect(str).Should(Equal("bar"))
				Expect(res.StatusCode).Should(Equal(200))
			})

			g.It("Should do a TRACE", func() {
				res, err := Request{Method: "TRACE", Uri: ts.URL + "/foo"}.Do()

				Expect(err).Should(BeNil())
				str, _ := res.Body.ToString()
				Expect(str).Should(Equal("bar"))
				Expect(res.StatusCode).Should(Equal(200))
			})

			g.It("Should do a custom method", func() {
				res, err := Request{Method: "FOOBAR", Uri: ts.URL + "/foo"}.Do()

				Expect(err).Should(BeNil())
				str, _ := res.Body.ToString()
				Expect(str).Should(Equal("bar"))
				Expect(res.StatusCode).Should(Equal(200))
			})

			g.Describe("Responses", func() {
				g.It("Should handle strings", func() {
					res, _ := Request{Method: "POST", Uri: ts.URL, Body: "foo bar"}.Do()

					str, _ := res.Body.ToString()
					Expect(str).Should(Equal("foo bar"))
				})

				g.It("Should handle io.ReaderCloser", func() {
					res, _ := Request{Method: "POST", Uri: ts.URL, Body: "foo bar"}.Do()

					body, _ := ioutil.ReadAll(res.Body)
					Expect(string(body)).Should(Equal("foo bar"))
				})

				g.It("Should handle parsing JSON", func() {
					res, _ := Request{Method: "POST", Uri: ts.URL, Body: `{"foo": "bar"}`}.Do()

					var foobar map[string]string

					res.Body.FromJsonTo(&foobar)

					Expect(foobar).Should(Equal(map[string]string{"foo": "bar"}))
				})
			})
			g.Describe("Redirects", func() {
				g.It("Should do not follow by default", func() {
					res, _ := Request{
						Uri: ts.URL + "/redirect_test/301",
					}.Do()
					Expect(res.StatusCode).Should(Equal(301))
				})
				g.It("Should follow only specified number of MaxRedirects", func() {
					res, _ := Request{
						Uri:          ts.URL + "/redirect_test/301",
						MaxRedirects: 1,
					}.Do()
					Expect(res.StatusCode).Should(Equal(302))
					res, _ = Request{
						Uri:          ts.URL + "/redirect_test/301",
						MaxRedirects: 2,
					}.Do()
					Expect(res.StatusCode).Should(Equal(303))
					res, _ = Request{
						Uri:          ts.URL + "/redirect_test/301",
						MaxRedirects: 3,
					}.Do()
					Expect(res.StatusCode).Should(Equal(307))
					res, _ = Request{
						Uri:          ts.URL + "/redirect_test/301",
						MaxRedirects: 4,
					}.Do()
					Expect(res.StatusCode).Should(Equal(200))
				})
			})
		})

		g.Describe("Timeouts", func() {

			g.Describe("Connection timeouts", func() {
				g.It("Should connect timeout after a default of 1000 ms", func() {
					start := time.Now()
					res, err := Request{Uri: "http://10.255.255.1"}.Do()
					elapsed := time.Since(start)

					Expect(elapsed).Should(BeNumerically("<", 1100*time.Millisecond))
					Expect(elapsed).Should(BeNumerically(">=", 1000*time.Millisecond))
					Expect(res).Should(BeNil())
					Expect(err.(*Error).Timeout()).Should(BeTrue())
				})
				g.It("Should connect timeout after a custom amount of time", func() {
					SetConnectTimeout(100 * time.Millisecond)
					start := time.Now()
					res, err := Request{Uri: "http://10.255.255.1"}.Do()
					elapsed := time.Since(start)

					Expect(elapsed).Should(BeNumerically("<", 150*time.Millisecond))
					Expect(elapsed).Should(BeNumerically(">=", 100*time.Millisecond))
					Expect(res).Should(BeNil())
					Expect(err.(*Error).Timeout()).Should(BeTrue())
				})
			})

			g.Describe("Request timeout", func() {
				var ts *httptest.Server
				stop := make(chan bool)

				g.Before(func() {
					ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						<-stop
						// just wait for someone to tell you when to end the request. this is used to simulate a slow server
					}))
				})
				g.After(func() {
					stop <- true
					ts.Close()
				})
				g.It("Should request timeout after a custom amount of time", func() {
					SetConnectTimeout(1000 * time.Millisecond)

					start := time.Now()
					res, err := Request{Uri: ts.URL, Timeout: 500 * time.Millisecond}.Do()
					elapsed := time.Since(start)

					Expect(elapsed).Should(BeNumerically("<", 550*time.Millisecond))
					Expect(elapsed).Should(BeNumerically(">=", 500*time.Millisecond))
					Expect(res).Should(BeNil())
					Expect(err.(*Error).Timeout()).Should(BeTrue())
				})
			})
		})

		g.Describe("Misc", func() {
			g.It("Should offer to set request headers", func() {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					Expect(r.Header.Get("User-Agent")).Should(Equal("foobaragent"))
					Expect(r.Host).Should(Equal("foobar.com"))
					Expect(r.Header.Get("Accept")).Should(Equal("application/json"))
					Expect(r.Header.Get("Content-Type")).Should(Equal("application/json"))
					Expect(r.Header.Get("X-Custom")).Should(Equal("foobar"))

					w.WriteHeader(200)
				}))
				defer ts.Close()

				req := Request{Uri: ts.URL, Accept: "application/json", ContentType: "application/json", UserAgent: "foobaragent", Host: "foobar.com"}
				req.AddHeader("X-Custom", "foobar")
				res, _ := req.Do()

				Expect(res.StatusCode).Should(Equal(200))
			})
			g.It("Should change transport TLS config if Request.Insecure is set", func() {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
				}))
				defer ts.Close()

				trans := &http.Transport{Dial: dialer.Dial}
				req := Request{
					transport: trans,
					Insecure:  true,
					Uri:       ts.URL,
					Host:      "foobar.com",
				}
				res, _ := req.Do()

				Expect(trans.TLSClientConfig.InsecureSkipVerify).Should(Equal(true))
				Expect(res.StatusCode).Should(Equal(200))
			})
		})

		g.Describe("Errors", func() {
			var ts *httptest.Server

			g.Before(func() {
				ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method == "POST" && r.URL.Path == "/" {
						w.Header().Add("Location", ts.URL+"/123")
						w.WriteHeader(201)
						io.Copy(w, r.Body)
					}
				}))
			})

			g.After(func() {
				ts.Close()
			})
			g.It("Should throw an error when FromJsonTo fails", func() {
				res, _ := Request{Method: "POST", Uri: ts.URL, Body: `{"foo": "bar"`}.Do()
				var foobar map[string]string

				err := res.Body.FromJsonTo(&foobar)
				Expect(err).Should(HaveOccurred())
			})
			g.It("Should handle Url parsing errors", func() {
				_, err := Request{Uri: ":"}.Do()

				Expect(err).ShouldNot(BeNil())
			})
			g.It("Should handle DNS errors", func() {
				_, err := Request{Uri: "http://*.localhost"}.Do()

				Expect(err).ShouldNot(BeNil())
			})
		})
	})
}
