package router

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gu-io/gu/notifications"
	"github.com/influx6/faux/pattern"
)

//==============================================================================

// WrapHandler returns a function of type PathHandler by wrapping the
// provided no argument function in one.
func WrapHandler(fx func()) func(PushEvent) {
	return func(_ PushEvent) {
		fx()
	}
}

// PushDirectiveEvent defines a event which is used to declare the switching
// of the route to another path provided.
type PushDirectiveEvent struct {
	To string
}

// PushEvent represent the current path and hash values.
//
//@notification:event
type PushEvent struct {
	Host   string
	Hash   string
	Path   string
	Rem    string
	To     string
	From   string
	Params map[string]string
}

// NewPushEvent returns PushEvent based on the path string provided.
func NewPushEvent(path string, useHash bool) (PushEvent, error) {
	ups, err := url.Parse(path)
	if err != nil {
		return PushEvent{}, err
	}

	hash := strings.TrimSpace(ups.Fragment)
	if hash == "" {
		hash = "/#"
	}

	var target string

	if useHash {
		target = hash
	} else {
		target = ups.Path
	}

	return PushEvent{
		Hash:   hash,
		Path:   ups.Path,
		Host:   ups.Host,
		Rem:    target,
		From:   ups.String(),
		Params: make(map[string]string),
	}, nil
}

// String returns the hash and path.
func (p PushEvent) String() string {
	return fmt.Sprintf("%s%s", pattern.TrimEndSlashe(p.Path), p.Hash)
}

// UseLocation returns the current Path associated with the provided path.
// Using the complete path has the Remaining path value.
func UseLocation(path string) PushEvent {
	event, err := NewPushEvent(path, false)
	if err != nil {
		return PushEvent{}
	}

	return event
}

// UseLocationHash returns the current Path associated with the provided path.
// Using the hash path has the Remaining path value.
func UseLocationHash(path string) PushEvent {
	event, err := NewPushEvent(path, true)
	if err != nil {
		return PushEvent{}
	}

	return event
}

// ListenAndResolve takes the giving pattern, matches it against changes provided by
// the current PathObserver, if the full URL(i.e Path+Hash) matches then fires
// the provided function.
func ListenAndResolve(pattern string, fx func(PushEvent), fail func(PushEvent)) Resolver {
	resolver := New(pattern)
	resolver.Done(fx)
	resolver.Failed(fail)

	notifications.Subscribe(NewPushEventHandler(resolver.Resolve))

	return resolver
}

// ListenFor takes the giving pattern, matches it against changes provided by
// the current PathObserver, if the full URL(i.e Path+Hash) matches then fires
// the provided function.
func ListenFor(hash bool, pattern string, fx func(PushEvent), fail func(PushEvent)) {
	matcher := URIMatcher(pattern)

	notifications.Subscribe(NewPushEventHandler(func(p PushEvent) {
		var target string

		if hash {
			target = p.Hash
		} else {
			target = p.String()
		}

		if params, rem, found := matcher.Validate(target); found {
			for key, val := range p.Params {
				if _, ok := params[key]; !ok {
					params[key] = val
				}
			}

			fx(PushEvent{
				Params: params,
				Rem:    rem,
				Path:   p.Path,
				Hash:   p.Hash,
				Host:   p.Host,
				From:   p.From,
			})

			return
		}

		if fail != nil {
			fail(p)
		}
	}))
}
