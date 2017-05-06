package router

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gu-io/gu/router/cache"
	"github.com/influx6/faux/pattern"
)

// Params defines a map type of key-value pairs to be sent as query parameters.
type Params map[string]string

// PreprocessHandler exposes a type which implements both the http.Handler and
// a method which pre-processes giving routes for use in a request.
type PreprocessHandler interface {
	Preprocess(string) string
}

// HTTPHandler defines a request interface which defines a type which will be used
// to service a http request.
type HTTPHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// HTTPCacheHandler defines a handler which implements a type which allows a
// handler to have access to a current request and response with the underline
// cache being used.
type HTTPCacheHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, cache.Cache)
}

//================================================================================

// server defines an interface used to provide a concrete method
// which returns a new path and request handler for that path.
type server interface {
	Match(string) (string, HTTPCacheHandler, error)
}

// Router exposes a struct which describes a multi-handler of request where
// it
type Router struct {
	cache cache.Cache
	sx    server
}

// NewRouter returns a new instance of a Router.
func NewRouter(handler interface{}, cache cache.Cache) *Router {
	var router Router
	router.cache = cache

	switch hl := handler.(type) {
	case Mux:
		router.sx = NewMultiplexer(hl)
	case HTTPCacheHandler, HTTPHandler:
		router.sx = NewHandleMux(hl)
	case Multiplexer:
		router.sx = hl
	default:
		panic("Unsupported handler type")
	}

	return &router
}

// Cache returns the internal cache used by the router.
func (r *Router) Cache() cache.Cache {
	return r.cache
}

// Patch retrieves the giving path and returns the response expected using a PATCH method.
func (r *Router) Patch(path string, params Params, body io.ReadCloser) (*http.Response, error) {
	return r.Do("PATCH", path, params, body)
}

// Put retrieves the giving path and returns the response expected using a PUT method.
func (r *Router) Put(path string, params Params, body io.ReadCloser) (*http.Response, error) {
	return r.Do("PUT", path, params, body)
}

// Post retrieves the giving path and returns the response expected using a POST method.
func (r *Router) Post(path string, params Params, body io.ReadCloser) (*http.Response, error) {
	return r.Do("POST", path, params, body)
}

// Options retrieves the giving path and returns the response expected using a OPTIONS method.
func (r *Router) Options(path string, params Params) (*http.Response, error) {
	return r.Do("OPTIONS", path, params, nil)
}

// Delete retrieves the giving path and returns the response expected using a DELETE method.
func (r *Router) Delete(path string, params Params) (*http.Response, error) {
	return r.Do("DELETE", path, params, nil)
}

// Head retrieves the giving path and returns the response expected using a HEAD method.
func (r *Router) Head(path string, params Params) (*http.Response, error) {
	return r.Do("HEAD", path, params, nil)
}

// Get retrieves the giving path and returns the response expected using a GET method.
func (r *Router) Get(path string, params Params) (*http.Response, error) {
	return r.Do("GET", path, params, nil)
}

// Do performs the giving requests for a giving path with the provided body and returns the
// response for that method.
func (r *Router) Do(method string, path string, params Params, body io.ReadCloser) (*http.Response, error) {
	path, handler, err := r.sx.Match(path)
	if err != nil {
		return nil, err
	}

	// Do we have parameters?
	if params != nil {
		parameters := WrapParams(params)

		// Does it already contain a query part?
		if strings.Contains(path, "?") {
			path = path + "&" + url.QueryEscape(parameters)
		} else {
			path = path + "?" + url.QueryEscape(parameters)
		}
	}

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	// Create a ResponseRecorder for the giving
	responseRecoder := httptest.NewRecorder()

	// TODO: Validate to ensure we don't need this here.
	if r.cache != nil {
		handler.ServeHTTP(responseRecoder, req, r.cache)
	}

	switch r.cache == nil {
	case true:
		handler.ServeHTTP(responseRecoder, req, nil)
	case false:
		if err := r.cache.Serve(responseRecoder, req); err != nil {
			handler.ServeHTTP(responseRecoder, req, r.cache)
		}
	}

	res := responseRecoder.Result()
	res.Request = req

	return res, nil
}

//================================================================================

// HandleMux defines a structure which handles the variaties in the supported request
// handlers of the router package and encapsulates the needed behaviour calls.
type HandleMux struct {
	normal     HTTPHandler
	caches     HTTPCacheHandler
	preprocess PreprocessHandler
}

// NewHandleMux returns a new instance of a HandleMux.
func NewHandleMux(handler interface{}) HandleMux {
	var hm HandleMux

	if ph, ok := handler.(PreprocessHandler); ok {
		hm.preprocess = ph
	}

	switch hl := handler.(type) {
	case HTTPCacheHandler:
		hm.caches = hl
	case HTTPHandler:
		hm.normal = hl
	default:
		panic("Unsupported handler type")
	}

	return hm
}

// Match examines the path and returns a new path, a Mux to handle the request
// else returns an error if one is not found.
func (m HandleMux) Match(path string) (string, HTTPCacheHandler, error) {
	if m.preprocess != nil {
		return m.preprocess.Preprocess(path), m, nil
	}

	return path, m, nil
}

// ServeAndCache attempts to service request with either the cache or normal handler found within
// itself.
func (m HandleMux) ServeHTTP(w http.ResponseWriter, r *http.Request, c cache.Cache) {

	// If we have no cache, then use the internal handler.
	switch c == nil {
	case true:
		m.normal.ServeHTTP(w, r)
		break
	case false:
		switch m.caches == nil {
		case true:
			m.normal.ServeHTTP(w, r)
			break
		case false:
			m.caches.ServeHTTP(w, r, c)
			break
		}
	}

}

//================================================================================

// Mux defines a structure which giving a Matcher and a handler which implements either the
// HTTPCacheHandler/BasicHandler/HTTPHandler and also optionally the PreprocessorHandler, where
// request being serviced will receive a new path when called and will be used to define,
// which part an external service will use to service a incoming relative path request.
// Below is demonstrated a scenario where Mux can be used for to handle request for specific
// namespace.
// Scenrio:
//   Mux has URI Matcher for github/*
// When:
//	Mux receives request for github/gu-io/buba
//
// If Mux has Preprocessor to to transform path to github.com/path:
//	Then: Path returned is github.com/gu-io/buba
//
// If Mux has No Preprocessor
// 	Then: Path returned is gu-io/buba.
type Mux struct {
	matcher pattern.URIMatcher
	handler HandleMux
}

// NewMux returns a new instance of a mux.
func NewMux(namespace string, handler interface{}) Mux {
	if !strings.HasSuffix(namespace, "/*") {
		namespace = strings.TrimSuffix(namespace, "/") + "/*"
	}

	var mx Mux
	mx.matcher = URIMatcher(namespace)
	mx.handler = NewHandleMux(handler)

	return mx
}

// Match validates that the giving Mux matches the wanted path and
// extracts the real path from the provided path, returning the true/false
// if it matched the path.
func (m *Mux) Match(path string) (string, HTTPCacheHandler, error) {
	_, rem, ok := m.matcher.Validate(path)
	if !ok {
		return "", nil, errors.New("Invalid Path")
	}

	return m.handler.Match(rem)
}

//================================================================================

// Multiplexer defines a struct which manages giving set of Mux and adequates
// calls the first to match a giving request about a incoming request.
type Multiplexer struct {
	mux []Mux
}

// NewMultiplexer returns a new instance of a Multiplexer.
func NewMultiplexer(mx ...Mux) Multiplexer {
	return Multiplexer{
		mux: mx,
	}
}

// Match examines the path and returns a new path, a Mux to handle the request
// else returns an error if one is not found.
func (m Multiplexer) Match(path string) (string, HTTPCacheHandler, error) {
	for _, item := range m.mux {
		newPath, handler, err := item.Match(path)
		if err != nil {
			continue
		}

		return newPath, handler, nil
	}

	return path, nil, fmt.Errorf("Mux not found for %s", path)
}

//================================================================================

// BasicHTTPHandler defines a request interface which defines a type which will be used
// to service a http request.
type BasicHTTPHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

// ErrorHandler defines a new HTTPHandler which is used to service a request and
// respond to a giving error that occurs.
type ErrorHandler struct {
	Handler BasicHTTPHandler
}

// NewErrorHandler returns a new instance of the HTTPHandler which is used to service a request
// and respond to a giving error that occurs.
func NewErrorHandler(bh BasicHTTPHandler) HTTPHandler {
	return ErrorHandler{Handler: bh}
}

// ServeHTTP services the incoming request to the underline Handler supplied for the
// basic handler.
func (e ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := e.Handler.ServeHTTP(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//================================================================================

// WrapParams transforms the params map into a string of &key=value pair.
func WrapParams(params Params) string {
	var q []string

	for k, v := range params {
		q = append(q, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(q, "&")
}

// ReadBody returns the body of the giving response.
func ReadBody(res *http.Response) ([]byte, error) {
	if res.Body == nil {
		return nil, errors.New("Response has no body/content")
	}

	var buf bytes.Buffer

	defer res.Body.Close()
	io.Copy(&buf, res.Body)

	return buf.Bytes(), nil
}
