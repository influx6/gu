package shell

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gu-io/gu/trees"
	"github.com/gu-io/gu/trees/elems"
)

// JSLink defines a struct to hold the Fetcher for "js-link" hook types.
// This hook embeds the css link as a link tag into the dom.
type JSLink struct{}

// Fetch returns the markup for the giving resource and where it should be inserted
// into the dom.
func (JSLink) Fetch(fetch Fetch, attr ManifestAttr) (*trees.Markup, bool, error) {
	script := trees.NewMarkup("script", false)
	trees.NewAttr("src", attr.Path).Apply(script)
	trees.NewAttr("type", "text/javascript").Apply(script)

	return script, true, nil
}

// CSSLink defines a struct to hold the Fetcher for "css-link" hook types.
// This hook embeds the css link as a link tag into the dom.
type CSSLink struct{}

// Fetch returns the markup for the giving resource and where it should be inserted
// into the dom.
func (CSSLink) Fetch(fetch Fetch, attr ManifestAttr) (*trees.Markup, bool, error) {
	style := trees.NewMarkup("link", false)
	trees.NewAttr("href", attr.Path).Apply(style)
	trees.NewAttr("ref", "stylesheet").Apply(style)
	trees.NewAttr("type", "text/css").Apply(style)

	return style, true, nil
}

// ImageLink defines a struct to hold the Fetcher for "css-link" hook types.
// This hook embeds the css link as a link tag into the dom.
type ImageLink struct{}

// Fetch returns the markup for the giving resource and where it should be inserted
// into the dom.
func (ImageLink) Fetch(fetch Fetch, attr ManifestAttr) (*trees.Markup, bool, error) {
	img := trees.NewMarkup("image", false)
	trees.NewAttr("src", attr.Path).Apply(img)

	return img, true, nil
}

//==============================================================================

var imageStyle = "data:image/%s;base64,%s"

// ImageEmbed defines a struct to hold the Fetcher for "img-css-embed" hook types.
// This hook retrieves the giving image data, transformed into
type ImageEmbed struct{}

// Fetch returns the markup for the giving resource and where it should be inserted
// into the dom.
func (m ImageEmbed) Fetch(fetch Fetch, attr ManifestAttr) (*trees.Markup, bool, error) {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(attr.Name), "."))

	img := trees.NewMarkup("img", false)

	if attr.Content != "" {
		decoded, err := attr.EncodeContentBase64()
		if err != nil {
			return nil, false, err
		}

		trees.NewAttr("src", fmt.Sprintf(imageStyle, ext, decoded)).Apply(img)
		return img, false, nil
	}

	res, err := fetch.Do(attr.WebRequest())
	if err != nil {
		return nil, false, err
	}

	decoded, err := res.EncodeContentBase64()
	if err != nil {
		return nil, false, err
	}

	trees.NewAttr("src", fmt.Sprintf(imageStyle, ext, decoded)).Apply(img)

	return img, false, nil
}

// ImageCSSEmbed defines a struct to hold the Fetcher for "img-css-embed" hook types.
// This hook retrieves the giving image data, transformed into
type ImageCSSEmbed struct{}

var cssImageStyle = `
	.%s-img{
		background: url('data:image/%s;base64,%s');
	}
`

// Fetch returns the markup for the giving resource and where it should be inserted
// into the dom.
func (m ImageCSSEmbed) Fetch(fetch Fetch, attr ManifestAttr) (*trees.Markup, bool, error) {
	if attr.Content != "" {
		style := trees.NewMarkup("style", false)

		ext := strings.TrimPrefix(filepath.Ext(attr.Name), ".")

		decoded, err := attr.EncodeContentBase64()
		if err != nil {
			return nil, false, err
		}

		elems.Text(cssImageStyle, getNewName(attr.Name), ext, decoded).Apply(style)

		return style, true, nil
	}

	res, err := fetch.Do(attr.WebRequest())
	if err != nil {
		return nil, false, err
	}

	style := trees.NewMarkup("style", false)
	ext := strings.TrimPrefix(filepath.Ext(attr.Name), ".")

	decoded, err := res.EncodeContentBase64()
	if err != nil {
		return nil, false, err
	}

	elems.Text(cssImageStyle, getNewName(attr.Name), ext, decoded).Apply(style)

	return style, true, nil
}

// getNewName returns an appropriate css name from the provided string.
func getNewName(val string) string {
	newVal := strings.Replace(val, filepath.Ext(val), "", 1)
	return strings.Replace(newVal, ".", "-", -1)
}

// CSSEmbed defines a struct to hold the Fetcher for "css-embed" hook types.
type CSSEmbed struct{}

// Fetch returns the markup for the giving resource and where it should be inserted
// into the dom.
func (CSSEmbed) Fetch(fetch Fetch, attr ManifestAttr) (*trees.Markup, bool, error) {
	if attr.Content != "" {
		style := trees.NewMarkup("style", false)

		decoded, err := attr.UnwrapBody()
		if err != nil {
			return nil, false, err
		}

		trees.NewText(string(decoded)).Apply(style)

		return style, true, nil
	}

	res, err := fetch.Do(attr.WebRequest())
	if err != nil {
		return nil, false, err
	}

	style := trees.NewMarkup("style", false)

	decoded, err := res.UnwrapBody()
	if err != nil {
		return nil, false, err
	}

	trees.NewText(string(decoded)).Apply(style)

	return style, true, nil
}

// JavascriptEmbed defines a struct to hold the Fetcher for "js-embed" hook types.
type JavascriptEmbed struct{}

// Fetch returns the markup for the giving resource and where it should be inserted
// into the dom.
func (JavascriptEmbed) Fetch(fetch Fetch, attr ManifestAttr) (*trees.Markup, bool, error) {
	if attr.Content != "" {
		style := trees.NewMarkup("script", false)

		decoded, err := attr.UnwrapBody()
		if err != nil {
			return nil, false, err
		}

		trees.NewText(string(decoded)).Apply(style)

		return style, true, nil
	}

	res, err := fetch.Do(attr.WebRequest())
	if err != nil {
		return nil, false, err
	}

	style := trees.NewMarkup("script", false)

	decoded, err := res.UnwrapBody()
	if err != nil {
		return nil, false, err
	}

	trees.NewText(string(decoded)).Apply(style)

	return style, true, nil
}

func init() {
	// Register CSS hooks.
	Register("css", CSSLink{})
	Register("css-embed", CSSEmbed{})

	// Register Javascript hooks.
	Register("js", JSLink{})
	Register("js-embed", JavascriptEmbed{})

	// Register Images hooks.
	Register("img", ImageCSSEmbed{})
	Register("img-embed", ImageEmbed{})
	Register("img-css-embed", ImageCSSEmbed{})
}
