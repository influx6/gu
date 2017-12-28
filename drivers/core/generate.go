// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

var (
	pkg = `// Package core contains javascript sources which are exported into specific
// drivers webview which exposes similar functionality needed to interoperate with
// the platform.

// Document is auto-generate and should not be modified by hand.

//go:generate go run generate.go

package core

// JavascriptDriverCore contains the javascript code to be injected into a webview
// to provide similar functionality desired to have it working with Gu.
`

	pkgVar = "var JavascriptDriverCore = `%s`\n"
)

func main() {
	js, err := ioutil.ReadFile("./core.js")
	if err != nil {
		panic(fmt.Sprintf("Unable to locate `core.js` file: %q", err.Error()))
	}

	goFile, err := os.OpenFile("./corejs.go", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(fmt.Sprintf("Unable to create `corejs.go` file: %q", err.Error()))
	}

	defer goFile.Close()

	if _, err := fmt.Fprint(goFile, pkg); err != nil {
		panic(fmt.Sprintf("Unable to write data to `corejs.go`: %q", err.Error()))
	}

	if _, err := fmt.Fprintf(goFile, pkgVar, js); err != nil {
		panic(fmt.Sprintf("Unable to write data to `corejs.go`: %q", err.Error()))
	}
}
