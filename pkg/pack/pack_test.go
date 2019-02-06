package pack_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/bilus/oya/pkg/pack"
	tu "github.com/bilus/oya/testutil"
)

func TestGithubPack_Vendor(t *testing.T) {
	vendorDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temp dir")
	defer os.RemoveAll(vendorDir)
	p, err := pack.New("github.com/tooploox/oya-fixtures", "v1.0.0")
	tu.AssertNoErr(t, err, "Error creating pack from uri")
	err = p.Vendor(vendorDir)
	tu.AssertNoErr(t, err, "Error vendoring pack")
	tu.AssertPathExists(t, filepath.Join(vendorDir, "github.com/tooploox/oya-fixtures/Oyafile"))
}
