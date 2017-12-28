package webcache

import (
	"bytes"
	"errors"
	"net/url"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gu-io/gu/router/cache"
)

// ObjectToRequest converts a object to a cache.Request.
func ObjectToRequest(o *js.Object) (cache.Request, error) {
	if o == nil || o == js.Undefined {
		return cache.Request{}, errors.New("Invalid/Undefined Object")
	}

	reqClone := o.Call("clone")

	body := make(chan []byte, 0)

	o.Call("text").Call("then", func(bd string) {
		go func() {
			body <- []byte(bd)
		}()
	})

	uri, _ := url.Parse(o.Get("url").String())

	return cache.Request{
		Underline: reqClone,
		Body:      *bytes.NewBuffer(<-body),
		URL:       uri,
		Path:      o.Get("url").String(),
		Method:    o.Get("method").String(),
		Headers:   ObjectToMap(o.Get("headers")),
	}, nil
}

// ObjectToResponse converts a object to a cache.Response.
func ObjectToResponse(o *js.Object) (cache.Response, error) {
	if o == nil || o == js.Undefined {
		return cache.Response{}, errors.New("Invalid/Undefined Object")
	}

	resClone := o.Call("clone")

	body := make(chan []byte, 0)

	o.Call("text").Call("then", func(bd string) {
		go func() {
			body <- []byte(bd)
		}()
	})

	return cache.Response{
		Body:      *bytes.NewBuffer(<-body),
		Underline: resClone,
		Status:    o.Get("status").Int(),
		Type:      o.Get("type").String(),
		Method:    o.Get("method").String(),
		Headers:   ObjectToMap(o.Get("headers")),
	}, nil
}

// AllObjectToResponse returns a slice of WebResponses from a slice
// of js.Object.
func AllObjectToResponse(o []*js.Object) []cache.Response {
	res := make([]cache.Response, 0)

	for _, ro := range o {
		if rq, err := ObjectToResponse(ro); err == nil {
			res = append(res, rq)
		}
	}

	return res
}

// AllObjectToRequest returns a slice of WebResponses from a slice
// of js.Object.
func AllObjectToRequest(o []*js.Object) []cache.Request {
	res := make([]cache.Request, 0)

	for _, ro := range o {
		if rq, err := ObjectToRequest(ro); err == nil {
			res = append(res, rq)
		}
	}

	return res
}

// ObjectToRequests returns a list of cache.Request objects from the provided object.
func ObjectToRequests(o *js.Object) []cache.Request {
	res := make([]cache.Request, 0)

	for i := 0; i < o.Length(); i++ {
		item := o.Index(i)

		if rq, err := ObjectToRequest(item); err == nil {
			res = append(res, rq)
		}
	}

	return res
}

// MapToHeaders converts a map into a js.Headers structure.
func MapToHeaders(res map[string]string) *js.Object {
	header := js.Global.Get("Headers").New()

	for name, val := range res {
		header.Call("set", name, val)
	}

	return header
}

// ResponseToJSResponse converts a object to a js.Response.
func ResponseToJSResponse(res *cache.Response) *js.Object {
	if res.Underline != nil {
		return res.Underline.(*js.Object)
	}

	body := js.NewArrayBuffer(res.Body.Bytes())
	// bodyBlob := js.Global.Get("Blob").New(body)

	res.Underline = js.Global.Get("Response").New(body, map[string]interface{}{
		"status":  res.Status,
		"headers": MapToHeaders(res.Headers),
	})

	return res.Underline.(*js.Object)
}

// RequestToJSRequest converts a object to a js.Request.
func RequestToJSRequest(res *cache.Request) *js.Object {
	if res.Underline != nil {
		return res.Underline.(*js.Object)
	}

	res.Underline = js.Global.Get("Request").New(res.URL, map[string]interface{}{
		"body":    res.Body.Bytes(),
		"method":  res.Method,
		"headers": MapToHeaders(res.Headers),
	})

	return res.Underline.(*js.Object)
}

// ObjectToMap returns a map from a giving object.
func ObjectToMap(o *js.Object) map[string]string {
	res := make(map[string]string)

	for i := 0; i < o.Length(); i++ {
		item := o.Index(i)
		itemName := item.String()
		res[itemName] = o.Get(itemName).String()
	}

	return res
}

// ObjectToStringList returns a map from a giving object.
func ObjectToStringList(o *js.Object) []string {
	res := make([]string, 0)

	for i := 0; i < o.Length(); i++ {
		item := o.Index(i)
		res = append(res, item.String())
	}

	return res
}

// ObjectToList returns a map from a giving object.
func ObjectToList(o *js.Object) []*js.Object {
	res := make([]*js.Object, 0)

	for i := 0; i < o.Length(); i++ {
		item := o.Index(i)
		res = append(res, item)
	}

	return res
}
