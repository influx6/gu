Assets
========
Gu comes with a asset bundling system prebuilt within it's CLI which is included
has part of the things generated when a new component or app is created using the
`gu app` and `gu component` commands respecively.


## Types of Assets

In Gu there are basically two types of assets: (a) Static Files Assets and (b) Static Markup Assets

- Static Files Assets

These types of assets represent anything that lies within either the projects `public`
directory or within a `components` directory, where all file extensions except `.static.html`
qualify. Any file seen will be turned into a go package which can be imported to either serve
these files through a `http.FileSystem` or by retrieving the individual contents by use of the
relative path of the file.

Such files in .css, .js, .html, .less, where each is processed by the packers registered to
handle those specific files. For example, the less asset packer will instead return a single
converted file depending on it's settings.

See more: https://github.com/gu-io/gu/tree/master/assets/packers


- Static Markup Assets

These types of assets are special in that then use the extension `.static.html` and when seen
will be transformed into a single go file within the same package where the html written within
them is transformed into Gu markup. These reduces alot of the runtime overhead of parsing markup
right in code, when it's far faster if we are dealing with pure go structures.
