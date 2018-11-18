package template_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
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

func TestRenderAll(t *testing.T) {
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	err = template.RenderAll("./fixtures/", outputDir, template.Scope{"foo": "bar"})
	tu.AssertNoErr(t, err, "Expected templates to render")

	tu.AssertFileContains(t, filepath.Join(outputDir, "good.txt.kasia"), "bar\n")
	tu.AssertFileContains(t, filepath.Join(outputDir, "subdir/nested.txt.kasia"), "bar\n")
}
