package gu

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/gu-io/gu/drivers/core"
	"github.com/gu-io/gu/hooks"
	"github.com/gu-io/gu/notifications"
	"github.com/gu-io/gu/router"
	"github.com/gu-io/gu/shell"
	"github.com/gu-io/gu/trees"
	"github.com/gu-io/gu/trees/elems"
	"github.com/influx6/faux/reflection"
)

// Resource defines any set of rendering links, scripts, styles needed by a view.
type Resource struct {
	Manifest shell.AppManifest
	body     []*trees.Markup
	head     []*trees.Markup
}

// AppAttr defines a struct for
type AppAttr struct {
	Name      string              `json:"name"`
	Title     string              `json:"title"`
	Manifests []shell.AppManifest `json:"manifests, omitempty"`
	Router    *router.Router      `json:"-"`
}

// NApp defines a struct which encapsulates all the core view management functions
// for views.
type NApp struct {
	uuid          string
	attr          AppAttr
	location      Location
	router        *router.Router
	notifications *notifications.AppNotification
	active        bool

	local       []shell.AppManifest
	views       []*NView
	activeViews []*NView

	globalResources []Resource

	tree *trees.Markup
}

// App creates a new app structure to rendering gu components.
func App(attr AppAttr) *NApp {
	var app NApp
	app.attr = attr
	app.uuid = NewKey()
	app.router = attr.Router

	// Add local notification channel for this giving app.
	app.notifications = notifications.New(app.uuid)

	// Sort Manifests file accordingly.
	if attr.Manifests != nil {

		for _, mani := range attr.Manifests {
			for _, attr := range mani.Manifests {
				if attr.IsGlobal {
					hooks.RegisterManifest(attr)
				}
			}

			if mani.GlobalScope {
				app.globalResources = append(app.globalResources, Resource{
					Manifest: mani,
				})

				continue
			}

			app.local = append(app.local, mani)
		}
	}

	// NOTE: Do this in the driver.
	// app.driver.OnReady(func() {
	// 	fmt.Printf("Running App: %q\n", app.attr.Name)
	// 	fmt.Printf("Running App Title: %q\n", app.attr.Name)
	//
	// 	app.active = false
	// 	app.ActivateRoute(app.driver.Location())
	//
	// 	app.driver.Render(&app)
	// 	app.active = true
	//
	// 	fmt.Printf("Ending App Run: %q\n", app.attr.Name)
	// })

	// NOTE: Do this in the driver.
	// app.notifications.Subscribe(func(directive router.PushDirectiveEvent) {
	// 	if !app.active || app.driver == nil {
	// 		return
	// 	}
	//
	// 	app.driver.Navigate(directive)
	// })

	// NOTE: Do this in the driver.
	// app.driver.OnRoute(&app)

	// NOTE: Do this in the driver.
	// app.driver.Ready()
	return &app
}

// Navigate sets the giving app location and also sets the location of the
// NOOPLocation which returns that always.
func (app *NApp) Navigate(pe router.PushDirectiveEvent) {
	app.initSanitCheck()
	app.location.Navigate(pe)
}

// Location returns the current route. It stores all set routes and returns the
// last route else returning a
func (app *NApp) Location() router.PushEvent {
	app.initSanitCheck()
	return app.location.Location()
}

// InitApp sets the Location to be used by the NApp and it's views and components.
func (app *NApp) InitApp(location Location) {
	app.location = location
}

// initSanitCheck will perform series of checks to ensure the needed features or
// structures required by the app is set else will panic.
func (app *NApp) initSanitCheck() {
	if app.location != nil {
		return
	}

	// Use the NoopLocation since we have not being set.
	app.location = NewNoopLocation(app)
}

// Notifications returns the underlying AppNotification pipeline for access.
func (app *NApp) Notifications() *notifications.AppNotification {
	return app.notifications
}

// Active returns true/false if the giving app is active and has already
// received rendering.
func (app *NApp) Active() bool {
	return app.active
}

// Mounted notifies all active views that they have been mounted.
func (app *NApp) Mounted() {
	for _, view := range app.activeViews {
		view.Mounted()
	}
}

// ActivateRoute actives the views which are to be rendered.
func (app *NApp) ActivateRoute(es interface{}) {
	var pe router.PushEvent

	switch esm := es.(type) {
	case string:
		tmp, err := router.NewPushEvent(esm, true)
		if err != nil {
			panic(fmt.Sprintf("Unable to create PushEvent for (URL: %q) -> %q\n", esm, err.Error()))
		}

		pe = tmp
	case router.PushEvent:
		pe = esm
	}

	app.activeViews = app.PushViews(pe)
}

// AppJSON defines a struct which holds the giving sets of tree changes to be
// rendered.
type AppJSON struct {
	AppID         string             `json:"AppId"`
	Name          string             `json:"Name"`
	Title         string             `json:"Title"`
	Head          []ViewJSON         `json:"Head"`
	Body          []ViewJSON         `json:"Body"`
	HeadResources []trees.MarkupJSON `json:"HeadResources"`
	BodyResources []trees.MarkupJSON `json:"BodyResources"`
}

// RenderJSON returns the giving rendered tree of the app respective of the path
// found as jons structure with markup content.
func (app *NApp) RenderJSON(es interface{}) AppJSON {
	if es != nil {
		app.ActivateRoute(es)
	}

	var tjson AppJSON
	tjson.Name = app.attr.Name
	tjson.Title = app.attr.Title

	toHead, toBody := app.Resources()

	for _, item := range toHead {
		tjson.HeadResources = append(tjson.HeadResources, item.TreeJSON())
	}

	for _, item := range toBody {
		tjson.BodyResources = append(tjson.BodyResources, item.TreeJSON())
	}

	var afterBody []ViewJSON

	for _, view := range app.activeViews {
		switch view.attr.Target {
		case HeadTarget:
			tjson.Head = append(tjson.Head, view.RenderJSON())
		case BodyTarget:
			tjson.Body = append(tjson.Body, view.RenderJSON())
		case AfterBodyTarget:
			afterBody = append(afterBody, view.RenderJSON())
		}

		viewHead, viewBody := view.Resources()
		for _, item := range viewHead {
			tjson.HeadResources = append(tjson.HeadResources, item.TreeJSON())
		}

		for _, item := range viewBody {
			tjson.BodyResources = append(tjson.BodyResources, item.TreeJSON())
		}
	}

	tjson.Body = append(tjson.Body, afterBody...)

	script := trees.NewMarkup("script", false)
	trees.NewAttr("type", "text/javascript").Apply(script)
	trees.NewText(core.JavascriptDriverCore).Apply(script)
	tjson.BodyResources = append(tjson.BodyResources, script.TreeJSON())

	return tjson
}

// Render returns the giving rendered tree of the app respective of the path
// found.
func (app *NApp) Render(es interface{}) *trees.Markup {
	if es != nil {
		app.ActivateRoute(es)
	}

	var html = trees.NewMarkup("html", false)
	var head = trees.NewMarkup("head", false)

	var body = trees.NewMarkup("body", false)
	trees.NewAttr("gu-app-id", app.uuid).Apply(body)

	// var app = trees.NewMarkup("app", false)
	// trees.NewAttr("gu-app-id", app.uuid).Apply(app)

	head.Apply(html)
	body.Apply(html)

	// Generate the resources according to the received data.
	toHead, toBody := app.Resources()
	head.AddChild(toHead...)

	var last = elems.Div()

	for _, view := range app.activeViews {
		switch view.attr.Target {
		case HeadTarget:
			view.Render().Apply(head)
		case BodyTarget:
			view.Render().Apply(body)
		case AfterBodyTarget:
			view.Render().Apply(last)
		}

		viewHead, viewBody := view.Resources()

		// Add the headers into the header so they load accordingly.
		head.AddChild(viewHead...)

		// Append the resources into the body has we need them last.
		toBody = append(toBody, viewBody...)
	}

	script := trees.NewMarkup("script", false)
	trees.NewAttr("type", "text/javascript").Apply(script)
	trees.NewText(core.JavascriptDriverCore).Apply(script)

	script.Apply(last)

	body.AddChild(last.Children()...)

	body.AddChild(toBody...)

	// Ensure to have this gc'ed.
	last = nil

	return html
}

// PushViews returns a slice of  views that match and pass the provided path.
func (app *NApp) PushViews(event router.PushEvent) []*NView {
	// fmt.Printf("Routing Path: %s\n", event.Rem)
	var active []*NView

	for _, view := range app.views {
		// _, rem, ok := view.router.Test(event.Rem)
		// fmt.Printf("Routing View: %s : %s -> %t -> %s\n", view.Attr().Name, view.router.Pattern(), ok, rem)

		if _, _, ok := view.router.Test(event.Rem); !ok {
			// Notify view to appropriate proper action when view does not match.
			view.router.Resolve(event)
			continue
		}

		view.propagateRoute(event)
		active = append(active, view)
	}

	return active
}

// Resources return the giving resource headers which relate with the
// view.
func (app *NApp) Resources() ([]*trees.Markup, []*trees.Markup) {
	var head, body []*trees.Markup

	head = append(head, elems.Title(elems.Text(app.attr.Title)))
	head = append(head, elems.Meta(trees.NewAttr("gu-app-id", app.uuid)))
	head = append(head, elems.Meta(trees.NewAttr("gu-app-name", app.attr.Name)))
	head = append(head, elems.Meta(trees.NewAttr("gu-app-title", app.attr.Title)))

	for _, def := range app.globalResources {
		if def.body != nil || def.head != nil {
			head = append(head, def.head...)
			body = append(body, def.body...)
			continue
		}

		if def.Manifest.Manifests == nil {
			continue
		}

		for _, manifest := range def.Manifest.Manifests {
			if !manifest.Init {
				continue
			}

			hook, err := hooks.Get(manifest.HookName)
			if err != nil {
				fmt.Printf("Hook[%q] does not exists: Resource[%q] unable to install\n", manifest.HookName, manifest.Name)
				continue
			}

			markup, toHead, err := hook.Fetch(app.router, manifest)
			if err != nil {
				fmt.Printf("Hook[%q] failed to retrieve Resource {Name: %q, Path: %q}\n", manifest.HookName, manifest.Name, manifest.Path)
				continue
			}

			trees.NewAttr("gu-resource", "true").Apply(markup)
			trees.NewAttr("gu-resource-view", app.uuid).Apply(markup)
			trees.NewAttr("gu-resource-from", manifest.Path).Apply(markup)
			trees.NewAttr("gu-resource-name", manifest.Name).Apply(markup)
			trees.NewAttr("gu-resource-id", manifest.ID).Apply(markup)
			trees.NewAttr("gu-resource-app-id", app.uuid).Apply(markup)

			if toHead {
				def.head = append(def.head, markup)
				head = append(head, markup)
				continue
			}

			def.body = append(def.body, markup)
			body = append(body, markup)
		}
	}

	return head, body
}

// UUID returns the uuid specific to the giving view.
func (app *NApp) UUID() string {
	return app.uuid
}

// ViewTarget defines a concrete type to define where the view should be rendered.
type ViewTarget int

const (

	// BodyTarget defines the view target where the view is rendered in the body.
	BodyTarget ViewTarget = iota

	// HeadTarget defines the view target where the view is rendered in the head.
	HeadTarget

	// AfterBodyTarget defines the view target where the view is rendered after
	// body views content. Generally the browser moves anything outside of the body
	// into the body as last elements. So this will be the last elements rendered
	// in the border accordingly in the order they are added into the respective app.
	AfterBodyTarget
)

// ViewAttr defines a structure to define a option values for setting up the appropriate
// settings for the view.
type ViewAttr struct {
	Name   string        `json:"name"`
	Route  string        `json:"route"`
	Target ViewTarget    `json:"target"`
	Base   *trees.Markup `json:"base"`
}

// View returns a new instance of the view object.
func (app *NApp) View(attr ViewAttr) *NView {
	app.initSanitCheck()

	if attr.Base == nil {
		attr.Base = trees.NewMarkup("view", false)
		trees.NewCSSStyle("display", "block").Apply(attr.Base)
	}

	var vw NView
	vw.attr = attr
	vw.root = app
	vw.uuid = NewKey()
	vw.appUUID = app.uuid
	// vw.location = app.location
	vw.notifications = app.notifications
	vw.Reactive = NewReactive()
	vw.local = app.local

	vw.router = router.New(attr.Route)

	vw.React(func() {

		// app.driver.Update(app, &vw)
		app.notifications.Dispatch(ViewUpdate{
			App:  app,
			View: &vw,
		})
	})

	// Register to listen for failure of route to match and
	// notify unmount call.
	vw.router.Failed(func(push router.PushEvent) {
		vw.disableView()
		vw.Unmounted()
	})

	vw.attr.Base.SwapUID(vw.uuid)

	app.views = append(app.views, &vw)

	return &vw
}

// RenderableData defines a struct which contains the name of a giving renderable
// and it's package.
type RenderableData struct {
	Name string
	Pkg  string
}

// NView defines a structure to encapsulates all rendering component for a given
// view.
type NView struct {
	Reactive
	root    *NApp
	uuid    string
	appUUID string
	active  bool
	attr    ViewAttr
	// location      Location
	router        router.Resolver
	notifications *notifications.AppNotification

	renderingData []RenderableData
	local         []shell.AppManifest

	localResources []Resource

	beginComponents []*Component
	anyComponents   []*Component
	lastComponents  []*Component
}

// UUID returns the uuid specific to the giving view.
func (v *NView) UUID() string {
	return v.uuid
}

// totalComponents returns the total component list.
func (v *NView) totalComponents() int {
	return len(v.beginComponents) + len(v.anyComponents) + len(v.lastComponents)
}

// ViewJSON defines a struct which holds the giving sets of view changes to be
// rendered.
type ViewJSON struct {
	AppID  string           `json:"AppID"`
	ViewID string           `json:"ViewID"`
	Tree   trees.MarkupJSON `json:"Tree"`
}

// RenderJSON returns the ViewJSON for the provided View and its current events and
// changes.
func (v *NView) RenderJSON() ViewJSON {
	return ViewJSON{
		AppID:  v.appUUID,
		ViewID: v.uuid,
		Tree:   v.Render().TreeJSON(),
	}
}

// Render returns the markup for the giving views.
func (v *NView) Render() *trees.Markup {
	fmt.Printf("Rendering View: %q Total Components: %d\n", v.attr.Name, v.totalComponents())

	base := v.attr.Base.Clone()

	// Update the base hash.
	base.UpdateHash()

	// Process the begin components and immediately add appropriately into base.
	for _, component := range v.beginComponents {
		if component.attr.Target == "" {
			component.Render().ApplyMorphers().Apply(base)
			continue
		}

		render := component.Render().ApplyMorphers()
		targets := trees.Query.QueryAll(base, component.attr.Target)
		for _, target := range targets {
			target.AddChild(render)
			target.UpdateHash()
		}
	}

	// Process the middle components and immediately add appropriately into base.
	for _, component := range v.anyComponents {
		if component.attr.Target == "" {
			component.Render().ApplyMorphers().Apply(base)
			continue
		}

		render := component.Render().ApplyMorphers()
		targets := trees.Query.QueryAll(base, component.attr.Target)
		for _, target := range targets {
			target.AddChild(render)
			target.UpdateHash()
		}
	}

	// Process the last components and immediately add appropriately into base.
	for _, component := range v.lastComponents {
		if component.attr.Target == "" {
			component.Render().ApplyMorphers().Apply(base)
			continue
		}

		render := component.Render().ApplyMorphers()
		targets := trees.Query.QueryAll(base, component.attr.Target)
		for _, target := range targets {
			target.AddChild(render)
			target.UpdateHash()
		}
	}

	return base
}

// Attr returns the views ViewAttr.
func (v *NView) Attr() ViewAttr {
	return v.attr
}

// propagateRoute supplies the needed route into the provided
func (v *NView) propagateRoute(pe router.PushEvent) {
	v.router.Resolve(pe)
}

// Resources return the giving resource headers which relate with the
// view.
func (v *NView) Resources() ([]*trees.Markup, []*trees.Markup) {
	var head, body []*trees.Markup

	for _, def := range v.localResources {
		if def.body != nil || def.head != nil {
			head = append(head, def.head...)
			body = append(body, def.body...)
			continue
		}

		if def.Manifest.Manifests == nil {
			continue
		}

		for _, manifest := range def.Manifest.Manifests {
			if !manifest.Init {
				continue
			}

			hook, err := hooks.Get(manifest.HookName)
			if err != nil {
				fmt.Printf("Hook[%q] does not exists: Resource[%q] unable to install\n", manifest.HookName, manifest.Name)
				continue
			}

			markup, toHead, err := hook.Fetch(v.root.router, manifest)
			if err != nil {
				fmt.Printf("Hook[%q] failed to retrieve Resource {Name: %q, Path: %q}\n", manifest.HookName, manifest.Name, manifest.Path)
				continue
			}

			trees.NewAttr("gu-resource", "true").Apply(markup)
			trees.NewAttr("gu-resource-view", v.uuid).Apply(markup)
			trees.NewAttr("gu-resource-from", manifest.Path).Apply(markup)
			trees.NewAttr("gu-resource-name", manifest.Name).Apply(markup)
			trees.NewAttr("gu-resource-id", manifest.ID).Apply(markup)

			if toHead {
				def.head = append(def.head, markup)
				head = append(head, markup)
				continue
			}

			def.body = append(def.body, markup)
			body = append(body, markup)
		}
	}

	return head, body
}

// Unmounted publishes changes notifications that the component is unmounted.
func (v *NView) Unmounted() {
	for _, component := range v.beginComponents {
		component.Unmounted.Publish()
	}
	for _, component := range v.anyComponents {
		component.Unmounted.Publish()
	}
	for _, component := range v.lastComponents {
		component.Unmounted.Publish()
	}
}

// Updated publishes changes notifications that the component is updated.
func (v *NView) Updated() {
	for _, component := range v.beginComponents {
		component.Updated.Publish()
	}
	for _, component := range v.anyComponents {
		component.Updated.Publish()
	}
	for _, component := range v.lastComponents {
		component.Updated.Publish()
	}
}

// Mounted publishes changes notifications that the component is mounted.
func (v *NView) Mounted() {
	for _, component := range v.beginComponents {
		component.Mounted.Publish()
	}
	for _, component := range v.anyComponents {
		component.Mounted.Publish()
	}
	for _, component := range v.lastComponents {
		component.Mounted.Publish()
	}
}

// RenderingOrder defines a type used to define the order which rendering is to be done for a resource.
type RenderingOrder int

const (
	// FirstOrder defines that rendering be first in order.
	FirstOrder RenderingOrder = iota

	// AnyOrder defines that rendering be middle in order.
	AnyOrder

	// LastOrder defines that rendering be last in order.
	LastOrder
)

// ComponentAttr defines a structure to define a component and its appropriate settings.
type ComponentAttr struct {
	Order     RenderingOrder `json:"order"`
	Tag       string         `json:"tag"`
	Target    string         `json:"target"`
	Route     string         `json:"route"`
	Base      interface{}    `json:"base"`
	Relations []string       `json:"relations"`
}

// Component adds the provided component into the selected view.
func (v *NView) Component(attr ComponentAttr) {
	if strings.TrimSpace(attr.Route) == "" {
		attr.Route = "*"
	}

	var c Component
	c.attr = attr
	c.uuid = NewKey()
	c.Reactive = NewReactive()
	c.Mounted = NewSubscriptions()
	c.Unmounted = NewSubscriptions()
	c.Rendered = NewSubscriptions()
	c.Updated = NewSubscriptions()
	c.Router = router.New(attr.Route)

	if attr.Tag == "" {
		attr.Tag = "component-element"
	}

	appServices := Services{
		AppUUID:       v.appUUID,
		Location:      v.root,
		Mounted:       c.Mounted,
		Updated:       c.Updated,
		Rendered:      c.Rendered,
		Unmounted:     c.Unmounted,
		ViewRouter:    c.Router,
		Router:        v.root.router,
		Notifications: v.notifications,
	}

	// Transform the base argument into the acceptable
	// format for the object.
	{
		switch mo := attr.Base.(type) {
		case func(Services) *trees.Markup:
			static := Static(mo(appServices))

			static.Morph = true
			c.Rendering = static

		case func() *trees.Markup:
			static := Static(mo())
			static.Morph = true
			c.Rendering = static

		case *trees.Markup:
			static := Static(mo)
			static.Morph = true
			c.Rendering = static
			break

		case string:
			parseTree := trees.ParseTree(mo)
			if len(parseTree) != 1 {
				section := elems.CustomElement(attr.Tag)
				section.AddChild(parseTree...)

				static := Static(section)
				static.Morph = true
				c.Rendering = static
				break
			}

			static := Static(parseTree[0])
			static.Morph = true
			c.Rendering = static
			break

		case Renderable:
			if service, ok := mo.(RegisterService); ok {
				service.RegisterService(appServices)
			}

			if renderField, _, err := reflection.StructAndEmbeddedTypeNames(mo); err == nil {
				v.renderingData = append(v.renderingData, RenderableData{
					Name: renderField.TypeName,
					Pkg:  renderField.Pkg,
				})
			}

			c.Rendering = mo
			break

		case func(Services) Renderable:
			rc := mo(appServices)

			if renderField, _, err := reflection.StructAndEmbeddedTypeNames(rc); err == nil {
				v.renderingData = append(v.renderingData, RenderableData{
					Name: renderField.TypeName,
					Pkg:  renderField.Pkg,
				})
			}

			c.Rendering = rc
			break

		case func() Renderable:
			rc := mo()

			if service, ok := rc.(RegisterService); ok {
				service.RegisterService(appServices)
			}

			if renderField, _, err := reflection.StructAndEmbeddedTypeNames(rc); err == nil {
				v.renderingData = append(v.renderingData, RenderableData{
					Name: renderField.TypeName,
					Pkg:  renderField.Pkg,
				})
			}

			c.Rendering = rc
			break

		default:
			panic(`
				Unknown markup or view processable type

					Accepted Markup Arguments:
						-	*trees.Markup
						- func() *trees.Markup

					Accepted View Arguments:
						-	[]Renderable
						-	Renderable
						-	func() []Renderable
						-	func() Renderable

				`)
		}
	}

	// Add the component into the right order.
	{
		switch attr.Order {
		case FirstOrder:
			v.beginComponents = append(v.beginComponents, &c)
		case LastOrder:
			v.lastComponents = append(v.lastComponents, &c)
		case AnyOrder:
			v.anyComponents = append(v.anyComponents, &c)
		}
	}

	// Connect the component into the rendering reactor if it has one.
	if rc, ok := c.Rendering.(Reactor); ok {
		rc.React(c.Reactive.Publish)
	}

	// Connect the view to react to a change from the component.
	c.React(v.Publish)

	// Register the component router into the views router.
	v.router.Register(c.Router)

	// Collect necessary app manifest that connect with rendering.
	{
		for _, relation := range v.renderingData {
			if app, err := shell.FindByRelation(v.local, relation.Name); err == nil {
				if v.hasRelation(app.Name) {
					continue
				}

				v.localResources = append(v.localResources, Resource{
					Manifest: app,
				})

				initRelation(v, app, v.local)
			}
		}
	}

	// Send call for view update.
	// v.driver.Update(v.root, v)
	v.Publish()
}

// Component defines a struct which
type Component struct {
	Reactive
	uuid      string
	attr      ComponentAttr
	Rendering Renderable
	Router    router.Resolver

	Mounted   Subscriptions
	Unmounted Subscriptions
	Rendered  Subscriptions
	Updated   Subscriptions
	live      *trees.Markup
}

// UUID returns the identification for the giving component.
func (c Component) UUID() string {
	return c.uuid
}

// Render returns the markup corresponding to the internal Renderable.
func (c *Component) Render() *trees.Markup {
	newTree := c.Rendering.Render()
	newTree.SwapUID(c.uuid)

	if c.live != nil {
		live := c.live
		live.EachEvent(func(e *trees.Event, _ *trees.Markup) {
			if e.Handle != nil {
				e.Handle.End()
			}
		})

		newTree.Reconcile(live)
		live.Empty()
	}

	c.live = newTree.ApplyMorphers()
	// fmt.Printf("Live: %s\n", c.live.HTML())

	c.Rendered.Publish()

	return c.live
}

// Disabled returns true/false if the giving view is disabled.
func (v *NView) Disabled() bool {
	// v.rl.RLock()
	// defer v.rl.RUnlock()

	return v.active
}

// enableView enables the active state of the view.
func (v *NView) enableView() {
	// v.rl.Lock()
	// {
	v.active = true
	// }
	// v.rl.Unlock()
}

// disableView disables the active state of the view.
func (v *NView) disableView() {
	// v.rl.Lock()
	// {
	v.active = false
	// }
	// v.rl.Unlock()
}

// hasRenderable returns true/false if a giving dependency name exists.
func (v *NView) hasRenderable(name string) bool {
	for _, rd := range v.renderingData {
		if rd.Name == name {
			return true
		}
	}

	return false
}

// hasRelation returns true/false if a giving manifests name exists.
func (v *NView) hasRelation(name string) bool {
	for _, rd := range v.localResources {
		if rd.Manifest.Name == name {
			return true
		}
	}

	return false
}

//==============================================================================

// StaticView defines a MarkupRenderer implementing structure which returns its Content has
// its markup.
type StaticView struct {
	uid      string
	Content  *trees.Markup
	Mounted  Subscriptions
	Rendered Subscriptions
	Morph    bool
}

// Static defines a toplevel function which returns a new instance of a StaticView using the
// provided markup as its content.
func Static(tree *trees.Markup) *StaticView {
	return &StaticView{
		Content: tree,
		uid:     NewKey(),
	}
}

// UUID returns the RenderGroup UUID for identification.
func (s *StaticView) UUID() string {
	return s.uid
}

// Dependencies returns the list of all internal dependencies of the given view.
// It returns the names of the structs and their internals composed values/fields
// to help conditional resource loading.
func (s *StaticView) Dependencies() []RenderableData {
	return nil
}

// Render returns the markup for the static view.
func (s *StaticView) Render() *trees.Markup {
	if s.Morph {
		return s.Content.ApplyMorphers()
	}

	return s.Content
}

// RenderHTML returns the html template version of the StaticView content.
func (s *StaticView) RenderHTML() template.HTML {
	return s.Content.EHTML()
}

//==============================================================================

// initRelation walks down the provided app relation adding the giving AppManifest
// which connect with this if not already in the list.
func initRelation(views *NView, app shell.AppManifest, relations []shell.AppManifest) {
	for _, relation := range app.Relation.Composites {
		if related, err := shell.FindByRelation(relations, relation); err == nil {
			if !views.hasRelation(related.Name) {

				views.localResources = append(views.localResources, Resource{
					Manifest: related,
				})

				initRelation(views, related, relations)
			}
		}
	}

	for _, field := range app.Relation.FieldTypes {
		if related, err := shell.FindByRelation(relations, field); err == nil {
			if !views.hasRelation(related.Name) {

				views.localResources = append(views.localResources, Resource{
					Manifest: related,
				})

				initRelation(views, related, relations)
			}
		}
	}
}
