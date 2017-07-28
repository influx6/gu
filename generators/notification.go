package generators

import (
	"fmt"
	"text/template"

	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
	"github.com/influx6/moz/gen/templates"
)

// NotificationTypeGenerator which defines a  function for generating a type for receiving a giving
//	struct type has a notification type which can then be wired as a notification.EventDistributor.
//
//	Usage:
//	We want users to be able to define a type within their source code where they can use an annotation to mark such
//	a type has a EventNotification type. More so, users will need a type that will cater to listening for such specific
//	struct type as an event, so they can register that type to listen specifically for such type.
//
//	Reason:
//	We want to remove the need for reflection but also provide flexibility and customization on the specifics of an event
//	users can get but also provide the low-risk approach of type assertion but that the user does not need to worry about.
//	If done this way we can get users to generate any event base type and get a handler to connect and handler type assertions
//	for that event without need to worry about that themselves.
//
//
func NotificationTypeGenerator(an ast.AnnotationDeclaration, str ast.StructDeclaration, pkg ast.PackageDeclaration) ([]gen.WriteDirective, error) {
	eventFileName := fmt.Sprintf("%s_event.go", str.Object.Name)

	typeGen := gen.Block(
		gen.Commentary(
			gen.Text(""),
		),
		gen.Package(
			gen.Name(pkg.Package),
			gen.Imports(),
			gen.Block(
				gen.SourceTextWith(
					string(templates.Must("notifications/eventtype.gen")),
					template.FuncMap{},
					struct {
						Struct  ast.StructDeclaration
						Package ast.PackageDeclaration
					}{
						Struct:  str,
						Package: pkg,
					},
				),
			),
		),
	)

	return []gen.WriteDirective{
		{
			Dir:          "./",
			DontOverride: true,
			Writer:       typeGen,
			FileName:     eventFileName,
		},
	}, nil
}
