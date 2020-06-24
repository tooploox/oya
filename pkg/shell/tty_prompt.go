package shell

import (
	"fmt"
	"io"
	"log"
	"reflect"

	"github.com/c-bata/go-prompt"
	"github.com/tooploox/oya/pkg/template"
)

type TTYPrompt struct {
	prompt         *prompt.Prompt
	stdin          io.Reader
	stdout, stderr io.Writer
	results        chan Result
}

func (p TTYPrompt) Stdin() io.Reader {
	return p.stdin
}

func (p TTYPrompt) Stdout() io.Writer {
	return p.stdout
}

func (p TTYPrompt) Stderr() io.Writer {
	return p.stderr
}

func (p TTYPrompt) Run() {
	p.prompt.Run()
}

func completer(scope template.Scope) prompt.Completer {
	return func(d prompt.Document) []prompt.Suggest {
		completions := make([]prompt.Suggest, 0, len(scope))
		for k, v := range scope {
			if reflect.TypeOf(v).Kind() != reflect.Func {
				completions = append(completions, completion(k, v))
			}
		}
		w := d.GetWordBeforeCursor()
		if w == "" {
			return nil
		}
		candidates := prompt.FilterHasPrefix(completions, "${", false) // Bash variable.
		return prompt.FilterFuzzy(candidates, w, false)
	}
}

func completion(k, v interface{}) prompt.Suggest {
	var description string

	return prompt.Suggest{
		Text:        fmt.Sprintf("${Oya[%v]}", k),
		Description: description,
	}
}

type ConsoleWriter struct {
	prompt.ConsoleWriter
	suppressed bool
}

// WriteRaw to write raw byte array.
func (w *ConsoleWriter) WriteRaw(data []byte) {
	if !w.suppressed {
		w.ConsoleWriter.WriteRaw(data)
	}
}

// Write to write safety byte array by removing control sequences.
func (w *ConsoleWriter) Write(data []byte) {
	if !w.suppressed {
		w.ConsoleWriter.Write(data)
	}
}

// WriteStr to write raw string.
func (w *ConsoleWriter) WriteRawStr(data string) {
	if !w.suppressed {
		w.ConsoleWriter.WriteRawStr(data)
	}
}

// WriteStr to write safety string by removing control sequences.
func (w *ConsoleWriter) WriteStr(data string) {
	if !w.suppressed {
		w.ConsoleWriter.WriteStr(data)
	}
}

func newTTYPrompt(scope template.Scope, results chan Result) Prompt {
	stdin, evalIn := io.Pipe()
	stdout := &ConsoleWriter{ConsoleWriter: prompt.NewStdoutWriter()}
	stderr := &ConsoleWriter{ConsoleWriter: prompt.NewStderrWriter()}

	exited := false
	return TTYPrompt{
		prompt: prompt.New(
			func(line string) {
				_, err := evalIn.Write([]byte(line))
				if err != nil {
					log.Fatalf("Internal error sending data to eval loop: %v", err)
				}
				_, err = evalIn.Write([]byte("\n"))
				if err != nil {
					log.Fatalf("Internal error sending data to eval loop: %v", err)
				}
				result := <-results // Synchronize with eval loop.
				stdout.Flush()
				stderr.Flush()
				if result.exited {
					// Prevent prompt from appearing.
					// See https://github.com/c-bata/go-prompt/issues/182
					stdout.suppressed = true
					stderr.suppressed = true
					exited = true
				}
			},
			completer(scope),
			prompt.OptionWriter(stdout),
			prompt.OptionSetExitCheckerOnInput(func(line string, breakline bool) bool { return exited }),
			prompt.OptionMaxSuggestion(10),
			prompt.OptionSuggestionBGColor(prompt.White),
			prompt.OptionSuggestionTextColor(prompt.Black),
		),
		stdin:   stdin,
		stdout:  writerAdapter{stdout},
		stderr:  writerAdapter{stderr},
		results: results,
	}
}

type writerAdapter struct {
	prompt.ConsoleWriter
}

func (w writerAdapter) Write(data []byte) (int, error) {
	w.ConsoleWriter.Write(data)
	return len(data), nil
}
