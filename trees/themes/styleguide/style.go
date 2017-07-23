package styleguide

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"text/template"

	"github.com/gu-io/gu/trees/css"
	"github.com/influx6/faux/colors"
	colorful "github.com/lucasb-eyer/go-colorful"
)

// contains different constants used within the package.
const (
	shadowLarge    = "0px 20px 40px 4px rgba(0, 0, 0, 0.51)"
	shadowPopDrops = "0px 9px 30px 2px rgba(0, 0, 0, 0.51)"
	shadowHovers   = "0px 13px 30px 5px rgba(0, 0, 0, 0.58)"
	shadowNormal   = "0px 13px 20px 2px rgba(0, 0, 0, 0.45)"

	AnimationCurveFastOutSlowIn   = "cubic-bezier(0.4, 0, 0.2, 1)"
	AnimationCurveLinearOutSlowIn = "cubic-bezier(0, 0, 0.2, 1)"
	AnimationCurveFastOutLinearIn = "cubic-bezier(0.4, 0, 1, 1)"
	AnimationCurveDefault         = AnimationCurveFastOutSlowIn

	smallBorderRadius  = 2
	mediumBorderRadius = 4
	largeBorderRadius  = 8

	AugmentedFourth = 1.414
	MinorSecond     = 1.067
	MajorSecond     = 1.125
	MinorThird      = 1.200
	MajorThird      = 1.250
	PerfectFourth   = 1.333
	PerfectFifth    = 1.500
	GoldenRatio     = 1.618

	LuminFlat     = 1.015
	LuminFat      = 1.200
	LuminFatThird = 1.245
)

var (
	helpers = template.FuncMap{
		"quote": func(b interface{}) string {
			switch bo := b.(type) {
			case string:
				return strconv.Quote(bo)
			case int:
				return strconv.Quote(strconv.Itoa(bo))
			case int64:
				return strconv.Quote(strconv.Itoa(int(bo)))
			case float32:
				mo := strconv.FormatFloat(float64(bo), 'f', 4, 32)
				return strconv.Quote(mo)
			case float64:
				mo := strconv.FormatFloat(bo, 'f', 4, 32)
				return strconv.Quote(mo)
			case rune:
				return strconv.QuoteRune(bo)
			default:
				mo, err := json.Marshal(b)
				if err != nil {
					return err.Error()
				}

				return strconv.Quote(string(mo))
			}
		},
		"prefixInt": func(prefix string, b int) string {
			return fmt.Sprintf("%s%d", prefix, b)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"lessThanEqual": func(a, b int) bool {
			return a <= b
		},
		"greaterThanEqual": func(a, b int) bool {
			return a >= b
		},
		"lessThan": func(a, b int) bool {
			return a < b
		},
		"greaterThan": func(a, b int) bool {
			return a > b
		},
		"len": func(a interface{}) int {
			switch real := a.(type) {
			case []interface{}:
				return len(real)
			case [][]byte:
				return len(real)
			case []byte:
				return len(real)
			case []float32:
				return len(real)
			case []float64:
				return len(real)
			case []string:
				return len(real)
			case []int:
				return len(real)
			default:
				return 0
			}
		},
		"multiply": func(a, b int) int {
			return a * b
		},
		"subtract": func(a, b int) int {
			return a - b
		},
		"divide": func(a, b int) int {
			return a / b
		},
		"perc": func(a, b float64) float64 {
			return (a / b) * 100
		},
		"textRhythmn": func(lineHeight int, capHeight int, fontSize int) int {
			return ((lineHeight - capHeight) * fontSize) / 2
		},
		"textRhythmnEM": func(lineHeight, capHeight, fontSize float64) float64 {
			return ((lineHeight - capHeight) * fontSize) / 2
		},
	}
)

// Attr defines different color and size of strings to define a specific brand.
type Attr struct {
	PrimaryWhite                  string
	SuccessColor                  string
	FailureColor                  string
	PrimaryColor                  string
	SecondaryColor                string
	PrimaryBrandColor             string
	SecondaryBrandColor           string
	AnimationCurveDefault         string
	AnimationCurveFastOutLinearIn string
	AnimationCurveFastOutSlowIn   string
	AnimationCurveLinearOutSlowIn string
	BaseScale                     float64             // BaseScale to use for generating expansion/detraction scale for font sizes.
	HeaderBaseScale               float64             // BaseScale to use for generating expansion/detraction scale for header h1-h6 tags.
	MinimumScaleCount             int                 // Total scale to generate small font sizes.
	MaximumScaleCount             int                 // Total scale to generate large font sizes
	MinimumHeadScaleCount         int                 // Total scale to generate small font sizes.
	MaximumHeadScaleCount         int                 // Total scale to generate large font sizes
	BaseFontSize                  int                 // BaseFontSize for typeface using the provide BaseScale.
	SmallBorderRadius             int                 // SmallBorderRadius for tiny components eg checkbox, radio buttons.
	MediumBorderRadius            int                 // MediaBorderRadius for buttons, inputs, etc
	LargeBorderRadius             int                 // LargeBorderRadius for components like cards, modals, etc.
	FloatingShadow                string              // shadow for floating icons, elements.
	HoverShadow                   string              // shadow for over dialog etc
	DropShadow                    string              // Useful for popovers/dropovers
	BaseShadow                    string              // Normal shadow of elemnts
	MaterialPalettes              map[string][]string `json:"palettes"`
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
	Brand            StyleColors `json:"brand"`
	SmallFontScale   []float64
	BigFontScale     []float64
	SmallHeaderScale []float64
	BigHeaderScale   []float64
	inited           bool
}

// Must returns the giving style or panics if it fails.
func Must(attr Attr) StyleGuide {
	style, err := New(attr)
	if err != nil {
		panic(err)
	}

	return style
}

// New returns a new StyleGuide object which generates the necessary css
// styles to utilize the defined style within any project.
func New(attr Attr) (StyleGuide, error) {
	var style StyleGuide

	if err := style.Init(); err != nil {
		return style, err
	}

	return style, nil
}

// Init initializes the style guide and all internal properties into
// appropriate defaults and states.
func (style *StyleGuide) Init() error {
	var err error

	style.Attr = initAttr(style.Attr)

	shm, bhm := GenerateValueScale(1, style.Attr.HeaderBaseScale, style.Attr.MinimumHeadScaleCount, style.Attr.MaximumHeadScaleCount)
	style.BigHeaderScale = bhm
	style.SmallHeaderScale = shm

	sm, bg := GenerateValueScale(1, style.Attr.BaseScale, style.Attr.MinimumScaleCount, style.Attr.MaximumScaleCount)
	style.BigFontScale = bg
	style.SmallFontScale = sm

	if style.Attr.PrimaryBrandColor != "" {
		style.Brand.PrimaryBrand, err = NewTones(style.Attr.PrimaryBrandColor)
		if err != nil {
			return errors.New("Invalid primary brand color")
		}
	}

	if style.Attr.SecondaryBrandColor != "" {
		style.Brand.SecondaryBrand, err = NewTones(style.Attr.SecondaryBrandColor)
		if err != nil {
			return errors.New("Invalid secondary brand color")
		}
	}

	style.Brand.Primary, err = NewTones(style.Attr.PrimaryColor)
	if err != nil {
		return errors.New("Invalid primary color")
	}

	style.Brand.Secondary, err = NewTones(style.Attr.SecondaryColor)
	if err != nil {
		return errors.New("Invalid secondary color")
	}

	style.Brand.White, err = NewTones(style.Attr.PrimaryWhite)
	if err != nil {
		return errors.New("Invalid primary white tone color")
	}

	style.Brand.Success, err = NewTones(style.Attr.SuccessColor)
	if err != nil {
		return errors.New("Invalid success color")
	}

	style.Brand.Failure, err = NewTones(style.Attr.FailureColor)
	if err != nil {
		return errors.New("Invalid failure color")
	}

	style.inited = true
	return err
}

// Ready returns true/false whether the giving style guide has being
// initialized.
func (style *StyleGuide) Ready() bool {
	return style.inited
}

// Stylesheet returns a css.Rule object which contains the styleguide style
// rules.
func (style *StyleGuide) Stylesheet() *css.Rule {
	return css.New(style.CSS(), nil)
}

// CSS returns a css style content for usage with a css stylesheet.
func (style *StyleGuide) CSS() string {
	tml, err := template.New("styleguide").Funcs(helpers).Parse(styleTemplate)
	if err != nil {
		return err.Error()
	}

	var buf bytes.Buffer
	if terr := tml.Execute(&buf, style); terr != nil {
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

func initAttr(attr Attr) Attr {
	if attr.MaterialPalettes == nil || len(attr.MaterialPalettes) == 0 {
		attr.MaterialPalettes = MaterialPalettes
	}

	if attr.PrimaryColor == "" {
		attr.PrimaryColor = MaterialPalettes["blue"][5]
	}

	if attr.SecondaryColor == "" {
		attr.SecondaryColor = MaterialPalettes["deep-purple"][5]
	}

	if attr.PrimaryWhite == "" {
		attr.PrimaryWhite = MaterialPalettes["white"][0]
	}

	if attr.SuccessColor == "" {
		attr.SuccessColor = MaterialPalettes["green"][5]
	}

	if attr.FailureColor == "" {
		attr.FailureColor = MaterialPalettes["red"][5]
	}

	if attr.BaseFontSize <= 0 {
		attr.BaseFontSize = 16
	}

	if attr.BaseScale <= 0 {
		attr.BaseScale = PerfectFourth
	}

	if attr.HeaderBaseScale <= 0 {
		attr.HeaderBaseScale = MajorThird
	}

	if attr.MinimumHeadScaleCount == 0 {
		attr.MinimumHeadScaleCount = 4
	}

	if attr.MaximumHeadScaleCount == 0 {
		attr.MaximumHeadScaleCount = 6
	}

	if attr.MinimumScaleCount == 0 {
		attr.MinimumScaleCount = 10
	}

	if attr.MaximumScaleCount == 0 {
		attr.MaximumScaleCount = 10
	}

	if attr.AnimationCurveDefault == "" {
		attr.AnimationCurveDefault = AnimationCurveDefault
	}

	if attr.AnimationCurveFastOutLinearIn == "" {
		attr.AnimationCurveFastOutLinearIn = AnimationCurveFastOutLinearIn
	}

	if attr.AnimationCurveFastOutSlowIn == "" {
		attr.AnimationCurveFastOutSlowIn = AnimationCurveFastOutSlowIn
	}

	if attr.AnimationCurveLinearOutSlowIn == "" {
		attr.AnimationCurveLinearOutSlowIn = AnimationCurveLinearOutSlowIn
	}

	if attr.FloatingShadow == "" {
		attr.FloatingShadow = shadowLarge
	}

	if attr.HoverShadow == "" {
		attr.HoverShadow = shadowHovers
	}

	if attr.BaseShadow == "" {
		attr.BaseShadow = shadowNormal
	}

	if attr.DropShadow == "" {
		attr.DropShadow = shadowPopDrops
	}

	if attr.SmallBorderRadius <= 0 {
		attr.SmallBorderRadius = smallBorderRadius
	}

	if attr.MediumBorderRadius <= 0 {
		attr.MediumBorderRadius = mediumBorderRadius
	}

	if attr.LargeBorderRadius <= 0 {
		attr.LargeBorderRadius = largeBorderRadius
	}

	return attr
}

//================================================================================================

// HamonicsFrom uses the above scale to return a slice of new Colors based on the provided
// HamonyScale set.
func HamonicsFrom(c Color) Tones {

	var scale []float64

	min, max := GenerateValueScale(0.1, LuminFatThird, 1, 10)

	lastItem := max[len(max)-1]
	_, inmax := GenerateValueScale(lastItem, LuminFlat, 0, 8)

	inmax = inmax[1:]

	reverse(len(min), func(index int) {
		scale = append(scale, min[index])
	})

	scale = append(scale, max...)
	scale = append(scale, inmax...)

	var colors []Color

	// TODO(alex): Should we have another scale for saturation?
	for _, scale := range scale {
		if scale > 1 {
			scale = 1
		}

		newColor := colorful.Hsl(c.Hue, c.Saturation, scale)
		h, s, l := newColor.Hsl()

		colors = append(colors, Color{
			C:          newColor,
			Hue:        h,
			Saturation: s,
			Luminosity: l,
			Alpha:      c.Alpha,
		})
	}

	var t Tones
	t.Base = c
	t.Grades = colors

	return t
}

// AdditiveSaturation adds the provided scale to the colors saturation value
// returning a new color suited to match.
func AdditiveSaturation(c Color, scale float64) Color {
	newLumen := c.Saturation + scale

	if newLumen > 1 {
		newLumen = 1
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
		newLumen = 1
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

// AdditiveLumination adds the provided scale to the colors Luminouse value
// returning a new color suited to match.
func AdditiveLumination(c Color, scale float64) Color {
	newLumen := c.Luminosity + scale

	if newLumen > 1 {
		newLumen = 1
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
		newLumen = 1
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

// GenerateValueScale returns a slice of values which are the a combination of
// a reducing + increasing scaled values of the provided scale generated from
// using the base initial 1.0 value against an ever incremental 1.0*(scale * n)
// or 1.0 / (scale *n) value, where n is the ever increasing index.
func GenerateValueScale(base float64, scale float64, minorCount int, majorCount int) ([]float64, []float64) {
	var major, minor []float64

	times(minorCount, func(index int) {
		if index > 1 {
			prevValue := minor[len(minor)-1]
			minor = append(minor, prevValue/scale)
			return
		}

		minor = append(minor, base/scale)
	})

	major = append(major, base)

	times(majorCount, func(index int) {
		if index > 1 {
			prevValue := major[index-1]
			major = append(major, prevValue*scale)
			return
		}

		major = append(major, base*scale)
	})

	return minor, major
}

func times(n int, fn func(int)) {
	for i := 0; i < n; i++ {
		fn(i + 1)
	}
}

func reverse(n int, fn func(int)) {
	for i := n; i > 0; i-- {
		fn(i - 1)
	}
}
