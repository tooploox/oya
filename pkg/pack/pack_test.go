package pack_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/bilus/oya/pkg/pack"
	tu "github.com/bilus/oya/testutil"
)

func TestGitPackage_Vendor(t *testing.T) {
	vendorDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temp dir")
	defer os.RemoveAll(vendorDir)
	p, err := pack.NewFromUri("github.com/bilus/oya", "fixtures")
	tu.AssertNoErr(t, err, "Error creating package from uri")
	err = p.Vendor(vendorDir)
	tu.AssertNoErr(t, err, "Error vendoring package")
	tu.AssertPathExists(t, filepath.Join(vendorDir, "github.com/bilus/oya/fixtures/ExampleGitPackage_Vendor/Oyafile"))
}
