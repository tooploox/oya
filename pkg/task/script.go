package task

import (
	"io"

	"github.com/tooploox/oya/pkg/shell"
	"github.com/tooploox/oya/pkg/template"
)

// OyaCmdOverride is used in tests, to override the path to the current oya executable.
// It is used to invoke other tasks from a task body.
// When tests are run, the current process executable path points to the test runner
// so it has to be overridden (with 'go run oya.go', roughly speaking).
var OyaCmdOverride *string

type Script struct {
	Script string
	Scope  *template.Scope
}

// Exec runs the script using the built-in shell interpreter.
func (s Script) Exec(workDir string, args []string, values template.Scope, stdout, stderr io.Writer) error {
	scope := values.Merge(*s.Scope)
	return shell.Exec(s.Script, workDir, args, scope, stdout, stderr, OyaCmdOverride)
}
