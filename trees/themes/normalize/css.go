// Package normalize defines a package which embeds all css files into a go file.
// This package is automatically generated and should not be modified by hand. 
// It provides a source which is used to build all css packages into a css.go 
// file which contains each allocated by name.

//go:generate go run generate.go

package normalize

import (
	"encoding/json"
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
  if err := json.Unmarshal([]byte("[\n\t{\n\t\t\"after\": null,\n\t\t\"before\": null,\n\t\t\"path\": \"normalize.css\",\n\t\t\"style\": \"    /*! normalize.css v6.0.0 | MIT License | github.com/necolas/normalize.css */\\n\\nhtml {\\n    line-height: 1.15;\\n    -ms-text-size-adjust: 100%;\\n    -webkit-text-size-adjust: 100%\\n}\\n\\narticle,\\naside,\\nfooter,\\nheader,\\nnav,\\nsection {\\n    display: block\\n}\\n\\nh1 {\\n    font-size: 2em;\\n    margin: .67em 0\\n}\\n\\nfigcaption,\\nfigure,\\nmain {\\n    display: block\\n}\\n\\nfigure {\\n    margin: 1em 40px\\n}\\n\\nhr {\\n    box-sizing: content-box;\\n    height: 0;\\n    overflow: visible\\n}\\n\\npre {\\n    font-family: monospace, monospace;\\n    font-size: 1em\\n}\\n\\na {\\n    background-color: transparent;\\n    -webkit-text-decoration-skip: objects; // Exposes a variable which contains the contents of the normalize css library.}abbr[title]{border-bottom:0;text-decoration:underline;text-decoration:underline dotted}b,strong{font-weight:inherit}b,strong{font-weight:bolder}code,kbd,samp{font-family:monospace,monospace;font-size:1em}dfn{font-style:italic}mark{background-color:#ff0;color:#000}small{font-size:80%}sub,sup{font-size:75%;line-height:0;position:relative;vertical-align:baseline}sub{bottom:-0.25em}sup{top:-0.5em}audio,video{display:inline-block}audio:not([controls]){display:none;height:0}img{border-style:none}svg:not(:root){overflow:hidden}button,input,optgroup,select,textarea{margin:0}button,input{overflow:visible}button,select{text-transform:none}button,html [type=\\\"button\\\"],[type=\\\"reset\\\"],[type=\\\"submit\\\"]{-webkit-appearance:button}button::-moz-focus-inner,[type=\\\"button\\\"]::-moz-focus-inner,[type=\\\"reset\\\"]::-moz-focus-inner,[type=\\\"submit\\\"]::-moz-focus-inner{border-style:none;padding:0}button:-moz-focusring,[type=\\\"button\\\"]:-moz-focusring,[type=\\\"reset\\\"]:-moz-focusring,[type=\\\"submit\\\"]:-moz-focusring{outline:1px dotted ButtonText}legend{box-sizing:border-box;color:inherit;display:table;max-width:100%;padding:0;white-space:normal}progress{display:inline-block;vertical-align:baseline}textarea{overflow:auto}[type=\\\"checkbox\\\"],[type=\\\"radio\\\"]{box-sizing:border-box;padding:0}[type=\\\"number\\\"]::-webkit-inner-spin-button,[type=\\\"number\\\"]::-webkit-outer-spin-button{height:auto}[type=\\\"search\\\"]{-webkit-appearance:textfield;outline-offset:-2px}[type=\\\"search\\\"]::-webkit-search-cancel-button,[type=\\\"search\\\"]::-webkit-search-decoration{-webkit-appearance:none}::-webkit-file-upload-button{-webkit-appearance:button;font:inherit}details,menu{display:block}summary{display:list-item}canvas{display:inline-block}template{display:none}[hidden]{display:none}\"\n\t}\n]"),&rules); err != nil {
  	fmt.Printf("Failed to unmarshal styles: %+q\n", err)
  }
}
