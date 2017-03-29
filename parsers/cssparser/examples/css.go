// Package examples defines a package which embeds all css files into a go file.
// This package is automatically generated and should not be modified by hand.
// It provides a source which is used to build all css packages into a css.go
// file which contains each allocated by name.

//go:generate go run generate.go

package examples

import (
	"encoding/json"
	"fmt"

	"github.com/gu-io/gu/trees/css"
)

var rules cssstyles

// Get returns the giving rules from the provided
func Get(dir string) *css.Rule {
	var target *cssstyle

	for _, item := range rules {
		if item.Path != dir {
			continue
		}

		target = &item
		break
	}

	if target == nil {
		return nil
	}

	return target.Rule(rules)
}

type cssstyles []cssstyle

// style defines a giving struct which contain the giving property style and dependencies.
type cssstyle struct {
	Style  string `json:"style"`
	Path   string `json:"path"`
	Before []int  `json:"before"`
	After  []int  `json:"after"`
}

// Rule retrieves the giving set of rules pertaining the giving style.
func (s *cssstyle) Rule(root []cssstyle) *css.Rule {
	var befores []*css.Rule

	for _, before := range s.Before {
		befores = append(befores, root[before].Rule(root))
	}

	self := css.New(s.Style, befores...)

	for _, after := range s.After {
		self = (root[after]).Rule(root).AddRoot(self)
	}

	return self
}

func init() {
	if err := json.Unmarshal([]byte("[\n\t{\n\t\t\"after\": null,\n\t\t\"before\": null,\n\t\t\"path\": \"base/base.css\",\n\t\t\"style\": \".base-component {\\n\\twidth: 100px;\\n}\"\n\t},\n\t{\n\t\t\"after\": null,\n\t\t\"before\": null,\n\t\t\"path\": \"base/ui/base-ui.css\",\n\t\t\"style\": \".base-ui{\\n\\tfont-size: 40px;\\n}\"\n\t},\n\t{\n\t\t\"after\": [\n\t\t\t1\n\t\t],\n\t\t\"before\": [\n\t\t\t3,\n\t\t\t4,\n\t\t\t0\n\t\t],\n\t\t\"path\": \"examples.css\",\n\t\t\"style\": \"/* #include ui/*:before, base/ui/base-ui.css:after, base/base.css */\\n\\n.examples {\\n\\twidth: 100px;\\n\\theight: 200px;\\n}\"\n\t},\n\t{\n\t\t\"after\": null,\n\t\t\"before\": null,\n\t\t\"path\": \"ui/button.css\",\n\t\t\"style\": \".ui-button {\\n\\tcolor: black;\\n}\"\n\t},\n\t{\n\t\t\"after\": null,\n\t\t\"before\": null,\n\t\t\"path\": \"ui/ui.css\",\n\t\t\"style\": \".ui {\\n\\tfont-family: Lato, Helvetica, sans-serif;\\n}\"\n\t}\n]"), &rules); err != nil {
		fmt.Printf("Failed to unmarshal styles: %+q\n", err)
	}
}
