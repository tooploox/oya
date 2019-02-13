package project_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/project"
	tu "github.com/bilus/oya/testutil"
)

func TestProject_Changeset(t *testing.T) {
	testCases := []struct {
		desc       string
		projectDir string
		workDir    string
		expected   []*oyafile.Oyafile
	}{
		{
			desc:       "Default changeset",
			projectDir: "./fixtures/project",
			workDir:    "./fixtures/project",
			expected:   tu.MustListOyafiles(t, "./fixtures/project"),
		},
		{
			desc:       "Default changeset from subdir",
			projectDir: "./fixtures/project",
			workDir:    "./fixtures/project/subdir/",
			expected:   tu.MustListOyafilesSubdir(t, "./fixtures/project", "./fixtures/project/subdir"),
		},
	}

	installDir, err := ioutil.TempDir("", "oya")
	tu.AssertNoErr(t, err, "Error creating temp dir")
	defer os.RemoveAll(installDir)

	for _, tc := range testCases {
		p, err := project.Detect(tc.projectDir, installDir)
		tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in %v (test case %q)", tc.projectDir, tc.desc)

		actual, err := p.Changeset(tc.workDir)
		tu.AssertNoErr(t, err, "Error calculating changeset (test case %q)", tc.desc)
		tu.AssertObjectsEqualMsg(t, tc.expected, actual, "Unexpected changeset (test case %q)", tc.desc)
	}
}
