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

	namebytes       = []byte("{{Name}}")
	pkgbytes        = []byte("{{PKG}}")
	sourcebytes     = []byte("{{SOURCE}}")
	goPathbytes     = []byte("{{GOPATH}}")
	pkgContentbytes = []byte("{{PKG_CONTENT}}")
	pkgNamebytes    = []byte("{{PKGNAME}}")
	dirNamebytes    = []byte("{{DIRNAME}}")
	nameLowerbytes  = []byte("{{Name_Lower}}")

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

			gendata = []byte(fmt.Sprintf("%q", gendata))
			cssgendata = bytes.Replace(cssgendata, pkgContentbytes, gendata, 1)
			cssgendata = bytes.Replace(cssgendata, dirNamebytes, []byte("css"), 1)
			cssgendata = bytes.Replace(cssgendata, pkgNamebytes, []byte("\""+cssDirName+"\""), 1)

			if err := writeFile(filepath.Join(cssDirPath, "generate.go"), cssgendata); err != nil {
				return err
			}

			fmt.Printf("- Adding project file: %q\n", filepath.Join(filepath.Base(cssDirPath), "generate.go"))

			return nil
		},
	})

	subcommands = append(subcommands, &cli.Command{
		Name:        "new",
		Usage:       "gu new <component-name>",
		Description: "Generates a new boiler code component file which can be set to be in it's own package or part of the component package ",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "base",
				Aliases: []string{"bs"},
				Usage:   "base=true",
				Value:   true,
			},
			&cli.BoolFlag{
				Name:    "flat",
				Aliases: []string{"fl"},
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
			args := ctx.Args()
			if args.Len() == 0 {
				return nil
			}

			cdir, err := os.Getwd()
			if err != nil {
				return err
			}

			flat := ctx.Bool("flat")
			base := ctx.Bool("base")

			gopath := os.Getenv("GOPATH")
			gup := filepath.Join(gopath, "src")
			gupkg := filepath.Join(gopath, "src", gupath)
			componentDir := filepath.Join(cdir, "components")

			componentName := ctx.String("component")

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
			newComponentDir := filepath.Join(componentDir, componentPkgName)

			cssDirName := "styles"
			newComponentCSSDir := filepath.Join(newComponentDir, cssDirName)
			newComponentCSSFilesDir := filepath.Join(newComponentCSSDir, "css")

			packagePath, err := filepath.Rel(gup, cdir)
			if err != nil {
				return err
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

				componentsPackagePath, coerr := findLower(packagePath, "components")
				if coerr != nil {
					return coerr
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

			cssbeforegendata, cerr := ioutil.ReadFile(filepath.Join(gupkg, "templates/css.template"))
			if cerr != nil {
				return cerr
			}

			cssgendata, merr := ioutil.ReadFile(filepath.Join(gupkg, "templates/cssgenerate.template"))
			if merr != nil {
				return merr
			}

			cssbeforegendata = []byte(fmt.Sprintf("%q", cssbeforegendata))
			cssgendata = bytes.Replace(cssgendata, pkgContentbytes, cssbeforegendata, 1)
			cssgendata = bytes.Replace(cssgendata, dirNamebytes, []byte("css"), 1)
			cssgendata = bytes.Replace(cssgendata, pkgNamebytes, []byte("\""+cssDirName+"\""), 1)

			if err = writeFile(filepath.Join(newComponentCSSDir, "generate.go"), cssgendata); err != nil {
				return err
			}

			fmt.Printf("- Adding project file: %q\n", filepath.Join("components", componentPkgName, "styles", "generate.go"))

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
			if err = os.Mkdir(manifestDirPath, 0777); err != nil {
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

			// Generate files for the project.
			switch driver {
			case "nodriver":
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
