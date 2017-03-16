CSSParser
==========
Although the name suggests that it parsers css but in actuality, it parsers css files for embedding.
CSSParser provides a baseline for running through all css files within a directory and it's subdirectories collecting all details and connections to generate an optimized structure for 
quick access which will be used to generate other structures as desired.



## Format
CSSParser uses the `trees/css` package included as part of `Gu` which uses `text/template` package 
and hence has no special format than what `text/template` and the general css style format rules.

## Inclusion Directive
CSSParser collects special inclusion directives which tell a giving css file to include another in it's
final output. This allows different css files to be bundled together with another file either due to shared styles.

It simply requires a comment line added at the top of the giving css file/files before any rule is defined, delimited by a new line, e.g


```css
/* #include ui/*:before, base/ui/base-ui:after, base/base.css */

.examples {
	width: 100px;
	height: 200px;
}
```

Where each included file or directory, has a `:modifier` ([after, before]) which defines where the content of the file, will be included in.