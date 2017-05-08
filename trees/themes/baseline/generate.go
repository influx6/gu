// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"path/filepath"

	"github.com/gu-io/gu/parsers/cssparser"
)

var pkgName = "baseline"
var pkg = "// Package {{PKG}} defines a package which embeds all css files into a go file.\n// This package is automatically generated and should not be modified by hand. \n// It provides a source which is used to build all css packages into a css.go \n// file which contains each allocated by name.\n\n//go:generate go run generate.go\n\npackage {{PKG}}\n\nimport (\n\t\"encoding/json\"\n\t\"strings\"\n\t\"fmt\"\n\n\t\"github.com/gu-io/gu/trees/css\"\n)\n\nvar rules cssstyles\n\n// Must returns the giving rules from the provided style rules else panics.\nfunc Must(dir string) *css.Rule {\n\tif rule := Get(dir); rule != nil {\n\t\treturn rule\n\t}\n\n\tpanic(fmt.Sprintf(\"Rule %s not found\", dir))\n}\n\n// GetSource returns the style contents of the stylesheet.\nfunc GetSource(dir string) string {\n\tfor _, item := range rules {\n\t\tif item.Path != dir {\n\t\t\tcontinue\n\t\t}\n\n\t\treturn item.RuleSource(rules)\n\t}\n\n\treturn \"\"\n}\n\n// Get returns the giving rules from the provided style rules.\nfunc Get(dir string) *css.Rule {\n\tvar target *cssstyle\n\n\tfor _, item := range rules {\n\t\tif item.Path != dir {\n\t\t\tcontinue\n\t\t}\n\n\t\ttarget = &item\n\t\tbreak\n\t}\n\n\tif target == nil {\n\t\treturn nil\n\t}\n\n\treturn target.Rule(rules)\n}\n\ntype cssstyles []cssstyle\n\n// style defines a giving struct which contain the giving property style and dependencies.\ntype cssstyle struct {\n\tStyle  string `json:\"style\"`\n\tPath   string `json:\"path\"`\n\tBefore []int  `json:\"before\"`\n\tAfter  []int  `json:\"after\"`\n}\n\n// RuleSource returns a string containing the giving rules and it's dependencies.\nfunc (s *cssstyle) RuleSource(root []cssstyle) string {\n\tvar befores []string\n\n\tfor _, before := range s.Before {\n\t\tbefores = append(befores, root[before].RuleSource(root))\n\t}\n\n\tbefores = append(befores, s.Style)\n\n\tfor _, after := range s.After {\n\t\tbefores = append(befores, root[after].RuleSource(root))\n\t}\n\n\treturn strings.Join(befores, \"\\n\")\n}\n\n// Rule retrieves the giving set of rules pertaining the giving style.\nfunc (s *cssstyle) Rule(root []cssstyle) *css.Rule {\n\tvar befores []*css.Rule\n\n\tfor _, before := range s.Before {\n\t\tbefores = append(befores, root[before].Rule(root))\n\t}\n\n\tself := css.New(s.Style, nil, befores...)\n\n\tfor _, after := range s.After {\n\t\tself = (root[after]).Rule(root).Add(self)\n\t}\n\n\treturn self\n}\n\nfunc init (){\n  if err := json.Unmarshal([]byte({{STYLES}}),&rules); err != nil {\n  \tfmt.Printf(\"Failed to unmarshal styles: %+q\\n\", err)\n  }\n}\n"

func main() {
	cdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	items, err := cssparser.ParseDir(filepath.Join(cdir, "css"))
	if err != nil {
		panic("Failed to walk CSS directories: "+ err.Error())
	}

	rd, err := json.MarshalIndent(items.Generate(), "", "\t")
	if err != nil {
		panic("Failed to Marshal CSS File struct: "+ err.Error())
	}


	file, err := os.Create(filepath.Join(cdir, "css.go"))
	if err != nil {
		panic("Failed to create css pkg file: "+ err.Error())
	}

	defer file.Close()

	quoted := fmt.Sprintf("%+q", rd)
	pkg = strings.Replace(pkg,"{{PKG}}", pkgName, -1)
	pkg = strings.Replace(pkg,"{{STYLES}}", quoted, -1)

	file.Write([]byte(pkg))
}