package localcache

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gu-io/gu/router/cache"
)

// API defines a structure which implements the cache.Cache interface.
type API struct {
	name  string
	pairs []cache.WebPair
}

// New returns a new instance of the API struct.
func New(name string) *API {
	return (&API{name: name}).init()
}

// String returns a json version of the internal array of pairs.
func (a *API) String() string {
	jsx, err := json.Marshal(a.pairs)
	if err != nil {
		return ""
	}

	return string(jsx)
}

// Empty deletes all giving requests from the underline cache.
func (a *API) Empty() error {
	a.pairs = nil
	return nil
}

// All returns all the pairs of requests which have been added into the cache in
func (a *API) All() ([]cache.WebPair, error) {
	return a.pairs[0:], nil
}

// DeleteRequest calls the underline cache.Cache.Delete.
func (a *API) DeleteRequest(w cache.Request) error {
	for index, pair := range a.pairs {
		if pair.Request.URL == w.URL {
			a.pairs = append(a.pairs[0:index], a.pairs[index+1:]...)
			return nil
		}
	}

	return errors.New("Request not found")
}

// Delete removes the giving path from the underline cache if found.
func (a *API) Delete(path string) error {
	return a.DeleteRequest(cache.Request{
		Path: path,
	})
}

// AddData adds the giving data object into the cache.
func (a *API) AddData(req string, res []byte) error {
	a.pairs = append(a.pairs, cache.WebPair{
		Request:  cache.Request{Path: req, Method: "GET"},
		Response: cache.Response{Method: "GET", Body: *bytes.NewBuffer(res)},
	})

	return nil
}

// Add adds the giving response object into the cache.
func (a *API) Add(req string, res *http.Response) error {
	resp, reqs := cache.HTTPResponseToResponse(res)
	if reqs == nil {
		reqs = &cache.Request{Path: req, Method: "GET"}
	}

	a.pairs = append(a.pairs, cache.WebPair{
		Request:  *reqs,
		Response: *resp,
	})

	return nil
}

// Serve attempts to find the request and serve the response into the provided
// http.ResponseWriter.
func (a *API) Serve(w http.ResponseWriter, r *http.Request) error {
	// 1. Attempt to get full URI
	// 2. Attempt to get only path

	res, _, err := a.Get(r.URL.String())
	if err != nil {
		res, _, err = a.Get(r.URL.Path)
		if err != nil {
			return err
		}
	}

	if res.Body.Len() == 0 {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res.Body.Bytes())
	return nil
}

// Put calls the internal caches.Cache.Put function matching against the
func (a *API) Put(req cache.Request, res cache.Response) error {
	a.pairs = append(a.pairs, cache.WebPair{
		Request:  req,
		Response: res,
	})

	a.sync()

	return nil
}

// PutPath calls the internal caches.Cache.Put function matching against the
func (a *API) PutPath(path string, res cache.Response) error {
	var req cache.Request
	req.Path = path

	uri, _ := url.Parse(path)
	req.URL = uri

	a.pairs = append(a.pairs, cache.WebPair{
		Request:  req,
		Response: res,
	})

	a.sync()

	return nil
}

// GetRequest calls CacheAPI.Match and passing in a default MatchAttr value.
func (a *API) GetRequest(w cache.Request) (cache.Response, error) {
	for _, pair := range a.pairs {
		if pair.Request.Path == w.Path {
			return pair.Response, nil
		}
	}

	return cache.Response{}, errors.New("Request not found")
}

// Get calls CacheAPI.MatchPath and passing in a default MatchAttr value.
func (a *API) Get(path string) (cache.Request, cache.Response, error) {
	for _, pair := range a.pairs {
		if pair.Request.Path == path {
			return pair.Request, pair.Response, nil
		}
	}

	return cache.Request{
		Path:   path,
		Method: "GET",
	}, cache.Response{}, errors.New("Request not found")
}

// init intializes the cache and its dependencies.
func (a *API) init() *API {
	a.unsync()
	return a
}

// sync stores the giving request and response pairs into a locastorage item
// using the provided api name.
func (a *API) sync() {
	if js.Global == nil || js.Global == js.Undefined {
		return
	}

	localStorage := js.Global.Get("localstorage")
	if localStorage == nil || localStorage == js.Undefined {
		return
	}

	localStorage.Call("setItem", a.name, a.String())
}

// unsync deletes the api's requests from the localstorage cache.
func (a *API) unsync() {
	if js.Global == nil || js.Global == js.Undefined {
		return
	}

	localStorage := js.Global.Get("localstorage")
	if localStorage == nil || localStorage == js.Undefined {
		return
	}

	localStorage.Call("removeItem", a.name)
}
