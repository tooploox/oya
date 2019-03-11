package task

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/tooploox/oya/pkg/template"
	"mvdan.cc/sh/interp"
	"mvdan.cc/sh/syntax"
)

type Script struct {
	Script string
	Shell  string
	Scope  *template.Scope
}

func (s Script) Exec(workDir string, values template.Scope, stdout, stderr io.Writer) error {
	scope := values.Merge(*s.Scope)
	defines := defines(scope)
	script := strings.Join(defines, "; ") + "\n" + s.Script

	file, err := syntax.NewParser().Parse(strings.NewReader(script), "")
	if err != nil {
		return err
	}
	r, _ := interp.New(interp.StdIO(nil, stdout, stderr),
		interp.Module(interp.DefaultExec),
		interp.Dir(workDir),
		interp.Env(nil),
		interp.Params(args(scope)...))
	ctx := context.Background()
StmtLoop:
	for _, stmt := range file.Stmts {
		err := r.Run(ctx, stmt)
		switch err.(type) {
		case nil:
		case interp.ExitStatus:
		case interp.ShellExitStatus:
			break StmtLoop
		default:
			// BUG(bilus): Add line error.
			return err // set -e
		}
	}
	return nil
}

func defines(scope template.Scope) []string {
	dfs := append([]string{}, "declare -A Oya=()")

	for k, v := range scope.Flat() {
		ks, ok := k.(string)
		if !ok {
			continue
		}
		dfs = append(dfs, define(ks, v))
	}
	return dfs
}

func define(k, v interface{}) string {
	switch vt := v.(type) {
	case string:
		return fmt.Sprintf("Oya[%v]='%v'", k, escapeQuotes(vt))
	default:
		return fmt.Sprintf("Oya[%v]='%v'", k, v)
	}
}

func escapeQuotes(s string) string {
	s1 := strings.Replace(s, "\\", "\\\\", -1)
	return strings.Replace(s1, "'", "\\'", -1)
}

func args(scope template.Scope) []string {
	v, ok := scope["Args"]
	if !ok || v == nil {
		return nil
	}
	arr, ok := v.([]string)
	if !ok {
		return nil
	}
	args := make([]string, 0)
	args = append(args, "--")
	args = append(args, arr...)
	return args
}
