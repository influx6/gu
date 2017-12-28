package trees_test

import (
	"testing"

	"github.com/gu-io/gu/trees"
)

// TestFinders validates the utility functions which help in retrieving underline
// markup.
func TestFinders(t *testing.T) {
	body := generateMarkup()

	if _, err := trees.GetStyle(body, "width"); err != nil {
		t.Fatalf("\t%s\t  Should have been able to find style in markup", failed)
	}
	t.Logf("\t%s\t  Should have been able to find style in markup", success)

	if _, err := trees.GetAttr(body, "id"); err != nil {
		t.Fatalf("\t%s\t  Should have been able to find attr in markup", failed)
	}
	t.Logf("\t%s\t  Should have been able to find attr in markup", success)

	if len(trees.GetStyles(body, "width", "auto")) < 1 {
		t.Fatalf("\t%s\t  Should have been able to find style in markup with value of auto and name of width", failed)
	}
	t.Logf("\t%s\t  Should have been able to find style in markup with value of auto and name of width", success)

	if len(trees.ElementsUsingStyle(body, "width", "")) > 2 {
		t.Fatalf("\t%s\t  Should have been able to found more than two elements using width style", failed)
	}
	t.Logf("\t%s\t  Should have been able to found more than two elements using width style", success)

	if len(trees.ElementsUsingStyle(body, "id", "")) > 2 {
		t.Fatalf("\t%s\t  Should have been able to found more than two elements using id attr", failed)
	}
	t.Logf("\t%s\t  Should have been able to found more than two elements using id attr", success)

	if len(trees.ElementsUsingStyle(body, "width", "200px")) > 2 {
		t.Fatalf("\t%s\t  Should have been able to found more than two elements using width style with value of 200px", failed)
	}
	t.Logf("\t%s\t  Should have been able to found more than two elements using width style with value of 200px", success)

	if len(trees.ElementsWithTag(body, "div")) < 2 {
		t.Fatalf("\t%s\t  Should have been able to found two div element", failed)
	}
	t.Logf("\t%s\t  Should have been able to found two div element", success)

	if len(trees.ElementsWithTag(body, "label")) != 1 {
		t.Fatalf("\t%s\t  Should have been able to find a label element", failed)
	}
	t.Logf("\t%s\t  Should have been able to find a label element", success)
}

func generateMarkup() *trees.Markup {
	body := trees.NewMarkup("body", false)
	trees.NewCSSStyle("width", "auto").Apply(body)
	trees.NewAttr("id", "main-wrapper").Apply(body)

	root := trees.NewMarkup("div", false)
	trees.NewAttr("id", "root-div").Apply(root)
	trees.NewAttr("class", "roots").Apply(root)
	trees.NewCSSStyle("width", "100px").Apply(root)
	trees.NewMarkup("section", false).Apply(root)
	root.Apply(body)

	root2 := trees.NewMarkup("div", false)
	trees.NewAttr("id", "root-div-2").Apply(root2)
	trees.NewAttr("class", "roots").Apply(root2)
	root2.Apply(body)

	label := trees.NewMarkup("label", false)
	trees.NewCSSStyle("width", "200px").Apply(label)
	label.Apply(root2)

	return body
}
