// Package examples defines a package which embeds all css files into a go file.
// This package is automatically generated and should not be modified by hand. 
// It provides a source which is used to build all css packages into a css.go 
// file which contains each allocated by name.

//go:generate go run generate.go

package examples

import (
	"strings"
	"fmt"

	"github.com/gu-io/gu/trees/css"
)

var rules cssstyles

// Must returns the giving rules from the provided style rules else panics.
func Must(dir string) *css.Rule {
	if rule := Get(dir); rule != nil {
		return rule
	}

	panic(fmt.Sprintf("Rule %s not found", dir))
}

// GetSource returns the style contents of the stylesheet.
func GetSource(dir string) string {
	for _, item := range rules {
		if item.Path != dir {
			continue
		}

		return item.RuleSource(rules)
	}

	return ""
}

// Get returns the giving rules from the provided style rules.
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

// RuleSource returns a string containing the giving rules and it's dependencies.
func (s *cssstyle) RuleSource(root []cssstyle) string {
	var befores []string

	for _, before := range s.Before {
		befores = append(befores, root[before].RuleSource(root))
	}

	befores = append(befores, s.Style)

	for _, after := range s.After {
		befores = append(befores, root[after].RuleSource(root))
	}

	return strings.Join(befores, "\n")
}

// Rule retrieves the giving set of rules pertaining the giving style.
func (s *cssstyle) Rule(root []cssstyle) *css.Rule {
	var befores []*css.Rule

	for _, before := range s.Before {
		befores = append(befores, root[before].Rule(root))
	}

	self := css.New(s.Style, nil, befores...)

	for _, after := range s.After {
		self = (root[after]).Rule(root).Add(self)
	}

	return self
}

func init (){

	
	rules = append(rules, cssstyle{
		Style: ".base-component {\n\twidth: 100px;\n}",
		Path: "base/base.css",
		After: []int{},
		Before: []int{},
	})


	
	rules = append(rules, cssstyle{
		Style: ".base-ui{\n\tfont-size: 40px;\n}",
		Path: "base/ui/base-ui.css",
		After: []int{},
		Before: []int{},
	})


	
	rules = append(rules, cssstyle{
		Style: "/* #include ui/*:before, base/ui/base-ui.css:after, base/base.css */\n\n.examples {\n\twidth: 100px;\n\theight: 200px;\n}",
		Path: "examples.css",
		After: []int{1,},
		Before: []int{3,4,0,},
	})


	
	rules = append(rules, cssstyle{
		Style: ".ui-button {\n\tcolor: black;\n}",
		Path: "ui/button.css",
		After: []int{},
		Before: []int{},
	})


	
	rules = append(rules, cssstyle{
		Style: ".ui {\n\tfont-family: Lato, Helvetica, sans-serif;\n}",
		Path: "ui/ui.css",
		After: []int{},
		Before: []int{},
	})


}
