package common

// Remover defines an interface which exposes a remove method.
type Remover interface {
	Remove()
}

// Settings defines a structure which contains fields that are used to contain
// specific user settings for the gu build system.
type Settings struct {
	Public Public `toml:"public"`
	Theme  Theme  `toml:"theme"`
}

// Public defines giving settings for the public assets folder which will be build
// to generate a embeddeable and servable assets package.
type Public struct {
	Path    string `toml:"path"`
	Package string `toml:"package"`
}

// Theme defines a struct whhich contains settings for generating a stylesheet of
// css rules.
type Theme struct {
	MinimumScaleCount             int     // Total scale to generate small font sizes.
	MaximumScaleCount             int     // Total scale to generate large font sizes
	MinimumHeadScaleCount         int     // Total scale to generate small font sizes.
	MaximumHeadScaleCount         int     // Total scale to generate large font sizes
	BaseFontSize                  int     // BaseFontSize for typeface using the provide BaseScale.
	SmallBorderRadius             int     // SmallBorderRadius for tiny components eg checkbox, radio buttons.
	MediumBorderRadius            int     // MediaBorderRadius for buttons, inputs, etc
	LargeBorderRadius             int     // LargeBorderRadius for components like cards, modals, etc.
	BaseScale                     float64 // BaseScale to use for generating expansion/detraction scale for font sizes.
	HeaderBaseScale               float64 // BaseScale to use for generating expansion/detraction scale for header h1-h6 tags.
	PrimaryWhite                  string
	SuccessColor                  string
	FailureColor                  string
	PrimaryColor                  string
	HoverShadow                   string
	DropShadow                    string
	BaseShadow                    string
	SecondaryColor                string
	FloatingShadow                string
	PrimaryBrandColor             string
	SecondaryBrandColor           string
	AnimationCurveDefault         string
	AnimationCurveFastOutLinearIn string
	AnimationCurveFastOutSlowIn   string
	AnimationCurveLinearOutSlowIn string
	MaterialPalettes              map[string][]string
}
