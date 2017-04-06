package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gu-io/gu/shell"
)

var (
	goPath    = os.Getenv("GOPATH")
	goSrcPath = filepath.Join(goPath, "src")
)

// Resource defines a giving component with specific resource metas written for
// it's type. It holds the lists of Resource meta gathered and parsed and contains
// detailed package and component data related to the resources and any other
// resources which connects with it.
type Resource struct {
	Global             bool
	File               string
	ComponentName      string
	PackageName        string
	Package            string
	CompositeTypeNames []string
	FieldTypeNames     []string
	Resources          []ResourceCollection
}

// GenManifests generates a manifest structure which mirrors with specific changes
// the and properly allocates the resources and attempts to load any local resource
// as a loaded resource.
func (r *Resource) GenManifests() (*shell.AppManifest, error) {
	appm := shell.NewAppManifest(fmt.Sprintf("%s.%s", r.PackageName, r.ComponentName))
	appm.GlobalScope = r.Global

	var relation shell.ComponentRelation
	relation.Name = r.ComponentName
	relation.Package = r.Package
	relation.FieldTypes = append(relation.FieldTypes, r.FieldTypeNames...)
	relation.Composites = append(relation.Composites, r.CompositeTypeNames...)

	appm.Relation = &relation

	for _, res := range r.Resources {
		relation.Composites = append(relation.Composites, res.Relations...)

		manifest, err := res.GenManifestAttr(r.Package)
		if err != nil && err != ErrNotLocalPath {
			return nil, err
		}

		if err == ErrNotLocalPath {
			manifest.Remote = true
		}

		appm.Manifests = append(appm.Manifests, manifest)
	}

	return appm, nil
}

// ResourceCollection holds the meta-information related to different resources
// gathered from a package list.
type ResourceCollection struct {
	Name          string
	Localize      bool
	Init          bool
	Global        bool
	Remote        bool
	Encode        bool
	Encoded       bool
	Base64Padding bool
	Path          string
	Cache         string
	HookName      string
	ContentSize   string
	Data          string
	Meta          map[string]string
	Relations     []string
}

// GenManifestAttr returns a new instance of a ManifestAttr generated from the giving
// ResourceCollection.
func (r *ResourceCollection) GenManifestAttr(pkg string) (shell.ManifestAttr, error) {
	var mattr shell.ManifestAttr

	mattr.Init = r.Init
	mattr.Name = r.Name
	mattr.Meta = r.Meta
	mattr.Content = r.Data
	mattr.Remote = r.Remote
	mattr.IsGlobal = r.Global
	mattr.B64Encode = r.Encode
	mattr.HookName = r.HookName
	mattr.Localize = r.Localize
	mattr.ContentB64 = r.Encoded
	mattr.Base64Padding = r.Base64Padding

	if size, err := strconv.Atoi(r.ContentSize); err == nil {
		mattr.Size = size
	}

	mattr.Path = strings.TrimSpace(r.Path)

	if mattr.Path != "" && !r.Remote && mattr.Content == "" {
		fmt.Printf("Retrieving File Content for %q: Base64: %t Base64Padding: %t\n", mattr.Path, mattr.B64Encode, mattr.Base64Padding)
		content, err := getFileContent(pkg, mattr.Path, mattr.B64Encode, mattr.Base64Padding)
		if err != nil {
			return mattr, err
		}

		if mattr.B64Encode {
			mattr.ContentB64 = true
		}

		mattr.Size = len(content)
		mattr.Content = string(content)
	}

	if mattr.Path != "" && r.Remote && mattr.Content == "" && r.Localize {
		fmt.Printf("Retrieving URL Content for %q: Base64: %t Base64Padding: %t\n", mattr.Path, mattr.B64Encode, mattr.Base64Padding)
		content, err := getURLContent(mattr.Path, mattr.B64Encode, mattr.Base64Padding)
		if err != nil {
			return mattr, err
		}

		if mattr.B64Encode {
			mattr.ContentB64 = true
		}

		mattr.Size = len(content)
		mattr.Content = string(content)
	}

	return mattr, nil
}

//==============================================================================

var (
	exceptions = []string{
		"github.com/gu-io/gu",
		"github.com/gu-io/gu/shell",
	}
)

// ShellResources parses a directoring creating a slice of Resources from
// all meta comments seen.
func ShellResources(dir string) ([]Resource, error) {
	var resources []Resource

	_, nodes, err := PackageDir(dir, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		for path, files := range node.Files {
			rel, _ := filepath.Rel(goSrcPath, path)
			rel = filepath.Dir(rel)

		importLoop:
			for _, imported := range files.Imports {
				val, _ := strconv.Unquote(imported.Path.Value)

				var skipImport bool

				{
				skipLog:
					for _, except := range exceptions {
						if except == val {
							skipImport = true
							break skipLog
						}

						urel, err := filepath.Rel(except, val)
						if err != nil {
							skipImport = true
							break skipLog
						}

						if strings.HasPrefix(urel, "..") {
							skipImport = true
							break skipLog
						}

						firstPath := strings.Replace(val, urel, "", 1)
						firstPath = strings.TrimSuffix(firstPath, "/")
						except = strings.TrimSuffix(except, "/")

						if firstPath == except {
							skipImport = true
							break skipLog
						}
					}
				}

				if skipImport {
					continue
				}

				pkgPath := filepath.Join(goSrcPath, val)

				if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
					continue importLoop
				}

				res, err := ShellResources(pkgPath)
				if err != nil {
					return nil, err
				}

				resources = append(resources, res...)
			}

			if strings.Contains(files.Doc.Text(), shellGlobalMarker) {
				resourceMaps, err := GrabShellDoc(files.Doc.Text())
				if err != nil {
					return nil, err
				}

				resd, err := toResources(resourceMaps)
				if err != nil {
					return nil, err
				}

				// Collect main fields from the maps into the resource
				// variable.
				var res Resource
				res.File = path
				res.Global = true
				res.Package = rel
				res.Resources = resd
				res.PackageName = node.Name
				res.ComponentName = files.Name.Name
				resources = append(resources, res)
			}

			for _, decl := range files.Decls {

				gdecl, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				if !strings.Contains(gdecl.Doc.Text(), shellComponentMarker) {
					continue
				}

				resourceMaps, err := GrabShellDoc(gdecl.Doc.Text())
				if err != nil {
					return nil, err
				}

				resd, err := toResources(resourceMaps)
				if err != nil {
					return nil, err
				}

				// Collect main fields from the maps into the resource
				// variable.
				var res Resource
				res.File = path
				res.Package = rel
				res.Resources = resd
				res.PackageName = node.Name

			innLoop:
				for _, spec := range gdecl.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue innLoop
					}

					res.ComponentName = ts.Name.Name

					switch tso := ts.Type.(type) {
					case *ast.StructType:

						{
						fieldLoop:
							for _, fld := range tso.Fields.List {
								if selectExpr, ok := fld.Type.(*ast.SelectorExpr); ok {
									res.CompositeTypeNames = append(res.CompositeTypeNames, selectExpr.Sel.Name)
									continue fieldLoop
								}

								if indentExpr, ok := fld.Type.(*ast.Ident); ok && indentExpr.Obj != nil {
									res.FieldTypeNames = append(res.FieldTypeNames, indentExpr.Name)
									continue fieldLoop
								}
							}
						}
					}

					break innLoop
				}

				resources = append(resources, res)
			}
		}
	}

	return resources, nil
}

// PackageDir turns a given go source file into a appropriate strucure which will be
// used to generate the needed manifests for a resource shell.
func PackageDir(file string, mode parser.Mode) (*token.FileSet, map[string]*ast.Package, error) {
	tokens := token.NewFileSet()
	nodes, err := parser.ParseDir(tokens, file, nil, mode)
	if err != nil {
		return nil, nil, err
	}

	return tokens, nodes, nil
}

// PackageFile turns a given go source file into a appropriate strucure which will be
// used to generate the needed manifests for a resource shell.
func PackageFile(file string, mode parser.Mode) (*token.FileSet, *ast.File, error) {
	tokens := token.NewFileSet()
	nodes, err := parser.ParseFile(tokens, file, nil, mode)
	if err != nil {
		return nil, nil, err
	}

	return tokens, nodes, nil
}
