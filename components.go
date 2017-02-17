package gu

import (
	"strings"

	"github.com/gu-io/gu/trees"
	"golang.org/x/net/html"
)

// ComponentMaker defines a function type which returns a Renderable based on
// a series of provied attributes.
type ComponentMaker func(fields map[string]string, template string) Renderable

// componentItem defines a struct which contains the tagName and maker corresponding
// to generating the giving tagName.
type componentItem struct {
	TagName string
	Maker   ComponentMaker
}

// ComponentRegistry defines a struct to manage all registered Component makers.
type ComponentRegistry struct {
	makers map[string]componentItem
}

// NewComponentRegistry returns a new instance of a ComponentRegistry.
func NewComponentRegistry() *ComponentRegistry {
	return &ComponentRegistry{
		makers: make(map[string]componentItem),
	}
}

// Generate returns a the component attribute with the provided markup as it's based.
// It provides a full complete set of all items in the list.
func (c *ComponentRegistry) Generate(markup string, attr ComponentAttr) ComponentAttr {
	attr.Base = ParseComponent(markup, c)
	return attr
}

// Parse returns a new Renderable from the giving markup.
func (c *ComponentRegistry) Parse(markup string) Renderable {
	return ParseComponent(markup, c)
}

// Has returns true/false if the giving tagName exists in the registry.
func (c *ComponentRegistry) Has(tag string) bool {
	tag = strings.ToLower(tag)
	_, ok := c.makers[tag]
	return ok
}

// ParseTag returns the giving Renderable for the giving markup.
func (c *ComponentRegistry) ParseTag(tag string, fields map[string]string, template string) Renderable {
	tag = strings.ToLower(tag)
	cm, ok := c.makers[tag]
	if !ok {
		return nil
	}

	return cm.Maker(fields, template)
}

// Register adds the giving tagName which will be used to create a new Renderable
// when found in the markup provided.
func (c *ComponentRegistry) Register(tagName string, maker ComponentMaker) {
	c.makers[strings.ToLower(tagName)] = componentItem{
		TagName: tagName,
		Maker:   maker,
	}
}

//================================================================================

// Treeset defines a structure which builds it's markup from a set of internal
// structures which returns a full complete markup including Components.
type Treeset struct {
	Registry      *ComponentRegistry
	DeferedTag    string
	DeferTemplate string
	Renderable    Renderable
	Children      []*Treeset
	Tree          *trees.Markup
	Attr          map[string]string
	Fields        map[string]string
}

// initialize any internal defered tag and template.
func (t *Treeset) init() {
	if t.DeferedTag != "" && t.Renderable == nil {
		t.Renderable = t.Registry.ParseTag(t.DeferedTag, t.Fields, t.DeferTemplate)
		t.Registry = nil
		t.Fields = nil
		t.DeferedTag = ""
		t.DeferTemplate = ""
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

	for name, val := range t.Attr {
		trees.NewAttr(name, val).Apply(base)
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

			attrs := make(map[string]string)
			fields := make(map[string]string)

			if hasAttr {
				{
				attrLoop:
					for {
						key, val, more := tokens.TagAttr()

						if string(key) != "" {
							keyName := strings.ToLower(string(key))
							if strings.HasPrefix(keyName, "component-") {
								fields[strings.TrimPrefix(keyName, "component-")] = string(val)
								continue
							}

							attrs[string(key)] = string(val)
						}

						if !more {
							break attrLoop
						}
					}
				}
			}

			var set Treeset

			if registery.Has(tag) {
				set = Treeset{
					Attr:       attrs,
					Fields:     fields,
					DeferedTag: tag,
					Registry:   registery,
					Tree:       trees.NewMarkup(string(tagName), token == html.SelfClosingTagToken),
				}
			} else {
				elem := trees.NewMarkup(string(tagName), token == html.SelfClosingTagToken)

				for name, val := range attrs {
					trees.NewAttr(name, val).Apply(elem)
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
