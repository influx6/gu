package gu

import (
	"bytes"
	"fmt"
	"strings"

	"html/template"

	"github.com/gu-io/gu/trees"
	"golang.org/x/net/html"
)

// DefaultComponentMakers provides a set of default components makers which can be
// readily used in creating markup.
var DefaultComponentMakers = []ComponentItem{
	{
		TagName: "css",
		Unwrap:  true,
		Maker: func(fields map[string]string, template string) Renderable {
			return Static(trees.CSSStylesheet(template, fields))
		},
	},
}

// ComponentMaker defines a function type which returns a Renderable based on
// a series of provied attributes.
type ComponentMaker func(fields map[string]string, template string) Renderable

// ComponentItem defines a struct which contains the tagName and maker corresponding
// to generating the giving tagName.
type ComponentItem struct {
	TagName string
	Maker   ComponentMaker
	Unwrap  bool
}

// ComponentRegistry defines a struct to manage all registered Component makers.
type ComponentRegistry struct {
	makers map[string]ComponentItem
}

// NewComponentRegistry returns a new instance of a ComponentRegistry.
func NewComponentRegistry() *ComponentRegistry {
	registry := &ComponentRegistry{
		makers: make(map[string]ComponentItem),
	}

	registry.Add(DefaultComponentMakers)

	return registry
}

// Generate returns a the component attribute with the provided markup as it's based.
// It provides a full complete set of all items in the list.
func (c *ComponentRegistry) Generate(markup string, attr ComponentAttr) ComponentAttr {
	attr.Base = ParseComponent(markup, c)
	return attr
}

// MustParseByTemplate returns a new Renderable from using text template to parse the provided
// markup.
func (c *ComponentRegistry) MustParseByTemplate(markup string, m interface{}) Renderable {
	renderable, err := c.ParseByTemplate(markup, m)
	if err != nil {
		panic(err)
	}
	return renderable
}

// ParseByTemplate returns a new Renderable from using text template to parse the provided
// markup.
func (c *ComponentRegistry) ParseByTemplate(markup string, m interface{}) (Renderable, error) {
	tmp, err := template.New("css").Parse(markup)
	if err != nil {
		return nil, err
	}

	var content bytes.Buffer
	if err := tmp.Execute(&content, m); err != nil {
		return nil, err
	}

	return ParseComponent(content.String(), c), nil
}

// Parse returns a new Renderable from the giving markup.
func (c *ComponentRegistry) Parse(markup string, m ...interface{}) Renderable {
	return ParseComponent(fmt.Sprintf(markup, m...), c)
}

// Add adds the giving set of possible item/items of the Acceptable type into
// the registry.
func (c *ComponentRegistry) Add(item interface{}) {
	switch realItem := item.(type) {
	case ComponentItem:
		realItem.TagName = strings.ToLower(realItem.TagName)
		c.makers[realItem.TagName] = realItem
	case []ComponentItem:
		for _, item := range realItem {
			item.TagName = strings.ToLower(item.TagName)
			c.makers[item.TagName] = item
		}
	case map[string]ComponentMaker:
		for name, item := range realItem {
			c.Register(name, item, false)
		}
	}
}

// Has returns true/false if the giving tagName exists in the registry.
func (c *ComponentRegistry) Has(tag string) bool {
	tag = strings.ToLower(tag)
	_, ok := c.makers[tag]
	return ok
}

// ParseTag returns the giving Renderable for the giving markup.
func (c *ComponentRegistry) ParseTag(tag string, fields map[string]string, template string) (Renderable, bool) {
	tag = strings.ToLower(tag)
	cm, ok := c.makers[tag]
	if !ok {
		return nil, false
	}

	return cm.Maker(fields, template), cm.Unwrap
}

// Register adds the giving tagName which will be used to create a new Renderable
// when found in the markup provided.
// Arguments:
//  TagName - the tagname of the component to be searched for.
//  Maker - the function which generates the Renderable
//  Unwrap - Boolean indicating if the content alone to be rendered or else be wrapped
//  					by the tag declared.
func (c *ComponentRegistry) Register(tagName string, maker ComponentMaker, unwrap bool) bool {
	c.makers[strings.ToLower(tagName)] = ComponentItem{
		TagName: tagName,
		Maker:   maker,
		Unwrap:  unwrap,
	}

	return true
}

//================================================================================

type attr struct {
	Name string
	Val  string
}

// Treeset defines a structure which builds it's markup from a set of internal
// structures which returns a full complete markup including Components.
type Treeset struct {
	Registry      *ComponentRegistry
	DeferedTag    string
	DeferTemplate string
	Renderable    Renderable
	Attr          []attr
	Children      []*Treeset
	Tree          *trees.Markup
	Fields        map[string]string
}

// initialize any internal defered tag and template.
func (t *Treeset) init() {
	if t.DeferedTag != "" && t.Renderable == nil {
		res, unwrap := t.Registry.ParseTag(t.DeferedTag, t.Fields, t.DeferTemplate)

		t.Renderable = res
		t.Registry = nil
		t.Fields = nil
		t.DeferedTag = ""
		t.DeferTemplate = ""

		if unwrap {
			t.Tree = nil
		}
	}
}

// Render returns a new markup instance for the Treeset.
func (t *Treeset) Render() *trees.Markup {
	// intialize internal renderable if required.
	t.init()

	var base *trees.Markup

	if t.Tree != nil {
		base = t.Tree.Clone()

		if t.Renderable != nil {
			t.Renderable.Render().Apply(base)
		}

		base.UpdateHash()
	} else {
		base = t.Renderable.Render()
	}

	for _, attr := range t.Attr {
		trees.NewAttr(attr.Name, attr.Val).Apply(base)
		base.UpdateHash()
	}

	for _, child := range t.Children {
		child.Render().Apply(base)
		base.UpdateHash()
	}

	return base
}

// ParseComponent returns a new Renderable from the giving markup generating the
// Renderable.
func ParseComponent(markup string, registry *ComponentRegistry) Renderable {
	var tree Treeset

	tokens := html.NewTokenizer(strings.NewReader(markup))
	parseTokens(tokens, &tree, registry)

	if len(tree.Children) == 1 {
		return tree.Children[0]
	}

	tree.Tree = trees.NewMarkup("section", false)
	return &tree
}

// parseTokens generates a new heirarchy of Treeset using the provided registery.
// Creating appropriate markup necessary to build.
func parseTokens(tokens *html.Tokenizer, parent *Treeset, registery *ComponentRegistry) {
	var templateEnable bool
	var template []string

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

			if templateEnable {
				template = append(template, text)
				continue
			}

			if token == html.CommentToken {
				text = "<!--" + text + "-->"
			}

			parent.Children = append(parent.Children, &Treeset{
				Tree: trees.NewText(text),
			})

			continue

		case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
			if token == html.EndTagToken {

				if templateEnable {
					parent.DeferTemplate = strings.Join(template, "")
					template = nil
					templateEnable = false
				}

				return
			}

			tagName, hasAttr := tokens.TagName()

			tag := strings.ToLower(string(tagName))
			if string(tag) == "root-template" {
				templateEnable = true
				continue
			}

			var unwrap bool

			attrs := make([]attr, 0)
			fields := make(map[string]string)

			if hasAttr {
				{
				attrLoop:
					for {
						key, val, more := tokens.TagAttr()

						if string(key) == "unwrap" {
							unwrap = true
						}

						if string(key) != "" {
							keyName := strings.ToLower(string(key))
							if strings.HasPrefix(keyName, "component-") {
								fields[strings.TrimPrefix(keyName, "component-")] = string(val)
								continue
							}

							attrs = append(attrs, attr{
								Name: string(key),
								Val:  string(val),
							})
						}

						if !more {
							break attrLoop
						}
					}
				}
			}

			var set Treeset

			if registery.Has(tag) {
				var wrap *trees.Markup

				if !unwrap {
					wrap = trees.NewMarkup(string(tagName), token == html.SelfClosingTagToken)
				}

				set = Treeset{
					Attr:       attrs,
					Fields:     fields,
					DeferedTag: tag,
					Registry:   registery,
					Tree:       wrap,
				}
			} else {
				elem := trees.NewMarkup(string(tagName), token == html.SelfClosingTagToken)

				for _, attr := range attrs {
					trees.NewAttr(attr.Name, attr.Val).Apply(elem)
				}

				for name, val := range fields {
					trees.NewAttr(name, val).Apply(elem)
				}

				set = Treeset{
					Tree: elem,
				}
			}

			parent.Children = append(parent.Children, &set)

			if token == html.SelfClosingTagToken {
				continue
			}

			parseTokens(tokens, &set, registery)
		}
	}
}
