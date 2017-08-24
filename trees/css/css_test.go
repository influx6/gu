package css_test

import (
	"testing"

	"github.com/gu-io/gu/trees/css"
	"github.com/influx6/faux/tests"
)

func TestBasicCSS(t *testing.T) {
	expected := "#galatica:hover {\n  color: red;\n}\n#galatica::before {\n  content: \"bugger\";\n}\n#galatica div a {\n  color: black;\n  font-family: Helvetica;\n}\n@media (max-width: 400px) {\n  #galatica:hover {\n    color: blue;\n    font-family: Helvetica;\n  }\n}"

	csr := css.New(`

    &:hover {
      color: red;
    }

    &::before {
      content: "bugger";
    }

    & div a {
      color: black;
      font-family: {{ .Font }};
    }

    @media (max-width: 400px){

      &:hover {
        color: blue;
        font-family: {{ .Font }};
      }

    }
`, nil)

	sheet, err := csr.Stylesheet(struct {
		Font string
	}{Font: "Helvetica"}, "#galatica")

	if err != nil {
		tests.Failed("Should have successfully processed stylesheet for rule")
	}
	tests.Passed("Should have successfully processed stylesheet for rule")

	if val := sheet.String(); val != expected {
		t.Logf("\t\tRecieved: %q\n", val)
		t.Logf("\t\tExpected: %q\n", expected)
		tests.Failed("Should have rendered expected stylesheet")
	}
	tests.Passed("Should have rendered expected stylesheet")
}

func TestExtensionCSS(t *testing.T) {
	expected := "#galatica::before {\n  content: \"bugger\";\n}\ndiv a {\n  font-family: Helvetica;\n  color: Pink;\n  border: 1px solid #000;\n}\n@media (max-width: 400px) {\n  #galatica:hover {\n    color: blue;\n    font-family: Helvetica;\n  }\n}"

	csr := css.New(`
    block {
      font-family: {{ .Font }};
      color: {{ .Color }};
    }
  `, nil)

	csx := css.New(`

    ::before {
      content: "bugger";
    }

    div a {
			{{ extend "block" }}
			border: 1px solid #000;
    }

    @media (max-width: 400px){

      :hover {
        color: blue;
        font-family: {{ .Font }};
      }

    }
`, csr)

	sheet, err := csx.Stylesheet(struct {
		Font  string
		Color string
	}{
		Font:  "Helvetica",
		Color: "Pink",
	}, "#galatica")

	if err != nil {
		tests.Failed("Should have successfully processed stylesheet for rule")
	}
	tests.Passed("Should have successfully processed stylesheet for rule")

	if res := sheet.String(); res != expected {
		t.Logf("\t\tRecieved: %q\n", res)
		t.Logf("\t\tExpected: %q\n", expected)
		tests.Failed("Should have rendered expected stylesheet")
	}
	tests.Passed("Should have rendered expected stylesheet")
}

func TestLinkedCSS(t *testing.T) {
	expected := "block {\n  font-family: Helvetica;\n  color: Pink;\n}\n#galatica::before {\n  content: \"bugger\";\n}\ndiv a {\n  color: black;\n  font-family: Helvetica;\n}\n@media (max-width: 400px) {\n  #galatica:hover {\n    color: blue;\n    font-family: Helvetica;\n  }\n}"

	csr := css.New(`
    block {
      font-family: {{ .Font }};
      color: {{ .Color }};
    }
  `, nil)

	csx := css.New(`

    ::before {
      content: "bugger";
    }

    div a {
      color: black;
      font-family: {{ .Font }};
    }

    @media (max-width: 400px){

      :hover {
        color: blue;
        font-family: {{ .Font }};
      }

    }
`, nil, csr)

	sheet, err := csx.Stylesheet(struct {
		Font  string
		Color string
	}{
		Font:  "Helvetica",
		Color: "Pink",
	}, "#galatica")

	if err != nil {
		tests.Failed("Should have successfully processed stylesheet for rule")
	}
	tests.Passed("Should have successfully processed stylesheet for rule")

	if res := sheet.String(); res != expected {
		t.Logf("\t\tRecieved: %q\n", res)
		t.Logf("\t\tExpected: %q\n", expected)
		tests.Failed("Should have rendered expected stylesheet")
	}
	tests.Passed("Should have rendered expected stylesheet")
}
