package packers

import "github.com/gu-io/gu/assets"

// StaticMarkupPacker defines a struct which implements the assets.Packer interface
// and will convert all .static files into go files with the file html content
// turned into type-safe trees.Markup structures(see github.com/gu-io/gu/tree/master/trees).
type StaticMarkupPacker struct {
	PackageName string
}

// Pack process all '.static.html' files present in the FileStatment slice and returns WriteDirectives
// which conta ins expected outputs for these files.
func (static StaticMarkupPacker) Pack(statements []assets.FileStatement, dir assets.DirStatement) ([]assets.WriteDirective, error) {
	var directives []assets.WriteDirective

	return directives, nil
}
