//+build ignore

package components

import (
	"testing"

	"github.com/gu-io/gu"
	"github.com/gu-io/gu/trees"
	"github.com/gu-io/gu/trees/elems"
	"github.com/influx6/faux/tests"
)

type hello struct {
	Name     string
	Template string
}

func (h *hello) Render() *trees.Markup {
	return elems.Div(
		elems.Text(h.Template, h.Name),
	)
}

func TestComponent(t *testing.T) {
	trees.SetMode(trees.Pretty)

	expected := "<div data-gen=\"gu\"  class=\"bomb\"  data-field=\"lexus\" style=\"\"><hello data-gen=\"gu\" style=\"\">Welcome to the world \"Alex Thunderbot\"</hello></div>"

	registry := gu.NewComponentRegistry()
	registry.Register("hello", func(fields map[string]string, template string) gu.Renderable {
		return &hello{Name: fields["name"], Template: template}
	}, false)

	component := registry.Parse(`
		<div class="bomb" data-field="lexus">
			<hello component-name="Alex Thunderbot">
				<root-template>
					Welcome to the world %q
				</root-template>
			</hello>
		</div>
	`, nil)

	val := component.Render().HTML()
	if val != expected {
		tests.Info("Recieved: %+q", val)
		tests.Info("Expected: %+q", expected)
		tests.Failed("Should have rendered expected markup")
	}
	tests.Passed("Should have rendered expected markup")
}
