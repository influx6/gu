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
)

// Params defines a map type of key-value pairs to be sent as query parameters.
type Params map[string]string

// HTTPHandler defines a request interface which defines a type which will be used
// to service a http request.
type HTTPHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// BasicHandler which defines a type which is used to service a request and returns an error
// if the request failed.
type BasicHandler interface {
	Serve(http.ResponseWriter, *http.Request) error
}

// PreprocessHandler exposes a type which implements both the http.Handler and
// a method which pre-processes giving routes for use in a request.
type PreprocessHandler interface {
	Preprocess(string) string
}

// CacheHandler defines a handler which implements a type which allows a
// handler to have access to a current request and response with the underline
// cache being used.
type CacheHandler interface {
	ServeAndCache(http.ResponseWriter, *http.Request, cache.Cache) error
}

// Router exposes a interface which describes a
type Router struct {
	cache    cache.Cache
	handler  HTTPHandler
	chandler CacheHandler
	phandler PreprocessHandler
}

// NewRouter returns a new instance of a Router.
func NewRouter(handler interface{}, cache cache.Cache) *Router {
	var router Router
	router.cache = cache

	if ph, ok := handler.(PreprocessHandler); ok {
		router.phandler = ph
	}

	switch mh := handler.(type) {
	case CacheHandler:
		router.chandler = mh
	case BasicHandler:
		router.handler = NewErrorHandler(mh)
	case HTTPHandler:
		router.handler = mh
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

	if r.phandler != nil {
		path = r.phandler.Preprocess(path)
	}

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	// Create a ResponseRecorder for the giving
	responseRecoder := httptest.NewRecorder()

	// TODO: Validate to ensure we don't need this here.
	// var wg sync.WaitGroup
	// wg.Add(1)

	// go func() {
	// 	defer wg.Done()

	// If we have no cache, then use the internal handler.
	if r.cache == nil {
		r.handler.ServeHTTP(responseRecoder, req)

		// return
		res := responseRecoder.Result()
		res.Request = req

		return res, nil
	}

	// Since we have a cache, attempt to serve the request, else use the
	// supplied http.Handler
	if err := r.cache.Serve(responseRecoder, req); err != nil {

		// If the CacheHandler is available then let it handle the request and pass
		// in the routers cache, incase it wishes to store the request into the cache.
		// If not, pass to normal http.handler.
		if r.chandler != nil {
			r.chandler.ServeAndCache(responseRecoder, req, r.cache)
		} else {
			r.handler.ServeHTTP(responseRecoder, req)
		}
	}
	// }()

	// wg.Wait()

	res := responseRecoder.Result()
	res.Request = req

	return res, nil
}

//================================================================================

// ErrorHandler defines a new HTTPHandler which is used to service a request and
// respond to a giving error that occurs.
type ErrorHandler struct {
	Handler BasicHandler
}

// NewErrorHandler returns a new instance of the HTTPHandler which is used to service a request
// and respond to a giving error that occurs.
func NewErrorHandler(bh BasicHandler) HTTPHandler {
	return ErrorHandler{Handler: bh}
}

// ServeHTTP services the incoming request to the underline Handler supplied for the
// basic handler.
func (e ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := e.Handler.Serve(w, r); err != nil {
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
