Components
==========

Creating components is the core reason Gu exists as a package. It's primary aim is to provide a base library that allows rendering these components easily and efficiently.

Gu takes a different approach to components and how they should work. Gu does not try to be a React version in Go, but instead it takes advantage of the simple concepts that makes the Go language very powerful.

-	Composition over Inheritance, where components compose each other to create larger components, rather than using a form of inheritance or inter-logic where components are separately rendered and communicate with each other.

-	Interfaces compliance for upgrades, whereby components provide the capability to expose themselves to higher functionality or gain access to objects such has the internal caching and resource request `Fetch` objects. Additionally, this appraoch allows components to declare themselves reactive and notify themselves and their views of change to be updated by the driver.

By sticking to such basic ideas and principles, it allows construction of components with the standard constructs provided by the Go language to the maximum capability allowed.

Basics
------

Creating a component is comparatively easy, in that you are only required to meet a single interface by which the rendering markup for the component is retrieved.

Gu provides a `Renderable` interface which exposes a single method:

```go
type Renderable interface {
	Render() *trees.Markup
}
```

Any `Type` which implements the `Renderable` type is considered a Component and will be called when attached to the Gu view.

```go

import (
	"github.com/gu-io/gu/trees"
	"github.com/gu-io/gu/eventx"
	"github.com/gu-io/gu/trees/elems/events"
	"github.com/gu-io/gu/trees/elems"
	"github.com/gu-io/gu/trees/property"
)

// Greeter takes a name and generates a greeting.
type Greeter struct {
	Name string
}

// change updates the greeters name field.
func (g *Greeting) change(name string) {
	g.Name = name
}

// Render returns the Gu's tree structures which declares the markup for
// the greeter.
func (g *Greeting) Render() *trees.Markup {
	return elems.Div(
		property.ClassAttr("greeter"),
		elems.Div(
			property.ClassAttr("greeting"),
			elems.Text("Welcome to the %s!", g.Name),
		),
		elems.Div(
			property.ClassAttr("box", "input"),
			elems.Input(
				property.PlaceholderAttr("Enter your Name"),
				property.TypeAttr("text"),
				events.ChangeEvent(func(ev trees.EventObject, root *trees.Markup) {
					changeEvent := ev.Underling.(*eventx.ChangeEvent)
					g.change(changeEvent.Value)
				}),
			),
		),
	)
}
```

Composed Components
-------------------

Gu favors `Composition` over complexity. In other words, if you have a two or more components which work as one, instead of rendering each individually within it's own view, it is preferable to compose the core types and let a master type handle their rendering calls. By following this basic principle, communication flow and functional flow is simplified.

As demonstrated by the example below:

```go
import (
	"github.com/gu-io/gu/trees"
	"github.com/gu-io/gu/trees/elems"
	"github.com/gu-io/gu/trees/property"
)

// MenuItem defines a component which displays an entry in a menu list.
type MenuItem struct {
	Name string
	URI  string
}

// Render returns the markup for a MenuItem.
func (m *MenuItem) Render() *trees.Markup {
	return elems.ListItem(
		elems.Anchor(elems.Text(m.Name), property.HrefAttr(m.URI)),
	)
}

// Menu defines a component which displays a menu list.
type Menu struct {
	Items []MenuItem
}

// Menu returns the markup for a Menu list.
func (m *Menu) Render() *trees.Markup {
	ul := elems.UnorderedList()

	for _, item := range m.Items {
		item.Render().Apply(ul)
	}

	return ul
}

```

By having the Menu Component logically encapsulate/compose it's internal list of items, we can easily provide a simple approach to higher and more complex relationships between components. Though not all relationships fit this pattern, the majority can be found to match the pattern perfectly.

Reactive Components
-------------------

Gu heavily depends on interfaces as a means of extending the capability of Component. By meeting the `Reactive` interface, a component type can be made reactive, allowing the Gu view system to listen for update signals to update the rendered output.

```go

import (
	"github.com/gu-io/gu/trees"
	"github.com/gu-io/gu/eventx"
	"github.com/gu-io/gu/trees/elems/events"
	"github.com/gu-io/gu/trees/elems"
	"github.com/gu-io/gu/trees/property"
)

// Greeter takes a name and generates a greeting.
type Greeter struct {
	gu.Reactive
	Name string
}

// New returns a new instance of a Greeter.
func New() *Greeter {
	return &Greeter{
		Reactive: gu.NewReactive(),
	}
}

// change updates the greeters name field.
func (g *Greeting) change(name string) {
	g.Name = name
	g.Publish()
}

// Render returns the Gu's tree structures which declares the markup for
// the greeter.
func (g *Greeting) Render() *trees.Markup {
	return elems.Div(
		property.ClassAttr("greeter"),
		elems.Div(
			property.ClassAttr("greeting"),
			elems.Text("Welcome to the %s!", g.Name),
		),
		elems.Div(
			property.ClassAttr("box", "input"),
			elems.Input(
				property.PlaceholderAttr("Enter your Name"),
				property.TypeAttr("text"),
				events.ChangeEvent(func(ev trees.EventObject, root *trees.Markup) {
					changeEvent := ev.Underling.(*eventx.ChangeEvent)
					g.change(changeEvent.Value)
				}),
			),
		),
	)
}
```

Complex Components
------------------

More complex components can be found in the [Components](https://github.com/gu-io/components) directory and other packages which demonstrate different structures and design to achieve the component's functionality.
