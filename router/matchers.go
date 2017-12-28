package router

import (
	"sync"

	"github.com/influx6/faux/pattern"
)

// matchers define a lists of pattern associated with url match validators.
var matchers = struct {
	m  map[string]pattern.URIMatcher
	ml sync.RWMutex
}{
	m: make(map[string]pattern.URIMatcher),
}

// URIMatcher returns a new uri matcher if it has not being already creatd.
func URIMatcher(path string) pattern.URIMatcher {
	matchers.ml.RLock()
	mk, ok := matchers.m[path]
	matchers.ml.RUnlock()

	if !ok {
		m := pattern.New(path)
		matchers.ml.Lock()
		matchers.m[path] = m
		matchers.ml.Unlock()
		return m
	}

	return mk
}
