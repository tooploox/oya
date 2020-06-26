package shell

import (
	"io"
	"log"
	"os"

	"github.com/tooploox/oya/pkg/template"
	"golang.org/x/crypto/ssh/terminal"
)

type Result struct {
	incomplete, exited bool
}

type Prompt interface {
	Run()
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
	Shutdown()
}

func NewPrompt(scope template.Scope, stdin io.Reader, stdout, stderr io.Writer, results chan Result) Prompt {
	_, haveTerminal := detectTerminal(stdin)
	if !haveTerminal {
		// TODO: Move it up.
		log.Println("WARN: No terminal detected")
		return newBasicPrompt(stdin, stdout, stderr, results)
	} else {
		return newTTYPrompt(scope, results) // Ignored stdin, stdout, stderr.
	}
}

// detectTerminal determines if the Reader passed as stdin is a file and a TTY.
func detectTerminal(stdin io.Reader) (*os.File, bool) {
	f, ok := stdin.(*os.File)
	if !ok {
		return nil, false
	}
	return f, terminal.IsTerminal(int(f.Fd()))
}
