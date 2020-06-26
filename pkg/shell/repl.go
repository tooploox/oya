package shell

import (
	"context"
	"io"

	"github.com/tooploox/oya/pkg/template"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// StartREPL starts a shell REPL, using TTY if available for better user
// experience with autocompletion and command history, falling back to simpler
// version if no TTY available, the function blocking until the shell exits.
func StartREPL(values template.Scope, workDir string, stdin io.Reader, stdout, stderr io.Writer, customPreamble *string) error {
	ctx, cancel := context.WithCancel(context.Background())

	// Accept up to 256 lines typed before TTY prompt appears,
	// preventing REPL getting stuck.
	results := make(chan Result, 256)
	prompt := NewPrompt(values, stdin, stdout, stderr, results)

	lastErr := make(chan error, 1)
	go func() {
		lastErr <- evalLoop(ctx, values, workDir, prompt.Stdin(), prompt.Stdout(), prompt.Stderr(), results, customPreamble)
	}()

	prompt.Run()
	cancel()
	prompt.Shutdown() // Close stdin read from in the Interactive call below.

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

	for {
		err = parser.Interactive(stdin, func(stmts []*syntax.Stmt) bool {
			var lastErr error
			if parser.Incomplete() {
				results <- Result{incomplete: true}
				return ctx.Err() != context.Canceled
			}
			for _, stmt := range stmts {
				lastErr = r.Run(ctx, stmt)
				if r.Exited() {
					results <- Result{exited: true}
					return false
				}
			}
			results <- Result{err: lastErr}
			return ctx.Err() != context.Canceled
		})
		// Ctrl-d pressed or 'exit'.
		if err == nil || ctx.Err() == context.Canceled {
			return nil
		}
		results <- Result{err: err}
	}

	return nil
}
