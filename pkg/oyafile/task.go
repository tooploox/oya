package oyafile

import (
	"io"

	"github.com/bilus/oya/pkg/template"
)

type Task interface {
	Exec(workDir string, params template.Scope, stdout, stderr io.Writer) error
}

type ScriptedTask struct {
	Script
	Shell string
	Scope *template.Scope
}

func (h ScriptedTask) Exec(workDir string, params template.Scope, stdout, stderr io.Writer) error {
	return h.Script.Exec(workDir, params.Merge(*h.Scope), stdout, stderr, h.Shell)
}

type BuiltinTask struct {
	OnExec func(stdout, stderr io.Writer) error
}

func (h BuiltinTask) Exec(workDir string, params template.Scope, stdout, stderr io.Writer) error {
	return h.OnExec(stdout, stderr)
}
