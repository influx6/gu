package cssparser_test

import (
	"testing"

	"github.com/gu-io/gu/parsers/cssparser/examples"
	"github.com/influx6/faux/tests"
)

func TestCSSFiles(t *testing.T) {
	if exm := examples.Get("examples.css"); exm == nil {
		tests.Failed("Should have successfully retrieved 'examples.css' file")
	}
	tests.Passed("Should have successfully retrieved 'examples.css' file")

	if base := examples.Get("base/base.css"); base == nil {
		tests.Failed("Should have successfully retrieved 'base/base.css' file")
	}
	tests.Passed("Should have successfully retrieved 'base/base.css' file")

	if baseui := examples.Get("base/ui/base-ui.css"); baseui == nil {
		tests.Failed("Should have successfully retrieved 'base/ui/base-ui.css' file")
	}
	tests.Passed("Should have successfully retrieved 'base/ui/base-ui.css' file")

	if ui := examples.Get("ui/ui.css"); ui == nil {
		tests.Failed("Should have successfully retrieved 'ui/ui.css' file")
	}
	tests.Passed("Should have successfully retrieved 'ui/ui.css' file")
}
