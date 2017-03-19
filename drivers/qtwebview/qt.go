// Package qtwebview implements the gu.Driver interface to provide rendering using the
// qt5 cross-platform GUI library and platform.
package qtwebview

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gu-io/gu"
	"github.com/gu-io/gu/router"
	"github.com/gu-io/gu/shell"
	"github.com/gu-io/gu/shell/cache/memorycache"
	"github.com/gu-io/gu/shell/fetch"
	"github.com/therecipe/qt/core"
)

// WebviewDriver provides a concrete implementation of the Gu.Driver interface.
type WebviewDriver struct {
	*QTApp
	attr          QTAttr
	readyHandlers []func()
	manifestData  []byte
}

// NewWebviewDriver returns a new instance of the qt WebviewDriver.
func NewWebviewDriver(attr QTAttr) *WebviewDriver {
	var driver WebviewDriver
	driver.attr = attr
	driver.QTApp = NewQTApp(attr)

	if attr.Manifest != "" {
		if err := driver.loadManifest(); err != nil {
			fmt.Printf("Error failed Loading manifest data: %q \n", err.Error())
		}
	}

	return &driver
}

// Ready tells the driver to initialize its operations has
// the registered app is ready for rendering.
func (wd WebviewDriver) Ready() {
	wd.QTApp.Init()

	wd.QTApp.view.ConnectUrlChanged(func(loc *core.QUrl) {

	})

	wd.QTApp.view.ConnectLoadFinished(func(ok bool) {

	})

	for _, ready := range wd.readyHandlers {
		ready()
	}
}

// Name returns the name of the driver.
func (WebviewDriver) Name() string {
	return "QTWebkit"
}

// OnReady registers the giving handle function to be called when the giving
// driver is ready and loaded.
func (wd WebviewDriver) OnReady(handle func()) {
	wd.readyHandlers = append(wd.readyHandlers, handle)
}

// Location returns the current location of the browser.
func (wd WebviewDriver) Location() router.PushEvent {
	var route router.PushEvent
	route.To = "/"
	route.Path = "/"

	// loc := wd.QTApp.view.Url()

	return route
}

// Navigate takes the provided route and navigates the rendering system to the
// desired page.
func (wd WebviewDriver) Navigate(route router.PushDirectiveEvent) {

}

// OnRoute registers the NApp instance for route changes and re-rendering.
func (wd WebviewDriver) OnRoute(app *gu.NApp) {

}

// Render issues a clean rendering of all content clearing out the current content
// of the browser to the one provided by the appliation.
func (wd WebviewDriver) Render(app *gu.NApp) {
	qt.view.SetHtml(app.Render(nil).HTML(), qt.baseURL)
}

// Update updates a giving view portion of a giving App within the designated
// rendering system(browser) provided by the driver.
func (wd WebviewDriver) Update(app *gu.NApp, view *gu.NView) {

}

// Services returns the Fetcher and Cache associated with the provided cacheName.
// Intercepting requests for usage.
func (wd WebviewDriver) Services(cacheName string, intercept bool) (shell.Fetch, shell.Cache) {
	mcache := memorycache.New(cacheName)
	fetcher := fetch.New(mcache)

	if wd.manifestData != nil {
		mcache.PutPath(wd.attr.Manifest, shell.WebResponse{
			Ok:         true,
			Status:     200,
			Body:       wd.manifestData,
			Type:       "application/json",
			StatusText: "OK",
			FinalURL:   wd.attr.Manifest,
		})
	}

	return fetcher, mcache
}

// loadManifest loads the giving manifest file data for easy manifest access on the cache.
func (wd WebviewDriver) loadManifest() error {
	if wd.manifestData != nil {
		return nil
	}

	cdir, err := os.Getwd()
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(filepath.Join(cdir, wd.attr.Manifest))
	if err != nil {
		return err
	}

	wd.manifestData = data
	return nil
}
