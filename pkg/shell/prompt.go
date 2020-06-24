package shell

import (
	"fmt"
	"io"
	"log"

	"github.com/c-bata/go-prompt"
)

type Prompt struct {
	*prompt.Prompt
	Stdin          io.Reader
	Stdout, Stderr io.Writer
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "${Oya[foo]}", Description: "bar"},
		{Text: "${Oya[somePassword]}", Description: "******"},
		{Text: "${Oya[first]}", Description: "111"},
	}
	w := d.GetWordBeforeCursor()
	if w == "" {
		return nil
	}
	return prompt.FilterHasPrefix(s, w, true)
}

func NewPrompt() *Prompt {
	stdin, evalIn := io.Pipe()
	stdout := prompt.NewStdoutWriter()
	stderr := prompt.NewStderrWriter()
	return &Prompt{
		Prompt: prompt.New(
			func(line string) {
				fmt.Println("EX:", line)
				_, err := evalIn.Write([]byte(line))
				if err != nil {
					log.Fatalf("Internal error sending data to eval loop: err", err)
				}
				_, err = evalIn.Write([]byte("\n"))
				if err != nil {
					log.Fatalf("Internal error sending data to eval loop: err", err)
				}
			},
			completer,
			prompt.OptionWriter(stdout),
		),
		Stdin:  stdin,
		Stdout: IOWriterAdapter{stdout},
		Stderr: IOWriterAdapter{stderr},
	}
}

type IOWriterAdapter struct {
	prompt.ConsoleWriter
}

func (w IOWriterAdapter) Write(data []byte) (int, error) {
	w.ConsoleWriter.Write(data)
	return len(data), nil
}
