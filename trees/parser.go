package trees

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

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

type counter interface {
	Next() int
}

type cn struct {
	ml  sync.Mutex
	val int
}

func (c *cn) Next() int {
	c.ml.Lock()
	c.val++
	c.ml.Unlock()
	return c.val
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

// ParseTreeToText takes a string markup and returns a *Markup which
// contains the full structure transpiled
// into the gutrees markup block structure.
func ParseTreeToText(markup string, withReturns bool) (io.WriterTo, error) {
	reader := strings.NewReader(markup)
	document, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer

	nameCounter := &cn{}

	doc := document.FirstChild

	if doc.FirstChild != nil && doc.FirstChild.FirstChild == nil {
		body := doc.LastChild

		if body.FirstChild == nil {
			return nil, errors.New("Body has no content nodes")
		}

		writeText(&buffer, "root := trees.NewMarkup(%q, %t)\n", "div", false)
		writeNodeSet(&buffer, body, "root", nameCounter)

		if withReturns {
			writeText(&buffer, `
			if len(root.Children()) == 1 {
				return root.Children()[0]
			}

			return root
		`)
		}
	} else {
		writeNode(&buffer, doc, "nil", "htmlNode")

		writeNodeSet(&buffer, doc, "root", nameCounter)

		if withReturns {
			writeText(&buffer, `
			return htmlNode
		`)
		}
	}

	return &buffer, nil
}

func writeNodeSet(w io.Writer, node *html.Node, parent string, count counter) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		elementName := fmt.Sprintf("elem%d", count.Next())

		writeNode(w, c, parent, elementName)

		writeNodeSet(w, c, elementName, count)
	}
}

func writeNode(w io.Writer, node *html.Node, parent string, elementName string) {
	switch node.Type {
	case html.ErrorNode:
		return
	case html.CommentNode:
		writeText(w, "trees.NewText(\"<---%s --->\").Apply(%s)", node.Data, parent)
		return
	case html.DoctypeNode:
		return
	case html.DocumentNode:
		return
	case html.ElementNode:
		writeText(w, "%s := trees.NewMarkup(%q, %t)\n%s.Apply(%s)", elementName, node.Data, false, elementName, parent)

		for _, attr := range node.Attr {
			if attr.Namespace != "" {
				writeText(w, "trees.NewAttr(\"%s:%s\", %q).Apply(%s)", attr.Namespace, attr.Key, attr.Val)
				continue
			}

			writeText(w, "trees.NewAttr(%+q, %q).Apply(%s)", attr.Key, attr.Val, elementName)
		}

		return
	case html.TextNode:
		text := strings.TrimSpace(node.Data)
		if text == "" {
			return
		}

		writeText(w, "trees.NewText(%+q).Apply(%s)", text, parent)
		return
	}
}

func writeText(w io.Writer, text string, vals ...interface{}) {
	fmt.Fprintf(w, text+"\n", vals...)
}
