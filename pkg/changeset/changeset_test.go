package changeset_test

import (
	"testing"

	"github.com/bilus/oya/pkg/changeset"
	"github.com/bilus/oya/pkg/oyafile"
	tu "github.com/bilus/oya/testutil"
)

// TestEmptyOyafile

func TestOneOyafile(t *testing.T) {
	rootDir := "./fixtures/TestOneOyafile"
	actual, err := changeset.Calculate(mustListOyafiles(t, rootDir))
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertEqual(t, 1, len(actual))
}

func TestNoChangesetTask(t *testing.T) {
	rootDir := "./fixtures/TestNoChangesetTask"
	allOyafiles, err := oyafile.List(rootDir)
	tu.AssertNoErr(t, err, "Error listing Oyafiles")
	expected := allOyafiles
	actual, err := changeset.Calculate(mustListOyafiles(t, rootDir))
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}

func TestEmptyChangeset(t *testing.T) {
	rootDir := "./fixtures/TestEmptyChangeset"
	var expected []*oyafile.Oyafile
	actual, err := changeset.Calculate(mustListOyafiles(t, rootDir))
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}

func TestMinimalChangeset(t *testing.T) {
	rootDir := "./fixtures/TestMinimalChangeset"
	expected := []*oyafile.Oyafile{mustLoadOyafile(t, rootDir, rootDir)}
	actual, err := changeset.Calculate(mustListOyafiles(t, rootDir))
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}

func TestFullChangeset(t *testing.T) {
	rootDir := "./fixtures/TestFullChangeset"
	allOyafiles, err := oyafile.List(rootDir)
	tu.AssertNoErr(t, err, "Error listing Oyafiles")
	expected := allOyafiles
	actual, err := changeset.Calculate(mustListOyafiles(t, rootDir))
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}

func TestLocalOverride(t *testing.T) {
	rootDir := "./fixtures/TestLocalOverride"
	expected := []*oyafile.Oyafile{mustLoadOyafile(t, "./fixtures/TestLocalOverride/project1", rootDir)}
	actual, err := changeset.Calculate(mustListOyafiles(t, rootDir))
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}

func mustListOyafiles(t *testing.T, dir string) []*oyafile.Oyafile {
	oyafiles, err := oyafile.List(dir)
	tu.AssertNoErr(t, err, "Error listing Oyafiles")
	tu.AssertTrue(t, len(oyafiles) > 0, "No Oyafiles found")
	return oyafiles
}

func mustLoadOyafile(t *testing.T, dir, rootDir string) *oyafile.Oyafile {
	o, found, err := oyafile.LoadFromDir(dir, rootDir)
	tu.AssertNoErr(t, err, "Error loading root Oyafile")
	tu.AssertTrue(t, found, "Root Oyafile not found")
	return o
}
