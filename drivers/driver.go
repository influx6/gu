package drivers

import (
	"github.com/gu-io/gu"
	"github.com/gu-io/gu/router"
)

// NOOPDriver provides a concrete implementation of the Gu.Driver interface.
type NOOPDriver struct {
	app     *gu.NApp
	started bool
}

// NewNOOPDriver returns a new instance of a js driver.
func NewNOOPDriver(app *gu.NApp) *NOOPDriver {
	var driver NOOPDriver
	driver.app = app
	return &driver
}

// Name returns the name of the driver.
func (driver *NOOPDriver) Name() string {
	return "NO-OP Driver"
}

// Location returns the current location of the browser.
func (driver *NOOPDriver) Location() router.PushEvent {
	return router.PushEvent{
		Host:   "noop.com",
		Path:   "/",
		Hash:   "#",
		From:   "/#",
		Params: make(map[string]string),
	}
}

// Navigate takes the provided route and navigates the rendering system to the
// desired page.
func (driver *NOOPDriver) Navigate(route router.PushDirectiveEvent) {
	var ps router.PushEvent

	ps, err := router.NewPushEvent(route.To, true)
	if err != nil {
		ps = router.PushEvent{
			Path: route.To,
		}
	}

	driver.app.ActivateRoute(ps)
}
