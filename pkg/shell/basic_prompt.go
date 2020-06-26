package shell

import (
	"fmt"
	"io"
)

type BasicPrompt struct {
	stdin          io.Reader
	stdout, stderr io.Writer
	results        chan Result
}

func newBasicPrompt(stdin io.Reader, stdout, stderr io.Writer, results chan Result) Prompt {
	return BasicPrompt{
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
		results: results,
	}
}

func (p BasicPrompt) Stdin() io.Reader {
	return p.stdin
}

func (p BasicPrompt) Stdout() io.Writer {
	return p.stdout
}

func (p BasicPrompt) Stderr() io.Writer {
	return p.stderr
}

func (p BasicPrompt) Shutdown() {}

func (p BasicPrompt) Run() {
	const prompt = "$ "
	const lineCont = "> "

	fmt.Fprint(p.stdout, prompt)
	for r := range p.results {
		if r.exited {
			break
		} else if r.incomplete {
			fmt.Fprint(p.stdout, lineCont)
		} else if r.err != nil {
			fmt.Fprintln(p.stderr, r.err)
			fmt.Fprint(p.stdout, prompt)
		} else {
			fmt.Fprint(p.stdout, prompt)
		}
	}
}
