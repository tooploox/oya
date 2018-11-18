package template_test

import (
	"bytes"
	"testing"

	"github.com/bilus/oya/pkg/template"
	tu "github.com/bilus/oya/testutil"
)

func TestLoad(t *testing.T) {
	_, err := template.Load("./fixtures/good.txt.kasia")
	tu.AssertNoErr(t, err, "Expected template to load")
}

func TestParse(t *testing.T) {
	_, err := template.Parse("$foo")
	tu.AssertNoErr(t, err, "Expected template to parse")
}

func TestTemplate_Render_Loaded(t *testing.T) {
	tpl, err := template.Load("./fixtures/good.txt.kasia")
	tu.AssertNoErr(t, err, "Expected template to load")
	output := new(bytes.Buffer)
	err = tpl.Render(output, map[string]string{"foo": "bar"})
	tu.AssertNoErr(t, err, "Expected template to render")
	tu.AssertEqual(t, "bar\n", output.String())
}

func TestTemplate_Render_Parsed(t *testing.T) {
	tpl, err := template.Parse("$foo")
	tu.AssertNoErr(t, err, "Expected template to parse")
	output := new(bytes.Buffer)
	err = tpl.Render(output, map[string]string{"foo": "bar"})
	tu.AssertNoErr(t, err, "Expected template to render")
	tu.AssertEqual(t, "bar", output.String())
}
