package raw_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/bilus/oya/pkg/raw"
	tu "github.com/bilus/oya/testutil"
)

func TestOyafile_AddRequire_NoRequire(t *testing.T) {
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	oyafilePath := filepath.Join(outputDir, "Oyafile")
	tu.MustCopyFile(t, "./fixtures/AddRequire/Oyafile", oyafilePath)

	raw, found, err := raw.Load(oyafilePath, oyafilePath)
	tu.AssertNoErr(t, err, "Error loading raw Oyafile")
	tu.AssertTrue(t, found, "No Oyafile found")

	pack := tu.MustMakeMockPack(t, "github.com/tooploox/foo", "v1.0.0")
	err = raw.AddRequire(pack)
	tu.AssertNoErr(t, err, "Error adding require")

	expectedContent := `Project: AddRequire
Require:
  github.com/tooploox/foo: v1.0.0
`
	tu.AssertFileContains(t, oyafilePath, expectedContent)
}

func TestOyafile_AddRequire_EmptyRequire(t *testing.T) {
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	oyafilePath := filepath.Join(outputDir, "Oyafile")
	tu.MustCopyFile(t, "./fixtures/AddRequire_EmptyRequire/Oyafile", oyafilePath)

	raw, found, err := raw.Load(oyafilePath, oyafilePath)
	tu.AssertNoErr(t, err, "Error loading raw Oyafile")
	tu.AssertTrue(t, found, "No Oyafile found")

	pack := tu.MustMakeMockPack(t, "github.com/tooploox/foo", "v1.0.0")
	err = raw.AddRequire(pack)
	tu.AssertNoErr(t, err, "Error adding require")

	expectedContent := `Project: AddRequire_EmptyRequire

Require:
  github.com/tooploox/foo: v1.0.0
`
	tu.AssertFileContains(t, oyafilePath, expectedContent)
}

func TestOyafile_AddRequire_ExistingRequire(t *testing.T) {
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	oyafilePath := filepath.Join(outputDir, "Oyafile")
	tu.MustCopyFile(t, "./fixtures/AddRequire_ExistingRequire/Oyafile", oyafilePath)

	raw, found, err := raw.Load(oyafilePath, oyafilePath)
	tu.AssertNoErr(t, err, "Error loading raw Oyafile")
	tu.AssertTrue(t, found, "No Oyafile found")

	pack := tu.MustMakeMockPack(t, "github.com/tooploox/bar", "v1.1.0")
	err = raw.AddRequire(pack)
	tu.AssertNoErr(t, err, "Error adding require")

	expectedContent := `Project: AddRequire_ExistingRequire

Require:
  github.com/tooploox/bar: v1.1.0
  github.com/tooploox/foo: v1.0.0
`
	tu.AssertFileContains(t, oyafilePath, expectedContent)
}

func TestOyafile_AddRequire_SameVersion(t *testing.T) {
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	oyafilePath := filepath.Join(outputDir, "Oyafile")
	tu.MustCopyFile(t, "./fixtures/AddRequire_ExistingRequire/Oyafile", oyafilePath)

	raw, found, err := raw.Load(oyafilePath, oyafilePath)
	tu.AssertNoErr(t, err, "Error loading raw Oyafile")
	tu.AssertTrue(t, found, "No Oyafile found")

	pack := tu.MustMakeMockPack(t, "github.com/tooploox/foo", "v1.0.0")
	err = raw.AddRequire(pack)
	tu.AssertNoErr(t, err, "Error adding require")

	expectedContent := `Project: AddRequire_ExistingRequire

Require:
  github.com/tooploox/foo: v1.0.0
`
	tu.AssertFileContains(t, oyafilePath, expectedContent)
}

func TestOyafile_AddRequire_DifferentVersion(t *testing.T) {
	outputDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temporary output dir")
	defer os.RemoveAll(outputDir)

	oyafilePath := filepath.Join(outputDir, "Oyafile")
	tu.MustCopyFile(t, "./fixtures/AddRequire_ExistingRequire/Oyafile", oyafilePath)

	raw, found, err := raw.Load(oyafilePath, oyafilePath)
	tu.AssertNoErr(t, err, "Error loading raw Oyafile")
	tu.AssertTrue(t, found, "No Oyafile found")

	pack := tu.MustMakeMockPack(t, "github.com/tooploox/foo", "v1.1.0")
	err = raw.AddRequire(pack)
	tu.AssertNoErr(t, err, "Error adding require")

	expectedContent := `Project: AddRequire_ExistingRequire

Require:
  github.com/tooploox/foo: v1.1.0
`
	tu.AssertFileContains(t, oyafilePath, expectedContent)
}

// BUG(bilus): When adding new line under Requires, assumes indentation is two spaces. Should detect it.
