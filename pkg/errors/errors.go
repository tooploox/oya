package errors

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/pkg/errors"
)

type Error struct {
	err   error
	cause error
	trace []Location
}

func (e Error) Error() string {
	return e.err.Error()
}

func (e Error) Trace() []Location {
	return e.trace
}

func (e Error) Print(out io.Writer) {
	fmt.Fprintln(out, "Error:", e.err)
	if len(e.trace) > 0 {
		fmt.Fprintln(out)
		for _, location := range e.trace {
			fmt.Fprintf(out, "  ")
			location.Print(out)
		}
	}
}

func (e Error) Cause() error {
	return e.cause
}

func New(err error, trace ...Location) error {
	return Error{
		err:   err,
		trace: trace,
	}
}
func Wrap(cause error, err error, trace ...Location) error {
	return Error{
		err:   err,
		cause: cause,
		trace: trace,
	}
}

func Wrapf(cause error, fmt string, args ...interface{}) error {
	return Error{
		err:   errors.Wrapf(cause, fmt, args...),
		cause: cause,
	}
}

func Errorf(fmt string, args ...interface{}) error {
	return errors.Errorf(fmt, args...)
}

type Location struct {
	Name        string // Optional.
	VerboseName string // Optional.
	Line, Col   uint   // Optional, 0 - unused.
	Snippet            // Optional.
}

func (l Location) Print(out io.Writer) {
	var name string
	if len(l.VerboseName) > 0 {
		name = l.VerboseName
	} else {
		name = l.Name
	}
	if len(name) > 0 {
		fmt.Fprintf(out, "%v", name)
		fmt.Fprintf(out, " ")
	}
	if l.Line != 0 {
		fmt.Fprintf(out, "at line %v", l.Line)
	}
	if l.Col != 0 {
		fmt.Fprintf(out, ", column %v", l.Col)
	}
	fmt.Fprintln(out)

	if !l.Snippet.IsEmpty() {
		fmt.Fprintln(out)
		l.Snippet.Print(out, l.Line, l.Col)
	}
}

type Trace interface {
	Trace() []Location
}

type Snippet struct {
	LineOffset uint
	Lines      []string
}

func (s Snippet) IsEmpty() bool {
	return len(s.Lines) == 0
}

func (r Snippet) Print(out io.Writer, line, col uint) {
	if r.IsEmpty() {
		return
	}

	var start uint = 1
	if line > 1 {
		start = line - 1
	}
	digits := int(math.Log10(float64(line)) + 1)
	lineFmt := fmt.Sprintf("%%v %%%vv%%s %%v\n", digits)
	lineSep := "|"
	for i := start; i <= line; i++ {
		var marker string
		if i == line {
			marker = ">"
		} else {
			marker = " "
		}
		li := i - r.LineOffset
		if li > 0 {
			fmt.Fprintf(out, lineFmt, marker, i, lineSep, r.Lines[li-1])
		}
	}
	if col > 0 {
		fmt.Fprintf(out, lineFmt, " ", " ", " ", colMarker(col))
	}
}

func colMarker(col uint) string {
	return strings.Repeat(" ", int(col)-1) + "^"
}
