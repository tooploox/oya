package template_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/tooploox/oya/pkg/template"
	tu "github.com/tooploox/oya/testutil"
)

func TestLoad(t *testing.T) {
	delimiters := template.Delimiters{"<%", "%>"}
	_, err := template.Load("./fixtures/good.txt.plush", delimiters)
	tu.AssertNoErr(t, err, "Expected template to load")
}

func TestParse(t *testing.T) {
	delimiters := template.Delimiters{"<%", "%>"}
	_, err := template.Parse("$foo", delimiters)
	tu.AssertNoErr(t, err, "Expected template to parse")
}

func TestTemplate_Render_Loaded(t *testing.T) {
	delimiters := template.Delimiters{"<%", "%>"}
	tpl, err := template.Load("./fixtures/good.txt.plush", delimiters)
	tu.AssertNoErr(t, err, "Expected template to load")
	output := new(bytes.Buffer)
	err = tpl.Render(output, template.Scope{"foo": "bar"})
	tu.AssertNoErr(t, err, "Expected template to render")
	tu.AssertEqual(t, "bar\n", output.String())
}

func TestTemplate_Render_Parsed(t *testing.T) {
	delimiters := template.Delimiters{"<%", "%>"}
	tpl, err := template.Parse("<%= foo %>", delimiters)
	tu.AssertNoErr(t, err, "Expected template to parse")
	output := new(bytes.Buffer)
	err = tpl.Render(output, template.Scope{"foo": "bar"})
	tu.AssertNoErr(t, err, "Expected template to render")
	tu.AssertEqual(t, "bar", output.String())
}

func TestTemplate_Render_MissingVariables(t *testing.T) {
	delimiters := template.Delimiters{"<%", "%>"}
	tpl, err := template.Parse("<%= noSuchVar %>", delimiters)
	tu.AssertNoErr(t, err, "Expected template to parse")
	err = tpl.Render(ioutil.Discard, template.Scope{"foo": "bar"})
	tu.AssertErr(t, err, "Expected template not to render")
}

func TestRenderAll_Directory(t *testing.T) {
	delimiters := template.Delimiters{"<%", "%>"}
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	err = template.RenderAll("./fixtures/", nil, outputDir,
		template.Scope{"foo": "bar", "baz": map[string]interface{}{"qux": "abc"}}, delimiters)
	tu.AssertNoErr(t, err, "Expected templates to render")

	tu.AssertFileContains(t, filepath.Join(outputDir, "good.txt.plush"), "bar\n")
	tu.AssertFileContains(t, filepath.Join(outputDir, "subdir/nested.txt.plush"), "barabc\n")
}

func TestRenderAll_SingleFile(t *testing.T) {
	delimiters := template.Delimiters{"<%", "%>"}
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	err = template.RenderAll("./fixtures/good.txt.plush", nil, outputDir, template.Scope{"foo": "bar"}, delimiters)
	tu.AssertNoErr(t, err, "Expected templates to render")

	tu.AssertFileContains(t, filepath.Join(outputDir, "good.txt.plush"), "bar\n")
}

func TestRenderAll_ExcludedPaths(t *testing.T) {
	delimiters := template.Delimiters{"<%", "%>"}
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	excludedPaths := []string{"good.txt.plush"}
	err = template.RenderAll("./fixtures/", excludedPaths, outputDir,
		template.Scope{"foo": "bar", "baz": map[string]interface{}{"qux": "abc"}}, delimiters)
	tu.AssertNoErr(t, err, "Expected templates to render")

	tu.AssertPathNotExists(t, filepath.Join(outputDir, "good.txt.plush"))
	tu.AssertFileContains(t, filepath.Join(outputDir, "subdir/nested.txt.plush"), "barabc\n")
}

func TestRenderAll_ExcludedPatterns(t *testing.T) {
	delimiters := template.Delimiters{"<%", "%>"}
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	excludedPaths := []string{"**.txt.plush"}
	err = template.RenderAll("./fixtures/", excludedPaths, outputDir, template.Scope{"foo": "bar"}, delimiters)
	tu.AssertNoErr(t, err, "Expected templates to render")

	tu.AssertPathNotExists(t, filepath.Join(outputDir, "good.txt.plush"))
	tu.AssertPathNotExists(t, filepath.Join(outputDir, "subdir/nested.txt.plush"))
}

func TestTemplate_Render_CustomDelimiters(t *testing.T) {
	delimiters := template.Delimiters{"{{", "}}"}
	tpl, err := template.Parse("{{= foo }}", delimiters)
	tu.AssertNoErr(t, err, "Expected template to parse")
	output := new(bytes.Buffer)
	err = tpl.Render(output, template.Scope{"foo": "bar"})
	tu.AssertNoErr(t, err, "Expected template to render")
	tu.AssertEqual(t, "bar", output.String())
}
