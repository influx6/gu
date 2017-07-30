package trees_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gu-io/gu/trees"
)

var success = "\u2713"
var failed = "\u2717"

func TestViewMarkup(t *testing.T) {
	trees.SetMode(trees.Pretty)
	defer trees.SetMode(trees.Normal)

	result := trees.ParseTree(`
		<header class="grid32 device">
			<div class="grid--block one-whole intro">
				<div class="grid--center one-half mascot">
					<a href="https://github.com/gu-io/gu"><img class="mascot-img" src="./assets/images/gu.png" /></a>
				</div>
				<div class="grid--center one-half title">
					<h1> GU Components</h1>
					<article class="griddesc">
						<p>This project demonstrates the different components built with the <a href="https://github.com/gu-io/gu">Gu</a> Library.</p>
					</article>
				</div>
			</div>
		</header>

		<div class="grid32 device components-demo">
			<div class="grid--block one-whole layout badges-demo"></div>
			<div class="grid--block one-whole layout layout-demo"></div>
			<div class="grid--block one-whole layout tooltip-demo"></div>
			<div class="grid--block one-whole layout avatar-demo"></div>
		</div>
  `)

	if len(result) != 2 {
		t.Fatalf("\t%s\t Should have parsed html markup properly: Found %d", failed, len(result))
	}
	t.Logf("\t%s\t Should have parsed html markup properly", success)

	var html []string
	for _, res := range result {
		html = append(html, res.HTML())
	}

	t.Logf("\t%s\t Parser should have produced markup for html: %q", success, strings.Join(html, ""))
}

func TestParserToText(t *testing.T) {
	input := `
		<header class="grid32 device">
			<div class="grid--block one-whole intro"></div>
		</header>
	`

	expected := "elemRoot := trees.NewMarkup(\"div\", false)\nelem1 := trees.NewMarkup(\"header\", false)\nelem1.Apply(elemRoot)\ntrees.NewAttr(\"class\", \"grid32 device\").Apply(elem1)\nelem2 := trees.NewMarkup(\"div\", false)\nelem2.Apply(elem1)\ntrees.NewAttr(\"class\", \"grid--block one-whole intro\").Apply(elem2)\n"

	var content bytes.Buffer
	trees.ParseTreeToText(input, false).WriteTo(&content)

	if content.String() != expected {
		t.Fatal("Parser should have successfully created markup in trees object format")
	}
	t.Log("Parser should have successfully created markup in trees object format")
}

func TestParser(t *testing.T) {
	trees.SetMode(trees.Pretty)
	defer trees.SetMode(trees.Normal)

	result := trees.ParseTree(`
		<!doctype html>
		<html>
			<head></head>
			<body>
				<div class="racket" id="racket-wrapper">
		      <a href="#" rel="bounce postive">Bounce +</a>
		    </div>

		    <!--thertorial words-->

				<div class="racket" id="racket-wrapper-2">
		      <a href="#" rel="bounce negative">Bounce -</a>
		    </div>
			</body>
		</html>
  `)

	var html []string
	for _, res := range result {
		html = append(html, res.HTML())
	}

	t.Logf("\t%s\t Parser should have produced markup for html: %q", success, strings.Join(html, ""))
}
