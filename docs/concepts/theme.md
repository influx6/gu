Style Guide
=======
Gu provides a [StyleGuide](https://github.com/gu-io/gu/tree/master/trees/themes/styleguide) package which provides a basic set of styles generated from the color sets provided and allows a quick startup in design. It's an optional piece of the package but can help to reduce alot of the metal load in creating a consistent set of design styles.

Example
-----------

Defining a style guide requires the selection of the primary and secondary color and brand which the project uses, with the selective color values for success and failures, which helps generate a giving stylesheet of property styles.

```go

var theme = styleguide.MustNew(styleguide.Attr{
	PrimaryColor: 			"#ffffff",
	SecondaryColor: 		"#ffffff",
	PrimaryBrandColor: 		"#ffffff",
	SecondaryBrandColor: 	"#ffffff",
	PrimaryWhite: 			"#ffffff",
	FailureColor: 			"#ffffff",
	SuccessColor: 			"#ffffff",
})


```

*Color fields must be in full else an error would be throw*

Executing `theme.CSS()` will produce: 

```css

html {
	font-size: 16px;
}

/*
____________ Base  classes ____________________________
 Letter spacing class go from highest to lowest where the highest is in positive 
   spacing value and the last in negative spacing values.
/*

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

/*
____________ Base border radius classes ____________________________

  These are classes for different border radius effect chosen specifically for 
  use with different components.

  .smallRadius: For basic border radius for small elements eg radio, checkboxes.
  .mediumRadius: For components like text, input, labels.
  .largeRadius: For components like cards, modal boxes, etc

*/

.border-radius-sm {
	-moz-border-radius: 2px;
	-webkit-border-radius: 2px;
	-o-border-radius: 2px;
	border-radius: 2px;
}

.border-radius-md {
	-moz-border-radius: 4px;
	-webkit-border-radius: 4px;
	-o-border-radius: 4px;
	border-radius: 4px;
}

.border-radius-lg {
	-moz-border-radius: 8px;
	-webkit-border-radius: 8px;
	-o-border-radius: 8px;
	border-radius: 8px;
}

.border-radius-cirle {
	-moz-border-radius: 50%;
	-webkit-border-radius: 50%;
	-o-border-radius: 50%;
	border-radius: 50%;
}

/*
____________ Base shadowdrop classes ____________________________

  These are classes for different shadow effect chosen specifically for 
  use with different components.

  .shadow: For basic shadows for normal elements.
  .shadow__drop: For shadow effects for dropdown/popovers type elements.
  .shadow__hover: For shadow effects for hovers. 
  .shadow__elevanted: For shadow effects for elevated modals, cards, etc
*/

.shadow {
	box-shadow: 0px 13px 20px 2px rgba(0, 0, 0, 0.45);
}

.shadow__dropdown {
	box-shadow: 0px 9px 30px 2px rgba(0, 0, 0, 0.51);
}

.shadow__hover {
	box-shadow: 0px 13px 30px 5px rgba(0, 0, 0, 0.58);
}

.shadow__elevated {
	box-shadow: 0px 20px 40px 4px rgba(0, 0, 0, 0.51);
}

/*
____________ Base font size classes ____________________________

  These are classes provide a simple set of font-scale font-size
  which allow you to use for scaling based on an initial font-size 
  set on a parent, they should scale well.

  font-size-sm: Defines font size for reducing sizes
  font-size-bg: Defines font size for increasing sizes using a scale eg MajorThirds.
*/


.font-size-sm__0 {
	font-size: 0.328em;
}

.font-size-sm__1 {
	font-size: 0.410em;
}

.font-size-sm__2 {
	font-size: 0.512em;
}

.font-size-sm__3 {
	font-size: 0.640em;
}

.font-size-sm__4 {
	font-size: 0.800em;
}



.font-size-bg__0 {
	font-size: 1.000em;
}

.font-size-bg__1 {
	font-size: 1.250em;
}

.font-size-bg__2 {
	font-size: 1.562em;
}

.font-size-bg__3 {
	font-size: 1.953em;
}

.font-size-bg__4 {
	font-size: 2.441em;
}

.font-size-bg__5 {
	font-size: 3.052em;
}

.font-size-bg__6 {
	font-size: 3.815em;
}

.font-size-bg__7 {
	font-size: 4.768em;
}

.font-size-bg__8 {
	font-size: 5.960em;
}

.font-size-bg__9 {
	font-size: 7.451em;
}

.font-size-bg__10 {
	font-size: 9.313em;
}


/*
____________ Color set ____________________________

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
	color: #ffffff;
}

.brand-border-color-primary {
	border-color: #ffffff;
}

.brand-background-color-primary {
	background: #ffffff;
}



.brand-background-color-primary-10 {
	background: #4d4d4d;
}

.brand-border-color-primary-10 {
	border-color: #4d4d4d;
}

.brand-primary-10 {
	color: #4d4d4d;
}


.brand-background-color-primary-20 {
	background: #4f4f4f;
}

.brand-border-color-primary-20 {
	border-color: #4f4f4f;
}

.brand-primary-20 {
	color: #4f4f4f;
}


.brand-background-color-primary-30 {
	background: #535353;
}

.brand-border-color-primary-30 {
	border-color: #535353;
}

.brand-primary-30 {
	color: #535353;
}


.brand-background-color-primary-40 {
	background: #595959;
}

.brand-border-color-primary-40 {
	border-color: #595959;
}

.brand-primary-40 {
	color: #595959;
}


.brand-background-color-primary-50 {
	background: #626262;
}

.brand-border-color-primary-50 {
	border-color: #626262;
}

.brand-primary-50 {
	color: #626262;
}



.color-primary {
	color: #ffffff;
}

.border-color-primary {
	border-color: #ffffff;
}

.background-color-primary {
	background: #ffffff;
}



.background-color-primary-10 {
	background: #4d4d4d;
}

.color-primary-10 {
	color: #4d4d4d;
}

.border-color-primary-10 {
	border-color: #4d4d4d;
}


.background-color-primary-20 {
	background: #4f4f4f;
}

.color-primary-20 {
	color: #4f4f4f;
}

.border-color-primary-20 {
	border-color: #4f4f4f;
}


.background-color-primary-30 {
	background: #535353;
}

.color-primary-30 {
	color: #535353;
}

.border-color-primary-30 {
	border-color: #535353;
}


.background-color-primary-40 {
	background: #595959;
}

.color-primary-40 {
	color: #595959;
}

.border-color-primary-40 {
	border-color: #595959;
}


.background-color-primary-50 {
	background: #626262;
}

.color-primary-50 {
	color: #626262;
}

.border-color-primary-50 {
	border-color: #626262;
}


/*____________ Secondary color set ____________________________*/

.color-secondary {
	color: #ffffff;
}

.background-color-secondary {
	background: #ffffff;
}

.border-color-secondary {
	background: #ffffff;
}



.background-color-secondary-10 {
	background: #4d4d4d;
}

.border-color-secondary-10 {
	border-color: #4d4d4d;
}

.color-secondary-10 {
	color: #4d4d4d;
}


.background-color-secondary-20 {
	background: #4f4f4f;
}

.border-color-secondary-20 {
	border-color: #4f4f4f;
}

.color-secondary-20 {
	color: #4f4f4f;
}


.background-color-secondary-30 {
	background: #535353;
}

.border-color-secondary-30 {
	border-color: #535353;
}

.color-secondary-30 {
	color: #535353;
}


.background-color-secondary-40 {
	background: #595959;
}

.border-color-secondary-40 {
	border-color: #595959;
}

.color-secondary-40 {
	color: #595959;
}


.background-color-secondary-50 {
	background: #626262;
}

.border-color-secondary-50 {
	border-color: #626262;
}

.color-secondary-50 {
	color: #626262;
}


.brand-color-secondary {
	color: #ffffff;
}

.brand-background-color-secondary {
	background: #ffffff;
}

.brand-border-color-secondary {
	background: #ffffff;
}



.brand-background-color-secondary-10 {
	background: #4d4d4d;
}

.brand-border-color-secondary-10 {
	border-color: #4d4d4d;
}

.brand-color-secondary-10 {
	color: #4d4d4d;
}


.brand-background-color-secondary-20 {
	background: #4f4f4f;
}

.brand-border-color-secondary-20 {
	border-color: #4f4f4f;
}

.brand-color-secondary-20 {
	color: #4f4f4f;
}


.brand-background-color-secondary-30 {
	background: #535353;
}

.brand-border-color-secondary-30 {
	border-color: #535353;
}

.brand-color-secondary-30 {
	color: #535353;
}


.brand-background-color-secondary-40 {
	background: #595959;
}

.brand-border-color-secondary-40 {
	border-color: #595959;
}

.brand-color-secondary-40 {
	color: #595959;
}


.brand-background-color-secondary-50 {
	background: #626262;
}

.brand-border-color-secondary-50 {
	border-color: #626262;
}

.brand-color-secondary-50 {
	color: #626262;
}



/*____________ Success color set ____________________________*/

.brand-success {
	color: #ffffff;
}

.background-color-success {
	background: #ffffff;
}

.border-color-success {
	border-color: #ffffff;
}



.background-color-success-10 {
	background: #4d4d4d;
}

.brand-success-10 {
	color: #4d4d4d;
}

.border-color-success-10 {
	border-color: #4d4d4d;
}


.background-color-success-20 {
	background: #4f4f4f;
}

.brand-success-20 {
	color: #4f4f4f;
}

.border-color-success-20 {
	border-color: #4f4f4f;
}


.background-color-success-30 {
	background: #535353;
}

.brand-success-30 {
	color: #535353;
}

.border-color-success-30 {
	border-color: #535353;
}


.background-color-success-40 {
	background: #595959;
}

.brand-success-40 {
	color: #595959;
}

.border-color-success-40 {
	border-color: #595959;
}


.background-color-success-50 {
	background: #626262;
}

.brand-success-50 {
	color: #626262;
}

.border-color-success-50 {
	border-color: #626262;
}



/*____________ White color set ____________________________*/

.background-color-white {
	background: #ffffff;
}

.brand-white {
	color: #ffffff;
}

.border-color-white {
	border-color: #ffffff;
}



.background-color-white-10 {
	background: #4d4d4d;
}

.brand-white-10 {
	color: #4d4d4d;
}

.border-color-white-10 {
	border-color: #4d4d4d;
}


.background-color-white-20 {
	background: #4f4f4f;
}

.brand-white-20 {
	color: #4f4f4f;
}

.border-color-white-20 {
	border-color: #4f4f4f;
}


.background-color-white-30 {
	background: #535353;
}

.brand-white-30 {
	color: #535353;
}

.border-color-white-30 {
	border-color: #535353;
}


.background-color-white-40 {
	background: #595959;
}

.brand-white-40 {
	color: #595959;
}

.border-color-white-40 {
	border-color: #595959;
}


.background-color-white-50 {
	background: #626262;
}

.brand-white-50 {
	color: #626262;
}

.border-color-white-50 {
	border-color: #626262;
}




/*____________ Failure color set ____________________________*/

.background-color-failure {
	background: #ffffff;
}

.brand-failure {
	color: #ffffff;
}

.border-color-failure {
	border-color: #ffffff;
}



.background-color-failure-10 {
	background: #4d4d4d;
}

.brand-failure-10 {
	color: #4d4d4d;
}

.border-color-failure-10 {
	border-color: #4d4d4d;
}


.background-color-failure-20 {
	background: #4f4f4f;
}

.brand-failure-20 {
	color: #4f4f4f;
}

.border-color-failure-20 {
	border-color: #4f4f4f;
}


.background-color-failure-30 {
	background: #535353;
}

.brand-failure-30 {
	color: #535353;
}

.border-color-failure-30 {
	border-color: #535353;
}


.background-color-failure-40 {
	background: #595959;
}

.brand-failure-40 {
	color: #595959;
}

.border-color-failure-40 {
	border-color: #595959;
}


.background-color-failure-50 {
	background: #626262;
}

.brand-failure-50 {
	color: #626262;
}

.border-color-failure-50 {
	border-color: #626262;
}


/*______________________________________________________________________*/
html {
	font-size: 16px;
}

/*
____________ Base  classes ____________________________
 Letter spacing class go from highest to lowest where the highest is in positive 
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

/*
____________ Base border radius classes ____________________________

  These are classes for different border radius effect chosen specifically for 
  use with different components.

  .smallRadius: For basic border radius for small elements eg radio, checkboxes.
  .mediumRadius: For components like text, input, labels.
  .largeRadius: For components like cards, modal boxes, etc

*/

.border-radius-sm {
	-moz-border-radius: 2px;
	-webkit-border-radius: 2px;
	-o-border-radius: 2px;
	border-radius: 2px;
}

.border-radius-md {
	-moz-border-radius: 4px;
	-webkit-border-radius: 4px;
	-o-border-radius: 4px;
	border-radius: 4px;
}

.border-radius-lg {
	-moz-border-radius: 8px;
	-webkit-border-radius: 8px;
	-o-border-radius: 8px;
	border-radius: 8px;
}

.border-radius-cirle {
	-moz-border-radius: 50%;
	-webkit-border-radius: 50%;
	-o-border-radius: 50%;
	border-radius: 50%;
}

/*
____________ Base shadowdrop classes ____________________________

  These are classes for different shadow effect chosen specifically for 
  use with different components.

  .shadow: For basic shadows for normal elements.
  .shadow__drop: For shadow effects for dropdown/popovers type elements.
  .shadow__hover: For shadow effects for hovers. 
  .shadow__elevanted: For shadow effects for elevated modals, cards, etc
*/

.shadow {
	box-shadow: 0px 13px 20px 2px rgba(0, 0, 0, 0.45);
}

.shadow__dropdown {
	box-shadow: 0px 9px 30px 2px rgba(0, 0, 0, 0.51);
}

.shadow__hover {
	box-shadow: 0px 13px 30px 5px rgba(0, 0, 0, 0.58);
}

.shadow__elevated {
	box-shadow: 0px 20px 40px 4px rgba(0, 0, 0, 0.51);
}

/*
____________ Base font size classes ____________________________

  These are classes provide a simple set of font-scale font-size
  which allow you to use for scaling based on an initial font-size 
  set on a parent, they should scale well.

  font-size-sm: Defines font size for reducing sizes
  font-size-bg: Defines font size for increasing sizes using a scale eg MajorThirds.
*/


.font-size-sm__0 {
	font-size: 0.328em;
}

.font-size-sm__1 {
	font-size: 0.410em;
}

.font-size-sm__2 {
	font-size: 0.512em;
}

.font-size-sm__3 {
	font-size: 0.640em;
}

.font-size-sm__4 {
	font-size: 0.800em;
}



.font-size-bg__0 {
	font-size: 1.000em;
}

.font-size-bg__1 {
	font-size: 1.250em;
}

.font-size-bg__2 {
	font-size: 1.562em;
}

.font-size-bg__3 {
	font-size: 1.953em;
}

.font-size-bg__4 {
	font-size: 2.441em;
}

.font-size-bg__5 {
	font-size: 3.052em;
}

.font-size-bg__6 {
	font-size: 3.815em;
}

.font-size-bg__7 {
	font-size: 4.768em;
}

.font-size-bg__8 {
	font-size: 5.960em;
}

.font-size-bg__9 {
	font-size: 7.451em;
}

.font-size-bg__10 {
	font-size: 9.313em;
}


/*
____________ Color set ____________________________

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
	color: #ffffff;
}

.brand-border-color-primary {
	border-color: #ffffff;
}

.brand-background-color-primary {
	background: #ffffff;
}



.brand-background-color-primary-10 {
	background: #4d4d4d;
}

.brand-border-color-primary-10 {
	border-color: #4d4d4d;
}

.brand-primary-10 {
	color: #4d4d4d;
}


.brand-background-color-primary-20 {
	background: #4f4f4f;
}

.brand-border-color-primary-20 {
	border-color: #4f4f4f;
}

.brand-primary-20 {
	color: #4f4f4f;
}


.brand-background-color-primary-30 {
	background: #535353;
}

.brand-border-color-primary-30 {
	border-color: #535353;
}

.brand-primary-30 {
	color: #535353;
}


.brand-background-color-primary-40 {
	background: #595959;
}

.brand-border-color-primary-40 {
	border-color: #595959;
}

.brand-primary-40 {
	color: #595959;
}


.brand-background-color-primary-50 {
	background: #626262;
}

.brand-border-color-primary-50 {
	border-color: #626262;
}

.brand-primary-50 {
	color: #626262;
}


.color-primary {
	color: #ffffff;
}

.border-color-primary {
	border-color: #ffffff;
}

.background-color-primary {
	background: #ffffff;
}



.background-color-primary-10 {
	background: #4d4d4d;
}

.color-primary-10 {
	color: #4d4d4d;
}

.border-color-primary-10 {
	border-color: #4d4d4d;
}


.background-color-primary-20 {
	background: #4f4f4f;
}

.color-primary-20 {
	color: #4f4f4f;
}

.border-color-primary-20 {
	border-color: #4f4f4f;
}


.background-color-primary-30 {
	background: #535353;
}

.color-primary-30 {
	color: #535353;
}

.border-color-primary-30 {
	border-color: #535353;
}


.background-color-primary-40 {
	background: #595959;
}

.color-primary-40 {
	color: #595959;
}

.border-color-primary-40 {
	border-color: #595959;
}


.background-color-primary-50 {
	background: #626262;
}

.color-primary-50 {
	color: #626262;
}

.border-color-primary-50 {
	border-color: #626262;
}


/*____________ Secondary color set ____________________________*/

.color-secondary {
	color: #ffffff;
}

.background-color-secondary {
	background: #ffffff;
}

.border-color-secondary {
	background: #ffffff;
}



.background-color-secondary-10 {
	background: #4d4d4d;
}

.border-color-secondary-10 {
	border-color: #4d4d4d;
}

.color-secondary-10 {
	color: #4d4d4d;
}


.background-color-secondary-20 {
	background: #4f4f4f;
}

.border-color-secondary-20 {
	border-color: #4f4f4f;
}

.color-secondary-20 {
	color: #4f4f4f;
}


.background-color-secondary-30 {
	background: #535353;
}

.border-color-secondary-30 {
	border-color: #535353;
}

.color-secondary-30 {
	color: #535353;
}


.background-color-secondary-40 {
	background: #595959;
}

.border-color-secondary-40 {
	border-color: #595959;
}

.color-secondary-40 {
	color: #595959;
}


.background-color-secondary-50 {
	background: #626262;
}

.border-color-secondary-50 {
	border-color: #626262;
}

.color-secondary-50 {
	color: #626262;
}


.brand-color-secondary {
	color: #ffffff;
}

.brand-background-color-secondary {
	background: #ffffff;
}

.brand-border-color-secondary {
	background: #ffffff;
}



.brand-background-color-secondary-10 {
	background: #4d4d4d;
}

.brand-border-color-secondary-10 {
	border-color: #4d4d4d;
}

.brand-color-secondary-10 {
	color: #4d4d4d;
}


.brand-background-color-secondary-20 {
	background: #4f4f4f;
}

.brand-border-color-secondary-20 {
	border-color: #4f4f4f;
}

.brand-color-secondary-20 {
	color: #4f4f4f;
}


.brand-background-color-secondary-30 {
	background: #535353;
}

.brand-border-color-secondary-30 {
	border-color: #535353;
}

.brand-color-secondary-30 {
	color: #535353;
}


.brand-background-color-secondary-40 {
	background: #595959;
}

.brand-border-color-secondary-40 {
	border-color: #595959;
}

.brand-color-secondary-40 {
	color: #595959;
}


.brand-background-color-secondary-50 {
	background: #626262;
}

.brand-border-color-secondary-50 {
	border-color: #626262;
}

.brand-color-secondary-50 {
	color: #626262;
}



/*____________ Success color set ____________________________*/

.brand-success {
	color: #ffffff;
}

.background-color-success {
	background: #ffffff;
}

.border-color-success {
	border-color: #ffffff;
}



.background-color-success-10 {
	background: #4d4d4d;
}

.brand-success-10 {
	color: #4d4d4d;
}

.border-color-success-10 {
	border-color: #4d4d4d;
}


.background-color-success-20 {
	background: #4f4f4f;
}

.brand-success-20 {
	color: #4f4f4f;
}

.border-color-success-20 {
	border-color: #4f4f4f;
}


.background-color-success-30 {
	background: #535353;
}

.brand-success-30 {
	color: #535353;
}

.border-color-success-30 {
	border-color: #535353;
}


.background-color-success-40 {
	background: #595959;
}

.brand-success-40 {
	color: #595959;
}

.border-color-success-40 {
	border-color: #595959;
}


.background-color-success-50 {
	background: #626262;
}

.brand-success-50 {
	color: #626262;
}

.border-color-success-50 {
	border-color: #626262;
}



/*____________ White color set ____________________________*/

.background-color-white {
	background: #ffffff;
}

.brand-white {
	color: #ffffff;
}

.border-color-white {
	border-color: #ffffff;
}



.background-color-white-10 {
	background: #4d4d4d;
}

.brand-white-10 {
	color: #4d4d4d;
}

.border-color-white-10 {
	border-color: #4d4d4d;
}


.background-color-white-20 {
	background: #4f4f4f;
}

.brand-white-20 {
	color: #4f4f4f;
}

.border-color-white-20 {
	border-color: #4f4f4f;
}


.background-color-white-30 {
	background: #535353;
}

.brand-white-30 {
	color: #535353;
}

.border-color-white-30 {
	border-color: #535353;
}


.background-color-white-40 {
	background: #595959;
}

.brand-white-40 {
	color: #595959;
}

.border-color-white-40 {
	border-color: #595959;
}


.background-color-white-50 {
	background: #626262;
}

.brand-white-50 {
	color: #626262;
}

.border-color-white-50 {
	border-color: #626262;
}




/*____________ Failure color set ____________________________*/

.background-color-failure {
	background: #ffffff;
}

.brand-failure {
	color: #ffffff;
}

.border-color-failure {
	border-color: #ffffff;
}



.background-color-failure-10 {
	background: #4d4d4d;
}

.brand-failure-10 {
	color: #4d4d4d;
}

.border-color-failure-10 {
	border-color: #4d4d4d;
}


.background-color-failure-20 {
	background: #4f4f4f;
}

.brand-failure-20 {
	color: #4f4f4f;
}

.border-color-failure-20 {
	border-color: #4f4f4f;
}


.background-color-failure-30 {
	background: #535353;
}

.brand-failure-30 {
	color: #535353;
}

.border-color-failure-30 {
	border-color: #535353;
}


.background-color-failure-40 {
	background: #595959;
}

.brand-failure-40 {
	color: #595959;
}

.border-color-failure-40 {
	border-color: #595959;
}


.background-color-failure-50 {
	background: #626262;
}

.brand-failure-50 {
	color: #626262;
}

.border-color-failure-50 {
	border-color: #626262;
}


/*______________________________________________________________________*/

```