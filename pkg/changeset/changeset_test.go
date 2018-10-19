package changeset_test

import (
	"testing"

	"github.com/bilus/oya/pkg/changeset"
	"github.com/bilus/oya/pkg/oyafile"
	tu "github.com/bilus/oya/testutil"
)

// TestEmptyOyafile

func mustLoadOyafile(t *testing.T, dir string) *oyafile.Oyafile {
	o, found, err := oyafile.LoadFromDir(dir)
	tu.AssertNoErr(t, err, "Error loading root Oyafile")
	tu.AssertTrue(t, found, "Root Oyafile not found")
	return o
}

func TestOneOyafile(t *testing.T) {
	rootDir := "./fixtures/TestOneOyafile"
	o := mustLoadOyafile(t, rootDir)
	actual, err := changeset.Calculate(o)
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertEqual(t, 1, len(actual))
}

func TestNoChangesetHook(t *testing.T) {
	rootDir := "./fixtures/TestNoChangesetHook"
	allOyafiles, err := oyafile.List(rootDir)
	tu.AssertNoErr(t, err, "Error listing Oyafiles")
	o := mustLoadOyafile(t, rootDir)
	expected := allOyafiles
	actual, err := changeset.Calculate(o)
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}

func TestEmptyChangeset(t *testing.T) {
	rootDir := "./fixtures/TestEmptyChangeset"
	var expected []*oyafile.Oyafile
	o := mustLoadOyafile(t, rootDir)
	actual, err := changeset.Calculate(o)
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}

func TestMinimalChangeset(t *testing.T) {
	rootDir := "./fixtures/TestMinimalChangeset"
	o := mustLoadOyafile(t, rootDir)
	expected := []*oyafile.Oyafile{o}
	actual, err := changeset.Calculate(o)
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}

func TestFullChangeset(t *testing.T) {
	rootDir := "./fixtures/TestFullChangeset"
	o := mustLoadOyafile(t, rootDir)
	allOyafiles, err := oyafile.List(rootDir)
	tu.AssertNoErr(t, err, "Error listing Oyafiles")
	expected := allOyafiles
	actual, err := changeset.Calculate(o)
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}

func TestLocalOverride(t *testing.T) {
	rootDir := "./fixtures/TestLocalOverride"
	o := mustLoadOyafile(t, rootDir)
	expected := []*oyafile.Oyafile{mustLoadOyafile(t, "./fixtures/TestLocalOverride/project1")}
	actual, err := changeset.Calculate(o)
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}
