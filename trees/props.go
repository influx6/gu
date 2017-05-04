package trees

import "strings"

//==============================================================================

// Property defines the interface for attributes in trees.
// It provides a apply and RenderAttribute which returns the key
// and value for that attribute.
type Property interface {
	Apply(*Markup)
	Clone() Property
	Render() (string, string)
}

// WhenElse checks if the first is a non-nil and returns else returns the second
// property instead.
func WhenElse(first Property, other Property) Property {
	if first != nil {
		return first
	}

	return other
}

// MarkupWhenElse returns the first markup if its not nil else the second.
func MarkupWhenElse(first, other *Markup) *Markup {
	if first != nil {
		return first
	}

	return other
}

// MarkupWhen returns the first or other markup when the giving state is false or true.
func MarkupWhen(state bool, first, other *Markup) *Markup {
	if state {
		return first
	}

	return other
}

// MarkupIf returns the markup when the giving state is true.
func MarkupIf(state bool, m *Markup) *Markup {
	if !state {
		return nil
	}

	return m
}

// If returns the property when the giving state is true.
func If(state bool, m Property) Property {
	if !state {
		return nil
	}

	return m
}

// When checks if the giving state is true and returns the first property else
// returns the second.
func When(state bool, first Property, other Property) Property {
	if state {
		return first
	}

	return other
}

//==============================================================================

// Attribute define the struct  for attributes
type Attribute struct {
	Name  string
	Value string
	After func(*Markup)
}

// NewAttr returns a new attribute instance.
func NewAttr(name, val string) *Attribute {
	a := Attribute{Name: strings.ToLower(name), Value: val}
	return &a
}

// NewAttrWith returns a new attribute instance with a provided function
// to call to provide a after effect to the markup.
func NewAttrWith(name, val string, after func(*Markup)) *Attribute {
	a := Attribute{Name: strings.ToLower(name), Value: val, After: after}
	return &a
}

// Render returns the key and value for this attribute rendered.
func (a *Attribute) Render() (string, string) {
	return a.Name, a.Value
}

// Apply applies a set change to the giving element attributes list
func (a *Attribute) Apply(e *Markup) {
	if e.allowAttributes {
		e.AddAttribute(a)

		if a.After != nil {
			a.After(e)
		}
	}
}

//Clone replicates the attribute into a unique instance
func (a *Attribute) Clone() Property {
	return &Attribute{Name: a.Name, Value: a.Value, After: a.After}
}

//==============================================================================

// CSSStyle define the style specification for element styles
type CSSStyle struct {
	Name  string
	Value string
}

// NewCSSStyle returns a new style instance
func NewCSSStyle(name, val string) *CSSStyle {
	s := CSSStyle{Name: strings.ToLower(name), Value: val}
	return &s
}

// Render returns the key and value for this style rendered.
func (s *CSSStyle) Render() (string, string) {
	return s.Name, s.Value
}

//Clone replicates the style into a unique instance
func (s *CSSStyle) Clone() Property {
	return &CSSStyle{Name: s.Name, Value: s.Value}
}

// Apply applies a set change to the giving element style list
func (s *CSSStyle) Apply(e *Markup) {
	if e.allowStyles {
		e.AddStyle(s)
	}
}

//==============================================================================

// ClassList defines the list type for class lists.
type ClassList struct {
	list []string
}

// NewClassList returns a new ClassList instance.
func NewClassList(c ...string) *ClassList {
	return &ClassList{list: c}
}

// Add adds a class name into the lists.
func (c *ClassList) Add(classes ...string) {
	c.list = append(c.list, classes...)
}

// Render returns the key and value for this style rendered.
func (c *ClassList) Render() (string, string) {
	return "class", strings.Join(c.list, " ")
}

// Apply checks for a class attribute
func (c *ClassList) Apply(em *Markup) {
	if em.allowAttributes {
		var old Property
		index := -1

		for ind, attr := range em.attrs {
			name, _ := attr.Render()
			if name != "class" {
				continue
			}

			old = attr
			index = ind
			break
		}

		if index == -1 {
			em.AddAttribute(c)
			return
		}

		if cold, ok := old.(*ClassList); ok {
			cold.Add(c.list...)
		} else {
			_, val := old.Render()
			c.Add(val)
			em.attrs[index] = c
		}

	}
}

// Clone replicates the lists of classnames.
func (c ClassList) Clone() Property {
	return &ClassList{
		list: c.list[:len(c.list)],
	}
}
