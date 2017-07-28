package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gu-io/gu/generators"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/sentries/stdout"
	"github.com/influx6/moz/annotations"
	"github.com/influx6/moz/ast"
	cli "gopkg.in/urfave/cli.v2"
)

var (
	version     = "0.0.1"
	defaultName = "manifests"
	commands    = []*cli.Command{}
	events      = metrics.New(stdout.Stdout{})

	namebytes         = []byte("{{Name}}")
	pkgbytes          = []byte("{{PKG}}")
	sourcebytes       = []byte("{{SOURCE}}")
	extbytes          = []byte("{{EXTENSIONS}}")
	goPathbytes       = []byte("{{GOPATH}}")
	pkgContentbytes   = []byte("{{PKG_CONTENT}}")
	pkgNamebytes      = []byte("{{PKGNAME}}")
	fileNamebytes     = []byte("{{FILENAME}}")
	filesDirNamebytes = []byte("{{FILESDIRNAME}}")
	dirNamebytes      = []byte("{{DIRNAME}}")
	nameLowerbytes    = []byte("{{Name_Lower}}")
	idLowerNamebytes  = []byte("{{Name_Lower_Id}}")

	gupath = "github.com/gu-io/gu"

	usage = `CLI to generate go code for the use in development with Gu.`
)

func main() {
	initCommands()

	app := &cli.App{}
	app.Name = "Gu"
	app.Version = version
	app.Commands = commands
	app.Usage = usage

	app.Run(os.Args)
}

func generateAddFile(name string, content []byte) string {
	return fmt.Sprintf(`
		addFile(%q, []byte(%+q))
`, name, content)
}

func capitalize(val string) string {
	return strings.ToUpper(val[:1]) + val[1:]
}

var badSymbols = regexp.MustCompile(`[(|\-|_|\W|\d)+]`)
var notAllowed = regexp.MustCompile(`[^(_|\w|\d)+]`)
var descore = regexp.MustCompile("-")

func validateName(val string) bool {
	return notAllowed.MatchString(val)
}

func initCommands() {
	var subcommands []*cli.Command

	subcommands = append(subcommands, &cli.Command{
		Name:        "css",
		Usage:       "gu css <css-dir-name>",
		Description: "Generates a styles package which builds all internal css files into a central go file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "name=hello",
			},
		},
		Action: func(ctx *cli.Context) error {
			args := ctx.Args()
			if args.Len() == 0 {
				return nil
			}

			cdir, err := os.Getwd()
			if err != nil {
				return err
			}

			cssDirName := ctx.String("name")
			if cssDirName == "" && args.Len() > 0 {
				cssDirName = args.First()
			}

			if cssDirName == "" && args.Len() == 0 {
				cssDirName = "gcss"
			}

			gopath := os.Getenv("GOPATH")
			gup := filepath.Join(gopath, "src")
			gupkg := filepath.Join(gup, gupath)
			cssDirPath := filepath.Join(cdir, cssDirName)

			if err = os.MkdirAll(cssDirPath, 0777); err != nil {
				return err
			}

			fmt.Printf("- Adding project directory: %q\n", filepath.Base(cssDirPath))

			if err = os.MkdirAll(filepath.Join(cssDirPath, "css"), 0777); err != nil {
				return err
			}

			fmt.Printf("- Adding project directory: %q\n", filepath.Join(filepath.Base(cssDirPath), "css"))

			gendata, err := ioutil.ReadFile(filepath.Join(gupkg, "templates/css.template"))
			if err != nil {
				return err
			}

			cssgendata, err := ioutil.ReadFile(filepath.Join(gupkg, "templates/cssgenerate.template"))
			if err != nil {
				return err
			}

			plainPKGData, err := ioutil.ReadFile(filepath.Join(gupkg, "templates/plain_generated_pkg.template"))
			if err != nil {
				return err
			}

			gendata = []byte(fmt.Sprintf("%q", gendata))
			cssgendata = bytes.Replace(cssgendata, pkgContentbytes, gendata, 1)
			cssgendata = bytes.Replace(cssgendata, dirNamebytes, []byte("css"), 1)
			cssgendata = bytes.Replace(cssgendata, pkgNamebytes, []byte("\""+cssDirName+"\""), 1)
			plainPKGData = bytes.Replace(plainPKGData, pkgNamebytes, []byte(cssDirName), -1)

			if err := writeFile(filepath.Join(cssDirPath, "generate.go"), cssgendata); err != nil {
				return err
			}

			fmt.Printf("- Adding project file: %q\n", filepath.Join(filepath.Base(cssDirPath), "generate.go"))

			if err := writeFile(filepath.Join(cssDirPath, "css.go"), plainPKGData); err != nil {
				return err
			}

			fmt.Printf("- Adding project file: %q\n", filepath.Join(filepath.Base(cssDirPath), "css.go"))

			return nil
		},
	})

	subcommands = append(subcommands, &cli.Command{
		Name:  "new",
		Usage: "gu new <component-name>",
		Description: `Generates a new boilerplate for giving component name.

		Options:
			- flat: When true will generate only a .go file for component. (Default: False)
			- stand: When true will generaly only go package with css style. (Default: False)
			- base:	When true will check for components dir first before generating go package. (Default: True)
		`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "base-package",
				Aliases: []string{"base"},
				Usage:   "base=true",
				Value:   true,
			},
			&cli.BoolFlag{
				Name:    "stand-alone",
				Aliases: []string{"stand"},
				Usage:   "stand=true",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "flat-file",
				Aliases: []string{"flat"},
				Usage:   "flat=true",
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "component",
				Aliases: []string{"c"},
				Usage:   "component=hello",
			},
		},
		Action: func(ctx *cli.Context) error {
			componentName := ctx.String("component")
			flat := ctx.Bool("flat-file")
			standAlone := ctx.Bool("stand-alone")
			base := ctx.Bool("base-of-components-package")

			args := ctx.Args()
			if args.Len() == 0 && componentName == "" {
				return nil
			}

			cdir, err := os.Getwd()
			if err != nil {
				return err
			}

			gopath := os.Getenv("GOPATH")
			gup := filepath.Join(gopath, "src")
			gupkg := filepath.Join(gopath, "src", gupath)
			componentDir := filepath.Join(cdir, "components")

			if componentName == "" && args.Len() > 0 {
				componentName = args.First()
			}

			// componentStructName := descore.ReplaceAllString(componentName, "_")
			componentStructName := badSymbols.ReplaceAllString(componentName, "")
			if validateName(componentStructName) {
				return errors.New("ComponentName does not meet go struct naming standards: " + componentStructName)
			}

			componentNameCap := capitalize(componentStructName)
			componentNameLower := strings.ToLower(componentStructName)

			componentPkgName := badSymbols.ReplaceAllString(componentName, "")
			newNoComponentDir := filepath.Join(cdir, componentPkgName)
			newComponentDir := filepath.Join(componentDir, componentPkgName)

			cssDirName := "styles"
			newComponentCSSDir := filepath.Join(newComponentDir, cssDirName)
			newComponentCSSFilesDir := filepath.Join(newComponentCSSDir, "css")
			newNoComponentCSSDir := filepath.Join(newNoComponentDir, cssDirName)
			newNoComponentCSSFilesDir := filepath.Join(newNoComponentCSSDir, "css")

			packagePath, err := filepath.Rel(gup, cdir)
			if err != nil {
				return err
			}

			cssbeforegendata, cerr := ioutil.ReadFile(filepath.Join(gupkg, "templates/css.template"))
			if cerr != nil {
				return cerr
			}

			cssgendata, merr := ioutil.ReadFile(filepath.Join(gupkg, "templates/cssgenerate.template"))
			if merr != nil {
				return merr
			}

			plainPKGData, err := ioutil.ReadFile(filepath.Join(gupkg, "templates/plain_generated_pkg.template"))
			if err != nil {
				return err
			}

			if standAlone {
				// fmt.Printf("Will to %q\n", newNoComponentDir)
				if err = os.Mkdir(newNoComponentDir, 0777); err != nil {
					return err
				}

				if err = os.MkdirAll(newNoComponentCSSFilesDir, 0777); err != nil {
					return err
				}

				fmt.Printf("- Adding project package: %q\n", filepath.Join("", componentPkgName))
				fmt.Printf("- Adding project directory: %q\n", filepath.Join("", componentPkgName, cssDirName))
				fmt.Printf("- Adding project directory: %q\n", filepath.Join("", componentPkgName, cssDirName, "css"))

				cssbeforegendata = []byte(fmt.Sprintf("%q", cssbeforegendata))
				cssgendata = bytes.Replace(cssgendata, pkgContentbytes, cssbeforegendata, 1)
				cssgendata = bytes.Replace(cssgendata, dirNamebytes, []byte("css"), 1)
				cssgendata = bytes.Replace(cssgendata, pkgNamebytes, []byte("\""+cssDirName+"\""), 1)
				plainPKGData = bytes.Replace(plainPKGData, pkgNamebytes, []byte(cssDirName), -1)

				if err = writeFile(filepath.Join(newNoComponentCSSDir, "generate.go"), cssgendata); err != nil {
					return err
				}

				fmt.Printf("- Adding project file: %q\n", filepath.Join("", componentPkgName, "styles", "generate.go"))

				if err := writeFile(filepath.Join(newNoComponentCSSDir, "css.go"), plainPKGData); err != nil {
					return err
				}

				fmt.Printf("- Adding project file: %q\n", filepath.Join("", componentPkgName, "styles", "css.go"))

				cpdata, cperr := ioutil.ReadFile(filepath.Join(gupkg, "templates/nopkgcomponent.template"))
				if cperr != nil {
					return cperr
				}

				cpdata = bytes.Replace(cpdata, pkgNamebytes, []byte(componentPkgName), -1)
				cpdata = bytes.Replace(cpdata, pkgbytes, []byte(packagePath), -1)
				cpdata = bytes.Replace(cpdata, namebytes, []byte(componentNameCap), -1)
				cpdata = bytes.Replace(cpdata, nameLowerbytes, []byte(componentNameLower), -1)
				cpdata = bytes.Replace(cpdata, idLowerNamebytes, []byte(componentNameLower[:2]), -1)

				componentFileName := fmt.Sprintf("%s.go", componentNameLower)
				cmdir := filepath.Join(newNoComponentDir, componentFileName)

				if err = writeFile(cmdir, cpdata); err != nil {
					return err
				}

				fmt.Printf("- Adding project file: %q\n", filepath.Join("", componentPkgName, componentFileName))
				return nil
			}

			if flat && base {
				baseDir := filepath.Base(componentDir)

				cpdata, cerr := ioutil.ReadFile(filepath.Join(gupkg, "templates/component-base.template"))
				if cerr != nil {
					return cerr
				}

				cpdata = bytes.Replace(cpdata, namebytes, []byte(componentNameCap), -1)
				cpdata = bytes.Replace(cpdata, nameLowerbytes, []byte(componentNameLower), -1)
				cpdata = bytes.Replace(cpdata, idLowerNamebytes, []byte(componentNameLower[:2]), -1)

				componentFileName := fmt.Sprintf("%s.go", componentNameLower)

				if _, merr := os.Stat(componentDir); merr != nil {
					componentDir, err = findLower(packagePath, "components")
					if err != nil {
						return err
					}

					componentDir = filepath.Join(gup, componentDir)
				}

				cmdir := filepath.Join(componentDir, componentFileName)
				if cerr := writeFile(cmdir, cpdata); cerr != nil {
					return cerr
				}

				fmt.Printf("- Adding project file: %q\n", filepath.Join(baseDir, componentFileName))
				return nil
			}

			if flat && !base {
				baseDir := filepath.Base(cdir)

				cpdata, cerr := ioutil.ReadFile(filepath.Join(gupkg, "templates/component.template"))
				if cerr != nil {
					return cerr
				}

				// Attempt to follow path down the stack to see if we can match it and
				// cheat.
				componentsPackagePath, coerr := findLower(packagePath, "components")
				if coerr != nil {

					// We couldn't cheat, so we follow the hard road and stat down the pipe
					// by attempting to see if we find a ./components dir down the tree.
					componentsPackagePath, coerr = findLowerByStat(gup, packagePath, "components", true)
					if coerr != nil {

						// We are still unable to find it, so just match if we are at the root directory
						// and we possibly just went stupid.
						possiblePath := filepath.Join(packagePath, "components")
						if _, err := os.Stat(possiblePath); err != nil {
							return fmt.Errorf("Error: %+q not found in %q", coerr, packagePath)
						}

						componentsPackagePath = possiblePath
					}
				}

				cpdata = bytes.Replace(cpdata, pkgNamebytes, []byte(baseDir), -1)
				cpdata = bytes.Replace(cpdata, pkgbytes, []byte(componentsPackagePath), -1)
				cpdata = bytes.Replace(cpdata, namebytes, []byte(componentNameCap), -1)
				cpdata = bytes.Replace(cpdata, nameLowerbytes, []byte(componentNameLower), -1)
				cpdata = bytes.Replace(cpdata, idLowerNamebytes, []byte(componentNameLower[:2]), -1)

				componentFileName := fmt.Sprintf("%s.go", componentNameLower)

				cmdir := filepath.Join(cdir, componentFileName)
				if cerr := writeFile(cmdir, cpdata); cerr != nil {
					return cerr
				}

				fmt.Printf("- Adding project file: %q\n", filepath.Join(baseDir, componentFileName))

				return nil
			}

			if err = os.Mkdir(newComponentDir, 0777); err != nil {
				return err
			}

			if err = os.MkdirAll(newComponentCSSFilesDir, 0777); err != nil {
				return err
			}

			fmt.Printf("- Adding project package: %q\n", filepath.Join("components", componentPkgName))
			fmt.Printf("- Adding project directory: %q\n", filepath.Join("components", componentPkgName, cssDirName))
			fmt.Printf("- Adding project directory: %q\n", filepath.Join("components", componentPkgName, cssDirName, "css"))

			cssbeforegendata = []byte(fmt.Sprintf("%q", cssbeforegendata))
			cssgendata = bytes.Replace(cssgendata, pkgContentbytes, cssbeforegendata, 1)
			cssgendata = bytes.Replace(cssgendata, dirNamebytes, []byte("css"), 1)
			cssgendata = bytes.Replace(cssgendata, pkgNamebytes, []byte("\""+cssDirName+"\""), 1)
			plainPKGData = bytes.Replace(plainPKGData, pkgNamebytes, []byte(cssDirName), -1)

			if err = writeFile(filepath.Join(newComponentCSSDir, "generate.go"), cssgendata); err != nil {
				return err
			}

			fmt.Printf("- Adding project file: %q\n", filepath.Join("components", componentPkgName, "styles", "generate.go"))

			if err := writeFile(filepath.Join(newComponentCSSDir, "css.go"), plainPKGData); err != nil {
				return err
			}

			fmt.Printf("- Adding project file: %q\n", filepath.Join("components", componentPkgName, "styles", "css.go"))

			cpdata, cperr := ioutil.ReadFile(filepath.Join(gupkg, "templates/pkgcomponent.template"))
			if cperr != nil {
				return cperr
			}

			cpdata = bytes.Replace(cpdata, pkgNamebytes, []byte(componentPkgName), -1)
			cpdata = bytes.Replace(cpdata, pkgbytes, []byte(packagePath), -1)
			cpdata = bytes.Replace(cpdata, namebytes, []byte(componentNameCap), -1)
			cpdata = bytes.Replace(cpdata, nameLowerbytes, []byte(componentNameLower), -1)
			cpdata = bytes.Replace(cpdata, idLowerNamebytes, []byte(componentNameLower[:2]), -1)

			componentFileName := fmt.Sprintf("%s.go", componentNameLower)
			cmdir := filepath.Join(newComponentDir, componentFileName)

			if err = writeFile(cmdir, cpdata); err != nil {
				return err
			}

			fmt.Printf("- Adding project file: %q\n", filepath.Join("components", componentPkgName, componentFileName))
			return nil
		},
	})

	commands = append(commands, &cli.Command{
		Name:        "files",
		Usage:       "gu files -dir=myfiles -extensions='.bo .loc .gob' <pkg-name>",
		Description: "Generates a package which builds all internal files that matches provided optional extension list into a go file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dirName",
				Aliases: []string{"dir"},
				Usage:   "dir=assets",
			},
			&cli.StringFlag{
				Name:    "extensions",
				Aliases: []string{"e"},
				Usage:   "extensions='.gob, .loc'",
			},
		},
		Action: func(ctx *cli.Context) error {
			mDirName := ctx.String("dir")
			vDirName := ctx.String("name")
			args := ctx.Args()
			if args.Len() == 0 && vDirName == "" {
				return nil
			}

			if mDirName == "" {
				mDirName = "files"
			}

			cdir, err := os.Getwd()
			if err != nil {
				return err
			}

			var extBu bytes.Buffer
			fmt.Fprintf(&extBu, "[]string{")

			extensions := strings.Split(ctx.String("extensions"), ",")
			totalLen := len(extensions) - 1
			for ind, ext := range extensions {
				fmt.Fprintf(&extBu, "%q", strings.TrimSpace(ext))
				if ind < totalLen {
					fmt.Fprintf(&extBu, ",")
				}
			}

			fmt.Fprintf(&extBu, "}")

			if vDirName == "" && args.Len() > 0 {
				vDirName = args.First()
			}

			if vDirName == "" && args.Len() == 0 {
				vDirName = "templates"
			}

			vDirFileName := strings.ToLower(vDirName) + ".go"

			gopath := os.Getenv("GOPATH")
			gup := filepath.Join(gopath, "src")
			gupkg := filepath.Join(gup, gupath)
			vDirPath := filepath.Join(cdir, vDirName)
			mDirPath := filepath.Join(cdir, vDirName, mDirName)

			if err = os.MkdirAll(vDirPath, 0777); err != nil {
				return err
			}

			if err = os.MkdirAll(mDirPath, 0777); err != nil {
				return err
			}

			fmt.Printf("- Adding project directory: %q\n", filepath.Base(vDirPath))

			gendata, err := ioutil.ReadFile(filepath.Join(gupkg, "templates/views.template"))
			if err != nil {
				return err
			}

			vgendata, err := ioutil.ReadFile(filepath.Join(gupkg, "templates/files.template"))
			if err != nil {
				return err
			}

			plainPKGData, err := ioutil.ReadFile(filepath.Join(gupkg, "templates/plain_generated_pkg.template"))
			if err != nil {
				return err
			}

			gendata = []byte(fmt.Sprintf("%q", gendata))
			vgendata = bytes.Replace(vgendata, extbytes, extBu.Bytes(), 1)
			vgendata = bytes.Replace(vgendata, pkgContentbytes, gendata, 1)
			vgendata = bytes.Replace(vgendata, filesDirNamebytes, []byte(mDirName), 1)

			vgendata = bytes.Replace(vgendata, pkgNamebytes, []byte("\""+vDirName+"\""), 1)
			vgendata = bytes.Replace(vgendata, fileNamebytes, []byte("\""+vDirFileName+"\""), 1)
			plainPKGData = bytes.Replace(plainPKGData, pkgNamebytes, []byte(vDirName), -1)

			if err := writeFile(filepath.Join(vDirPath, "generate.go"), vgendata); err != nil {
				return err
			}

			fmt.Printf("- Adding project file: %q\n", filepath.Join(filepath.Base(vDirPath), "generate.go"))

			if err := writeFile(filepath.Join(vDirPath, vDirFileName), plainPKGData); err != nil {
				return err
			}

			fmt.Printf("- Adding project file: %q\n", filepath.Join(filepath.Base(vDirPath), vDirFileName))

			return nil
		},
	})

	commands = append(commands, &cli.Command{
		Name:        "templates",
		Usage:       "gu templates --dir=layouts --name=mytemplates",
		Description: "Generates a package to builds internal files [.html|.xhtml|.xml|.gml|.ghtml|.tml] as a go file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "name=mytemplates",
			},
			&cli.StringFlag{
				Name:    "dirName",
				Aliases: []string{"dir"},
				Usage:   "dir=assets",
			},
		},
		Action: func(ctx *cli.Context) error {
			mDirName := ctx.String("dirName")
			vDirName := ctx.String("name")
			args := ctx.Args()
			if args.Len() == 0 && vDirName == "" {
				return nil
			}

			cdir, err := os.Getwd()
			if err != nil {
				return err
			}

			if mDirName == "" {
				mDirName = "layouts"
			}

			if vDirName == "" && args.Len() > 0 {
				vDirName = args.First()
			}

			if vDirName == "" && args.Len() == 0 {
				vDirName = "templates"
			}

			vDirFileName := strings.ToLower(vDirName) + ".go"

			gopath := os.Getenv("GOPATH")
			gup := filepath.Join(gopath, "src")
			gupkg := filepath.Join(gup, gupath)
			vDirPath := filepath.Join(cdir, vDirName)
			mDirPath := filepath.Join(cdir, vDirName, mDirName)

			if err = os.MkdirAll(vDirPath, 0777); err != nil {
				return err
			}

			if err = os.MkdirAll(mDirPath, 0777); err != nil {
				return err
			}

			fmt.Printf("- Adding project directory: %q\n", filepath.Base(vDirPath))

			gendata, err := ioutil.ReadFile(filepath.Join(gupkg, "templates/views.template"))
			if err != nil {
				return err
			}

			vgendata, err := ioutil.ReadFile(filepath.Join(gupkg, "templates/viewsgenerate.template"))
			if err != nil {
				return err
			}

			plainPKGData, err := ioutil.ReadFile(filepath.Join(gupkg, "templates/plain_generated_pkg.template"))
			if err != nil {
				return err
			}

			gendata = []byte(fmt.Sprintf("%q", gendata))
			vgendata = bytes.Replace(vgendata, pkgContentbytes, gendata, 1)
			vgendata = bytes.Replace(vgendata, filesDirNamebytes, []byte(mDirName), 1)

			vgendata = bytes.Replace(vgendata, pkgNamebytes, []byte("\""+vDirName+"\""), 1)
			vgendata = bytes.Replace(vgendata, fileNamebytes, []byte("\""+vDirFileName+"\""), 1)
			plainPKGData = bytes.Replace(plainPKGData, pkgNamebytes, []byte(vDirName), -1)

			if err := writeFile(filepath.Join(vDirPath, "generate.go"), vgendata); err != nil {
				return err
			}

			fmt.Printf("- Adding project file: %q\n", filepath.Join(filepath.Base(vDirPath), "generate.go"))

			if err := writeFile(filepath.Join(vDirPath, vDirFileName), plainPKGData); err != nil {
				return err
			}

			fmt.Printf("- Adding project file: %q\n", filepath.Join(filepath.Base(vDirPath), vDirFileName))

			return nil
		},
	})

	commands = append(commands, &cli.Command{
		Name:        "components",
		Usage:       "gu components <sub-comand>",
		Description: "This provides subcommands which are used in the development of components",
		Subcommands: subcommands,
	})

	commands = append(commands, &cli.Command{
		Name:        "new",
		Usage:       "gu new <PackageName>",
		Description: "Generates a new gu component package for a gu app",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "driver",
				Aliases: []string{"dri"},
				Usage:   "driver=js|nodriver",
				Value:   "js",
			},
			&cli.StringFlag{
				Name:  "dir",
				Usage: "dir=path-to-dir",
			},
			&cli.StringFlag{
				Name:    "packageName",
				Aliases: []string{"pkg"},
				Usage:   "pkg=hello",
			},
		},
		Action: func(ctx *cli.Context) error {
			args := ctx.Args()
			if args.Len() == 0 {
				return nil
			}

			cdir, err := os.Getwd()
			if err != nil {
				return err
			}

			gopath := os.Getenv("GOPATH")
			gupsrc := filepath.Join(gopath, "src")
			gup := filepath.Join(gupsrc, gupath)

			// packagePath, err := filepath.Rel(gup, cdir)
			// if err != nil {
			// 	return err
			// }

			indir := ctx.String("dir")

			if indir != "" {
				if strings.HasPrefix(indir, ".") || !strings.HasPrefix(indir, "/") {
					indir = filepath.Join(cdir, indir)
				}
			} else {
				indir = cdir
			}

			packageName := ctx.String("packageName")

			if packageName == "" && args.Len() > 0 {
				packageName = args.First()
			}

			driver := ctx.String("driver")
			appDir := filepath.Join(indir, packageName)

			// appPackagePath := filepath.Join(packagePath, packageName)

			appPackagePath, err := filepath.Rel(gupsrc, appDir)
			if err != nil {
				return err
			}

			fmt.Printf("- Creating new project: %q\n", packageName)
			fmt.Printf("- Using driver template: %q\n", driver)

			// Generate dirs for the project.
			if err = os.MkdirAll(appDir, 0777); err != nil && err != os.ErrExist {
				return err
			}

			fmt.Printf("- Creating project directory: %q\n", packageName)

			if err = os.MkdirAll(filepath.Join(appDir, "components"), 0777); err != nil && err != os.ErrExist {
				return err
			}

			fmt.Printf("- Creating project directory: %q\n", filepath.Join(packageName, "components"))

			if err = os.MkdirAll(filepath.Join(appDir, "assets"), 0777); err != nil && err != os.ErrExist {
				return err
			}

			fmt.Printf("- Creating project directory: %q\n", filepath.Join(packageName, "assets"))

			registrydata, rerr := ioutil.ReadFile(filepath.Join(gup, "templates/registry.template"))
			if rerr != nil {
				return rerr
			}

			if err = writeFile(filepath.Join(indir, packageName, "components/components.go"), registrydata); err != nil {
				return err
			}

			fmt.Printf("- Adding project file: %q\n", "components/components.go")

			// Generate files for the project.
			switch driver {
			case "package", "nomain", "empty", "plain":
				appdata, err := ioutil.ReadFile(filepath.Join(gup, "templates/app_empty.template"))
				if err != nil {
					return err
				}

				appdata = bytes.Replace(appdata, goPathbytes, []byte(gopath), -1)
				appdata = bytes.Replace(appdata, pkgbytes, []byte(appPackagePath), -1)
				appdata = bytes.Replace(appdata, pkgNamebytes, []byte(packageName), -1)
				appdata = bytes.Replace(appdata, namebytes, []byte(packageName), -1)

				if err := writeFile(filepath.Join(indir, packageName, "app.go"), appdata); err != nil {
					return err
				}

				fmt.Printf("- Adding project file: %q\n", "app.go")
			case "nodriver", "no-driver":
				appdata, err := ioutil.ReadFile(filepath.Join(gup, "templates/app_nodriver.template"))
				if err != nil {
					return err
				}

				appdata = bytes.Replace(appdata, goPathbytes, []byte(gopath), -1)
				appdata = bytes.Replace(appdata, pkgbytes, []byte(appPackagePath), -1)
				appdata = bytes.Replace(appdata, pkgNamebytes, []byte(packageName), -1)
				appdata = bytes.Replace(appdata, namebytes, []byte(packageName), -1)

				if err := writeFile(filepath.Join(indir, packageName, "app.go"), appdata); err != nil {
					return err
				}

				fmt.Printf("- Adding project file: %q\n", "app.go")

			case "js":
				appdata, err := ioutil.ReadFile(filepath.Join(gup, "templates/app_js.template"))
				if err != nil {
					return err
				}

				apphtmldata, err := ioutil.ReadFile(filepath.Join(gup, "templates/app_js_html.template"))
				if err != nil {
					return err
				}

				appdata = bytes.Replace(appdata, goPathbytes, []byte(gopath), -1)
				appdata = bytes.Replace(appdata, pkgbytes, []byte(appPackagePath), -1)
				appdata = bytes.Replace(appdata, namebytes, []byte(packageName), -1)
				apphtmldata = bytes.Replace(apphtmldata, namebytes, []byte(packageName), -1)

				if err := writeFile(filepath.Join(indir, packageName, "app.go"), appdata); err != nil {
					return err
				}

				fmt.Printf("- Adding project file: %q\n", "app.go")

				if err := writeFile(filepath.Join(indir, packageName, "index.html"), apphtmldata); err != nil {
					return err
				}

				fmt.Printf("- Adding project file: %q\n", "index.html")
			}

			// Change to new app directory.
			if err := os.Chdir(filepath.Join(indir, packageName)); err != nil {
				return nil
			}

			return nil
		},
	})

	commands = append(commands, &cli.Command{
		Name:        "generate",
		Usage:       "gu generate",
		Description: "Generate will call needed code generators to create project assets and files as declared by the project and it's sources",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"dir"},
				Usage:   "dir=./my-gu-project",
			},
		},
		Action: func(ctx *cli.Context) error {
			args := ctx.Args()
			if args.Len() == 0 {
				return nil
			}

			indir := ctx.String("dir")

			if indir == "" {
				cdir, err := os.Getwd()
				if err != nil {
					return err
				}

				indir = cdir
			}

			register := ast.NewAnnotationRegistry()

			generators.RegisterGenerators(register)

			// Register @assets annotation for our registery as well.
			register.Register("assets", annotations.AssetsAnnotationGenerator)

			pkgs, err := ast.ParseAnnotations(events, indir)
			if err != nil {
				return err
			}

			if err := ast.Parse(events, register, pkgs...); err != nil {
				return err
			}

			return nil
		},
	})

}

// FindLowerByStat searches the path line down until it's roots to find the directory with the giving
// dirName matching else returns an error.
func findLowerByStat(root string, path string, dirName string, dirOnly bool) (string, error) {
	path = filepath.Clean(path)

	if path == "." {
		return "", errors.New("'" + dirName + "' path not found")
	}

	// Let's attempt to see if there is a dirName in this path and if it's a
	// directory.
	possiblePath := filepath.Join(root, path, dirName)
	possibleStat, err := os.Stat(possiblePath)
	if err == nil {
		if dirOnly && !possibleStat.IsDir() {
			return findLower(filepath.Join(path, ".."), dirName)
		}

		return filepath.Join(path, dirName), nil
	}

	return findLowerByStat(root, filepath.Join(path, ".."), dirName, dirOnly)
}

// Searches the path line down until it's roots to find the directory with the giving
// dirName matching else returns an error.
func findLower(path string, dirName string) (string, error) {
	path = filepath.Clean(path)

	if path == "." {
		return "", errors.New("'" + dirName + "' path not found")
	}

	if filepath.Base(path) == dirName {
		return path, nil
	}

	return findLower(filepath.Join(path, ".."), dirName)
}

func writeFile(targetFile string, data []byte) error {
	file, err := os.Create(targetFile)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}
