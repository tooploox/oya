package task

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/tooploox/oya/pkg/errors"
	"github.com/tooploox/oya/pkg/template"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type ErrScriptFail struct {
	ExitCode int
	Message  string
}

func (e ErrScriptFail) Error() string {
	return e.Message
}

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
	head := strings.Join(defines, "; ") + ";"
	if OyaCmdOverride != nil {
		head = *OyaCmdOverride + "; " + head
	}
	headLines := uint(strings.Count(head, "\n") + 1)
	script := head + "\n" + s.Script
	file, err := syntax.NewParser().Parse(strings.NewReader(script), "")
	if err != nil {
		switch err := err.(type) {
		case syntax.ParseError:
			return s.errScriptFail(err.Pos, headLines, err, -1)
		default:
			return err
		}
	}
	r, err := interp.New(interp.StdIO(nil, stdout, stderr),
		interp.Dir(workDir),
		interp.Env(nil),
		interp.Params(toParams(args)...))
	if err != nil {
		return err
	}

	var lastFailure error
	ctx := context.Background()
LOOP:
	for _, stmt := range file.Stmts {
		err := r.Run(ctx, stmt)
		switch err := err.(type) {
		case nil:
			// Continue with next statement in script.
			lastFailure = nil
		case interp.ExitStatus:
			// Record to use if it's the last command.
			errCode := int(err)
			if errCode != 0 {
				lastFailure = s.errScriptFail(stmt.Pos(), headLines, err, errCode)
			}
		case interp.ShellExitStatus:
			// Shell interpreter exited.
			// Either early return (due to "exit" or "set -e"), or task is finished.
			errCode := int(err)
			if errCode != 0 {
				lastFailure = s.errScriptFail(stmt.Pos(), headLines, err, errCode)
			}
			break LOOP

		default:
			return s.errScriptFail(stmt.Pos(), headLines, err, -1)
		}
	}
	return lastFailure
}

func (s Script) errScriptFail(pos syntax.Pos, headLines uint, err error, exitCode int) error {
	return errors.New(
		ErrScriptFail{
			ExitCode: exitCode,
			Message:  err.Error(), // Simplify error trace.
		},
		errors.Location{
			Line: pos.Line() - headLines,
			Col:  pos.Col(),
			Snippet: errors.Snippet{
				Lines: strings.Split(s.Script, "\n"),
			},
		},
	)
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
