package oyafile_test

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/template"
	tu "github.com/bilus/oya/testutil"
)

func init() {
	oyaCmdOverride := "go run github.com/bilus/oya"
	oyafile.OyaCmdOverride = &oyaCmdOverride
}

// TODO: More complete test coverage.

func TestRunningTasks(t *testing.T) {
	rootDir := fullPath("./fixtures/TestRunningTasks")
	installDir := filepath.Join(rootDir, ".packs")
	o, found, err := oyafile.LoadFromDir(rootDir, rootDir)
	tu.AssertTrue(t, found, "Oyafile not found")
	tu.AssertNoErr(t, err, "Error loading Oyafile")
	out := strings.Builder{}
	err = o.Build(installDir)
	tu.AssertNoErr(t, err, "Error building Oyafile")
	found, err = o.RunTask("bar", template.Scope{}, &out, &out)
	tu.AssertTrue(t, found, "Task 'bar' not found")
	tu.AssertNoErr(t, err, "Error running task 'bar'")
	tu.AssertRegexpMatch(t, ".*foo\nbar\n$", out.String())
}

func fullPath(relPath string) string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(filename), relPath)
}
