Drivers
=======

Drivers are a system developed in the redesign of Gu to provide capability to render the result of a giving app to any supportable platform. Due to the rise in rendering targets such as Mobile and Desktop, we eagerly believe that GopherJS will not be the only means to deliver applications to users, hence there was a need to provide a means by which the same Gu applications can be rendered on such platforms with ease.

GopherJS does provide a convenient and powerful means to take Go Apps into the JS world, but we desired to allow flexibility in the way anyone can take the any app built on Gu, to easily be rendered and usable on different targets/systems and drivers met that need.

By easily abstracting out the rendering details for each platform and making the core Gu package concentrate on organization and structures, we easily allow flexibility for more larger systems which can easily embed any app built with Gu.

*We hope that developers will take this and push the boundaries further to allow easy deployment of Gu apps to other platforms e.g QT, Android, iOS,...etc*

Examples of Drivers:
--------------------

Below are the list of drivers being actively developed or are already usable. We hope this list can increase the more.

-	GopherJS Driver(https://github.com/gu-io/gopherjs/) (Stable)
-	QT Driver(https://github.com/gu-io/qt) (Experiemental)

Gu's Driver Interface
---------------------

This is the interface which all `Drivers` implementation must match and through this, Gu can easily be rendered to the target platform.

```go

// Driver defines an interface which provides an Interface to render Apps and views
// to a display system. eg. GopherJS, wkWebview.
type Driver interface {
	// Method to subscribe to ready state for the driver.
	OnReady(func())

	// Name of the driver.
	Name() string

	// Navigate the Driver to the provided path.
	Navigate(router.PushDirectiveEvent)

	// Current Location of the driver path.
	Location() router.PushEvent

	// OnRoute registers for route change notifications to react accordingly.
	// This does a totally re-write of the whole display.
	OnRoute(*NApp)

	// Render app and it's content.
	Render(*NApp)

	// Update app's view and it's target content.
	Update(*NApp, ...*NView)

	// Fetcher returns a new resource fetcher which can be used to retrieve Resources.
	// Fetcher requires the cacheName and a boolean indicating if it should intercept
	// all requests for resources.
	Fetcher(cacheName string, interceptRequests bool) (shell.Fetch, shell.Cache)
}

```

It's job is to provide an interface by which it's selected platform renders out appropriately the designed application built using GU and to provide the flexibility of not being tied to a specific rendering endpoint. It also exposes means by which, requests for resources can be made through and caches by which response can be sorted to reduce network usage.
