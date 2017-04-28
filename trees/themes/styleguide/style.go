package styleguide

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/influx6/faux/colors"
	colorful "github.com/lucasb-eyer/go-colorful"
)

var (
	helpers = template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"multiply": func(a, b int) int {
			return a * b
		},
		"subtract": func(a, b int) int {
			return b - a
		},
	}
)

// Attr defines different color and size of strings to define a specific brand.
type Attr struct {
	PrimaryColor        string
	SecondaryColor      string
	PrimaryBrandColor   string
	SecondaryBrandColor string
	PrimaryWhite        string
	SuccessColor        string
	FailureColor        string

	BaseFontSize int // BaseFontSize for typeface using MajorThird.

	SmallBorderRadius  int // SmallBorderRadius for tiny components eg checkbox, radio buttons.
	MediumBorderRadius int // MediaBorderRadius for buttons, inputs, etc
	LargeBorderRadius  int // LargeBorderRadius for components like cards, modals, etc.

	// Optional shadow values, need not set
	FloatingShadow string // shadow for floating icons, elements.
	HoverShadow    string // shadow for over dialog etc
	DropShadow     string // Useful for popovers/dropovers
	BaseShadow     string // Normal shadow of elemnts
}

// StyleColors defines a struct which holds all possible giving brand colors utilized
// for the project.
type StyleColors struct {
	Primary        Tones `json:"primary"`
	Secondary      Tones `json:"secondary"`
	Success        Tones `json:"success"`
	Failure        Tones `json:"failure"`
	White          Tones `json:"white"`
	PrimaryBrand   Tones `json:"primary_support"`
	SecondaryBrand Tones `json:"secondary_support"`
}

// TypeSize defines a type for a float64 with a special
// string method to return it in 3 decimal places.
type TypeSize float64

// String returns the version of the giving TypeSize;
func (t TypeSize) String() string {
	return fmt.Sprintf("%.3f", t)
}

// StyleGuide represent the fullset of brand style properties for the project.
type StyleGuide struct {
	Attr
	Brand          StyleColors `json:"brand"`
	BigFontScale   []TypeSize
	SmallFontScale []TypeSize
}

// NewStyleGuide returns a new StyleGuide object which generates the necessary css
// styles to utilize the defined style within any project.
func NewStyleGuide(attr Attr) (StyleGuide, error) {
	var style StyleGuide

	var err error
	style.Brand.Primary, err = NewTones(attr.PrimaryColor)
	if err != nil {
		return style, err
	}

	style.Brand.PrimaryBrand, err = NewTones(attr.PrimaryBrandColor)
	if err != nil {
		return style, err
	}

	style.Brand.Secondary, err = NewTones(attr.SecondaryColor)
	if err != nil {
		return style, err
	}

	style.Brand.SecondaryBrand, err = NewTones(attr.SecondaryBrandColor)
	if err != nil {
		return style, err
	}

	style.Brand.White, err = NewTones(attr.PrimaryWhite)
	if err != nil {
		return style, err
	}

	style.Brand.Success, err = NewTones(attr.SuccessColor)
	if err != nil {
		return style, err
	}

	style.Brand.Failure, err = NewTones(attr.FailureColor)
	if err != nil {
		return style, err
	}

	style.init()
	return style, nil
}

const (
	shadowLarge    = "0px 20px 40px 4px rgba(0, 0, 0, 0.51)"
	shadowPopDrops = "0px 9px 30px 2px rgba(0, 0, 0, 0.51)"
	shadowHovers   = "0px 13px 30px 5px rgba(0, 0, 0, 0.58)"
	shadowNormal   = "0px 13px 20px 2px rgba(0, 0, 0, 0.45)"

	smallBorderRadius  = 2
	mediumBorderRadius = 4
	largeBorderRadius  = 8
)

// init initializes the style guide properties which require values.
func (sc *StyleGuide) init() {
	if sc.Attr.BaseFontSize <= 0 {
		sc.BaseFontSize = 16
	}

	if sc.FloatingShadow == "" {
		sc.FloatingShadow = shadowLarge
	}

	if sc.HoverShadow == "" {
		sc.HoverShadow = shadowHovers
	}

	if sc.BaseShadow == "" {
		sc.BaseShadow = shadowNormal
	}

	if sc.DropShadow == "" {
		sc.DropShadow = shadowPopDrops
	}

	if sc.SmallBorderRadius <= 0 {
		sc.SmallBorderRadius = smallBorderRadius
	}

	if sc.MediumBorderRadius <= 0 {
		sc.MediumBorderRadius = mediumBorderRadius
	}

	if sc.LargeBorderRadius <= 0 {
		sc.LargeBorderRadius = largeBorderRadius
	}

	bg, sm := GenerateValueScale(MajorThird, 1)

	for _, item := range bg {
		sc.BigFontScale = append(sc.BigFontScale, TypeSize(item))
	}

	for _, item := range sm {
		sc.SmallFontScale = append(sc.SmallFontScale, TypeSize(item))
	}
}

// CSS returns a css style content for usage with a css stylesheet.
func (sc *StyleGuide) CSS() string {
	tml, err := template.New("styleguide").Funcs(helpers).Parse(styleTemplate)
	if err != nil {
		return err.Error()
	}

	var buf bytes.Buffer
	if terr := tml.Execute(&buf, sc); terr != nil {
		return terr.Error()
	}

	return buf.String()
}

//================================================================================================

// Tones defines the set of color tones generated for a base color using the Hamonic tone
// sets, it provides a very easily set of color variations for use in styles.
type Tones struct {
	Base   Color   `json:"base"`
	Grades []Color `json:"tones"`
}

// NewTones returns a new Tones object representing the provided color tones if
// the value provided is a valid color.
func NewTones(base string) (Tones, error) {
	c, err := ColorFrom(base)
	if err != nil {
		return Tones{}, err
	}

	return HamonicsFrom(c), nil
}

// JSON returns the string representation of the provided tone.
func (t Tones) JSON() string {
	return fmt.Sprintf(`{
	"base": %q,
	"grades": %+q
  }`, t.Base, t.Grades)
}

// String returns the string representation of the provided tone.
func (t Tones) String() string {
	return fmt.Sprintf(`%q %q`, t.Base, t.Grades)
}

//================================================================================================

var (
	lnScale     = []float64{0.19, 0.25, 0.30}
	normalScale = []float64{0.09, 0.15, 0.30}
	midScale    = []float64{0.035, 0.13, 0.23}

	lowScale    = []float64{0.009, 0.016, 0.025, 0.035}
	darkerScale = []float64{-0.05, -0.70}

	satuScale = []float64{0, 0, 0.00, 0.00}
)

// HamonicsFrom uses the above scale to return a slice of new Colors based on the provided
// HamonyScale set.
func HamonicsFrom(c Color) Tones {
	var colors []Color

	var scale []float64

	switch {
	case c.Luminosity < 0.3:
		scale = lnScale
	case c.Luminosity > 0.9:
		scale = lowScale
	case c.Luminosity > 0.5:
		scale = midScale
	default:
		scale = normalScale
	}

	var darkers []Color
	for _, scale := range darkerScale {
		darkers = append(darkers, MultiplicativeLumination(c, scale))
	}

	for i := len(darkers) - 1; i > 0; i-- {
		colors = append(colors, darkers[i])
	}

	for index, scale := range scale {
		var next Color

		if len(colors) > 0 {
			next = colors[len(colors)-1]
		} else {
			next = c
		}

		saturateScale := satuScale[index]

		colors = append(colors, AdditiveSaturation(AdditiveLumination(next, scale), saturateScale))
	}

	var t Tones
	t.Base = c
	t.Grades = colors

	return t
}

// AdditiveSaturation adds the provided scale to the colors saturation value
// returning a new color suited to match.
func AdditiveSaturation(c Color, scale float64) Color {
	// fmt.Printf("AS: H: %.4f S: %.4f L: %.4f \n", c.Hue, c.Saturation, c.Luminosity)
	newLumen := c.Saturation + scale
	if newLumen > 1 {

		// Use the difference to reduce the lightness.
		diff := 1 - c.Luminosity
		if diff > 0 {
			newLum := c.Luminosity + (diff / 2)

			// fmt.Printf("Diff: %.4f : %.4f -> %.4f : %.4f \n", diff, diff/2, c.Luminosity, newLum)

			if newLum > 1 {
				newLumen = 0.999
			} else {
				newLumen = newLum
			}
		} else {
			newLumen = 0.999
		}
	}

	if newLumen < 0 {
		newLumen = 0
	}

	newColor := colorful.Hsl(c.Hue, newLumen, c.Luminosity)

	h, s, l := newColor.Hsl()

	return Color{
		C:          newColor,
		Hue:        h,
		Saturation: s,
		Luminosity: l,
		Alpha:      c.Alpha,
	}
}

// MultiplicativeSaturation multiples the scale to the colors saturation value
// using the returned value as a addition to the current saturation value,
// Creating a gradual change in saturation for the returned color.
func MultiplicativeSaturation(c Color, scale float64) Color {
	newLuma := (c.Saturation * scale)
	newLumen := c.Saturation + newLuma
	if newLumen > 1 {

		// Use the difference to reduce the lightness.
		diff := 1 - c.Luminosity
		if diff > 0 {
			newLum := c.Luminosity + (diff / 2)

			// fmt.Printf("Diff: %.4f : %.4f -> %.4f : %.4f \n", diff, diff/2, c.Luminosity, newLum)

			if newLum > 1 {
				newLumen = 0.999
			} else {
				newLumen = newLum
			}
		} else {
			newLumen = 0.999
		}
	}

	newColor := colorful.Hsl(c.Hue, c.Saturation, newLumen)
	h, s, l := newColor.Hsl()

	return Color{
		C:          newColor,
		Hue:        h,
		Saturation: s,
		Luminosity: l,
		Alpha:      c.Alpha,
	}
}

// AdditiveLumination adds the provided scale to the colors Luminouse value
// returning a new color suited to match.
func AdditiveLumination(c Color, scale float64) Color {
	// fmt.Printf("AL: H: %.4f S: %.4f L: %.4f  -> S: %.4f\n", c.Hue, c.Saturation, c.Luminosity, scale)

	newLumen := c.Luminosity + scale
	if newLumen > 1 {

		// Use the difference to reduce the lightness.
		diff := 1 - c.Luminosity
		if diff > 0 {
			newLum := c.Luminosity + (diff / 2)

			// fmt.Printf("Diff: %.4f : %.4f -> %.4f : %.4f \n", diff, diff/2, c.Luminosity, newLum)

			if newLum > 1 {
				newLumen = 0.999
			} else {
				newLumen = newLum
			}
		} else {
			newLumen = 0.999
		}
	}

	if newLumen < 0 {
		newLumen = 0
	}

	newColor := colorful.Hsl(c.Hue, c.Saturation, newLumen)

	h, s, l := newColor.Hsl()

	return Color{
		C:          newColor,
		Hue:        h,
		Saturation: s,
		Luminosity: l,
		Alpha:      c.Alpha,
	}
}

// MultiplicativeLumination multiples the scale to the colors Luminouse value
// using the returned value as addition to the current Luminouse value.
// Creating a gradual change in luminousity for the returned color.
func MultiplicativeLumination(c Color, scale float64) Color {
	// fmt.Printf("ML: H: %.4f S: %.4f L: %.4f \n", c.Hue, c.Saturation, c.Luminosity)
	newLum := (c.Luminosity * scale)
	newLumen := c.Luminosity + newLum

	if newLumen > 1 {

		// Use the difference to reduce the lightness.
		diff := 1 - c.Luminosity
		if diff > 0 {
			newLum := c.Luminosity + (diff / 2)

			if newLum > 1 {
				newLumen = 0.999
			} else {
				newLumen = newLum
			}
		} else {
			newLumen = 0.999
		}
	}

	if newLumen < 0 {
		newLumen = 0
	}

	newColor := colorful.Hsl(c.Hue, c.Saturation, newLumen)
	h, s, l := newColor.Hsl()

	return Color{
		C:          newColor,
		Hue:        h,
		Saturation: s,
		Luminosity: l,
		Alpha:      c.Alpha,
	}
}

//================================================================================================

// Color defines a basic struct which expresses the color values provided
// a struct containing HSL points.
type Color struct {
	C          colorful.Color
	Hue        float64 `json:"hue"`
	Luminosity float64 `json:"luminosity"`
	Saturation float64 `json:"saturation"`
	Alpha      float64 `json:"alpha"`
}

// String returns the Hex representation of the color.
func (c Color) String() string {
	return c.C.Hex()
}

// ColorFrom returns a Color instance representing the valid
// color values provided else returning error if the color value
// is not a valid color presentation i.e (rgb,rgba, hsl, hex).
func ColorFrom(value string) (Color, error) {
	var c colorful.Color

	alpha := float64(1)

	switch {
	case colors.IsHex(value):
		c, _ = colorful.Hex(value)
		break
	case colors.IsHSL(value):
		h, s, l := colors.ParseHSL(value)
		c = colorful.Hsl(h, s, l)
		break
	case colors.IsRGB(value):
		var red, green, blue int
		red, green, blue, alpha = colors.ParseRGB(value)
		c = colorful.Color{R: float64(red) / 255, G: float64(green) / 255, B: float64(blue) / 255}
		break
	case colors.IsRGBA(value):
		var red, green, blue int
		red, green, blue, alpha = colors.ParseRGB(value)
		c = colorful.Color{R: float64(red) / 255, G: float64(green) / 255, B: float64(blue) / 255}
		break
	default:
		return Color{}, errors.New("Invalid color value received")
	}

	h, s, l := c.Hsl()

	return Color{
		C:          c,
		Hue:        h,
		Saturation: s,
		Luminosity: l,
		Alpha:      alpha,
	}, nil
}

//==================================================================================

// contains different constants of different accepted sclaes.
const (
	AugmentedFourth = 1.414
	MinorSecond     = 1.067
	MajorSecond     = 1.125
	MinorThird      = 1.200
	MajorThird      = 1.250
	PerfectFourth   = 1.333
	PerfectFifth    = 1.500
	GoldenRation    = 1.618
)

// GenerateValueScale returns a value scale which is produced from generating
// a slice of n values representing the given scale value and are multipled
// by the provided base values.
func GenerateValueScale(scale float64, base float64) ([]float64, []float64) {

	// Generate scale based on 1.0 scale using the provided scale.
	max, min := GenerateScale(scale, 5, 10)

	// Multiply all scale value by the provided base.
	Times(len(max), func(index int) {
		max[index-1] = max[index-1] * base
	})

	// Multiply all scale value by the provided base.
	Times(len(min), func(index int) {
		min[index-1] = min[index-1] * base
	})

	return max, min
}

// GenerateScale returns a slice of values which are the a combination of
// a reducing + increasing scaled values of the provided scale generated from
// using the base initial 1.0 value against an ever incremental 1.0*(scale * n)
// or 1.0 / (scale *n) value, where n is the ever increasing index.
func GenerateScale(scale float64, minorCount int, majorCount int) ([]float64, []float64) {
	var scales []float64

	minorScales := make([]float64, minorCount)

	Times(minorCount, func(index int) {
		scaled := 1.0

		Times(index, func(_ int) {
			scaled *= scale
		})

		minorLen := len(minorScales)
		minorScales[minorLen-index] = 1.0 / scaled
	})

	scales = append(scales, 1.0)

	Times(majorCount, func(index int) {
		scaled := 1.0

		Times(index, func(_ int) {
			scaled *= scale
		})

		scales = append(scales, scaled)
	})

	return scales, minorScales
}

// Times using the provided count, runs the function (n-1) number of times, since
// it starts from zero.
func Times(n int, fn func(int)) {
	for i := 0; i < n; i++ {
		fn(i + 1)
	}
}
