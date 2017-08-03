package gu

import (
	"fmt"
	"html/template"

	"github.com/gu-io/gu/drivers/core"
	"github.com/gu-io/gu/notifications"
	"github.com/gu-io/gu/router"
	"github.com/gu-io/gu/trees"
	"github.com/gu-io/gu/trees/elems"
)

// NApp defines a struct which encapsulates all the core view management functions
// for views.
type NApp struct {
	active         bool
	title          string
	uuid           string
	location       Location
	views          []*NView
	activeViews    []*NView
	tree           *trees.Markup
	notifications  *notifications.AppEventNotification
	router         *router.Router
	resourceHeader []*trees.Markup
	resourceBody   []*trees.Markup
}

// App creates a new app structure to rendering gu components.
func App(title string, router *router.Router) *NApp {
	var app NApp
	app.title = title
	app.uuid = NewKey()
	app.router = router
	app.notifications = notifications.AppNotification(app.uuid)

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

// Do calls the giving function providing it with the NApp instance.
func (app *NApp) Do(appFun func(*NApp)) *NApp {
	if appFun != nil {
		appFun(app)
		return app
	}

	return app
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
	tjson.Name = app.title

	toHead, toBody := app.Resources()

	for _, item := range toHead {
		tjson.HeadResources = append(tjson.HeadResources, item.TreeJSON())
	}

	for _, item := range toBody {
		tjson.BodyResources = append(tjson.BodyResources, item.TreeJSON())
	}

	var afterBody []ViewJSON

	for _, view := range app.activeViews {
		switch view.target {
		case HeadTarget:
			tjson.Head = append(tjson.Head, view.RenderJSON())
		case BodyTarget:
			tjson.Body = append(tjson.Body, view.RenderJSON())
		case AfterBodyTarget:
			afterBody = append(afterBody, view.RenderJSON())
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
		switch view.target {
		case HeadTarget:
			view.Render().Apply(head)
		case BodyTarget:
			view.Render().Apply(body)
		case AfterBodyTarget:
			view.Render().Apply(last)
		}
	}

	script := trees.NewMarkup("script", false)
	trees.NewAttr("type", "text/javascript").Apply(script)
	trees.NewText(core.JavascriptDriverCore).Apply(script)

	script.Apply(last)

	body.AddChild(last.Children()...)

	body.AddChild(toBody...)

	return html
}

// PushViews returns a slice of  views that match and pass the provided path.
func (app *NApp) PushViews(event router.PushEvent) []*NView {
	// fmt.Printf("Routing Path: %s\n", event.Rem)
	var active []*NView

	for _, view := range app.views {
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
	if app.resourceHeader != nil && app.resourceBody != nil {
		return app.resourceHeader, app.resourceBody
	}

	var head, body []*trees.Markup

	head = append(head, elems.Title(elems.Text(app.title)))
	head = append(head, elems.Meta(trees.NewAttr("app-id", app.uuid)))
	head = append(head, elems.Meta(trees.NewAttr("charset", "utf-8")))

	app.resourceHeader = head
	app.resourceBody = body

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

// View returns a new instance of the view object.
func (app *NApp) View(renderable interface{}, route string, target ViewTarget) *NView {
	app.initSanitCheck()

	if route == "" {
		route = "*"
	}

	var base Renderable

	switch rnb := renderable.(type) {
	case Renderable:
		base = rnb
		break
	case *trees.Markup:
		base = Static(rnb)
		break
	case trees.Appliable:
		base = ApplyStatic(rnb)
		break
	default:
		panic("Only Renderable/trees.Markup allowed")
	}

	var vw NView
	vw.root = app
	vw.target = target
	vw.base = base
	vw.uuid = NewKey()
	vw.appUUID = app.uuid
	vw.Reactive = NewReactive()

	vw.router = router.NewResolver(route)

	// app.driver.Update(app, &vw)
	vw.React(func() {
		notifications.Dispatch(ViewUpdate{
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
	base    Renderable
	target  ViewTarget

	// location      Location
	router router.Resolver

	mounted   Subscriptions
	rendered  Subscriptions
	updated   Subscriptions
	unmounted Subscriptions

	beginComponents []*Component
	anyComponents   []*Component
	lastComponents  []*Component
}

// UUID returns the uuid specific to the giving view.
func (v *NView) UUID() string {
	return v.uuid
}

// Do calls the giving function providing it with the NApp instance.
func (v *NView) Do(viewFun func(*NView)) *NView {
	if viewFun != nil {
		viewFun(v)
		return v
	}

	return v
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

// Target returns the associated view target.
func (v *NView) Target() ViewTarget {
	return v.target
}

// Render returns the markup for the giving views.
func (v *NView) Render() *trees.Markup {
	base := v.base.Render()

	// Process the begin components and immediately add appropriately into base.
	for _, component := range v.beginComponents {
		if component.Target == "" {
			component.Render().ApplyMorphers().Apply(base)
			continue
		}

		render := component.Render().ApplyMorphers()
		targets := trees.Query.QueryAll(base, component.Target)
		for _, target := range targets {
			target.AddChild(render)
			target.UpdateHash()
		}
	}

	// Process the middle components and immediately add appropriately into base.
	for _, component := range v.anyComponents {
		if component.Target == "" {
			component.Render().ApplyMorphers().Apply(base)
			continue
		}

		render := component.Render().ApplyMorphers()
		targets := trees.Query.QueryAll(base, component.Target)
		for _, target := range targets {
			target.AddChild(render)
			target.UpdateHash()
		}
	}

	// Process the last components and immediately add appropriately into base.
	for _, component := range v.lastComponents {
		if component.Target == "" {
			component.Render().ApplyMorphers().Apply(base)
			continue
		}

		render := component.Render().ApplyMorphers()
		targets := trees.Query.QueryAll(base, component.Target)
		for _, target := range targets {
			target.AddChild(render)
			target.UpdateHash()
		}
	}

	if len(base.Children()) == 1 {
		child := base.Children()[0]
		child.SwapUID(v.uuid)
		child.UpdateHash()

		return child
	}

	base.SwapUID(v.uuid)
	base.UpdateHash()

	return base
}

// propagateRoute supplies the needed route into the provided
func (v *NView) propagateRoute(pe router.PushEvent) {
	v.router.Resolve(pe)
}

// Unmounted publishes changes notifications that the view is unmounted.
func (v *NView) Unmounted() {
	v.unmounted.Publish()
}

// Updated publishes changes notifications that the view is updated.
func (v *NView) Updated() {
	v.updated.Publish()
}

// Rendered publishes changes notifications that the view is rendered.
func (v *NView) Rendered() {
	v.rendered.Publish()
}

// Mounted publishes changes notifications that the view is mounted.
func (v *NView) Mounted() {
	v.mounted.Publish()
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

// Services return s a Service instance which contains fields used by the
// Components of a view to gain access to the specific functionality of it's app root.
func (v *NView) Services() Services {
	return Services{
		AppUUID:   v.appUUID,
		Location:  v.root,
		ViewRoute: v.router,
		Router:    v.root.router,
		Mounted:   v.mounted,
		Unmounted: v.unmounted,
		Updated:   v.updated,
		Rendered:  v.rendered,
	}
}

// Component adds the provided component into the selected view.
func (v *NView) Component(renderable interface{}, order RenderingOrder, route string, target string) {
	var base Renderable

	switch rnb := renderable.(type) {
	case Renderable:
		base = rnb
		break
	case *trees.Markup:
		base = Static(rnb)
		break
	case trees.Appliable:
		base = ApplyStatic(rnb)
		break
	default:
		panic("Only Renderable/trees.Markup allowed")
	}

	var c Component
	c.uuid = NewKey()
	c.Target = target
	c.Rendering = base
	c.Reactive = NewReactive()
	c.Router = router.NewResolver(route)

	// if the renderable can push reactions then listen.
	if rr, ok := base.(Reactor); ok {
		rr.React(c.Reactive.Publish)
	}

	// Connect the view to react to a change from the component.
	c.React(v.Publish)

	// Register the component router into the views router.
	v.router.Register(c.Router)

	// format for the object.
	// Add the component into the right order.
	{
		switch order {
		case FirstOrder:
			v.beginComponents = append(v.beginComponents, &c)
		case LastOrder:
			v.lastComponents = append(v.lastComponents, &c)
		case AnyOrder:
			v.anyComponents = append(v.anyComponents, &c)
		}
	}
}

// Component defines a struct which
type Component struct {
	Reactive
	uuid   string
	Target string

	Rendering Renderable
	Router    router.Resolver

	live *trees.Markup
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
			if e.Remove != nil {
				e.Remove.Remove()
			}
		})

		newTree.Reconcile(live)
		live.Empty()
	}

	c.live = newTree.ApplyMorphers()

	return c.live
}

// Disabled returns true/false if the giving view is disabled.
func (v *NView) Disabled() bool {
	return v.active
}

// enableView enables the active state of the view.
func (v *NView) enableView() {
	v.active = true
}

// disableView disables the active state of the view.
func (v *NView) disableView() {
	v.active = false
}

//==============================================================================

// ApplyView defines a MarkupRenderer implementing structure which returns its Content has
// its markup.
type ApplyView struct {
	uid      string
	Morph    bool
	Mounted  Subscriptions
	Rendered Subscriptions
	Content  trees.Appliable
	base     *trees.Markup
}

// ApplyStatic defines a toplevel function which returns a new instance of a StaticView using the
// provided markup as its content.
func ApplyStatic(tree trees.Appliable) *ApplyView {
	return &ApplyView{
		Content: tree,
		uid:     NewKey(),
		base:    trees.NewMarkup("div", false),
	}
}

// UUID returns the RenderGroup UUID for identification.
func (a *ApplyView) UUID() string {
	return a.uid
}

// Render returns the markup for the static view.
func (a *ApplyView) Render() *trees.Markup {
	a.Content.Apply(a.base)

	children := a.base.Children()

	if len(children) == 0 {
		return a.base
	}

	a.base.Empty()

	root := children[0]
	if a.Morph {
		return root.ApplyMorphers()
	}

	return root
}

// RenderHTML returns the html template version of the StaticView content.
func (a *ApplyView) RenderHTML() template.HTML {
	return a.Render().EHTML()
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
