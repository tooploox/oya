package task

import (
	"io"

	"github.com/tooploox/oya/pkg/template"
)

type Task interface {
	Exec(workDir string, params template.Scope, stdout, stderr io.Writer) error
}
