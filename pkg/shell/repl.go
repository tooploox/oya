package shell

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tooploox/oya/pkg/template"
	"golang.org/x/crypto/ssh/terminal"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// StartREPL starts a shell REPL, using TTY if available for better user
// experience with autocompletion and command history, falling back to simpler
// version if no TTY available, the function blocking until the shell exits.
func StartREPL(values template.Scope, workDir string, stdin io.Reader, stdout, stderr io.Writer, customPreamble *string) error {
	ctx := context.Background()
	// f, ok := stdin.(*os.File)
	// fmt.Println("file?", ok)
	// if ok {
	// 	log.Debugf("Terminal detected")
	// } else {
	// 	log.Debugf("WARN: No terminal detected")
	// }
	//
	_, haveTerminal := detectTerminal(stdin)
	if !haveTerminal {
		log.Println("WARN: No terminal detected")
	}

	var (
		evalStdin              io.Reader
		evalStdout, evalStderr io.Writer

		prompt *Prompt
	)

	if haveTerminal {
		prompt = NewPrompt()
		evalStdin = prompt.Stdin
		evalStdout = prompt.Stdout
		evalStderr = prompt.Stderr

	} else {
		evalStdin = stdin
		evalStdout = stdout
		evalStderr = stderr
	}

	lastErr := make(chan error, 1)
	go func() {
		lastErr <- evalLoop(ctx, values, workDir, evalStdin, evalStdout, evalStderr, customPreamble, prompt == nil)
	}()

	if prompt != nil {
		prompt.Run()
	}

	return <-lastErr
}

// detectTerminal determines if the Reader passed as stdin is a file and a TTY.
func detectTerminal(stdin io.Reader) (*os.File, bool) {
	f, ok := stdin.(*os.File)
	if !ok {
		return nil, false
	}
	return f, terminal.IsTerminal(int(f.Fd()))
}

// evalLoop runns the eval part of REPL.
func evalLoop(ctx context.Context, values template.Scope, workDir string, stdin io.Reader, stdout, stderr io.Writer, customPreamble *string, showPrompt bool) error {
	const promptStr = "$ "
	const lineCont = "> "

	r, err := interp.New(interp.StdIO(nil, stdout, stderr),
		interp.Dir("."),
		interp.Env(nil))
	if err != nil {
		return err
	}

	parser := syntax.NewParser()
	if showPrompt {
		fmt.Fprint(stdout, promptStr)
	}

	if err := addPreamble(ctx, r, parser, values, customPreamble); err != nil {
		return err
	}

	var lastErr error

	err = parser.Interactive(stdin, func(stmts []*syntax.Stmt) bool {
		if parser.Incomplete() {
			if showPrompt {
				fmt.Fprint(stdout, lineCont)
			}
			return true
		}
		for _, stmt := range stmts {
			lastErr = r.Run(ctx, stmt)
			stdout.(IOWriterAdapter).ConsoleWriter.Flush()
			if r.Exited() {
				return false
			}
		}
		if showPrompt {
			fmt.Fprint(stdout, promptStr)
		}
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
