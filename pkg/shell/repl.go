package shell

import (
	"context"
	"fmt"
	"io"

	"github.com/tooploox/oya/pkg/template"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// StartREPL starts a shell REPL, using TTY if available for better user
// experience with autocompletion and command history, falling back to simpler
// version if no TTY available, the function blocking until the shell exits.
func StartREPL(values template.Scope, workDir string, stdin io.Reader, stdout, stderr io.Writer, customPreamble *string) error {
	ctx := context.Background()

	results := make(chan Result)
	prompt := NewPrompt(values, stdin, stdout, stderr, results)

	lastErr := make(chan error, 1)
	go func() {
		lastErr <- evalLoop(ctx, values, workDir, prompt.Stdin(), prompt.Stdout(), prompt.Stderr(), results, customPreamble)
	}()

	prompt.Run()

	return <-lastErr
}

// evalLoop runns the eval part of REPL.
func evalLoop(ctx context.Context, values template.Scope, workDir string, stdin io.Reader, stdout, stderr io.Writer, results chan Result, customPreamble *string) error {
	r, err := interp.New(interp.StdIO(nil, stdout, stderr),
		interp.Dir("."),
		interp.Env(nil))
	if err != nil {
		return err
	}

	parser := syntax.NewParser()

	if err := addPreamble(ctx, r, parser, values, customPreamble); err != nil {
		return err
	}

	var lastErr error

	err = parser.Interactive(stdin, func(stmts []*syntax.Stmt) bool {
		if parser.Incomplete() {
			results <- Result{incomplete: true}
			return true
		}
		for _, stmt := range stmts {
			lastErr = r.Run(ctx, stmt)
			if r.Exited() {
				results <- Result{exited: true}
				return false
			}
		}
		results <- Result{}
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
