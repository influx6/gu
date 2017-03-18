// Package qtwebview implements the gu.Driver interface to provide rendering using the
// qt5 cross-platform GUI library and platform.
package qtwebview

import (
	"github.com/gu-io/gu"
	"github.com/gu-io/gu/router"
	"github.com/gu-io/gu/shell"
	"github.com/gu-io/gu/shell/cache/memorycache"
)

// WebviewDriver provides a concrete implementation of the Gu.Driver interface.
type WebviewDriver struct {
	*QTApp
	attr          *QTAttr
	readyHandlers []func()
}

// NewWebviewDriver returns a new instance of the qt WebviewDriver.
func NewWebviewDriver(attr *QTAttr) *WebviewDriver {
	var driver WebviewDriver
	driver.attr = attr

	driver.QTApp = NewQTApp(attr)

	return &driver
}

// Ready tells the driver to initialize its operations has
// the registered app is ready for rendering.
func (wd WebviewDriver) Ready() {
	wd.QTApp.Init()
}

// Name returns the name of the driver.
func (WebviewDriver) Name() string {
	return "QTWebkit"
}

// OnReady registers the giving handle function to be called when the giving
// driver is ready and loaded.
func (WebviewDriver) OnReady(handle func()) {

}

// Location returns the current location of the browser.
func (WebviewDriver) Location() router.PushEvent {
	var route router.PushEvent
	return route
}

// Navigate takes the provided route and navigates the rendering system to the
// desired page.
func (WebviewDriver) Navigate(route router.PushDirectiveEvent) {

}

// OnRoute registers the NApp instance for route changes and re-rendering.
func (WebviewDriver) OnRoute(app *gu.NApp) {

}

// Render issues a clean rendering of all content clearing out the current content
// of the browser to the one provided by the appliation.
func (WebviewDriver) Render(app *gu.NApp) {

}

// Update updates a giving view portion of a giving App within the designated
// rendering system(browser) provided by the driver.
func (WebviewDriver) Update(app *gu.NApp, view *gu.NView) {

}

// Services returns the Fetcher and Cache associated with the provided cacheName.
// Intercepting requests for usage.
func (WebviewDriver) Services(cacheName string, intercept bool) (shell.Fetch, shell.Cache) {
	inMemoryCache := memorycache.New(cacheName)
	fetcher := fetch.New(inMemoryCache)

	return fetcher, inMemoryCache
}
