package oyafile

import "io"

type Hook interface {
	Exec(env map[string]string, stdout, stderr io.Writer, shell string) error
}

type ScriptedHook struct {
	Name string
	Script
}

func (h ScriptedHook) Exec(env map[string]string, stdout, stderr io.Writer, shell string) error {
	return h.Script.Exec(env, stdout, stderr, shell)
}

type BuiltinHook struct {
	Name   string
	OnExec func(env map[string]string, stdout, stderr io.Writer, shell string) error
}

func (h BuiltinHook) Exec(env map[string]string, stdout, stderr io.Writer, shell string) error {
	return h.OnExec(env, stdout, stderr, shell)
}
