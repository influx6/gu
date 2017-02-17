package gu

// ComponentMaker defines a function type which returns a Renderable based on
// a series of provied attributes.
type ComponentMaker func(attrs map[string]string) Renderable

// ComponentRegistry defines a struct to manage all registered Component makers.
type ComponentRegistry struct {
	makers []componentItem
}

// componentItem defines a struct which contains the tagName and maker corresponding
// to generating the giving tagName.
type componentItem struct {
	TagName string
	Maker   ComponentMaker
}

// Generate returns a the component attribute with the provided markup as it's based.
// It provides a full complete set of all items in the list.
func (c *ComponentRegistry) Generate(markup string, attr ComponentAttr) ComponentAttr {
	return attr
}

// Register adds the giving tagName which will be used to create a new Renderable
// when found in the markup provided.
func (c *ComponentRegistry) Register(tagName string, maker ComponentMaker) {
	c.makers = append(c.makers, componentItem{
		TagName: tagName,
		Maker:   maker,
	})
}

// func parseTokens(tokens *html.Tokenizer, tree []trees.Applier) {
//
// 	for {
// 		token := tokens.Next()
//
// 		switch token {
// 		case html.ErrorToken:
// 			return
//
// 		case html.TextToken, html.CommentToken, html.DoctypeToken:
// 			text := strings.TrimSpace(string(tokens.Text()))
//
// 			if text == "" {
// 				continue
// 			}
//
// 			if token == html.CommentToken {
// 				text = "<!--" + text + "-->"
// 			}
//
// 			if node != nil {
// 				// NewText(text).Apply(node)
// 				continue
// 			}
//
// 			// NewText(text).Apply(root)
// 			continue
//
// 		case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
// 			if token == html.EndTagToken {
// 				return
// 			}
//
// 			tagName, hasAttr := tokens.TagName()
// 			attrs := make(map[string]string)
//
// 			// node = NewMarkup(string(tagName), token == html.SelfClosingTagToken)
// 			// node.Apply(root)
//
// 			if hasAttr {
// 			attrLoop:
// 				for {
// 					key, val, more := tokens.TagAttr()
// 					if !more {
// 						break attrLoop
// 					}
//
// 					attrs[string(key)] = string(val)
// 				}
// 			}
//
// 			if token == html.SelfClosingTagToken {
// 				continue
// 			}
//
// 			parseTokens(tokens, tree)
// 		}
// 	}
// }
