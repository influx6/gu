package jigsaw

import (
	"github.com/gu-io/gu"
	"github.com/gu-io/gu/trees"
	"github.com/gu-io/gu/trees/elems"
	"github.com/gu-io/gu/trees/property"

	componentsbase "github.com/gu-io/gu/examples/boxes/components"
)

_ = componentsbase.Components.Register("jigsaw", func(attr map[string]string, template string) gu.Renderable {
	return NewJigsaw()
}, false)


// Jigsaw defines a component which implements the gu.Renderable interface.
type Jigsaw struct{
	gu.Reactive
}

// NewJigsaw returns a new instance of Jigsaw component.
func NewJigsaw() *Jigsaw {
  return &Jigsaw{
  	Reactive: gu.NewReactive(),
  }
}

// Render returns the markup for this Jigsaw component.
func (c Jigsaw) Render() *trees.Markup {
	return elems.Div(property.Class("component"))
}