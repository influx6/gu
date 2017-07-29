package generators

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/influx6/faux/fmtwriter"

	"github.com/gu-io/gu/generators/data"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

// GuPackageGenerator which defines a  function for generating a type for receiving a giving
//	struct type has a notification type which can then be wired as a notification.EventDistributor.
//
func GuPackageGenerator(an ast.AnnotationDeclaration, pkg ast.PackageDeclaration) ([]gen.WriteDirective, error) {
	if len(an.Arguments) == 0 {
		return nil, errors.New("Expected atleast one argument for annotation as component name")
	}

	componentName := an.Arguments[0]
	componentNameLower := strings.ToLower(componentName)

	typeGen := gen.Block(
		gen.Package(
			gen.Name(componentName),
			gen.Imports(
				gen.Import("github.com/gu-io/gu", ""),
				gen.Import("github.com/gu-io/gu/trees", ""),
				gen.Import("github.com/gu-io/gu/trees/elems", ""),
				gen.Import("github.com/gu-io/gu/trees/property", ""),
			),
			gen.Block(
				gen.SourceTextWith(
					string(data.Must("scaffolds/base.gen")),
					template.FuncMap{},
					struct {
						Name string
					}{
						Name: componentName,
					},
				),
			),
		),
	)

	settingsGen := gen.Block(
		gen.SourceText(
			string(data.Must("scaffolds/settings.gen")),
			nil,
		),
	)

	pipeGen := gen.Block(
		gen.SourceText(
			string(data.Must("scaffolds/settings.toml.gen")),
			nil,
		),
	)

	return []gen.WriteDirective{
		{
			DontOverride: false,
			Dir:          componentNameLower,
		},
		{
			DontOverride: false,
			Dir:          filepath.Join(componentNameLower, "public"),
		},
		{
			DontOverride: false,
			Dir:          filepath.Join(componentNameLower, "settings.toml"),
			Writer:       pipeGen,
		},
		{
			DontOverride: false,
			Dir:          componentNameLower,
			FileName:     "generator.go",
			Writer:       fmtwriter.New(settingsGen, true, true),
		},
		{
			DontOverride: false,
			Dir:          filepath.Join(componentNameLower, "app"),
			FileName:     fmt.Sprintf("%s.go", componentNameLower),
			Writer:       fmtwriter.New(typeGen, true, true),
		},
	}, nil
}
