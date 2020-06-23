package shell

import (
	"context"
	"io"
	"strings"

	"github.com/tooploox/oya/pkg/errors"
	"github.com/tooploox/oya/pkg/template"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type ErrExecFail struct {
	ExitCode int
	Message  string
}

func (e ErrExecFail) Error() string {
	return e.Message
}

func Exec(script string, workDir string, args []string, values template.Scope, stdout, stderr io.Writer, customPreamble *string) error {
	ctx := context.Background()

	r, err := interp.New(interp.StdIO(nil, stdout, stderr),
		interp.Dir(workDir),
		interp.Env(nil),
		interp.Params(toParams(args)...))
	if err != nil {
		return err
	}

	parser := syntax.NewParser()

	if err := addPreamble(ctx, r, parser, values, customPreamble); err != nil {
		return err
	}

	file, err := syntax.NewParser().Parse(strings.NewReader(script), "")
	if err != nil {
		switch err := err.(type) {
		case syntax.ParseError:
			return errExecFail(err.Pos, script, err, -1)
		default:
			return err
		}
	}

	var lastErr error

	for _, stmt := range file.Stmts {
		lastErr = r.Run(ctx, stmt)
		if lastErr != nil {
			if r.Exited() {
				exitStatus, ok := interp.IsExitStatus(lastErr)
				if !ok {
					exitStatus = 0
				}

				lastErr = errExecFail(stmt.Pos(), script, lastErr, int(exitStatus))
				break
			}
		}
		if r.Exited() {
			break
		}
	}
	return lastErr

}

func errExecFail(pos syntax.Pos, script string, err error, exitCode int) error {
	return errors.New(
		ErrExecFail{
			ExitCode: exitCode,
			Message:  err.Error(), // Simplify error trace.
		},
		errors.Location{
			Line: pos.Line(),
			Col:  pos.Col(),
			Snippet: errors.Snippet{
				Lines: strings.Split(script, "\n"),
			},
		},
	)
}

func toParams(taskArgs []string) []string {
	args := make([]string, 0, len(taskArgs)+1)
	args = append(args, "--")
	args = append(args, taskArgs...)
	return args
}
