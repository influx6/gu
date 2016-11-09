# Getting Started 
Gu is fundamentally a library built to handle view rendering, nothing more nor less.
When creating Gu, my main focus was on creating a solution that did not bind itself 
tightly to the perculiarities of either the frontend or backend, but allow each content to 
be renderable regardless of where it is loaded. 

The concepts in Gu are practially simple and rely majorly on a functional and interface based system,
where the structures which define the markup to be rendered are combined functionally and components 
are required to match/implement specific interfaces to enable specific features eg Reaction for page/view update.

Gu is in no way a Flux-like framework or library, nor does it provide complex structures and 
layed down paths by which such can be attained, It simply provides a baseline to render the desire output and 
gives the freedom for the developer to determine how his/her application data flow  works.


## The Guide
In this guide, we simply will be looking at Gu is its only two possible representation 
or in essense, the two ways you will ever use its constructs, which are as a `Page` and as `Component`. 
This allows you to grasp the available spectrum and capabilities provided by Gu and 
to evolve and define for yourself how you wish to mix and match Gu into your projects and workflow.

*Gu has evolved over its development lifetime alot, loosing and adding new ideas into the way it works but 
I have ensured to keep the tenant of simplicity and a non-intrusive architecture, that greatly meet the rendering 
needs of the developer and also grant freedom in how it should be used with in projects.*

### A Page
The Page approach is pretty much as its called, a "Page".  Combining ideas from [Goa](https://goa.design/) in 
the way its structures are defined and generated, Gu's [Design Package](./designs) defines structures which 
unlike [Goa](https://goa.design/)  do not provide any form of code generation but produces the desired effect of defining with intent 
the content expected on a page. In a sense, it is the organizational structure when wanting to define multipage or 
single page applications with gu.

Basically a page is a resource amongs many resources which depending on the availability
of a predefined route to be validated against, will be rendered. 
On the client side this means any resource which matches the current URL hash which gets
rendered onto the browser DOM, while on the client it is the URL provided to generated a complete 
html response to render a full page with its differents parts (i.e head, links, scripts, styles and body tags).

Pages as well listen for update requests from the things they render if this matches
the library reaction interface, which effectively provides a easy but simply way 
for developers to build some form of reactivity to changes in structures.

Below is an example of a single Resource rendered by a Resources Manager. The 
structures created where made to provide self description in the intent of what is 
expected, so that anyone could easily grasp the intent and expected result for that page.


```go
package main

import (
	"github.com/gopherjs/gopherjs/js"
	. "github.com/influx6/gu/design"
	. "github.com/influx6/gu/trees"
	. "github.com/influx6/gu/trees/attrs"
	. "github.com/influx6/gu/trees/elems"
	events "github.com/influx6/gu/trees/events"
	redom "github.com/influx6/gu/dom"
	"honnef.co/go/js/dom"
)

var _ = Resource(func() {

	UseRoute("/hello")

	DoTitle("Hello App")
	
	DoMarkup(func() *Markup {
		return Div(
			CSS(`
				${
					width:100%;
					height: 100%;
				}
			`, struct{ Size string }{Size: "130px"}),
			ID("hello"),
			Header1(
				Text("Hello"),
				events.Click(func(ev EventObject, tree *Markup) {
					js.Global.Call("alert", "I just got clicked, Yaay!!!")
				}, ""),
			),
			Span(Text("Click me")),
		)
	}, "", false)
})

func main(){
	New(&redom.DOMRenderer{
		Document: dom.GetWindow().Document(),
	}).Init()
}

```

The code flurishes with intentional declarations of its intent which one needs 
a surface detailing of what it does. Within Gu a resource is a single page, that 
is the very basic unit and can not be broken down any further not can resources 
call each other has they define seperate entities that could be rendered together 
if the route criteria matches, which is defaultly the case if no route is set.

By calling the gu design package `Resource` function we create a Resource which 
gets taken into the care of the created `ResourceManager` defined in the main function.
The `Resource` function returns the index for which the resource exists within 
its `ResourceManager` but simply provides a reference point but if required as 
a means of gaining access to the raw `Resource` object created by its manager.

```go
var _ = Resource(func() {})

```

The main function handles the creation of the `ResourceManager` by calling the `New`
function from the `gu/design` package which expects as an optional argument a 
renderer which will be called when on the client to handle page rendering during 
initial load and change of browsers URi. Gu already provides this and can be provided 
at load to get the client rendering adequately.

```go
	New(&redom.DOMRenderer{
		Document: dom.GetWindow().Document(),
	}).Init()
```

The internals of the resource definition are as simple, by using the design package 
functions we quickly define the title expected for that page and the route criteria 
which is matched by looking at the hash if on the client and through the url supplied 
through calling a render function which builds the adequate response of the complete 
page when dealing on the server.

```go
	UseRoute("/hello")

	DoTitle("Hello App")
```

The final piece and meat of the page is the intention declared about the markup to be 
rendered. As Gu provides a functional approach which creates structures underline 
that describes to the renderer what it expects on the page as its content, this is something
which will be heavily used has its the core of what makes it possible to render both 
on either the client or backend seperating a tight coupling between either ground

```go
	DoMarkup(func() *Markup {
		return Div(
			CSS(`
				${
					width:100%;
					height: 100%;
				}
			`, struct{ Size string }{Size: "130px"}),
			ID("hello"),
			Header1(
				Text("Hello"),
				events.Click(func(ev EventObject, tree *Markup) {
					js.Global.Call("alert", "I just got clicked, Yaay!!!")
				}, ""),
			),
			Span(Text("Click me")),
		)
	}, "", false)

```

The `DoMarkup` function accepts either a `Markup` structure or a function which 
returns a markup or a lsits of markup points, which are then organized with in 
the resource to make up its content. The reason we use this approach is because 
by the very name of the function, an intent of contents gets declared and allows 
us to either target or vary the organizational architecture to which a page is to 
be describe it.

The `DoMarkup` as well as an alternative function within the `gu/design` package,
called the `DoView`. In Gu

The `Markup` architecture found in the `gu/trees` package, defines structures which
varying and provides a complete set of HTML/HTML5 sturctures through functions 
that easily define either the CSS to be include and the attributes and markup content.

## A Component