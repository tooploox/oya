package oyafile

import (
	"io"
	"strings"

	"github.com/bilus/oya/pkg/template"
)

type Task interface {
	Exec(workDir string, params template.Scope, stdout, stderr io.Writer) error
	GetName() string
	IsBuiltIn() bool
}

type ScriptedTask struct {
	Name string
	Script
	Shell string
	Scope *template.Scope
}

func (h ScriptedTask) Exec(workDir string, params template.Scope, stdout, stderr io.Writer) error {
	return h.Script.Exec(workDir, params.Merge(*h.Scope), stdout, stderr, h.Shell)
}

func (t ScriptedTask) GetName() string {
	return t.Name
}

func (t ScriptedTask) IsBuiltIn() bool {
	firstChar := t.Name[0:1]
	return firstChar == strings.ToUpper(firstChar)
}

type BuiltinTask struct {
	Name   string
	OnExec func(stdout, stderr io.Writer) error
}

func (h BuiltinTask) Exec(workDir string, params template.Scope, stdout, stderr io.Writer) error {
	return h.OnExec(stdout, stderr)
}

func (t BuiltinTask) GetName() string {
	return t.Name
}

func (t BuiltinTask) IsBuiltIn() bool {
	return true
}
