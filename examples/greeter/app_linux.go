// +build !js

package main

import (
	"github.com/gu-io/gu"
)

// AppSettings defines the attributes for the app for using the GopherJS
// driver.
var AppSettings = gu.AppAttr{
	InterceptRequests: true,
	Name:              "{{ .Name }}",
	Mode:              gu.DevelopmentMode,
	Title:             "{{ .Name }} Gu App",
	Manifests:         "assets/manifests.json",
}
