package oyafile

import (
	"io"

	"github.com/bilus/oya/pkg/template"
)

type Task interface {
	Exec(stdout, stderr io.Writer) error
}

type ScriptedTask struct {
	Name string
	Script
	Shell string
	Scope *template.Scope
}

func (h ScriptedTask) Exec(stdout, stderr io.Writer) error {
	return h.Script.Exec(*h.Scope, stdout, stderr, h.Shell)
}

type BuiltinTask struct {
	Name   string
	OnExec func(stdout, stderr io.Writer) error
}

func (h BuiltinTask) Exec(stdout, stderr io.Writer) error {
	return h.OnExec(stdout, stderr)
}
