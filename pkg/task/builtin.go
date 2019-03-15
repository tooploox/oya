package task

import (
	"io"

	"github.com/tooploox/oya/pkg/template"
)

type Builtin struct {
	OnExec func(stdout, stderr io.Writer) error
}

func (h Builtin) Exec(workDir string, args []string, params template.Scope, stdout, stderr io.Writer) error {
	return h.OnExec(stdout, stderr)
}
