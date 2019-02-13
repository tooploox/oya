package pack_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/semver"
	tu "github.com/bilus/oya/testutil"
)

func TestGithubPack_Vendor(t *testing.T) {
	installDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temp dir")
	defer os.RemoveAll(installDir)
	l, err := pack.OpenLibrary("github.com/tooploox/oya-fixtures")
	tu.AssertNoErr(t, err, "Error opening pack library")
	p, err := l.Version(semver.MustParse("v1.0.0"))
	tu.AssertNoErr(t, err, "Error getting pack")
	err = p.Install(installDir)
	tu.AssertNoErr(t, err, "Error vendoring pack")
	tu.AssertPathExists(t, filepath.Join(installDir, "github.com/tooploox/oya-fixtures@v1.0.0/Oyafile"))
}
