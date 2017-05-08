package trees

import (
	"strings"

	"bytes"
	"text/template"

	"golang.org/x/net/html"
)

// ParseTemplateInto parses the provided string has a template which
// is processed with the provided binding and passed into the root.
func ParseTemplateInto(root *Markup, markup string, binding interface{}) {
	var bu bytes.Buffer

	tmpl := template.Must(template.New("Parsed").Parse(markup))
	if err := tmpl.Execute(&bu, binding); err != nil {
		return
	}

	ParseToRoot(root, bu.String())
}

// ParseTemplate parses the provided string has a template which
// is processed with the provided binding.
func ParseTemplate(markup string, binding interface{}) []*Markup {
	var bu bytes.Buffer

	tmpl := template.Must(template.New("Parsed").Parse(markup))
	if err := tmpl.Execute(&bu, binding); err != nil {
		return nil
	}

	return ParseTree(bu.String())
}

// ParseFirstOrMakeRoot attempts to parse the giving markup and returns the
// element if only one else creates a div and adds all children as part of div.
func ParseFirstOrMakeRoot(markup string) *Markup {
	trees := ParseTree(markup)
	if len(trees) == 1 {
		return trees[0]
	}

	root := NewMarkup("div", false)
	root.AddChild(trees...)
	return root
}

// ParseToRoot passes the markup generated from the markup added to the provided
// root.
func ParseToRoot(root *Markup, markup string) {
	trees := ParseTree(markup)
	for _, child := range trees {
		child.Apply(root)
	}
}

// ParseAndFirst expects the markup provided to only have one root element which
// will be returned.
func ParseAndFirst(markup string) *Markup {
	trees := ParseTree(markup)
	if len(trees) != 1 {
		panic("Markup must only returned single item in tree")
	}

	return trees[0]
}

// ParseAsRoot returns the markup generated from the provided markup,
// returning them as children of the provided root.
func ParseAsRoot(root string, markup string) *Markup {
	tokens := html.NewTokenizer(strings.NewReader(markup))

	var sel *Selector
	if sels := Query.ParseSelector(root); sels != nil {
		sel = sels[0]
	} else {
		sel.Tag = root
	}

	rootElem := NewMarkup(sel.Tag, false)

	if sel.ID != "" {
		NewAttr("id", sel.ID).Apply(rootElem)
	}

	if sel.Classes != nil {
		(&ClassList{list: sel.Classes}).Apply(rootElem)
	}

	pullNode(tokens, rootElem)

	return rootElem
}

// ParseTree takes a string markup and returns a *Markup which
// contains the full structure transpiled
// into the gutrees markup block structure.
func ParseTree(markup string) []*Markup {
	tokens := html.NewTokenizer(strings.NewReader(markup))

	rootElem := NewMarkup("div", false)
	pullNode(tokens, rootElem)

	return rootElem.Children()
}

func pullNode(tokens *html.Tokenizer, root *Markup) {
	for {
		token := tokens.Next()

		switch token {
		case html.ErrorToken:
			return

		case html.TextToken, html.CommentToken, html.DoctypeToken:
			text := strings.TrimSpace(string(tokens.Text()))

			if text == "" {
				continue
			}

			if token == html.CommentToken {
				text = "<!--" + text + "-->"
			}

			// if node != nil {
			// 	NewText(text).Apply(node)
			// 	continue
			// }

			NewText(text).Apply(root)
			continue

		case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
			tagName, hasAttr := tokens.TagName()

			// fmt.Printf("Token: %#v -> %+q -> %q -> %t\n", token, token, tagName, token == html.SelfClosingTagToken)

			if token == html.EndTagToken && string(tagName) == root.tagname {
				return
			}

			node := NewMarkup(string(tagName), token == html.SelfClosingTagToken)
			node.Apply(root)

			if hasAttr {
			attrLoop:
				for {
					key, val, more := tokens.TagAttr()

					if string(key) != "" {
						NewAttr(string(key), string(val)).Apply(node)
					}

					if !more {
						break attrLoop
					}
				}
			}

			if token == html.SelfClosingTagToken {
				continue
			}

			pullNode(tokens, node)
		}
	}
}
