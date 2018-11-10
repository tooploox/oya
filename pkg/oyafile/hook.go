package oyafile

import "io"

type Hook interface {
	Exec(values Scope, stdout, stderr io.Writer) error
}

type ScriptedHook struct {
	Name string
	Script
	Shell string
}

func (h ScriptedHook) Exec(values Scope, stdout, stderr io.Writer) error {
	return h.Script.Exec(values, stdout, stderr, h.Shell)
}

type BuiltinHook struct {
	Name   string
	OnExec func(values Scope, stdout, stderr io.Writer) error
}

func (h BuiltinHook) Exec(values Scope, stdout, stderr io.Writer) error {
	return h.OnExec(values, stdout, stderr)
}
