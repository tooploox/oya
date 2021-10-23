package task

import (
	"io"

	"github.com/tooploox/oya/pkg/template"
)

type MockTask struct{}

func (MockTask) Exec(workDir string, args []string, scope template.Scope, stdout, stderr io.Writer) error {
	return nil
}
