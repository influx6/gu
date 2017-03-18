package memorycache

import (
	"encoding/json"
	"errors"
)

// API defines a structure which implements the shell.Cache interface.
type API struct {
	name  string
	pairs []shell.WebPair
}

// New returns a new instance of the API struct.
func New(name string) *API {
	return &API{name: name}
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
func (a *API) All() ([]shell.WebPair, error) {
	return a.pairs[0:], nil
}

// Delete calls the underline CacheAPI.DeleteRequest.
func (a *API) Delete(w shell.WebRequest) error {
	for index, pair := range a.pairs {
		if pair.Request.URL == w.URL {
			a.pairs = append(a.pairs[0:index], a.pairs[index+1:]...)
			return nil
		}
	}

	return errors.New("Request not found")
}

// Put calls the internal caches.Cache.Put function matching against the
func (a *API) Put(req shell.WebRequest, res shell.WebResponse) error {
	a.pairs = append(a.pairs, shell.WebPair{
		Request:  req,
		Response: res,
	})

	a.sync()

	return nil
}

// PutPath calls the internal caches.Cache.Put function matching against the
func (a *API) PutPath(path string, res shell.WebResponse) error {
	var req shell.WebRequest
	req.URL = path
	req.Cache = shell.DefaultStrategy

	a.pairs = append(a.pairs, shell.WebPair{
		Request:  req,
		Response: res,
	})

	return nil
}

// GetRequest calls CacheAPI.Match and passing in a default MatchAttr value.
func (a *API) GetRequest(w shell.WebRequest) (shell.WebResponse, error) {
	for _, pair := range a.pairs {
		if pair.Request.URL == w.URL {
			return pair.Response, nil
		}
	}

	return shell.WebResponse{}, errors.New("Request not found")
}

// GetPath calls CacheAPI.MatchPath and passing in a default MatchAttr value.
func (a *API) GetPath(path string) (shell.WebRequest, shell.WebResponse, error) {
	for _, pair := range a.pairs {
		if pair.Request.URL == path {
			return pair.Request, pair.Response, nil
		}
	}

	return shell.WebRequest{
		URL:    path,
		Method: "GET",
	}, shell.WebResponse{}, errors.New("Request not found")
}

// GetManifest returns the giving request, response associated with the given
func (a *API) GetManifest(mn shell.ManifestAttr) (shell.WebRequest, shell.WebResponse, error) {
	req := mn.WebRequest()

	for _, pair := range a.pairs {
		if pair.Request.URL == req.URL {
			return pair.Request, pair.Response, nil
		}
	}

	return req, shell.WebResponse{}, errors.New("Request not found")
}
