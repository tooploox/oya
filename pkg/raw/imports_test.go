package raw_test

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/bilus/oya/pkg/raw"
	tu "github.com/bilus/oya/testutil"
)

func TestOyafile_AddImport_NoImport(t *testing.T) {
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	oyafilePath := filepath.Join(outputDir, "Oyafile")
	copyFile("./fixtures/AddImport/Oyafile", oyafilePath)

	raw, found, err := raw.Load(oyafilePath, oyafilePath)
	tu.AssertNoErr(t, err, "Error loading raw Oyafile")
	tu.AssertTrue(t, found, "No Oyafile found")

	err = raw.AddImport("foo", "github.com/tooploox/foo")
	tu.AssertNoErr(t, err, "Error adding import")

	expectedContent := `Project: AddImport
Import:
  foo: github.com/tooploox/foo
`
	tu.AssertFileContains(t, oyafilePath, expectedContent)
}

func TestOyafile_AddImport_ToExisting(t *testing.T) {
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	oyafilePath := filepath.Join(outputDir, "Oyafile")
	copyFile("./fixtures/AddImport_ToExisting/Oyafile", oyafilePath)

	raw, found, err := raw.Load(oyafilePath, oyafilePath)
	tu.AssertNoErr(t, err, "Error loading raw Oyafile")
	tu.AssertTrue(t, found, "No Oyafile found")

	err = raw.AddImport("bar", "github.com/tooploox/bar")
	tu.AssertNoErr(t, err, "Error adding import")

	expectedContent := `Project: AddImport_ToExisting

Import:
  bar: github.com/tooploox/bar
  foo: github.com/tooploox/foo
`
	tu.AssertFileContains(t, oyafilePath, expectedContent)
}

func TestOyafile_AddImport_MoreKeys(t *testing.T) {
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	oyafilePath := filepath.Join(outputDir, "Oyafile")
	copyFile("./fixtures/AddImport_MoreKeys/Oyafile", oyafilePath)

	raw, found, err := raw.Load(oyafilePath, oyafilePath)
	tu.AssertNoErr(t, err, "Error loading raw Oyafile")
	tu.AssertTrue(t, found, "No Oyafile found")

	err = raw.AddImport("bar", "github.com/tooploox/bar")
	tu.AssertNoErr(t, err, "Error adding import")

	expectedContent := `Project: AddImport_MoreKeys
Import:
  bar: github.com/tooploox/bar

Values:
  baz: qux
`
	tu.AssertFileContains(t, oyafilePath, expectedContent)
}

func TestOyafile_AddImport_Twice(t *testing.T) {
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	oyafilePath := filepath.Join(outputDir, "Oyafile")
	copyFile("./fixtures/AddImport/Oyafile", oyafilePath)

	raw, found, err := raw.Load(oyafilePath, oyafilePath)
	tu.AssertNoErr(t, err, "Error loading raw Oyafile")
	tu.AssertTrue(t, found, "No Oyafile found")

	err = raw.AddImport("foo", "github.com/tooploox/foo")
	tu.AssertNoErr(t, err, "Error adding import")
	err = raw.AddImport("bar", "github.com/tooploox/bar")
	tu.AssertNoErr(t, err, "Error adding import")

	expectedContent := `Project: AddImport
Import:
  bar: github.com/tooploox/bar
  foo: github.com/tooploox/foo
`
	tu.AssertFileContains(t, oyafilePath, expectedContent)
}

func copyFile(fromPath, toPath string) error {
	from, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	return err
}
