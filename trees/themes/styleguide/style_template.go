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


/* Letter spacing class go from highest to lowest where the highest is in positive 
   spacing value and the last in negative spacing values.
*/
.letter-spacing-1 {
	letter-spacing: 1px;
}

.letter-spacing-2 {
	letter-spacing: -0.5px;
}

.letter-spacing-3 {
	letter-spacing: -1px;
}

.letter-spacing-4 {
	letter-spacing: -2px;
}


.wrap {
  text-wrap: wrap;
  white-space: -moz-pre-wrap;
  white-space: pre-wrap;
  word-wrap: break-word;
}

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

.border-radius-cirle {
	-moz-border-radius: 50%;
	-webkit-border-radius: 50%;
	-o-border-radius: 50%;
	border-radius: 50%;
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

/*____________ Color set ____________________________*/
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

  Other color tones are graded from 10...nth where n is a multiple of 10 * index.
  The lowest grade of 10 is where the color is close to it's darkest version while 
  the highest means a continous increase in luminousity.

  We further generate classes for Color, Border-Color and Background based on the division and 
  subdivision.

*/

.brand-color-primary {
	color: {{.Brand.PrimaryBrand.Base}};
}

.brand-border-color-primary {
	border-color: {{.Brand.PrimaryBrand.Base}};
}

.brand-background-color-primary {
	background: {{.Brand.PrimaryBrand.Base}};
}

{{ range $index, $item := .Brand.PrimaryBrand.Grades }}
{{ $rn := add $index 1 }}
.brand-primary-{{ multiply $rn 10}} {
	color: {{$item}};
}
{{ end }}

{{ range $index, $item := .Brand.PrimaryBrand.Grades }}
{{ $rn := add $index 1 }}
.brand-background-color-primary-{{ multiply $rn 10}} {
	background: {{$item}};
}
{{ end }}


{{ range $index, $item := .Brand.PrimaryBrand.Grades }}
{{ $rn := add $index 1 }}
.brand-border-color-primary-{{ multiply $rn 10}} {
	border-color: {{$item}};
}
{{ end }}

.color-primary {
	color: {{.Brand.Primary.Base}};
}

.border-color-primary {
	border-color: {{.Brand.Primary.Base}};
}

.background-color-primary {
	background: {{.Brand.Primary.Base}};
}

{{ range $index, $item := .Brand.Primary.Grades }}
{{ $rn := add $index 1 }}
.color-primary-{{ multiply $rn 10}} {
	color: {{$item}};
}
{{ end }}

{{ range $index, $item := .Brand.Primary.Grades }}
{{ $rn := add $index 1 }}
.background-color-primary-{{ multiply $rn 10}} {
	background: {{$item}};
}
{{ end }}


{{ range $index, $item := .Brand.Primary.Grades }}
{{ $rn := add $index 1 }}
.border-color-primary-{{ multiply $rn 10}} {
	border-color: {{$item}};
}
{{ end }}


*____________ Base background color set ____________________________*/


.color-secondary {
	color: {{.Brand.Secondary.Base}};
}

.background-color-secondary {
	background: {{.Brand.Secondary.Base}};
}

.border-color-secondary {
	background: {{.Brand.Secondary.Base}};
}

{{ range $index, $item := .Brand.Secondary.Grades }}
{{ $rn := add $index 1 }}
.color-secondary-{{ multiply $rn 10}} {
	color: {{$item}};
}
{{ end }}

{{ range $index, $item := .Brand.Secondary.Grades }}
{{ $rn := add $index 1 }}
.background-color-secondary-{{ multiply $rn 10}} {
	background: {{$item}};
}
{{ end }}

{{ range $index, $item := .Brand.Secondary.Grades }}
{{ $rn := add $index 1 }}
.border-color-secondary-{{ multiply $rn 10}} {
	border-color: {{$item}};
}
{{ end }}


.brand-color-secondary {
	color: {{.Brand.SecondaryBrand.Base}};
}

.brand-background-color-secondary {
	background: {{.Brand.SecondaryBrand.Base}};
}

.brand-border-color-secondary {
	background: {{.Brand.SecondaryBrand.Base}};
}

{{ range $index, $item := .Brand.SecondaryBrand.Grades }}
{{ $rn := add $index 1 }}
.brand-color-secondary-{{ multiply $rn 10}} {
	color: {{$item}};
}
{{ end }}

{{ range $index, $item := .Brand.SecondaryBrand.Grades }}
{{ $rn := add $index 1 }}
.brand-background-color-secondary-{{ multiply $rn 10}} {
	background: {{$item}};
}
{{ end }}

{{ range $index, $item := .Brand.SecondaryBrand.Grades }}
{{ $rn := add $index 1 }}
.brand-border-color-secondary-{{ multiply $rn 10}} {
	border-color: {{$item}};
}
{{ end }}

*____________ Base secondary color set ____________________________*/

.brand-success {
	color: {{.Brand.Success.Base}};
}

.background-color-success {
	background: {{.Brand.Success.Base}};
}

.border-color-success {
	border-color: {{.Brand.Success.Base}};
}

{{ range $index, $item := .Brand.Success.Grades }}
{{ $rn := add $index 1 }}
.brand-success-{{ multiply $rn 10}} {
	color: {{$item}};
}
{{ end }}

{{ range $index, $item := .Brand.Success.Grades }}
{{ $rn := add $index 1 }}
.background-color-success-{{ multiply $rn 10}} {
	background: {{$item}};
}
{{ end }}

{{ range $index, $item := .Brand.Success.Grades }}
{{ $rn := add $index 1 }}
.border-color-success-{{ multiply $rn 10}} {
	border-color: {{$item}};
}
{{ end }}

/*____________ Base border primary color set ____________________________*/

.background-color-failure {
	background: {{.Brand.Failure.Base}};
}

.brand-failure {
	color: {{.Brand.Failure.Base}};
}

.border-color-failure {
	border-color: {{.Brand.Failure.Base}};
}

{{ range $index, $item := .Brand.Failure.Grades }}
{{ $rn := add $index 1 }}
.background-color-failure-{{ multiply $rn 10}} {
	background: {{$item}};
}
{{ end }}

{{ range $index, $item := .Brand.Failure.Grades }}
{{ $rn := add $index 1 }}
.brand-failure-{{ multiply $rn 10}} {
	color: {{$item}};
}
{{ end }}

{{ range $index, $item := .Brand.Failure.Grades }}
{{ $rn := add $index 1 }}
.border-color-failure-{{ multiply $rn 10}} {
	border-color: {{$item}};
}
{{ end }}

/*______________________________________________________________________*/
`
