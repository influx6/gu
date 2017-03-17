package qtwebview

import (
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/webengine"
	"github.com/therecipe/qt/widgets"
)

//===========================================================================================

// QTAttr defines a attribute struct which provides specific configurations
// for a QTApp.
type QTAttr struct {
	URL    string
	Width  int
	Height int
}

// QTApp defines a struct which creates a QTWindow with a QTWebkit for loading singularly
// the provided gu.NApp.
type QTApp struct {
	attr    QTAttr
	baseURL *core.QUrl
	window  *widgets.QMainWindow
	fm      *widgets.QWidget
	view    *webengine.QWebEngineView
}

// NewQTApp returns a new instance of the QTApp.
func NewQTApp(attr QTAttr) *QTApp {
	var app QTApp
	app.attr = attr
	app.baseURL = core.NewQUrl3(attr.URL, core.QUrl__DecodedMode)
	return &app
}

// Init initializes the QTWindow and sets the desired
// widgets and webview.
func (qt *QTApp) Init() {

	// Create a new widget window.
	qt.window = widgets.NewQMainWindow(nil, 0)

	if qt.attr.Width > 0 {
		qt.window.SetMinimumWidth(qt.attr.Width)
	}

	if qt.attr.Height > 0 {
		qt.window.SetMaximumHeight(qt.attr.Height)
	}

	// create a widget group later.
	qt.fm = widgets.NewQWidget(nil, 0)
	qt.fm.SetLayout(widgets.NewQVBoxLayout())

	// Create the view we wish to render with.
	qt.view = webengine.NewQWebEngineView(nil)
	qt.fm.Layout().AddWidget(qt.view)

	// Add widget to the window.
	qt.window.SetCentralWidget(qt.fm)
}

// URI returns the QTApp internal URL.
func (qt *QTApp) URI() *core.QUrl {
	return qt.baseURL
}

// View returns the underline webview for the view.
func (qt *QTApp) View() *webengine.QWebEngineView {
	return qt.view
}

// Run initializes the window to show and calls the needed
// methods to intialize and block the routine call till exit.
func (qt *QTApp) Run() {

	// Set the base view for the webview.
	qt.view.SetHtml("", qt.baseURL)

	// Ask the window to show.
	qt.window.Show()

	// Block till exit.
	widgets.QApplication_Exec()
}

//===========================================================================================

var qtInit bool

// InitQTApplication initializes the QTApplication to be used for
// the initialization of the application and must be called in the
// main() function.
func InitQTApplication() {
	if !qtInit {
		widgets.NewQApplication(len(os.Args), os.Args)
		webengine.QtWebEngine_Initialize()
		qtInit = true
	}
}
