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

// OyaCmdOverride is used in tests, to override the path to the current oya executable.
// It is used to invoke other tasks from a task body.
// When tests are run, the current process executable path points to the test runner
// so it has to be overridden (with 'go run oya.go', roughly speaking).
var OyaCmdOverride *string

type Script struct {
	Script string
	Shell  string
	Scope  *template.Scope
}

func (s Script) Exec(workDir string, args []string, values template.Scope, stdout, stderr io.Writer) error {
	scope := values.Merge(*s.Scope)
	defines := defines(scope)
	script := strings.Join(defines, "; ") + ";\n" + s.Script
	if OyaCmdOverride != nil {
		script = *OyaCmdOverride + "; " + script
	}

	file, err := syntax.NewParser().Parse(strings.NewReader(script), "")
	if err != nil {
		return err
	}
	r, _ := interp.New(interp.StdIO(nil, stdout, stderr),
		interp.Module(interp.DefaultExec),
		interp.Dir(workDir),
		interp.Env(nil),
		interp.Params(toParams(args)...))
	ctx := context.Background()
	for _, stmt := range file.Stmts {
		err := r.Run(ctx, stmt)
		switch err.(type) {
		case nil:
		case interp.ExitStatus:
			errCode := err.(interp.ExitStatus)
			if errCode != 0 {
				return fmt.Errorf("task exited with code %d", errCode)
			}
			return nil
		case interp.ShellExitStatus:
			errCode := err.(interp.ShellExitStatus)
			if errCode != 0 {
				return fmt.Errorf("task exited with code %d", errCode)
			}
			return nil
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

func toParams(taskArgs []string) []string {
	args := make([]string, 0, len(taskArgs)+1)
	args = append(args, "--")
	args = append(args, taskArgs...)
	return args
}
