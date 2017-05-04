// +build ignore
// Credit to Richard Musiol (https://github.com/neelance/dom)
// His code was crafted to fit gu's use

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var elemNameMap = map[string]string{
	"g":                "Group",
	"font-face-face":   "Fontface",
	"font-face-format": "FontFaceFormat",
	"font-face-name":   "FontfaceName",
	"font-face-src":    "FontFaceSrc",
	"font-face-uri":    "FontfaceURI",
	"missing-glyph":    "MissingGlyph",
	"a":                "Anchor",
	"article":          "Article",
	"aside":            "Aside",
	"area":             "Area",
	"abbr":             "Abbreviation",
	"b":                "Bold",
	"base":             "Base",
	"bdi":              "BidirectionalIsolation",
	"bdo":              "BidirectionalOverride",
	"blockquote":       "BlockQuote",
	"br":               "Break",
	"cite":             "Citation",
	"col":              "Column",
	"colgroup":         "ColumnGroup",
	"datalist":         "DataList",
	"dialog":           "Dialog",
	"details":          "Details",
	"dd":               "Description",
	"del":              "DeletedText",
	"dfn":              "Definition",
	"Def":              "Definition",
	"dl":               "DescriptionList",
	"dt":               "DefinitionTerm",
	"G":                "Group",
	"em":               "Emphasis",
	"embed":            "Embed",
	"footer":           "Footer",
	"figure":           "Figure",
	"figcaption":       "FigureCaption",
	"fieldset":         "FieldSet",
	"h1":               "Header1",
	"h2":               "Header2",
	"h3":               "Header3",
	"h4":               "Header4",
	"h5":               "Header5",
	"h6":               "Header6",
	"hgroup":           "HeadingsGroup",
	"header":           "Header",
	"hr":               "HorizontalRule",
	"i":                "Italic",
	"iframe":           "InlineFrame",
	"img":              "Image",
	"ins":              "InsertedText",
	"kbd":              "KeyboardInput",
	"keygen":           "KeyGen",
	"li":               "ListItem",
	"meta":             "Meta",
	"menuitem":         "MenuItem",
	"nav":              "Navigation",
	"noframes":         "NoFrames",
	"noscript":         "NoScript",
	"ol":               "OrderedList",
	"option":           "Option",
	"optgroup":         "OptionsGroup",
	"p":                "Paragraph",
	"param":            "Parameter",
	"pre":              "Preformatted",
	"q":                "Quote",
	"rp":               "RubyParenthesis",
	"Ref":              "Reference",
	"rt":               "RubyText",
	"s":                "Strikethrough",
	"samp":             "Sample",
	"source":           "Source",
	"section":          "Section",
	"sub":              "Subscript",
	"sup":              "Superscript",
	"tbody":            "TableBody",
	"textarea":         "TextArea",
	"td":               "TableData",
	"tfoot":            "TableFoot",
	"th":               "TableHeader",
	"thead":            "TableHead",
	"tr":               "TableRow",
	"u":                "Underline",
	"ul":               "UnorderedList",
	"var":              "Variable",
	"track":            "Track",
	"wbr":              "WordBreakOpportunity",
}

//list of self closing tags
var autoclosers = map[string]bool{
	"use":     true,
	"area":    true,
	"base":    true,
	"col":     true,
	"command": true,
	"embed":   true,
	"hr":      true,
	"input":   true,
	"keygen":  true,
	"meta":    true,
	"param":   true,
	"source":  true,
	"track":   true,
	"wbr":     true,
	"br":      true,
}

func pullDoc(url string, fx func(doc *goquery.Document)) error {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}

	fx(doc)
	return nil
}

func main() {
	file, err := os.Create("elems.gen.go")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	fmt.Fprint(file, `// Package elems contains definition for different html element types and some custom node types

//go:generate go run generate.go

// Documentation source: "HTML element reference" by Mozilla Contributors, https://developer.mozilla.org/en-US/docs/Web/HTML/Element, licensed under CC-BY-SA 2.5.
// Documentation for custom types lies within the  "github.com/influx6/faux/domtrees" package

package elems

import (
	"strings"

	"github.com/gu-io/gu/trees"
	"github.com/gu-io/gu/trees/css"
	"github.com/russross/blackfriday"
)

// SpaceCharacter provides text markup which contains the '&nbsp' text for
// a space element.
func SpaceCharacter(count int) *trees.Markup {
	if count < 1 {
		count = 0
	}

	var spaces []string

	for i := 0; i < count; i++ {
		spaces = append(spaces, "&nbsp;")
	}

	return trees.NewText(strings.Join(spaces, ""))
}

// Markdown takes the giving string which contains markdown written contents 
// and parses to html, which then is used to generate a new markup.
func Markdown(md string) *trees.Markup {
	hml := blackfriday.MarkdownCommon([]byte(md))
	return Parse(string(hml))
}

// MarkdownWithTemplate takes the giving template string which contains markdown  
// go template format strings. These will be parsed with text/template and using 
// blackfriday to parse the final output to html, which then is used to 
// generate a new markup. Any error is returned as a <error> tag markup.
func MarkdownWithTemplate(md string, bind interface{}) *trees.Markup {
	return trees.MarkdownTemplate(md, bind)
}

// CustomElement defines a type which returns a custom element type provided by
// the tagname.
func CustomElement(tag string, markup ...trees.Appliable) *trees.Markup {
	e := trees.NewMarkup(tag,false)
	trees.NewCSSStyle("display", "block").Apply(e)

	for _, m := range markup {
		if m == nil { continue }
		m.Apply(e)
	}
	return e
}

// Text provides custom type for defining text nodes with the trees markup.
func Text(content string, dl ...interface{}) *trees.Markup {
	return trees.NewText(content, dl...)
}

// ParseTemplate returns the giving markup structure generated from the string
// through the template used with the binding provided.
func ParseTemplate(markup string, bind interface{}, ms ...trees.Appliable) *trees.Markup {
	tms := trees.ParseTemplate(markup, bind)
	if len(tms) > 1 {
		sec := trees.NewMarkup("section", false)
		for _, el := range tms {
			el.Apply(sec)
		}

		for _, m := range ms {
			if m == nil {
				continue
			}

			m.Apply(sec)
		}

		return sec
	}

	root := tms[0]

	for _, m := range ms {
		if m == nil {
			continue
		}
		m.Apply(root)
	}

	return root
}

// Parse returns the giving markup structure generated from the string.
func Parse(markup string, ms ...trees.Appliable) *trees.Markup {
	tms := trees.ParseTree(markup)
	if len(tms) > 1 {
		sec := trees.NewMarkup("section",false)
		for _, el := range tms {
			el.Apply(sec)
		}

		for _, m := range ms {
			if m == nil { continue }
			m.Apply(sec)
		}

		return sec
	}

	root := tms[0]

	for _, m := range ms {
		if m == nil { continue }
		m.Apply(root)
	}

	return root
}

// ParseIn returns the giving markup structure generated from the string.
func ParseIn(root string,markup string, mo ...trees.Appliable) *trees.Markup {
	mroot := trees.ParseAsRoot(root, markup)

	for _, m := range mo {
		if m == nil { continue }
		m.Apply(mroot)
	}

	return mroot
}

// CSS provides a function that takes style rules which returns a stylesheet embeded into
// the provided element parent and is built on the gu/css package which collects
// necessary details from its parent to only target where it gets mounted.
func CSS(styles interface{}, bind interface{}, ext *css.Rule) *trees.Markup {
	return trees.CSSStylesheet(styles, bind, ext)
}
`)

	code := regexp.MustCompile("</?code>")
	unwanted := regexp.MustCompile("[^\\w\\d-]+")

	doneSvg := make(map[string]bool)
	err = pullDoc("https://developer.mozilla.org/en-US/docs/Web/SVG/Element", func(doc *goquery.Document) {
		doc.Find(".index ul li a").Each(func(i int, s *goquery.Selection) {
			link, _ := s.Attr("href")

			if !strings.HasPrefix(link, "/en-US/docs/Web/SVG/Element/") {
				return
			}

			if s.Parent().Find(".icon-trash, .icon-thumbs-down-alt, .icon-warning-sign").Length() > 0 {
				return
			}

			desc, _ := s.Attr("title")

			text := code.ReplaceAllString(s.Text(), "")

			name := text[1 : len(text)-1]

			// for key, item := range elemNameMap {
			// 	if strings.HasPrefix(name, key) || strings.HasSuffix(name, key) {
			// 		name = strings.Replace(name, key, item, 1)
			// 	}
			// }

			if doneSvg[name] || unwanted.MatchString(name) {
				return
			}

			writeSVGElem(file, name, desc, link)
			doneSvg[name] = true
		})
	})

	if err != nil {
		log.Fatalf("Unable to pull SVG ELEMENTS: %s", err)
	}

	doneHtml := make(map[string]bool)
	err = pullDoc("https://developer.mozilla.org/en-US/docs/Web/HTML/Element", func(doc *goquery.Document) {
		doc.Find(".quick-links a").Each(func(i int, s *goquery.Selection) {
			link, _ := s.Attr("href")
			if !strings.HasPrefix(link, "/en-US/docs/Web/HTML/Element/") {
				return
			}

			if s.Parent().Find(".icon-trash, .icon-thumbs-down-alt, .icon-warning-sign").Length() > 0 {
				return
			}

			desc, _ := s.Attr("title")
			text := s.Text()

			if text == "Heading elements" || text == "<h1>–<h6>" {
				// fmt.Printf("Write with %q\n", text)
				writeElem(file, "h1", desc, link)
				writeElem(file, "h2", desc, link)
				writeElem(file, "h3", desc, link)
				writeElem(file, "h4", desc, link)
				writeElem(file, "h5", desc, link)
				writeElem(file, "h6", desc, link)
				return
			}

			name := text[1 : len(text)-1]

			if name == "html" || name == "head" || name == "body" || unwanted.MatchString(name) {
				return
			}

			if doneHtml[name] {
				return
			}

			writeElem(file, name, desc, link)
			doneHtml[name] = true
		})
	})

	if err != nil {
		log.Fatalf("Unable to pull HTML ELEMENTS: %s", err)
	}

}

var badsymbs = regexp.MustCompile("-(.+)")

func writeSVGElem(w io.Writer, name, desc, link string) {
	var autocloser = autoclosers[name]
	funName := elemNameMap[name]
	// fmt.Printf("Header with %q\n", funName)

	funName = restruct(funName)
	// fmt.Printf("Header2 with %q\n", funName)

	if funName == "" {
		funName = restruct(name)

		for badsymbs.MatchString(funName) {
			if simbs := badsymbs.FindStringSubmatch(funName); len(simbs) > 0 {
				item := capitalize(simbs[1])
				funName = badsymbs.ReplaceAllString(funName, item)
			}
		}

		funName = capitalize(funName)
	}

	if funName != "Svg" {
		funName = "Svg" + funName
	}

	fmt.Fprintf(w, `
// %s provides the following for SVG XML elements ->
// %s
// https://developer.mozilla.org%s
func %s(markup ...trees.Appliable) *trees.Markup {
	e := trees.NewMarkup("%s",%t)
	for _, m := range markup {
		if m == nil { continue }
		m.Apply(e)
	}
	return e
}
`, funName, desc, link, funName, name, autocloser)
}

func writeElem(w io.Writer, name, desc, link string) {
	var autocloser = autoclosers[name]
	funName := elemNameMap[name]
	// fmt.Printf("FuncName with %q\n", name)

	funName = restruct(funName)
	// fmt.Printf("FuncName2 with %q\n", funName)

	if funName == "" {
		funName = restruct(name)

		for badsymbs.MatchString(funName) {
			if simbs := badsymbs.FindStringSubmatch(funName); len(simbs) > 0 {
				item := capitalize(simbs[1])
				funName = badsymbs.ReplaceAllString(funName, item)
			}
		}

		funName = capitalize(funName)
	}

	fmt.Fprintf(w, `
// %s provides the following for HTML elements ->
// %s
// https://developer.mozilla.org%s
func %s(markup ...trees.Appliable) *trees.Markup {
	e := trees.NewMarkup("%s",%t)
	for _, m := range markup {
		if m == nil { continue }
		m.Apply(e)
	}
	return e
}
`, funName, desc, link, funName, name, autocloser)
}

// capitalize capitalizes the first character in a string
func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func restruct(s string) string {
	if strings.Contains(s, "-") {
		mo := strings.Split(s, "-")
		for index, mi := range mo {
			if index == 0 {
				continue
			}

			mo[index] = capitalize(mi)
		}

		return strings.Join(mo, "")
	}

	return s
}
