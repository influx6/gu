package styleguide_test

import (
	"testing"

	"github.com/gu-io/gu/trees/themes/styleguide"
	"github.com/influx6/faux/tests"
)

func TestColor(t *testing.T) {
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

// func TestColorScale(t *testing.T) {
// 	color, err := styleguide.ColorFrom("#f4f4f4")
// 	if err != nil {
// 		tests.Failed("Should have successfully returned hsl value for hashed color: %+q.", err)
// 	}
// 	tests.Passed("Should have successfully returned hsl value for hashed color.")

// 	tones := styleguide.HamonicsFrom(color)
// 	fmt.Printf("Tone: %s -> %d\n", tones, len(tones.Grades))
// }

// func TestStyleGuide(t *testing.T) {
// 	sl, err := styleguide.New(styleguide.Attr{
// 		PrimaryColor:        "#000000",
// 		SecondaryColor:      "#000000",
// 		PrimaryBrandColor:   "#000000",
// 		SecondaryBrandColor: "#000000",
// 		PrimaryWhite:        "#ffffff",
// 		FailureColor:        "#000000",
// 		SuccessColor:        "#000000",
// 	})
// 	if err != nil {
// 		tests.Failed("Should have successfully returned new styleguide: %+q.", err)
// 	}
// 	tests.Passed("Should have successfully returned new styleguide.")

// 	fmt.Println(sl.CSS())
// }
