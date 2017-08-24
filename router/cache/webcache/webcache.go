package webcache

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gu-io/gu/router/cache"
)

// ErrInvalidState is returned when the cache api does not exists in the global
// context.
var ErrInvalidState = errors.New("Cache Not Found")

// NewCacheResponse defines the response returned when a cache is request from the
// API.
type newCacheResponse struct {
	Error error     `json:"error"`
	Cache *CacheAPI `json:"cache"`
}

// ErrCacheNotFound defines the error returned when the cached desired is not found.
var ErrCacheNotFound = errors.New("Cache Not Found")

// New returns a new instance of the CacheAPI using the browsers cache API.
// API found in webkit browsers. https://developer.mozilla.org/en-US/docs/Web/API/Cache.
func New(cacheName string) (*CacheAPI, error) {
	if js.Global == nil || js.Global == js.Undefined {
		return nil, ErrInvalidState
	}

	caches := js.Global.Get("caches")

	if caches == nil || caches == js.Undefined {
		return nil, ErrInvalidState
	}

	openReq := caches.Call("open", cacheName)
	if openReq == js.Undefined || openReq == nil {
		return nil, ErrCacheNotFound
	}

	var opVal newCacheResponse
	res := make(chan newCacheResponse, 0)

	openReq.Call("then", func(o *js.Object) {
		go func() {
			res <- newCacheResponse{Cache: NewCacheAPI(o), Error: nil}
		}()
	}, func(o *js.Object) {
		go func() {
			res <- newCacheResponse{Error: errors.New(o.String())}
		}()
	})

	opVal = <-res
	return opVal.Cache, opVal.Error
}

//================================================================================

// CacheAPI defines individual cache item received from a call to the internal cache API.
type CacheAPI struct {
	*js.Object
}

// NewCacheAPI returns a new instance of the CacheAPI.
func NewCacheAPI(o *js.Object) *CacheAPI {
	var c CacheAPI
	c.Object = o
	return &c
}

// AddAll calls the internal caches.Cache.Add and caches.Cache.AddAll function matching against the
// request string/strings.
func (c *CacheAPI) AddAll(request ...string) error {
	resChan := make(chan error, 0)

	var isList bool

	if len(request) > 1 {
		isList = true
	}

	if !isList {
		c.Call("add", request[0]).Call("then", func() {
			go func() {
				close(resChan)
			}()
		}).Call("catch", func(err *js.Object) {
			go func() {
				resChan <- errors.New(err.String())
			}()
		})

		return <-resChan
	}

	c.Call("addAll", request).Call("then", func() {
		go func() {
			close(resChan)
		}()
	}).Call("catch", func(err *js.Object) {
		go func() {
			resChan <- errors.New(err.String())
		}()
	})

	return <-resChan
}

// AddData adds the giving data object into the cache.
func (c *CacheAPI) AddData(path string, data []byte) error {
	req := cache.Request{Path: path, Method: "GET"}
	res := cache.Response{Method: "GET", Body: *bytes.NewBuffer(data)}

	return c.Put(req, res)
}

// Add the giving path and response into the cache.
func (c *CacheAPI) Add(request string, resp *http.Response) error {
	resChan := make(chan error, 0)

	res, reqs := cache.HTTPResponseToResponse(resp)
	if reqs == nil {
		reqs = &cache.Request{Path: request, Method: "GET"}
	}

	resObj := ResponseToJSResponse(res)

	c.Call("put", reqs, resObj).Call("catch", func(err *js.Object) {
		go func() {
			resChan <- errors.New(err.String())
		}()
	})

	return <-resChan
}

// Put calls the internal caches.Cache.Put function matching against the
// request string.
func (c *CacheAPI) Put(request cache.Request, res cache.Response) error {
	resChan := make(chan error, 0)

	resObj := ResponseToJSResponse(&res)
	reqObj := RequestToJSRequest(&request)

	c.Call("put", reqObj, resObj).Call("then", func(_ *js.Object) {
		go func() {
			close(resChan)
		}()
	}).Call("catch", func(err *js.Object) {
		go func() {
			resChan <- errors.New(err.String())
		}()
	})

	return <-resChan
}

// Get calls CacheAPI.MatchPath and passing in a default MatchAttr value.
func (c *CacheAPI) Get(request string) (cache.Request, cache.Response, error) {
	wr := cache.Request{
		Method: "GET",
		Path:   request,
	}

	res, err := c.MatchPath(request, nil)
	if err != nil {
		return wr, res, err
	}

	return wr, res, nil
}

// CacheResponse defines the response returned when a cache request is made.
type cacheResponse struct {
	Error    error          `json:"error"`
	Response cache.Response `json:"response"`
}

// CacheResponseChannel defines a channel type for the response to be received
// over a request.
type cacheResponseChannel chan cacheResponse

// MatchPath calls the internal caches.Cache.Match function matching against the
// request string.
func (c *CacheAPI) MatchPath(request string, attr map[string]interface{}) (cache.Response, error) {
	resChn := make(cacheResponseChannel, 0)

	c.Call("match", request, attr).Call("then", func(response *js.Object) {
		res, err := ObjectToResponse(response)
		if err != nil {
			go func() {
				resChn <- cacheResponse{Error: err}
			}()
			return
		}

		go func() {
			resChn <- cacheResponse{
				Response: res,
			}
		}()
	}).Call("catch", func(err *js.Object) {
		go func() {
			resChn <- cacheResponse{Error: errors.New(err.String())}
		}()
	})

	opVal := <-resChn
	return opVal.Response, opVal.Error
}

// Serve attempts to find the request and serve the response into the provided
// http.ResponseWriter.
func (c *CacheAPI) Serve(w http.ResponseWriter, r *http.Request) error {
	// 1. Attempt to get full URI
	// 2. Attempt to get only path

	res, err := c.MatchPath(r.URL.String(), nil)
	if err != nil {
		res, err = c.MatchPath(r.URL.Path, nil)
		if err != nil {
			return err
		}
	}

	w.WriteHeader(res.Status)
	w.Write(res.Body.Bytes())
	return nil
}

// Match calls the internal caches.Cache.Match function matching against the
// request string for a js.Cache.Match.
func (c *CacheAPI) Match(request cache.Request, attr map[string]interface{}) (cache.Response, error) {
	resChan := make(cacheResponseChannel, 0)

	c.Call("match", RequestToJSRequest(&request), attr).Call("then", func(response *js.Object) {
		res, err := ObjectToResponse(response)
		if err != nil {
			go func() {
				resChan <- cacheResponse{Error: err}
			}()
			return
		}

		go func() {
			resChan <- cacheResponse{
				Response: res,
				Error:    nil,
			}
		}()
	}).Call("catch", func(err *js.Object) {
		go func() {
			resChan <- cacheResponse{Error: errors.New(err.String())}
		}()
	})

	opVal := <-resChan
	return opVal.Response, opVal.Error
}

// CacheAllResponse defines the slice of responses returned when a cache request
// is made to matchAll.
type cacheAllResponse struct {
	Error    error            `json:"error"`
	Response []cache.Response `json:"response"`
}

// CacheAllResponseChannel defines a channel type for the response to be received
// over a request to js.Cache.MatchAll.
type cacheAllResponseChannel chan cacheAllResponse

// MatchAllPath calls the internal caches.Cache.MatchAll function matching against the
// request string for a js.Cache.Match.
func (c *CacheAPI) MatchAllPath(request string, attr map[string]interface{}) ([]cache.Response, error) {
	resChan := make(cacheAllResponseChannel, 0)

	c.Call("match", request, attr).Call("then", func(responses *js.Object) {
		go func() {
			resChan <- cacheAllResponse{
				Response: AllObjectToResponse(ObjectToList(responses)),
			}
		}()
	}).Call("catch", func(err *js.Object) {
		go func() {
			resChan <- cacheAllResponse{Error: errors.New(err.String())}
		}()
	})

	opVal := <-resChan
	return opVal.Response, opVal.Error
}

// MatchAll calls the internal caches.Cache.MatchAll function matching against the
// request string for a js.Cache.Match.
func (c *CacheAPI) MatchAll(request cache.Request, attr map[string]interface{}) ([]cache.Response, error) {
	resChan := make(cacheAllResponseChannel, 0)

	c.Call("match", RequestToJSRequest(&request), attr).Call("then", func(responses *js.Object) {
		go func() {
			resChan <- cacheAllResponse{
				Response: AllObjectToResponse(ObjectToList(responses)),
			}
		}()
	}).Call("catch", func(err *js.Object) {
		go func() {
			resChan <- cacheAllResponse{Error: errors.New(err.String())}
		}()
	})

	opVal := <-resChan
	return opVal.Response, opVal.Error
}

// Delete calls the underline CacheAPI.DeleteRequest.
func (c *CacheAPI) Delete(request string) error {
	return c.DeletePath(request, nil)
}

// ErrRequestDoesNotExists is returned when the giving request does not exists in
// the catch.
var ErrRequestDoesNotExists = errors.New("Request not existing")

// DeletePath deletes the given path from the cache if found.
func (c *CacheAPI) DeletePath(request string, attr map[string]interface{}) error {
	resChn := make(chan error, 0)

	c.Call("delete", request, attr).Call("then", func(response *js.Object) {
		if response.Bool() {
			go func() {
				close(resChn)
			}()
			return
		}

		go func() {
			resChn <- ErrRequestDoesNotExists
		}()

	}).Call("catch", func(err *js.Object) {
		go func() {
			resChn <- errors.New(err.String())
		}()
	})

	return <-resChn
}

// DeleteRequest calls the internal caches.Cache.Delete function matching against the
// request for a js.Cache.Delete.
func (c *CacheAPI) DeleteRequest(request cache.Request, attr map[string]interface{}) error {
	resChn := make(chan error, 0)

	if request.Underline != nil {
		c.Call("delete", request.Underline, attr).Call("then", func(response *js.Object) {
			if response.Bool() {
				go func() {
					close(resChn)
				}()
				return
			}

			go func() {
				resChn <- ErrRequestDoesNotExists
			}()

		}).Call("catch", func(err *js.Object) {
			go func() {
				resChn <- errors.New(err.String())
			}()
		})

		return <-resChn
	}

	c.Call("delete", request, attr).Call("then", func(response *js.Object) {
		if response.Bool() {
			go func() {
				close(resChn)
			}()
			return
		}

		go func() {
			resChn <- ErrRequestDoesNotExists
		}()

	}).Call("catch", func(err *js.Object) {
		go func() {
			resChn <- errors.New(err.String())
		}()
	})

	return <-resChn
}

// cacheKeys defines the response received when all keys of the cache is retrieved
// or when filtered by a request value.
type cacheKeys struct {
	Error    error    `json:"error"`
	Response []string `json:"response"`
}

type cacheAllRequest struct {
	Error    error           `json:"error"`
	Response []cache.Request `json:"response"`
}

// Keys returns a slice of all cache keys for all request added in the order they
// were added.
func (c *CacheAPI) Keys(request interface{}, attr map[string]interface{}) ([]cache.Request, error) {
	resChn := make(chan cacheAllRequest, 0)

	if request == nil {
		c.Call("keys").Call("then", func(response *js.Object) {
			go func() {
				resChn <- cacheAllRequest{
					Response: ObjectToRequests(response),
				}
			}()
		}).Call("catch", func(err *js.Object) {
			go func() {
				resChn <- cacheAllRequest{
					Error: errors.New(err.String()),
				}
			}()
		})

		opVal := <-resChn
		return opVal.Response, opVal.Error
	}

	switch ro := request.(type) {
	case string:
		c.Call("keys", ro, attr).Call("then", func(response *js.Object) {
			go func() {
				resChn <- cacheAllRequest{
					Response: ObjectToRequests(response),
				}
			}()
		}).Call("catch", func(err *js.Object) {
			go func() {
				resChn <- cacheAllRequest{
					Error: errors.New(err.String()),
				}
			}()
		})
	case cache.Request:
		c.Call("keys", RequestToJSRequest(&ro), attr).Call("then", func(response *js.Object) {
			go func() {
				resChn <- cacheAllRequest{
					Response: ObjectToRequests(response),
				}
			}()
		}).Call("catch", func(err *js.Object) {
			go func() {
				resChn <- cacheAllRequest{
					Error: errors.New(err.String()),
				}
			}()
		})
	case *cache.Request:
		c.Call("keys", RequestToJSRequest(ro), attr).Call("then", func(response *js.Object) {
			go func() {
				resChn <- cacheAllRequest{
					Response: ObjectToRequests(response),
				}
			}()
		}).Call("catch", func(err *js.Object) {
			go func() {
				resChn <- cacheAllRequest{
					Error: errors.New(err.String()),
				}
			}()
		})
	default:
		go func() {
			resChn <- cacheAllRequest{
				Error: errors.New("Request type not supported in cache"),
			}
		}()
	}

	opVal := <-resChn
	return opVal.Response, opVal.Error
}

// All returns all the pairs of requests which have been added into the cache in
// the order they were added.
func (c *CacheAPI) All() ([]cache.WebPair, error) {
	keys, err := c.Keys(nil, nil)
	if err != nil {
		return nil, err
	}

	var pairs []cache.WebPair

	for _, req := range keys {
		//
		res, err := c.MatchPath(req.Path, nil)
		if err != nil {
			return pairs, err
		}

		pairs = append(pairs, cache.WebPair{
			Request:  req,
			Response: res,
		})
	}

	return pairs, nil
}

// Empty deletes all giving requests from the underline cache.
func (c *CacheAPI) Empty() error {
	keys, err := c.Keys(nil, nil)

	if err != nil {
		return err
	}

	for _, key := range keys {
		if err := c.DeleteRequest(key, nil); err != nil && err != ErrRequestDoesNotExists {
			fmt.Printf("Delete Error: %q\n", err.Error())
		}
	}

	return nil
}
