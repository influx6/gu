package router

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
)

// Cache defines a interface which exposes a cache like structure for retrieving
// requests.
type Cache interface {
	Empty() error
	Delete(string) error
	Add(string, *http.Response)
	Serve(http.ResponseWriter, *http.Request) error
}

// Router exposes a interface which describes a
type Router struct {
	handler http.Handler
	cache   Cache
}

// NewRouter returns a new instance of a Router.
func NewRouter(handler http.Handler, cache Cache) *Router {
	var router Router
	router.cache = cache
	router.handler = handler

	return &router
}

// Patch retrieves the giving path and returns the response expected using a PATCH method.
func (r *Router) Patch(path string, body io.ReadCloser) (*http.Response, error) {
	return r.Do("PATCH", path, body)
}

// Put retrieves the giving path and returns the response expected using a PUT method.
func (r *Router) Put(path string, body io.ReadCloser) (*http.Response, error) {
	return r.Do("PUT", path, body)
}

// Post retrieves the giving path and returns the response expected using a POST method.
func (r *Router) Post(path string, body io.ReadCloser) (*http.Response, error) {
	return r.Do("POST", path, body)
}

// Options retrieves the giving path and returns the response expected using a OPTIONS method.
func (r *Router) Options(path string) (*http.Response, error) {
	return r.Do("OPTIONS", path, nil)
}

// Delete retrieves the giving path and returns the response expected using a DELETE method.
func (r *Router) Delete(path string) (*http.Response, error) {
	return r.Do("DELETE", path, nil)
}

// Head retrieves the giving path and returns the response expected using a HEAD method.
func (r *Router) Head(path string) (*http.Response, error) {
	return r.Do("HEAD", path, nil)
}

// Get retrieves the giving path and returns the response expected using a GET method.
func (r *Router) Get(path string) (*http.Response, error) {
	return r.Do("GET", path, nil)
}

// Do performs the giving requests for a giving path with the provided body and returns the
// response for that method.
func (r *Router) Do(method string, path string, body io.ReadCloser) (*http.Response, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	// Create a ResponseRecorder for the giving
	responseRecoder := httptest.NewRecorder()

	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		defer wg.Done()

		// If we have no cache, then use the internal handler.
		if r.cache == nil {
			r.handler.ServeHTTP(responseRecoder, req)
			return
		}

		// Since we have a cache, attempt to serve the request, else use the
		// supplied http.Handler
		if err := r.cache.Serve(responseRecoder, req); err != nil {
			r.handler.ServeHTTP(responseRecoder, req)
		}
	}()

	wg.Wait()

	res := responseRecoder.Result()
	res.Request = req

	return res, nil
}
