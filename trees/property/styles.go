package property

import "github.com/gu-io/gu/trees"

// CustomStyle returns a new Style with the giving name and
// value.
func CustomStyle(name, value string) trees.Property {
	return trees.NewCSSStyle(name, value)
}

// BackgroundStyle provides the color style value
func BackgroundStyle(value string) trees.Property {
	return trees.NewCSSStyle("background", value)
}

// ColorStyle provides the color style value
func ColorStyle(value string) trees.Property {
	return trees.NewCSSStyle("color", value)
}

// DisplayStyle provides the style setter that sets the css display value.
func DisplayStyle(ops string) trees.Property {
	return trees.NewCSSStyle("display", ops)
}

// HeightStyle provides the height style value
func HeightStyle(size string) trees.Property {
	return trees.NewCSSStyle("height", size)
}

// FontstringStyle provides the margin style value
func FontstringStyle(size string) trees.Property {
	return trees.NewCSSStyle("font-size", size)
}

// PaddingStyle provides the margin style value
func PaddingStyle(size string) trees.Property {
	return trees.NewCSSStyle("padding", size)
}

// MarginStyle provides the margin style value
func MarginStyle(size string) trees.Property {
	return trees.NewCSSStyle("margin", size)
}

// WidthStyle provides the width style value
func WidthStyle(size string) trees.Property {
	return trees.NewCSSStyle("width", size)
}
