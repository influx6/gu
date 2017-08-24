package property

import "github.com/gu-io/gu/trees"

// CustomAttr returns a new Attribute with the giving name and
// value.
func CustomAttr(name, val string) trees.Property {
	return &trees.Attribute{Name: name, Value: val}
}

// NameAttr defines attributes of type "Name" for html element types.
func NameAttr(val string) trees.Property {
	return &trees.Attribute{Name: "name", Value: val}
}

// CheckedAttr defines attributes of type "Checked" for html element types.
func CheckedAttr(val string) trees.Property {
	return &trees.Attribute{Name: "checked", Value: val}
}

// ClassAttr defines attributes of type "Class" for html element types.
func ClassAttr(vals ...string) trees.Property {
	return trees.NewClassList(vals...)
}

// ClassNameAttr defines attributes of type "ClassName" for html element types.
func ClassNameAttr(val string) trees.Property {
	return &trees.Attribute{Name: "className", Value: val}
}

// AutofocusAttr defines attributes of type "Autofocus" for html element types.
func AutofocusAttr(val string) trees.Property {
	return &trees.Attribute{Name: "autofocus", Value: val}
}

// IDAttr defines attributes of type "Id" for html element types..
// This is rather special in that, it uses the after effect function field
// to set the markups internal Id field to its value. This is intentional
// has it helps the css system provide a more accurate selector for the element.
func IDAttr(val string) trees.Property {
	return &trees.Attribute{Name: "id", Value: val, After: func(owner *trees.Markup) { owner.ID = val }}
}

// HTMLForAttr defines attributes of type "HtmlFor" for html element types.
func HTMLForAttr(val string) trees.Property {
	return &trees.Attribute{Name: "htmlFor", Value: val}
}

// SrcAttr defines attributes of type "Src" for html element types.
func SrcAttr(val string) trees.Property {
	return &trees.Attribute{Name: "src", Value: val}
}

// HrefAttr defines attributes of type "Href" for html element types.
func HrefAttr(val string) trees.Property {
	return &trees.Attribute{Name: "href", Value: val}
}

// RelAttr defines attributes of type "Rel" for html element types.
func RelAttr(val string) trees.Property {
	return &trees.Attribute{Name: "rel", Value: val}
}

// TypeAttr defines attributes of type "Type" for html element types.
func TypeAttr(val string) trees.Property {
	return &trees.Attribute{Name: "type", Value: val}
}

// PlaceholderAttr defines attributes of type "Placeholder" for html element types.
func PlaceholderAttr(val string) trees.Property {
	return &trees.Attribute{Name: "placeholder", Value: val}
}

// ValueAttr defines attributes of type "Value" for html element types.
func ValueAttr(val string) trees.Property {
	return &trees.Attribute{Name: "value", Value: val}
}
