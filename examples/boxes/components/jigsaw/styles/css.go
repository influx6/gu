// Package styles defines a package which embeds all css files into a go file.
// This package is automatically generated and should not be modified by hand.
// It provides a source which is used to build all css packages into a css.go
// file which contains each allocated by name.

//go:generate go run generate.go

package styles

import (
	"encoding/json"
	"fmt"

	"github.com/gu-io/gu/trees/css"
)

var rules styles

// Get returns the giving rules from the provided
func Get(dir string) *css.Rule {
	var target *style

	for _, item := range styles {
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

type styles []style

// style defines a giving struct which contain the giving property style and dependencies.
type style struct {
	Rule   string `json:"rule"`
	Path   string `json:"path"`
	Before []int  `json:"before"`
	After  []int  `json:"after"`
}

// Rule retrieves the giving set of rules pertaining the giving style.
func (s *style) Rule(root []styles) *css.Rule {
	var befores []*css.Rule

	for _, before := range s.Before {
		befores = append(befores, root[before].Rule(root))
	}

	self := css.New(s.Rule, befores...)

	for _, after := range s.After {
		self = (root[after]).Rule(root).AddRoot(self)
	}

	return self
}

func init() {
	if err := json.Unmarshal(&rules, "null"); err != nil {
		fmt.Printf("Failed to unmarshal styles: %+q\n", err)
	}
}
