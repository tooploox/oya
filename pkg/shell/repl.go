package shell

import (
	"context"
	"fmt"
	"io"

	"github.com/tooploox/oya/pkg/template"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

const prompt = "$ "
const lineCont = "> "

func StartREPL(workDir string, values template.Scope, stdin io.Reader, stdout, stderr io.Writer, customPreamble *string) error {
	ctx := context.Background()

	r, err := interp.New(interp.StdIO(nil, stdout, stderr),
		interp.Dir("."),
		interp.Env(nil))
	if err != nil {
		return err
	}

	parser := syntax.NewParser()
	fmt.Fprint(stdout, prompt)

	if err := addPreamble(ctx, r, parser, values, customPreamble); err != nil {
		return err
	}

	var lastErr error

	err = parser.Interactive(stdin, func(stmts []*syntax.Stmt) bool {
		if parser.Incomplete() {
			fmt.Fprint(stdout, lineCont)
			return true
		}
		for _, stmt := range stmts {
			lastErr = r.Run(ctx, stmt)
			if r.Exited() {
				return false
			}
		}
		fmt.Fprint(stdout, prompt)
		return true
	})
	if err != nil {
		switch err := err.(type) {
		case syntax.ParseError:
			return fmt.Errorf("Error: %v", err)
			// TODO: Better error reporting?
			// return s.errScriptFail(err.Pos, preambleLines, err, -1)
		default:
			return err
		}
	}

	return lastErr
}
