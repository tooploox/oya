package oyafile

import (
	"io"

	"github.com/bilus/oya/pkg/template"
)

type Task interface {
	Exec(workDir string, stdout, stderr io.Writer) error
}

type ScriptedTask struct {
	Name string
	Script
	Shell string
	Scope *template.Scope
}

func (h ScriptedTask) Exec(workDir string, stdout, stderr io.Writer) error {
	return h.Script.Exec(workDir, *h.Scope, stdout, stderr, h.Shell)
}

type BuiltinTask struct {
	Name   string
	OnExec func(stdout, stderr io.Writer) error
}

func (h BuiltinTask) Exec(workDir string, stdout, stderr io.Writer) error {
	return h.OnExec(stdout, stderr)
}
