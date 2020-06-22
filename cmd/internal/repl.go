package internal

import (
	"context"
	"fmt"
	"io"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

const prompt = "$ "
const lineCont = "> "

func REPL(stdin io.Reader, stdout, stderr io.Writer) error {
	r, err := interp.New(interp.StdIO(nil, stdout, stderr),
		interp.Dir("."),
		interp.Env(nil))
	if err != nil {
		return err
	}

	var lastErr error

	ctx := context.Background()

	parser := syntax.NewParser()
	fmt.Fprint(stdout, prompt)

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
			// return s.errScriptFail(err.Pos, headLines, err, -1)
		default:
			return err
		}
	}

	return lastErr
}
