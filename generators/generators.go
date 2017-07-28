package generators

import "github.com/influx6/moz/ast"

// RegisterGenerators will add all generator functions from this package
// into the provided registry.
func RegisterGenerators(house *ast.AnnotationRegistry) {
	house.Register("notification:event", NotificationTypeGenerator)
}
