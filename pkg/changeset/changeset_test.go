package changeset_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/bilus/oya/pkg/changeset"
	"github.com/bilus/oya/pkg/oyafile"
	tu "github.com/bilus/oya/testutil"
)

// TestEmptyOyafile

func TestOneOyafile(t *testing.T) {
	rootDir := fullPath("./fixtures/TestOneOyafile")
	actual, err := changeset.Calculate(tu.MustListOyafiles(t, rootDir))
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertEqual(t, 1, len(actual))
}

func TestNoChangesetTask(t *testing.T) {
	rootDir := fullPath("./fixtures/TestNoChangesetTask")
	allOyafiles := tu.MustListOyafiles(t, rootDir)
	expected := allOyafiles
	actual, err := changeset.Calculate(tu.MustListOyafiles(t, rootDir))
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}

func TestEmptyChangeset(t *testing.T) {
	rootDir := fullPath("./fixtures/TestEmptyChangeset")
	var expected []*oyafile.Oyafile
	actual, err := changeset.Calculate(tu.MustListOyafiles(t, rootDir))
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}

func TestMinimalChangeset(t *testing.T) {
	rootDir := fullPath("./fixtures/TestMinimalChangeset")
	expected := []*oyafile.Oyafile{tu.MustLoadOyafile(t, rootDir, rootDir)}
	actual, err := changeset.Calculate(tu.MustListOyafiles(t, rootDir))
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}

func TestFullChangeset(t *testing.T) {
	rootDir := fullPath("./fixtures/TestFullChangeset")
	allOyafiles := tu.MustListOyafiles(t, rootDir)
	expected := allOyafiles
	actual, err := changeset.Calculate(tu.MustListOyafiles(t, rootDir))
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}

func TestLocalOverride(t *testing.T) {
	rootDir := fullPath("./fixtures/TestLocalOverride")
	expected := []*oyafile.Oyafile{tu.MustLoadOyafile(t, filepath.Join(rootDir, "./project1"), rootDir)}
	actual, err := changeset.Calculate(tu.MustListOyafiles(t, rootDir))
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}

func fullPath(relPath string) string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(filename), relPath)
}
