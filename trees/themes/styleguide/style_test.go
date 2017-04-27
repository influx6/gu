package styleguide_test

import (
	"fmt"
	"testing"

	"github.com/gu-io/gu/trees/themes/styleguide"
	"github.com/influx6/faux/tests"
)

func TestColor(t *testing.T) {
	style, _ := styleguide.NewStyleGuide(styleguide.Attr{
		PrimaryBrandColor:   "#7fffd4",
		SecondaryBrandColor: "#7fffd4",
		SuccessColor:        "#7fffd4",
		FailureColor:        "#7fffd4",
	})

	fmt.Printf("%+s\n", style.CSS())

	color, err := styleguide.ColorFrom("#7fffd4")
	if err != nil {
		tests.Failed("Should have successfully returned hsl value for hashed color: %+q.", err)
	}
	tests.Passed("Should have successfully returned hsl value for hashed color.")

	if color.C.Hex() != "#7fffd4" {
		tests.Failed("Should have successfully matched hex value.")
	}
	tests.Passed("Should have successfully matched hex value.")

	hslColor, err := styleguide.ColorFrom("hsl(160, 100%, 75%)")
	if err != nil {
		tests.Failed("Should have successfully returned hsl value for hashed color: %+q.", err)
	}
	tests.Passed("Should have successfully returned hsl value for hashed color.")

	if hslColor.C.DistanceRgb(color.C) > 0.1 {
		tests.Failed("Should have successfully matched hex value.")
	}
	tests.Passed("Should have successfully matched hex value.")
}
