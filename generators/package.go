package generators

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

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
		gen.SourceText(
			string(data.Must("scaffolds/base.gen")),
			struct {
				Name string
			}{
				Name: componentName,
			},
		),
	)

	publicStandInGen := gen.Block(
		gen.Package(
			gen.Name(componentNameLower),
			gen.SourceText(
				string(data.Must("scaffolds/bundle.gen")),
				struct {
					Name    string
					Package string
				}{
					Name:    componentName,
					Package: componentNameLower,
				},
			),
		),
	)

	publicGen := gen.Block(
		gen.SourceText(
			string(data.Must("scaffolds/pack-bundle.gen")),
			struct {
				Name      string
				LessFile  string
				Package   string
				TargetDir string
			}{
				TargetDir: "public",
				Name:      componentName,
				Package:   componentNameLower,
				LessFile:  fmt.Sprintf("less/%s.less", componentNameLower),
			},
		),
	)

	settingsGen := gen.Block(
		gen.SourceText(
			string(data.Must("scaffolds/settings.gen")),
			struct {
				Name    string
				Package string
			}{
				Name:    componentName,
				Package: componentNameLower,
			},
		),
	)

	lessGen := gen.Block(
		gen.SourceText(
			string(data.Must("scaffolds/main.less.gen")),
			struct {
				Name    string
				Package string
			}{
				Name:    componentName,
				Package: componentNameLower,
			},
		),
	)

	tomlGen := gen.Block(
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
			Dir:          filepath.Join(componentNameLower, "public/less"),
		},
		{
			DontOverride: false,
			Writer:       lessGen,
			Dir:          filepath.Join(componentNameLower, "public/less"),
			FileName:     fmt.Sprintf("%s.less", componentNameLower),
		},
		{
			DontOverride: true,
			Writer:       tomlGen,
			Dir:          componentNameLower,
			FileName:     "settings.toml",
		},
		{
			DontOverride: false,
			Dir:          componentNameLower,
			FileName:     "settings_bundle.go",
			Writer:       fmtwriter.New(settingsGen, true, true),
		},
		{
			DontOverride: true,
			Dir:          componentNameLower,
			FileName:     "public_bundle.go",
			Writer:       fmtwriter.New(publicGen, true, true),
		},
		{
			DontOverride: true,
			Dir:          filepath.Join(componentNameLower, "public"),
			FileName:     fmt.Sprintf("%s_bundle.go", componentNameLower),
			Writer:       fmtwriter.New(publicStandInGen, true, true),
		},
		{
			DontOverride: true,
			Dir:          componentNameLower,
			FileName:     fmt.Sprintf("%s.go", componentNameLower),
			Writer:       fmtwriter.New(typeGen, true, true),
		},
	}, nil
}
