// +build !js

package generators

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gu-io/gu/common"
	"github.com/gu-io/gu/generators/data"
	"github.com/influx6/faux/fmtwriter"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

// JSDriverGenerator which defines a  function for generating a type for receiving a giving
//	struct type has a notification type which can then be wired as a notification.EventDistributor.
//
func JSDriverGenerator(an ast.AnnotationDeclaration, pkg ast.PackageDeclaration) ([]gen.WriteDirective, error) {
	var config common.Settings

	// Load settings into configuration.
	if _, err := toml.DecodeFile("./settings.toml", &config); err != nil {
		return nil, fmt.Errorf("Please execute command where settings.toml file is located: %+q", err)
	}

	if err := config.Public.Validate(); err != nil {
		return nil, err
	}

	jsFileName := fmt.Sprintf("%s_app_bundle.js", strings.ToLower(config.App))

	var jsPath string

	if config.Static.IndexDir == "./" {
		jsPath = filepath.Join("public", "js")
	} else {
		jsPath = "js"
	}

	htmlGen := gen.Block(
		gen.SourceText(
			string(data.Must("scaffolds/base.html.gen")),
			struct {
				Name   string
				Path   string
				JSFile string
			}{
				Name:   config.App,
				Path:   config.Public.Path,
				JSFile: fmt.Sprintf("%s/%s", jsPath, jsFileName),
			},
		),
	)

	jsGen := gen.Block(
		gen.SourceText(
			string(data.Must("scaffolds/jsdriver.gen")),
			struct {
				Name    string
				Package string
				Path    string
				JSFile  string
			}{
				Name:    config.App,
				Package: config.Package,
				Path:    config.Public.Path,
				JSFile:  filepath.Join("../../", config.Public.Path, "js", jsFileName),
			},
		),
	)

	return []gen.WriteDirective{
		{
			DontOverride: false,
			FileName:     "main.go",
			Dir:          "./driver/js",
			Writer:       fmtwriter.New(jsGen, true, true),
		},
		{
			DontOverride: false,
			Writer:       htmlGen,
			FileName:     "index.html",
			Dir:          config.Static.IndexDir,
		},
	}, nil
}
