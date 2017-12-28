package cache

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// WebPair defines a struct which contains the request object and response object
// received for that request.
type WebPair struct {
	Request  Request
	Response Response
}

// Cache defines a interface which exposes a cache like structure for retrieving
// requests.
type Cache interface {
	Empty() error
	Delete(string) error
	AddData(string, []byte) error
	Add(string, *http.Response) error
	Get(string) (Request, Response, error)
	Serve(http.ResponseWriter, *http.Request) error
}

// Request defines the structure which holds request fields for a given
// request.
type Request struct {
	Body      bytes.Buffer      `json:"body"`
	Method    string            `json:"method"`
	Path      string            `json:"path"`
	URL       *url.URL          `json:"url"`
	Headers   map[string]string `json:"headers"`
	Underline interface{}       `json:"underline"`
}

// Response defines a struct to hold a response for a giving resource.
type Response struct {
	Status    int               `json:"status"`
	Method    string            `json:"method"`
	Type      string            `json:"type"`
	Body      bytes.Buffer      `json:"body"`
	Headers   map[string]string `json:"headers"`
	Cookies   []string          `json:"cookies"`
	Underline interface{}       `json:"underline"`
}

// HTTPRequestToRequest transforms a giving request object into a cache.Request
// object.
func HTTPRequestToRequest(req *http.Request) *Request {
	var rq *Request
	rq.URL = req.URL
	rq.Path = req.URL.String()
	rq.Method = req.Method
	rq.Headers = headerToMap(req.Header)

	return rq
}

// HTTPResponseToResponse transforms a giving response object into a Response
// object.
func HTTPResponseToResponse(res *http.Response) (*Response, *Request) {
	var buf bytes.Buffer

	if res.Body != nil {
		io.Copy(&buf, res.Body)
		res.Body.Close()
	}

	var rq *Request
	var wq *Response

	if res.Request != nil {
		rq.URL = res.Request.URL
		rq.Path = res.Request.URL.String()
		rq.Method = res.Request.Method
		rq.Headers = headerToMap(res.Request.Header)

		wq.Method = res.Request.Method
	}

	wq.Body = buf
	wq.Status = res.StatusCode
	wq.Headers = headerToMap(res.Header)
	wq.Cookies = cookies(res.Cookies())

	return wq, rq
}

func cookies(cookies []*http.Cookie) []string {
	var co []string

	for _, cookie := range cookies {
		co = append(co, cookie.String())
	}

	return co
}

func headerToMap(hl http.Header) map[string]string {
	ma := make(map[string]string)

	for key, vals := range hl {
		ma[key] = strings.Join(vals, ";")
	}

	return ma
}
