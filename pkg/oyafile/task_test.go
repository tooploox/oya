package oyafile_test

import (
	"testing"

	"github.com/bilus/oya/pkg/oyafile"
	tu "github.com/bilus/oya/testutil"
)

func TestTask_SplitName(t *testing.T) {
	testCases := []struct {
		taskName    string
		importAlias string
		baseName    string
	}{
		{
			taskName:    "",
			importAlias: "",
			baseName:    "",
		},
		{
			taskName:    "build",
			importAlias: "",
			baseName:    "build",
		},
		{
			taskName:    "docker.build",
			importAlias: "docker",
			baseName:    "build",
		},
		{
			taskName:    "foo.bar.baz",
			importAlias: "foo",
			baseName:    "bar.baz",
		},
	}
	for _, tc := range testCases {
		task := oyafile.ScriptedTask{
			Name: tc.taskName,
		}
		a, n := task.SplitName()
		tu.AssertEqual(t, tc.importAlias, a)
		tu.AssertEqual(t, tc.baseName, n)
	}
}
