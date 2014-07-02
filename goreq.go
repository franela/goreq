package goreq

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type Request struct {
	transport    *http.Transport
	headers      []headerTuple
	Method       string
	Uri          string
	Body         interface{}
	QueryString  interface{}
	Timeout      time.Duration
	ContentType  string
	Accept       string
	Host         string
	UserAgent    string
	Insecure     bool
	MaxRedirects int
}

type Response struct {
	StatusCode int
	Body       Body
	Header     http.Header
}

type headerTuple struct {
	name  string
	value string
}

type Body struct {
	io.ReadCloser
}

type Error struct {
	timeout bool
	Err     error
}

func (e *Error) Timeout() bool {
	return e.timeout
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (b *Body) FromJsonTo(o interface{}) error {
	if body, err := ioutil.ReadAll(b); err != nil {
		return err
	} else if err := json.Unmarshal(body, o); err != nil {
		return err
	}

	return nil
}

func (b *Body) ToString() (string, error) {
	body, err := ioutil.ReadAll(b)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func paramParse(query interface{}) (string, error) {
	var (
		v = &url.Values{}
		s = reflect.ValueOf(query)
		t = reflect.TypeOf(query)
	)


        switch query.(type) {
        case url.Values:
                return query.(url.Values).Encode(), nil
        default:
                for i := 0; i < s.NumField(); i++ {
                        v.Add(strings.ToLower(t.Field(i).Name), fmt.Sprintf("%v", s.Field(i).Interface()))
                }
                return v.Encode(), nil
        }
}

func prepareRequestBody(b interface{}) (io.Reader, error) {
	switch b.(type) {
	case string:
		// treat is as text
		return strings.NewReader(b.(string)), nil
	case io.Reader:
		// treat is as text
		return b.(io.Reader), nil
	case []byte:
		//treat as byte array
		return bytes.NewReader(b.([]byte)), nil
	default:
		// try to jsonify it
		j, err := json.Marshal(b)
		if err == nil {
			return bytes.NewReader(j), nil
		}
		return nil, err
	}
}

func newResponse(res *http.Response) *Response {
	return &Response{StatusCode: res.StatusCode, Header: res.Header, Body: Body{res.Body}}
}

var defaultDialer = &net.Dialer{Timeout: 1000 * time.Millisecond}
var defaultTransport = &http.Transport{Dial: defaultDialer.Dial}
var	defaultClient = &http.Client{Transport: defaultTransport}

func SetConnectTimeout(duration time.Duration) {
	defaultDialer.Timeout = duration
}

func (r *Request) AddHeader(name string, value string) {
	if r.headers == nil {
		r.headers = []headerTuple{}
	}
	r.headers = append(r.headers, headerTuple{name: name, value: value})
}

func (r Request) Do() (*Response, error) {
  var req *http.Request
  var er error
  var client *http.Client
  var transport *http.Transport
  customTransport := false

  if r.transport != nil {
    transport = r.transport
    customTransport = true
  } else {
    transport = defaultTransport
  }

	if r.Insecure {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	} else if transport.TLSClientConfig != nil {
		// the default TLS client (when transport.TLSClientConfig==nil) is
		// already set to verify, so do nothing in that case
		transport.TLSClientConfig.InsecureSkipVerify = false
	}

  if customTransport {
    client = &http.Client{Transport: r.transport}
  } else {
    client = defaultClient
  }

	b, e := prepareRequestBody(r.Body)

	if e != nil {
		// there was a problem marshaling the body
		return nil, &Error{Err: e}
	}

	if strings.EqualFold(r.Method, "GET") || strings.EqualFold(r.Method, "") {
		if r.QueryString != nil {
			param, e := paramParse(r.QueryString)
			if e != nil {
				return nil, &Error{Err: e}
			}
			r.Uri = r.Uri + "?" + param
		}
		req, er = http.NewRequest(r.Method, r.Uri, nil)
	} else {
		req, er = http.NewRequest(r.Method, r.Uri, b)
	}

	if er != nil {
		// we couldn't parse the URL.
		return nil, &Error{Err: er}
	}

	// add headers to the request
	req.Host = r.Host
	req.Header.Add("User-Agent", r.UserAgent)
	req.Header.Add("Content-Type", r.ContentType)
	req.Header.Add("Accept", r.Accept)
	if r.headers != nil {
		for _, header := range r.headers {
			req.Header.Add(header.name, header.value)
		}
	}

	timeout := false
	var timer *time.Timer
	if r.Timeout > 0 {
		timer = time.AfterFunc(r.Timeout, func() {
			transport.CancelRequest(req)
			timeout = true
		})
	}

	res, err := client.Do(req)
	if timer != nil {
		timer.Stop()
	}

	if err != nil {
		if op, ok := err.(*net.OpError); !timeout && ok {
			timeout = op.Timeout()
		}
		return nil, &Error{timeout: timeout, Err: err}
	}

	if isRedirect(res.StatusCode) && r.MaxRedirects > 0 {
		loc, _ := res.Location()
		r.MaxRedirects--
		r.Uri = loc.String()
		return r.Do()
	}

	return newResponse(res), nil
}

func isRedirect(status int) bool {
	switch status {
	case http.StatusMovedPermanently:
		return true
	case http.StatusFound:
		return true
	case http.StatusSeeOther:
		return true
	case http.StatusTemporaryRedirect:
		return true
	default:
		return false
	}
}
