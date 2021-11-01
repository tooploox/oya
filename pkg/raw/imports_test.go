package raw_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/tooploox/oya/pkg/raw"
	tu "github.com/tooploox/oya/testutil"
)

func TestOyafile_AddImport_NoImport(t *testing.T) {
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	oyafilePath := filepath.Join(outputDir, "Oyafile")
	tu.MustCopyFile(t, "./fixtures/AddImport/Oyafile", oyafilePath)

	raw, found, err := raw.Load(oyafilePath, oyafilePath)
	tu.AssertNoErr(t, err, "Error loading raw Oyafile")
	tu.AssertTrue(t, found, "No Oyafile found")

	err = raw.AddImport("foo", "github.com/tooploox/foo")
	tu.AssertNoErr(t, err, "Error adding import")

	tu.AssertNoErr(t, raw.ApplyChanges(), "Error applying changes")

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
	tu.MustCopyFile(t, "./fixtures/AddImport_ToExisting/Oyafile", oyafilePath)

	raw, found, err := raw.Load(oyafilePath, oyafilePath)
	tu.AssertNoErr(t, err, "Error loading raw Oyafile")
	tu.AssertTrue(t, found, "No Oyafile found")

	err = raw.AddImport("bar", "github.com/tooploox/bar")
	tu.AssertNoErr(t, err, "Error adding import")

	tu.AssertNoErr(t, raw.ApplyChanges(), "Error applying changes")

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
	tu.MustCopyFile(t, "./fixtures/AddImport_MoreKeys/Oyafile", oyafilePath)

	raw, found, err := raw.Load(oyafilePath, oyafilePath)
	tu.AssertNoErr(t, err, "Error loading raw Oyafile")
	tu.AssertTrue(t, found, "No Oyafile found")

	err = raw.AddImport("bar", "github.com/tooploox/bar")
	tu.AssertNoErr(t, err, "Error adding import")

	tu.AssertNoErr(t, raw.ApplyChanges(), "Error applying changes")

	expectedContent := `Project: AddImport_MoreKeys
Import:
  bar: github.com/tooploox/bar

Values:
  baz: qux
`
	tu.AssertFileContains(t, oyafilePath, expectedContent)
}

func TestOyafile_AddImport_TwoPacks(t *testing.T) {
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	oyafilePath := filepath.Join(outputDir, "Oyafile")
	tu.MustCopyFile(t, "./fixtures/AddImport/Oyafile", oyafilePath)

	raw, found, err := raw.Load(oyafilePath, oyafilePath)
	tu.AssertNoErr(t, err, "Error loading raw Oyafile")
	tu.AssertTrue(t, found, "No Oyafile found")

	err = raw.AddImport("foo", "github.com/tooploox/foo")
	tu.AssertNoErr(t, err, "Error adding import")
	err = raw.AddImport("bar", "github.com/tooploox/bar")
	tu.AssertNoErr(t, err, "Error adding import")

	tu.AssertNoErr(t, raw.ApplyChanges(), "Error applying changes")

	expectedContent := `Project: AddImport
Import:
  bar: github.com/tooploox/bar
  foo: github.com/tooploox/foo
`
	tu.AssertFileContains(t, oyafilePath, expectedContent)
}

func TestOyafile_AddImport_AlreadyImported(t *testing.T) {
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	oyafilePath := filepath.Join(outputDir, "Oyafile")
	tu.MustCopyFile(t, "./fixtures/AddImport_ToExisting/Oyafile", oyafilePath)

	raw, found, err := raw.Load(oyafilePath, oyafilePath)
	tu.AssertNoErr(t, err, "Error loading raw Oyafile")
	tu.AssertTrue(t, found, "No Oyafile found")

	err = raw.AddImport("foo", "github.com/tooploox/foo")
	tu.AssertErr(t, err, "Expected an error trying to import the same pack twice")

	tu.AssertNoErr(t, raw.ApplyChanges(), "Error applying changes")

	expectedContent := `Project: AddImport_ToExisting

Import:
  foo: github.com/tooploox/foo
`
	tu.AssertFileContains(t, oyafilePath, expectedContent)
}

func TestOyafile_Expose_FirstTime(t *testing.T) {
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	oyafilePath := filepath.Join(outputDir, "Oyafile")
	tu.MustCopyFile(t, "./fixtures/Expose_FirstTime/Oyafile", oyafilePath)

	raw, found, err := raw.Load(oyafilePath, oyafilePath)
	tu.AssertNoErr(t, err, "Error loading raw Oyafile")
	tu.AssertTrue(t, found, "No Oyafile found")

	err = raw.Expose("foo")
	tu.AssertNoErr(t, err, "Error adding import")

	tu.AssertNoErr(t, raw.ApplyChanges(), "Error applying changes")

	expectedContent := `Project: Expose_FirstTime
Expose: foo
Import:
  foo: github.com/tooploox/foo
`
	tu.AssertFileContains(t, oyafilePath, expectedContent)

}
