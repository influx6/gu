Gu
==

[![Go Report Card](https://goreportcard.com/badge/github.com/gu-io/gu)](https://goreportcard.com/report/github.com/gu-io/gu)
[![Build Status](https://travis-ci.org/gu-io/gu.svg?branch=master)](https://travis-ci.org/gu-io/gu)

A component rendering library for Go. It efficiently renders standard HTML both on the frontend and backend.

Install
-------

```
go get -u github.com/gu-io/gu/...
```

Example
-------

![GopherJS Example](./media/greeter.png)

The example above can be found in the [Example Repo](https://github.com/gu-io/examples/tree/master/greeter).

CLI
---

Gu provides an cli tooling which is installed when `go get` is done for this package, the tooling provides easier means of generating a project and components using the gu project, among other features.

It is provided to both improve the workflow of the user, as well as to provide quick setup of your project. Provided below are examples of workflows which are generally done when developing with Gu.

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

-	Creating a new component as part of a existing package.

```bash
> cd sonar
> cd components/tableui/
> gu components new --flat=true --base=false tables
- Adding project file: "tableui/tables.go"
```

-	Creating a component as a self contained package.

*This types of package generated won't call out to a root package to register themselves with the projects `Components` registry*

```bash
>  gu components new --stand=true tooltip
- Adding project package: "tooltip"
- Adding project directory: "tooltip/styles"
- Adding project directory: "tooltip/styles/css"
- Adding project file: "tooltip/styles/generate.go"
- Adding project file: "tooltip/styles/css.go"
- Adding project file: "tooltip/tooltip.go"
```

-	Creating a component as part of the main components package for a project.

```bash
> cd sonar
>  gu components new --flat=true menubar
- Adding project file: "components/menubar.go"
```

-	Add a css package which internal will generate a new go file containing all css files in it.

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

Concepts
--------

Gu is fundamentally a library built to provide rendering capabilities with simple principles in building components that make up your application. There exists certain concepts which should be grasped due to the architecture and these do make it easier to reason and thinking when using the library. To fully grasp these concept, there is set below a series of short explanations about the different core pieces that make up the libray and I hope these will help in the use of this libray and it's examples.

-	[Virtual DOM](./docs/concepts/dom.md)

-	[Notifications](./docs/concepts/notifications.md)

-	[Routing](./docs/concepts/routing.md)

-	[Components](./docs/concepts/components.md)

-	[Drivers](./docs/concepts/drivers.md)

-	[App, Views, Component](./docs/concepts/app.md)

-	[Embeddable Resources](./docs/concepts/embedded-resources.md)

-	[Style Guides](./docs/concepts/theme.md)

How to Contribute
-----------------

Please read the contribution guidelines [Contribution Guidelines](./docs/concepts/contributing.md)

Limitations
-----------

Gu by it's very design and architecture is "Simple". It lacks the bells and whistles of similar frameworks and libraries. It's geared to solve your rendering needs and due to this, certain limitations exists with it.

-	Gu provides no react like flux structure.

-	Gu only focuses on providing you sets of structures able to work on the client and server for HTML/HTML5 rendering markup rendering with diffing support.

-	Gu component are simply Go types that implements Gu's set of interfaces and nothing else.

Once these limitations are not a problem, I believe using the library should help in achieving the end product and design you wish to build.

Last Note
---------

Please feel free to make issues on suggestions, questions, changes, bugs or improvements for the library. They all will be gladly received with much fan-fare.

God bless.
