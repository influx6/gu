Routing
=======

Gu provides a simplified routing system, which does not provide many bells and whistles found in routing solution these days. This is intended has no major requirement are required for more complex system (Maybe you mena to say "This is intentional, as complex routing is not expected to be needed").

Gu provides `Resolvers` which is part of the `routing` package.

Gu's `Resolvers` combine a pubsub system and a pattern matcher, which checks the validity of a provided route and then notifies any registered callbacks ("callbacks" is too generic. might need a better word here...) of either the success or failure of the route matching. If the route passed the matcher, the callback receives a structure (might be better to say "data struct") which contains parameter information and other details as related to the route which was matched. This then allows components and views to react accordingly.

Examples
--------

Below are examples of using the `routing` package to create Resolvers and how multiple Resolvers can be used to create nested routing where one level is nested below another to form part of the expected route/path.

-	Simple Routes. This example demonstrates the creation of a Resolver for a given route and how paths can be tested against the resolver's internal matcher. It also demonstrates the  usage of the pubsub capability of a Resolver in resolving a route path structure called a `PushEvent`.

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

-	Level Routes.  This example demonstrates how two independent routes can be connected together to create a level routing where the first Resolver feeds the next registered Resolvers it's successfully matched `PushEvent` which then is passed with the `Remaining` path left of the successfully matched path. (Huh ? This sentence is highly confusing ! )

```go

import "github.com/gu-io/gu/router"

func main() {
	home := router.New("/home/*") // the /* tells the router to allow more paths.
	rx := router.New("/:id")

	home.Register(rx)

	home.Done(func(px router.PushEvent) {
		// px.Params{}, px.Rem: /12
		//...
	})

	rx.Done(func(px router.PushEvent) {
		// px.Params{id:12}, px.Rem: /12
		//...
	})

	rx.Failed(func(px router.PushEvent) {
		//...
	})

	home.Resolve(router.UseLocation("home/12"))
}
```

```go

import "github.com/gu-io/gu/router"

func main() {
	home := router.New("/home/*") // the /* tells the router to allow more paths.
	rx := home.Only("/:id")

	rx.Done(func(px router.PushEvent) {
		// px.Params{id:12}
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

By combining these simple concepts, it should provide a flexible approach in routing for components, views and any other use case. As times does go by and more complex and complicated needs arise, the `router` package will be updated.
