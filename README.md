Gu
==

A component rendering library for Go. It efficiently renders standard HTML both on the frontend and backend.

Install
-------

```
go install github.com/gu-io/gu/...
```

Gu CLI
------

Gu provides an cli tooling which is installed when `go get` is done for this package, the tooling provides easier means of generating a project and components using the gu project, among other features. It is provided to both improve the workflow of the user but also provide quick setup of your project. We will delve into simple commands which allows us create a simple project with a series of components for development.

-	Creating a gu golang project using the [GopherJS](https://github.com/gu-io/gopherjs) driver

```bash
> gu new --driver=js sonar
- Creating new project: "sonar"
- Using driver template: "js"
	- Creating project directory: "sonar"
	- Creating project directory: "sonar/components"
	- Creating project directory: "sonar/assets"
	- Adding project file: "components/components.go"
	- Adding project file: "app.go"
```

-	Creating a component with it's assets for a project as a unique package.

	*Continuing from the command above*

```bash
> cd sonar
> gu components new tableui
- Adding project package: "components/tableui"
- Adding project directory: "components/tableui/styles"
- Adding project directory: "components/tableui/styles/css"
- Adding project file: "components/tableui/styles/generate.go"
- Adding project file: "components/tableui/tableui.go"

```

-	Creating a component file as part of a giving component package.

```bash
> cd sonar
> cd components/tableui/
> gu components new --flat=true --base=false tables
- Adding project file: "tableui/tables.go"
```

-	Creating a component file as part of the base components package.

```bash
> cd sonar
>  gu components new --flat=true menubar
- Adding project file: "components/menubar.go"
```

-	Adding css go package to turn css files into a go file into any directory.

```bash
> cd sonar
> cd sonar/components
>  gu components css menucss
- Adding project directory: "menucss"
- Adding project directory: "menucss/css"
- Adding project file: "menucss/generate.go"
```

Goals
-----

-	Dead Simple API.
-	Embeddable Resources.
-	Simplicity and Flexibility as core philosophies.
-	Able to render on both front and back end.
-	Quickly craft your UI without touching HTML.
-	Share code between backend and frontend.

Advantages
----------

-	Complex component libraries can be built up and shared as Golang packages.
-	Components are hierarchical allowing further reuse.
-	Event handling is simple and strongly typed.
-	Compile time safety

Examples
--------

The github repo [Examples](https://github.com/gu-io/examples), provides examples demonstrating the usage of the Gu library in building applications.

Documentation
-------------

Gu is fundamentally a library built to provide a component rendering package which exposes the means to effectively and efficiently render HTML/HTML5 content as needed. It provides different concepts and packages which help fulfill this in the most idiomatic form possible. It was built with the philosophy that simplicity is far better than complexity, and carries this principle in the way it's structures are built.

It offers a driver based system which allows the package to be used to render to different output system (e.g Browser , Headless Webkit system, ...etc). Though some of these features and drivers are still in works, Gu currently provides a [GopherJS](https://github.com/gopherjs) driver that showcases the possibility of the provided system.

*Gu is in no way a Flux-like framework. It is just a library that simply provides a baseline to render the desire output and gives the freedom for the developer to determine how their application data flow should works.*

### The Concepts

In grasping the examples and approach Gu takes, there exists certain concepts which need be introduced and you can quickly run down through them, has each concept tries to be short but informative about how that part of the Gu library works.

-	[Virtual DOM](./docs/concepts/dom.md)

-	[Notifications](./docs/concepts/notifications.md)

-	[Routing](./docs/concepts/routing.md)

-	[Components](./docs/concepts/components.md)

-	[Drivers](./docs/concepts/drivers.md)

-	[App, Views, Component](./docs/concepts/app.md)

-	[Embeddable Resource](./docs/concepts/embedded-resources.md)

How to Contribute
-----------------

Please read the contribution guidelines [Contribution Guidelines](./docs/concepts/contributing.md)

Limitations
-----------

Gu by it's very design and architecture was constructed to be "Simple". It lacks the bells and whistles of similar frameworks and libraries. It's geared to solve your rendering needs and due to this, certain limitations exists with it.

-	Gu provides no react like flux structure.

-	Gu only focuses on providing you sets of structures able to work on the client and server for HTML/HTML5 rendering.

-	Gu component are simply Go types that implements Gu's set of interfaces and nothing else.

Once these limitations are not a problem, I believe using the library should help in achieving the end product and design you wish to build.

Last Note
---------

Please feel free to make issues on suggestions, questions, changes, bugs or improvements for the library. They all will be gladly received with much fan-fare.

God bless.
