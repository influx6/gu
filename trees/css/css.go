package css

import (
	"bytes"
	"strings"
	"sync"
	"text/template"

	bcss "github.com/aymerick/douceur/css"
	"github.com/aymerick/douceur/parser"
)

// Rule defines the a single css rule which will be transformed and
// converted into a usable stylesheet during rendering.
type Rule struct {
	template *template.Template
	rx       sync.RWMutex
	depends  []*Rule
}

// New returns a new instance of a Rule which provides capability to parse
// and extrapolate the giving content using the provided binding.
func New(rules string, rs ...*Rule) *Rule {
	tmp, err := template.New("css").Parse(rules)
	if err != nil {
		panic(err)
	}

	return &Rule{
		template: tmp,
		depends:  rs,
	}
}

// AddRoot adds the giving Rule as a style which this Rule depends on.
func (r *Rule) AddRoot(rx *Rule) *Rule {
	r.rx.Lock()
	{
		r.depends = append(r.depends, rx)
	}
	r.rx.Unlock()

	return r
}

// Stylesheet returns the provided styles using the binding as the argument for the
// provided css template.
func (r *Rule) Stylesheet(bind interface{}, parentNode string) (*bcss.Stylesheet, error) {
	var stylesheet bcss.Stylesheet

	r.rx.RLock()
	{
		for _, rule := range r.depends {
			sheet, err := rule.Stylesheet(bind, parentNode)
			if err != nil {
				return nil, err
			}

			stylesheet.Rules = append(stylesheet.Rules, sheet.Rules...)
		}
	}
	r.rx.RUnlock()

	var content bytes.Buffer
	if err := r.template.Execute(&content, bind); err != nil {
		return nil, err
	}

	sheet, err := parser.Parse(content.String())
	if err != nil {
		return nil, err
	}

	for _, rule := range sheet.Rules {
		r.morphRule(rule, parentNode)
	}

	stylesheet.Rules = append(stylesheet.Rules, sheet.Rules...)

	return &stylesheet, nil
}

// adjustName adjust the provided name according to the set rules of for specific
// css selectors.
func (r *Rule) adjustName(sel string, parentNode string) string {
	sel = strings.TrimSpace(sel)

	switch {
	case strings.Contains(sel, "&"):
		return strings.Replace(sel, "&", parentNode, -1)

	case strings.HasPrefix(sel, ":"):
		return parentNode + "" + sel

	default:
		return sel
	}
}

// morphRules adjusts the provided rules with the parent selector.
func (r *Rule) morphRule(base *bcss.Rule, parentNode string) {
	for index, sel := range base.Selectors {
		base.Selectors[index] = r.adjustName(sel, parentNode)
	}

	for _, rule := range base.Rules {
		if rule.Kind == bcss.AtRule {
			r.morphRule(rule, parentNode)
			continue
		}

		for index, sel := range rule.Selectors {
			rule.Selectors[index] = r.adjustName(sel, parentNode)
		}
	}
}
