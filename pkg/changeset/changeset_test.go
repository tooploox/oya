package changeset_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/tooploox/oya/pkg/changeset"
	"github.com/tooploox/oya/pkg/oyafile"
	tu "github.com/tooploox/oya/testutil"
)

func TestNoChangesetTask(t *testing.T) {
	// No Changeset: directive means no changeset at this level.
	// project.Changeset takes care of adding a default changeset task
	// to root Oyafile if it's missing but here we don't worry about it.
	rootDir := fullPath("./fixtures/TestOneOyafile")
	actual, err := changeset.Calculate(tu.MustListOyafiles(t, rootDir))
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertEqual(t, 0, len(actual))
}

func TestEmptyChangeset(t *testing.T) {
	rootDir := fullPath("./fixtures/TestEmptyChangeset")
	expected := make([]*oyafile.Oyafile, 0)
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

func TestUniqueness(t *testing.T) {
	rootDir := fullPath("./fixtures/TestUniqueness")
	allOyafiles := tu.MustListOyafiles(t, rootDir)
	expected := allOyafiles
	actual, err := changeset.Calculate(tu.MustListOyafiles(t, rootDir))
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}

func fullPath(relPath string) string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(filename), relPath)
}
