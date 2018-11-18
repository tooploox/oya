package oyafile

import (
	"io"

	"github.com/bilus/oya/pkg/template"
)

type Hook interface {
	Exec(stdout, stderr io.Writer) error
}

type ScriptedHook struct {
	Name string
	Script
	Shell string
	Scope *template.Scope
}

func (h ScriptedHook) Exec(stdout, stderr io.Writer) error {
	return h.Script.Exec(*h.Scope, stdout, stderr, h.Shell)
}

type BuiltinHook struct {
	Name   string
	OnExec func(stdout, stderr io.Writer) error
}

func (h BuiltinHook) Exec(stdout, stderr io.Writer) error {
	return h.OnExec(stdout, stderr)
}
