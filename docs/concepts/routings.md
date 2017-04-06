Routing
=======

Gu provides a simplified routing system, which does not provide many bells and whistles found in routing solution these days. This is intentional, as complex routing is not expected to be needed.

Gu provides two routing concepts for the library:

-	View Routers The `View Routers`, also called `Resolvers` callback style, chaining strucuture where higher chains can effect the visibility of lower chains and also feed the lower chains pieces of routers which are left from their own path conditions. With this views can inform internal markup to hide/display themselves based on the supplied routers. This provides a clean approach to dealing with views and how the current paths affects those views.

-	Request Routers The `Request Routers` are the defactor means by which views and components can make request to retrieve resources from remote endpoints.

Request Router
==============

Router expresses a new system to allow components make requests for resources like database records, contents and assets from either the backend or frontend without much change of code. By exposing a structure which implements the Go `http` package `http.Handler`, this can be used by the router to service all request.

It is special in that for a App, only one ever exists and uses the supplied `http.Handler` and `router.Cache` implementing structure to resolve requests. This allows us to drastically move apps offline by providing a `http.Handler` that services requests from some offline store or the supplied cache, or implements the processes in making requests to the remote http endpoint for the resources.

One major benefit of this is, the fact we easily are able to use such a system on the server without much code change, since we can swap the supplied `http.Handler`, that passes all made requests to the running server without any actually use of a `http.Client`.

This was done to provide the flexibile and massive compatibility in both usage for either client or server codebase.

*Note: Now the `Cache` supplied is never updated by the router but is used to respond to request first before using the provided `http.Handler`, this approach safe guards the user has full control on how the cache operates and how it validates and invalidates requests, before allowing the router to proceed to the `http.Handler` to handle the request.*

Example
-------

The `gu/router` package lets you initialize a new `router.Router` which will use the supplied `http.handler` like below:

```go

import (
	"github.com/gu-io/gu/router"
	"github.com/gu-io/gu/router/cache/memorycache"
)


type serviceProvider struct{}

func (serviceProvider) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "reset":
		w.WriteHeader(http.StatusNoContent)
	case "count":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("1"))
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

var mainCache := memorycache.New("in-memory-store")
var mainRouter := router.NewRouter(serviceProvider{}, mainCache)

res, _ := mainRouter.Get("/count") // res.Status == http.StatusOK
res, _ := mainRouter.Get("/reset") // res.Status == http.StatusNoContent
res, _ := mainRouter.Get("/users") // res.Status == http.StatusBadRequest


```

All views and components will recieve access to the provided router through the implementation of the `RegisterService` interface.

View Routers
------------

View Routers are a construct built out in providing a means of chaining multiple path matchers which affect each other based on a callback system. Each router is restricted by the supplied path provided to it. These form allows us to use this type of routers to condition specific pieces of a components rendered output to either hide or show itself based on the validity of it's router to the current path. More so, others can use this to perform specific actions when this routers are trigger.

This provides a simple but powerful construct for components and views to interact with the external display easily.

Below are two example demonstrating the creation of a `View Reouter`:

1.	Demonstrate the usage of a given route and how paths can be tested against the resolver's internal matcher. It also demonstrates the usage of the pubsub capability of a Resolver in resolving a route path supplied by a `PushEvent`.

```go

import "github.com/gu-io/gu/router"

func main() {
	rx := router.New("/:id")

	// Test if the route matches specific path.
	params, rem, state := rx.Test("12")
	// Where:
	// params => are the parameters extracted from the test. {id: 12}
	// rem => remaining path if this route allows extensive routes.
	// state => boolean value which declares if the path matches.

	// Register callbacks for the success of the a match.
	rx.Done(func(px router.PushEvent) {
		// ....
	})

	// Register callbacks for the failure of the a match.
	rx.Failed(func(px router.PushEvent) {
		// ....
	})

	// Request the Resolver to resolve the provided route PushEvent.
	rx.Resolve(router.UseLocation("/12"))
}
```

1.	Demonstrate the usage of a chained routers and how they can be combined to create a reactive chain, where the parent route can pass values and remaining path's down to a lower router to resolve accordingly.

```go

import "github.com/gu-io/gu/router"

func main() {
	home := router.New("/home/*") // the /* tells the router to allow more paths.
	rx := router.New("/:id")

	home.Register(rx)

	home.Done(func(px router.PushEvent) {
		// px.Params{}, px.Rem: /12
		// DO something, we we passed
		//...
	})

	rx.Done(func(px router.PushEvent) {
		// DO something, we got a id
		// px.Params{id:12}, px.Rem: /12
		//...
	})

	rx.Failed(func(px router.PushEvent) {
		//...
	})

	home.Resolve(router.UseLocation("home/12"))
}
```

Conclusion
----------

By combining these simple concepts, it should provide a flexible approach in routing for components, views and requesting resources using the Gu library.
