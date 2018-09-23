package changeset_test

import (
	"testing"

	"github.com/bilus/oya/pkg/changeset"
	"github.com/bilus/oya/pkg/fixtures"
	"github.com/bilus/oya/pkg/oyafile"
	tu "github.com/bilus/oya/testutil"
)

func TestNoHook(t *testing.T) {
	rootDir := "/tmp/"
	oyafiles := []*oyafile.Oyafile{
		fixtures.Oyafile(
			rootDir,
		),
	}
	expected := oyafiles
	actual, err := changeset.Calculate(oyafiles)
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}

func TestEmptyChangeset(t *testing.T) {
	rootDir := "/tmp/"
	oyafiles := []*oyafile.Oyafile{
		fixtures.Oyafile(
			rootDir,
			"Changeset", "cat /dev/null",
		),
	}
	expected := []*oyafile.Oyafile{}
	actual, err := changeset.Calculate(oyafiles)
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}

func TestMinimalChangeset(t *testing.T) {
	rootDir := "/tmp/"
	oyafiles := []*oyafile.Oyafile{
		fixtures.Oyafile(
			rootDir,
			"Changeset", "echo '+/tmp/'",
		),
	}
	expected := oyafiles
	actual, err := changeset.Calculate(oyafiles)
	tu.AssertNoErr(t, err, "Error calculating changeset")
	tu.AssertObjectsEqual(t, expected, actual)
}
