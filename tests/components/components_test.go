package components

import (
	"testing"

	"github.com/gu-io/gu"
	"github.com/gu-io/gu/tests"
	"github.com/gu-io/gu/trees"
	"github.com/gu-io/gu/trees/elems"
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

	expected := "<div data-gen=\"gu\"  data-field=\"lexus\"  class=\"bomb\"><hello data-gen=\"gu\"><div data-gen=\"gu\">Welcome to the world \"Alex Thunderbot\"</div></hello></div>"
	expected2 := "<div data-gen=\"gu\"  class=\"bomb\"  data-field=\"lexus\"><hello data-gen=\"gu\"><div data-gen=\"gu\">Welcome to the world \"Alex Thunderbot\"</div></hello></div>"

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
	`)

	// Component prints:
	// <div data-gen="gu" class="bomb" data-field="lexus">
	// 	<hello data-gen="gu">
	// 		<div data-gen="gu">Welcome to the world "Alex Thunderbot"</div>
	// 	</hello>
	// </div>

	if val := component.Render().HTML(); val != expected && val != expected2 {
		t.Logf("\t\tRecieved: %q\n", val)
		t.Logf("\t\tExpected: %q\n", expected)
		tests.Failed(t, "Should have rendered expected markup")
	}
	tests.Passed(t, "Should have rendered expected markup")
}
