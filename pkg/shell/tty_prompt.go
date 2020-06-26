package shell

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/tooploox/oya/pkg/template"
)

const maxPreviewLen = 40

type TTYPrompt struct {
	prompt         *prompt.Prompt
	stdin          io.Reader
	stdout, stderr io.Writer
	results        chan Result

	shutdown func()
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
	printWelcome(p.stdout)
	p.prompt.Run()
}

func (p TTYPrompt) Shutdown() {
	p.shutdown()
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
		stdin:    stdin,
		stdout:   writerAdapter{stdout},
		stderr:   writerAdapter{stderr},
		results:  results,
		shutdown: func() { stdin.Close() },
	}
}

func printWelcome(stdout io.Writer) {
	fmt.Fprintln(stdout, "Welcome to Oya REPL! To exit, press Ctrl-D.")
	fmt.Fprintln(stdout, "Type ${ to auto-complete Oya values.")
}

func completer(scope template.Scope) prompt.Completer {
	const trigger = "${"

	return func(d prompt.Document) []prompt.Suggest {
		values := scope.Flat()

		completions := make([]prompt.Suggest, 0, len(values))
		for k, v := range values {
			if reflect.TypeOf(v).Kind() != reflect.Func {

				completions = append(completions, completion(k, v))
			}
		}
		sort.Slice(completions,
			func(i, j int) bool {
				return completions[i].Text < completions[j].Text
			})
		w := d.GetWordBeforeCursor()
		if !strings.HasPrefix(w, trigger) {
			return nil
		}

		return prompt.FilterFuzzy(completions, strings.TrimPrefix(w, trigger), false)
	}
}

func previewValue(v interface{}) string {
	s, ok := v.(string)
	if ok {
		return shorten(firstLine(s), maxPreviewLen)
	} else {
		return "(object)"
	}
}

func completion(k, v interface{}) prompt.Suggest {
	return prompt.Suggest{
		Text:        substitution(k),
		Description: previewValue(v),
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

type writerAdapter struct {
	prompt.ConsoleWriter
}

func (w writerAdapter) Write(data []byte) (int, error) {
	w.ConsoleWriter.Write(data)
	return len(data), nil
}
