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
	// SplitName splits task name into import alias and base task name. For instance "docker.build" becomes ["docker", "build"]. If the task name is not aliased, import alias is empty. For instance "build" becomes ["", "build"]
	SplitName() (importAlias string, baseName string)
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

func (t ScriptedTask) SplitName() (string, string) {
	return splitTaskName(t.Name)
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

func (t BuiltinTask) SplitName() (string, string) {
	return splitTaskName(t.Name)
}

func splitTaskName(name string) (string, string) {
	parts := strings.Split(name, ".")
	switch len(parts) {
	case 0:
		return "", ""
	case 1:
		return "", parts[0]
	default:
		return parts[0], strings.Join(parts[1:], ".")
	}
}
