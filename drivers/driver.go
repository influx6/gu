package drivers

import (
	"github.com/gu-io/gu"
	"github.com/gu-io/gu/router"
	"github.com/gu-io/gu/shell"
)

// NOOPDriver provides a concrete implementation of the Gu.Driver interface.
type NOOPDriver struct {
	readyHandlers []func()
	started       bool
}

// NewNOOPDriver returns a new instance of a js driver.
func NewNOOPDriver() *NOOPDriver {
	var driver NOOPDriver
	return &driver
}

// Ready is called to intialize the driver and load the page.
func (driver *NOOPDriver) Ready() {
	for _, ready := range driver.readyHandlers {
		ready()
	}
}

// Name returns the name of the driver.
func (driver *NOOPDriver) Name() string {
	return "NO-OP Driver"
}

// OnReady registers the giving handle function to be called when the giving
// driver is ready and loaded.
func (driver *NOOPDriver) OnReady(handle func()) {
	driver.readyHandlers = append(driver.readyHandlers, handle)
}

// Location returns the current location of the browser.
func (driver *NOOPDriver) Location() router.PushEvent {
	return router.PushEvent{
		// Host:   host,
		// Path:   path,
		// Hash:   hash,
		// Rem:    hash,
		// From:   location,
		Params: make(map[string]string),
	}
}

// Navigate takes the provided route and navigates the rendering system to the
// desired page.
func (driver *NOOPDriver) Navigate(route router.PushDirectiveEvent) {
	// SetLocationByPushEvent(route)
}

// OnRoute registers the NApp instance for route changes and re-rendering.
func (driver *NOOPDriver) OnRoute(app *gu.NApp) {
}

// Render issues a clean rendering of all content clearing out the current content
// of the browser to the one provided by the appliation.
func (driver *NOOPDriver) Render(app *gu.NApp) {
	driver.started = true
}

// Update updates a giving view portion of a giving App within the designated
// rendering system(browser) provided by the driver.
func (driver *NOOPDriver) Update(app *gu.NApp, view *gu.NView) {
}

// Services returns the Fetcher and Cache associated with the provided cacheName.
// Intercepting requests for usage.
func (driver *NOOPDriver) Services(cacheName string, intercept bool) (shell.Fetch, shell.Cache) {
	return nil, nil
}
