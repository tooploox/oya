package oyafile

import "io"

type Hook interface {
	Exec(env map[string]string, stdout, stderr io.Writer) error
}

type ScriptedHook struct {
	Name string
	Script
	Shell string
}

func (h ScriptedHook) Exec(env map[string]string, stdout, stderr io.Writer) error {
	return h.Script.Exec(env, stdout, stderr, h.Shell)
}

type BuiltinHook struct {
	Name   string
	OnExec func(env map[string]string, stdout, stderr io.Writer) error
}

func (h BuiltinHook) Exec(env map[string]string, stdout, stderr io.Writer) error {
	return h.OnExec(env, stdout, stderr)
}
