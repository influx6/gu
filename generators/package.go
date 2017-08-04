// +build !js

package generators

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/influx6/faux/fmtwriter"

	"github.com/gu-io/gu/generators/data"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

var (
	inGOPATH    = os.Getenv("GOPATH")
	inGOPATHSrc = filepath.Join(inGOPATH, "src")
)

// GuPackageGenerator which defines a  function for generating a type for receiving a giving
//	struct type has a notification type which can then be wired as a notification.EventDistributor.
//
func GuPackageGenerator(an ast.AnnotationDeclaration, pkg ast.PackageDeclaration) ([]gen.WriteDirective, error) {
	if len(an.Arguments) == 0 {
		return nil, errors.New("Expected atleast one argument for annotation as component name")
	}

	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve current directory path: %+q", err)
	}

	guSrc := filepath.Join(inGOPATHSrc, "github.com/gu-io/gu")
	gutheme := filepath.Join(guSrc, "common/themes")
	guGridCSS := filepath.Join(gutheme, "/grids/grid.css")
	guNormCSS := filepath.Join(gutheme, "/normalize/normalize.css")

	gridCSSData, err := ioutil.ReadFile(guGridCSS)
	if err != nil {
		fmt.Printf("Failed to retrieve /grid/grid.css file from gu src dir: %q -> %+q\n", guGridCSS, err)
		return nil, err
	}

	gridNormCSS, err := ioutil.ReadFile(guNormCSS)
	if err != nil {
		fmt.Printf("Failed to retrieve /normalize/normalize.css file from gu src dir: %q -> %+q\n", guNormCSS, err)
		return nil, err
	}

	packageDir, err := filepath.Rel(inGOPATHSrc, workDir)
	if err != nil {
		fmt.Printf("Failed to retrieve package directory path in go src: %+q\n", err)
	}

	componentName := an.Arguments[0]
	componentNameLower := strings.ToLower(componentName)

	componentPackageDir := filepath.Join(packageDir, componentNameLower)

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
			string(data.Must("scaffolds/pack-bundle-public.gen")),
			struct {
				Name          string
				LessFile      string
				Package       string
				TargetDir     string
				TargetPackage string
			}{
				TargetDir:     "public",
				TargetPackage: "public",
				Name:          componentName,
				Package:       componentNameLower,
				LessFile:      fmt.Sprintf("less/%s.less", componentNameLower),
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
			struct {
				Name    string
				Package string
			}{
				Name:    componentName,
				Package: componentPackageDir,
			},
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
			Dir:          filepath.Join(componentNameLower, "public/css"),
			FileName:     "normalize.css",
			Writer:       bytes.NewBuffer(gridNormCSS),
		},
		{
			DontOverride: false,
			Dir:          filepath.Join(componentNameLower, "public/css"),
			FileName:     "grid.css",
			Writer:       bytes.NewBuffer(gridCSSData),
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
