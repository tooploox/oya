package task

import (
	"io"

	"github.com/tooploox/oya/pkg/template"
)

type Task interface {
	Exec(workDir string, args []string, scope template.Scope, stdout, stderr io.Writer) error
}
