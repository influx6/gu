package router

import (
	"github.com/gu-io/gu/trees"
	"github.com/influx6/faux/pattern"
)

// Handler defines a function type for a Resolver subcriber.
type Handler func(PushEvent)

// Resolvable defines an interface for a resolvable type object.
type Resolvable interface {
	Resolve(PushEvent)
}

// Resolver defines an interface for a type that resolves a
// provided instance of a dispatch.PushEvent.
type Resolver interface {
	Resolvable

	Flush()
	Pattern() string
	Only(string, ...trees.SwitchMorpher) ResolveMorpher
	Register(Resolver)
	Done(Handler) Resolver
	Failed(Handler) Resolver
	Test(string) (map[string]string, string, bool)
}

// ResolveMorpher defines an interface for a Resolver which can morph a trees.Markup
// state.
type ResolveMorpher interface {
	Resolver
	trees.Morpher
}

// NewResolver returns a new instance of a structure that matches
// the Resolver interface.
func NewResolver(path string) Resolver {
	var br basicResolver
	if path != "" {
		br.matcher = URIMatcher(path)
	}

	return &br
}

// basicResolver defines a struct that implements
type basicResolver struct {
	children []Resolver
	fails    []Handler
	subs     []Handler
	matcher  pattern.URIMatcher
}

// Flush resets the subscriptions and children lists to empty.
func (b *basicResolver) Flush() {
	b.subs = nil
	b.fails = nil
	b.children = nil
}

// Pattern returns the giving path pattern used by this resolver.
func (b *basicResolver) Pattern() string {
	return b.matcher.Pattern()
}

// Only returns a new instance of a ResolveMorpher which can change the rendering
// state of the tree markup.
func (b *basicResolver) Only(pattern string, m ...trees.SwitchMorpher) ResolveMorpher {
	var morpher trees.SwitchMorpher

	if len(m) != 0 {
		morpher = m[0]
	} else {
		morpher = &trees.RemoveMorpher{}
	}

	morph := NewRouting(pattern, morpher)
	b.Register(morph)
	return morph
}

// Register adds a resolver into the list which will get triggerd
// when this resolver gets triggered, they will recieve a new PushEvent
// made out of the remaining path from the PushEvent received by this.
func (b *basicResolver) Register(r Resolver) {
	b.children = append(b.children, r)
}

// Test evaluates the giving path returning a map of parameters, a string of path left unmatch if
// possible(incase the resolver allows extra path lines within its pattern using the /*) and a boolean
// indicating when the path was indeed matched succesfully or if it was a total failure.
func (b *basicResolver) Test(path string) (map[string]string, string, bool) {
	return b.matcher.Validate(path)
}

// Resolve takes a `dispatch.PushEvent` instance, matches the content
// against it's own matcher, if its a match, it calls its subscribers
// then takes the remaining path left after removing its own matching
// piece and sends that off to its children, if there exists any remaining
// path that is.
func (b *basicResolver) Resolve(path PushEvent) {
	if b.matcher == nil {

		// Notify the subscribers.
		for _, sub := range b.subs {
			sub(path)
		}

		// Notify the kids with what is left in the PushEvent.
		for _, child := range b.children {
			child.Resolve(path)
		}

		return
	}

	params, rem, ok := b.matcher.Validate(path.Rem)
	if !ok {

		// Notify the fail subscribers.
		for _, sub := range b.fails {
			sub(path)
		}

		return
	}

	// Copy over the parameter left from the previous path.
	for key, val := range path.Params {
		params[key] = val
	}

	newPath := PushEvent{
		Rem:    rem,
		Params: params,
		Hash:   path.Hash,
		Host:   path.Host,
		Path:   path.Path,
		From:   path.Rem,
		To:     rem,
	}

	// Notify the subscribers.
	for _, sub := range b.subs {
		sub(newPath)
	}

	// Notify the kids with what is left in the PushEvent.
	for _, child := range b.children {
		child.Resolve(newPath)
	}
}

// Failed adds a function to the failed subscription list for this
// resolver.
func (b *basicResolver) Failed(sub Handler) Resolver {
	b.fails = append(b.fails, sub)
	return b
}

// Done adds a function to the subscription list for this
// resolver.
func (b *basicResolver) Done(sub Handler) Resolver {
	b.subs = append(b.subs, sub)
	return b
}
