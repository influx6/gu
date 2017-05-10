package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gu-io/gu/shell"
	"github.com/gu-io/gu/shell/parse"
	cli "gopkg.in/urfave/cli.v2"
)

var (
	version     = "0.0.1"
	defaultName = "manifests"
	commands    = []*cli.Command{}

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

	gupath = "github.com/gu-io/gu"

	usage = `Provides a CLi tool which allows deployment and generation of project files for use in development.`

	aferoTemplate = `// Package %s is auto-generated and should not be modified by hand.
// This package contains a virtual file system for generate resources which are not accessed
// through a remote endpoint (i.e those resources generated from the manifests that are local in the
// filesystem and are not marked as remote in access).
package %s

import (
	"path/filepath"

	"github.com/spf13/afero"
)

// AppFS defines the global handler for which all assets generated from manifests
// files which are not remote resources are provided as binary embedded assets.
var AppFS = afero.NewMemMapFs()

// addFile adds a giving file name into the file system.
func addFile(path string, content []byte){
	dir, _ := filepath.Split(path)
	if dir != "" {
		AppFS.MkdirAll(dir,0755)
	}

	afero.WriteFile(AppFS, path, content, 0644)
}

func init(){
%+s
}

`
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
		Description: `Generates a new boiler code component package or file, which can be set to be in it's 
		own package or part of the current directory.

		Options:
			- flat: This option when true, will indicate that only a .go file of that component is to be generated in the app's components package.

			- base:	This option when false, will force that component file or package to be generated right in the directory where the command was called and not in the components package.
		`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "base-of-components-package",
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
		Usage:       "gu templates --dirName layouts templates",
		Description: "Generates a package which builds all internal [.html|.xhtml|.xml|.gml|.ghtml|.tml] files into a go file",
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

			manifestDirPath := filepath.Join(appDir, "manifests")
			if err = os.MkdirAll(manifestDirPath, 0777); err != nil {
				return err
			}

			fmt.Printf("- Adding project directory: %q\n", filepath.Base(manifestDirPath))

			manifestGendata, err := ioutil.ReadFile(filepath.Join(gup, "templates/manifest-generate.template"))
			if err != nil {
				return err
			}

			manifestSource, err := ioutil.ReadFile(filepath.Join(gup, "templates/manifests.template"))
			if err != nil {
				return err
			}

			manifestSource = []byte(fmt.Sprintf("%q", manifestSource))
			manifestGendata = bytes.Replace(manifestGendata, sourcebytes, manifestSource, 1)

			if err := writeFile(filepath.Join(manifestDirPath, "generate.go"), manifestGendata); err != nil {
				return err
			}

			fmt.Printf("- Adding project file: %q\n", filepath.Join(filepath.Base(manifestDirPath), "generate.go"))

			plainPKGData, err := ioutil.ReadFile(filepath.Join(gup, "templates/plain_generated_pkg.template"))
			if err != nil {
				return err
			}

			plainPKGData = bytes.Replace(plainPKGData, pkgNamebytes, []byte(filepath.Base(manifestDirPath)), -1)

			if err := writeFile(filepath.Join(manifestDirPath, "manifests.go"), plainPKGData); err != nil {
				return err
			}

			fmt.Printf("- Adding project file: %q\n", filepath.Join(filepath.Base(manifestDirPath), "manifests.go"))

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

			case "osx":
				// read the full qt template and write into the file.
				data, err := ioutil.ReadFile(filepath.Join(gup, "templates/app_osx.template"))
				if err != nil {
					return err
				}

				data = bytes.Replace(data, namebytes, []byte(packageName), -1)
				data = bytes.Replace(data, pkgbytes, []byte(appPackagePath), -1)

				if err := writeFile(filepath.Join(indir, packageName, "app.go"), data); err != nil {
					return err
				}

				fmt.Printf("- Adding project file: %q\n", "app.go")

			case "win", "edge":
				// read the full qt template and write into the file.
				data, err := ioutil.ReadFile(filepath.Join(gup, "templates/app_win.template"))
				if err != nil {
					return err
				}

				data = bytes.Replace(data, namebytes, []byte(packageName), -1)
				data = bytes.Replace(data, pkgbytes, []byte(appPackagePath), -1)

				if err := writeFile(filepath.Join(indir, packageName, "app.go"), data); err != nil {
					return err
				}

				fmt.Printf("- Adding project file: %q\n", "app.go")

			case "linux", "gtk":
				// read the full qt template and write into the file.
				data, err := ioutil.ReadFile(filepath.Join(gup, "templates/app_gtk.template"))
				if err != nil {
					return err
				}

				data = bytes.Replace(data, namebytes, []byte(packageName), -1)
				data = bytes.Replace(data, pkgbytes, []byte(appPackagePath), -1)

				if err := writeFile(filepath.Join(indir, packageName, "app.go"), data); err != nil {
					return err
				}

				fmt.Printf("- Adding project file: %q\n", "app.go")

			case "qt":
				// read the full qt template and write into the file.
				data, err := ioutil.ReadFile(filepath.Join(gup, "templates/app_qt.template"))
				if err != nil {
					return err
				}

				data = bytes.Replace(data, namebytes, []byte(packageName), -1)
				data = bytes.Replace(data, pkgbytes, []byte(appPackagePath), -1)

				if err := writeFile(filepath.Join(indir, packageName, "app.go"), data); err != nil {
					return err
				}

				fmt.Printf("- Adding project file: %q\n", "app.go")
			}

			// Change to new app directory.
			if err := os.Chdir(filepath.Join(indir, packageName)); err != nil {
				return nil
			}

			return nil
		},
	})

	commands = append(commands, &cli.Command{
		Name:        "manifests",
		Usage:       "gu manifests",
		Description: "Generate a manifest.json file that contains all resources from meta-comments within the package to be embedded",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "input-dir",
				Aliases: []string{"indir"},
				Usage:   "in-dir=path-to-dir-to-scan",
			},
			&cli.StringFlag{
				Name:    "output-dir",
				Aliases: []string{"outdir"},
				Usage:   "out-dir=path-to-store-manifest-file",
			},
		},
		Action: func(ctx *cli.Context) error {
			cdir, err := os.Getwd()
			if err != nil {
				return err
			}

			indir := ctx.String("input-dir")
			outdir := ctx.String("output-dir")

			if indir != "" {
				if strings.HasPrefix(indir, ".") || !strings.HasPrefix(indir, "/") {
					indir = filepath.Join(cdir, indir)
				}
			} else {
				indir = cdir
			}

			if outdir != "" {
				if strings.HasPrefix(outdir, ".") || !strings.HasPrefix(outdir, "/") {
					outdir = filepath.Join(cdir, outdir)
				}
			} else {
				outdir = cdir
			}

			res, err := parse.ShellResources(indir)
			if err != nil {
				return err
			}

			var manifests []*shell.AppManifest

			for _, rs := range res {
				manifest, merr := rs.GenManifests()
				if merr != nil {
					return merr
				}

				manifests = append(manifests, manifest)
			}

			manifestJSON, err := json.MarshalIndent(manifests, "", "\t")
			if err != nil {
				return err
			}

			if bytes.Equal(manifestJSON, []byte("null")) {
				manifestJSON = []byte("{}")
			}

			manifestFile, err := os.Create(filepath.Join(outdir, "manifest.json"))
			if err != nil {
				return err
			}

			defer manifestFile.Close()

			total, err := manifestFile.Write(manifestJSON)
			if err != nil {
				return err
			}

			if total != len(manifestJSON) {
				return errors.New("Data written is incomplete")
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
