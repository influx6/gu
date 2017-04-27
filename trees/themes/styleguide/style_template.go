// Package style is generated to contain the content of a css template for generating a styleguide for use in projects.

// Document is auto-generate and should not be modified by hand.

//go:generate go run generate.go

package styleguide

// styleTemplate contains the text template used to generated the full set of 
// css template for a giving styleguide.
var styleTemplate = `
html {
	font-size: {{ .BaseFontSize }}px;
}

/*____________ Base  classes ____________________________*/

/*____________ Base border radius classes ____________________________*/
/**
  These are classes for different border radius effect chosen specifically for 
  use with different components.

  .smallRadius: For basic border radius for small elements eg radio, checkboxes.
  .mediumRadius: For components like text, input, labels.
  .largeRadius: For components like cards, modal boxes, etc

*/

.border-radius-sm {
	-moz-border-radius: {{ .SmallBorderRadius }}px;
	-webkit-border-radius: {{ .SmallBorderRadius }}px;
	-o-border-radius: {{ .SmallBorderRadius }}px;
	border-radius: {{ .SmallBorderRadius }}px;
}

.border-radius-md {
	-moz-border-radius: {{ .MediumBorderRadius }}px;
	-webkit-border-radius: {{ .MediumBorderRadius }}px;
	-o-border-radius: {{ .MediumBorderRadius }}px;
	border-radius: {{ .MediumBorderRadius }}px;
}

.border-radius-lg {
	-moz-border-radius: {{ .LargeBorderRadius }}px;
	-webkit-border-radius: {{ .LargeBorderRadius }}px;
	-o-border-radius: {{ .LargeBorderRadius }}px;
	border-radius: {{ .LargeBorderRadius }}px;
}

/*____________ Base shadowdrop classes ____________________________*/
/**
  These are classes for different shadow effect chosen specifically for 
  use with different components.

  .shadow: For basic shadows for normal elements.
  .shadow__drop: For shadow effects for dropdown/popovers type elements.
  .shadow__hover: For shadow effects for hovers. 
  .shadow__elevanted: For shadow effects for elevated modals, cards, etc
*/

.shadow {
	box-shadow: {{ .BaseShadow }};
}

.shadow__dropdown {
	box-shadow: {{ .DropShadow }};
}

.shadow__hover {
	box-shadow: {{ .HoverShadow }};
}

.shadow__elevated {
	box-shadow: {{ .FloatingShadow }};
}

/*______________________________________________________________________*/


/*____________ Base font size classes ____________________________*/
/**
  These are classes provide a simple set of font-scale font-size
  which allow you to use for scaling based on an initial font-size 
  set on a parent, they should scale well.

  font-size-sm: Defines font size for reducing sizes
  font-size-bg: Defines font size for increasing sizes using a scale eg MajorThirds.
*/

{{ range $key, $item := .SmallFontScale }}
.font-size-sm__{{ $key }} {
	font-size: {{$item}}em;
}
{{ end }}

{{ range $key, $item := .BigFontScale }}
.font-size-bg__{{ $key }} {
	font-size: {{$item}}em;
}
{{ end }}

/*______________________________________________________________________*/

/*____________ Base primary color set ____________________________*/
/**
  These are classes provide a color set based on specific brand colors 
  provided, these allows us to easily generate color, background and border-color
  classes suited to provide a consistent color design for use in project.

  These brand colors are divided into:

  primary: Main color for the giving project's brand
  secondary: Secondary brand color for project.
  success: Color for successful operation or messages.
  failure: Color for failed operations or messages.


  These are further subdivided into these diffent tones/shades:

  base: The original color without any modification.
  lumen: The original color version with a bit of increase in lumination.
  lumin: The original color version as close to white or 50/60% illumination.
  light: The original color version as close to white or 100% illumination.
  Lush: The original color with more black added.
  Dark: The original color with 50% dark harded.

  We further generate classes for Color, Border-Color and Background based on the division and 
  subdivision.

*/

.brand-primary__base, .brand-primary {
	color: {{.Brand.Primary.Base}};
}

.brand-primary__lumen {
	color: {{.Brand.Primary.Lumen}};
}

.brand-primary__lumin {
	color: {{.Brand.Primary.Lumin}};
}

.brand-primary__light {
	color: {{.Brand.Primary.Light}};
}

.brand-primary__lush {
	color: {{.Brand.Primary.Lush}};
}

.brand-primary__dark {
	color: {{.Brand.Primary.Dark}};
}


.brand-secondary__base, .brand-secondary {
	color: {{.Brand.Secondary.Base}};
}

.brand-secondary__lumen {
	color: {{.Brand.Secondary.Lumen}};
}

.brand-secondary__lumin {
	color: {{.Brand.Secondary.Lumin}};
}

.brand-secondary__light {
	color: {{.Brand.Secondary.Light}};
}

.brand-secondary__lush {
	color: {{.Brand.Secondary.Lush}};
}

.brand-secondary__dark {
	color: {{.Brand.Secondary.Dark}};
}


.brand-success__base, .brand-success {
	color: {{.Brand.Success.Base}};
}

.brand-success__lumen {
	color: {{.Brand.Success.Lumen}};
}

.brand-success__lumin {
	color: {{.Brand.Success.Lumin}};
}

.brand-success__light {
	color: {{.Brand.Success.Light}};
}

.brand-success__lush {
	color: {{.Brand.Success.Lush}};
}

.brand-success__dark {
	color: {{.Brand.Success.Dark}};
}


.brand-failure__base, .brand-failure {
	color: {{.Brand.Failure.Base}};
}

.brand-failure__lumen {
	color: {{.Brand.Failure.Lumen}};
}

.brand-failure__lumin {
	color: {{.Brand.Failure.Lumin}};
}

.brand-failure__light {
	color: {{.Brand.Failure.Light}};
}

.brand-failure__lush {
	color: {{.Brand.Failure.Lush}};
}

.brand-failure__dark {
	color: {{.Brand.Failure.Dark}};
}

/*______________________________________________________________________*/


*____________ Base background color set ____________________________*/

.background-color-primary__base, .background-color-primary {
	color: {{.Brand.Primary.Base}};
}

.background-color-primary__lumen {
	color: {{.Brand.Primary.Lumen}};
}

.background-color-primary__lumin {
	color: {{.Brand.Primary.Lumin}};
}

.background-color-primary__light {
	color: {{.Brand.Primary.Light}};
}

.background-color-primary__lush {
	color: {{.Brand.Primary.Lush}};
}

.background-color-primary__dark {
	color: {{.Brand.Primary.Dark}};
}


.background-color-secondary__base, .background-color-secondary {
	color: {{.Brand.Secondary.Base}};
}

.background-color-secondary__lumen {
	color: {{.Brand.Secondary.Lumen}};
}

.background-color-secondary__lumin {
	color: {{.Brand.Secondary.Lumin}};
}

.background-color-secondary__light {
	color: {{.Brand.Secondary.Light}};
}

.background-color-secondary__lush {
	color: {{.Brand.Secondary.Lush}};
}

.background-color-secondary__dark {
	color: {{.Brand.Secondary.Dark}};
}

.background-color-success__base, .background-color-success {
	color: {{.Brand.Success.Base}};
}

.background-color-success__lumen {
	color: {{.Brand.Success.Lumen}};
}

.background-color-success__lumin {
	color: {{.Brand.Success.Lumin}};
}

.background-color-success__light {
	color: {{.Brand.Success.Light}};
}

.background-color-success__lush {
	color: {{.Brand.Success.Lush}};
}

.background-color-success__dark {
	color: {{.Brand.Success.Dark}};
}


.background-color-failure__base, .background-color-failure {
	color: {{.Brand.Failure.Base}};
}

.background-color-failure__lumen {
	color: {{.Brand.Failure.Lumen}};
}

.background-color-failure__lumin {
	color: {{.Brand.Failure.Lumin}};
}

.background-color-failure__light {
	color: {{.Brand.Failure.Light}};
}

.background-color-failure__lush {
	color: {{.Brand.Failure.Lush}};
}

.background-color-failure__dark {
	color: {{.Brand.Failure.Dark}};
}

/*______________________________________________________________________*/


/*____________ Base border primary color set ____________________________*/

.border-color-primary__base, .border-color-primary {
	border-color: {{.Brand.Primary.Base}};
}

.border-color-primary__lumen {
	border-color: {{.Brand.Primary.Lumen}};
}

.border-color-primary__lumin {
	border-color: {{.Brand.Primary.Lumin}};
}

.border-color-primary__light {
	border-color: {{.Brand.Primary.Light}};
}

.border-color-primary__lush {
	border-color: {{.Brand.Primary.Lush}};
}

.border-color-primary__dark {
	border-color: {{.Brand.Primary.Dark}};
}


.border-color-secondary__base, .border-color-secondary {
	border-color: {{.Brand.Secondary.Base}};
}

.border-color-secondary__lumen {
	border-color: {{.Brand.Secondary.Lumen}};
}

.border-color-secondary__lumin {
	border-color: {{.Brand.Secondary.Lumin}};
}

.border-color-secondary__light {
	border-color: {{.Brand.Secondary.Light}};
}

.border-color-secondary__lush {
	border-color: {{.Brand.Secondary.Lush}};
}

.border-color-secondary__dark {
	border-color: {{.Brand.Secondary.Dark}};
}


.border-color-success__base, .border-color-success {
	border-color: {{.Brand.Success.Base}};
}

.border-color-success__lumen {
	border-color: {{.Brand.Success.Lumen}};
}

.border-color-success__lumin {
	border-color: {{.Brand.Success.Lumin}};
}

.border-color-success__light {
	border-color: {{.Brand.Success.Light}};
}

.border-color-success__lush {
	border-color: {{.Brand.Success.Lush}};
}

.border-color-success__dark {
	border-color: {{.Brand.Success.Dark}};
}


.border-color-failure__base, .border-color-failure {
	border-color: {{.Brand.Failure.Base}};
}

.border-color-failure__lumen {
	border-color: {{.Brand.Failure.Lumen}};
}

.border-color-failure__lumin {
	border-color: {{.Brand.Failure.Lumin}};
}

.border-color-failure__light {
	border-color: {{.Brand.Failure.Light}};
}

.border-color-failure__lush {
	border-color: {{.Brand.Failure.Lush}};
}

.border-color-failure__dark {
	border-color: {{.Brand.Failure.Dark}};
}

/*______________________________________________________________________*/
`
