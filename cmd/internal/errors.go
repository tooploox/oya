package internal

import (
	"fmt"
	"io"
	"strings"

	"github.com/tooploox/oya/pkg/errors"
	"github.com/tooploox/oya/pkg/task"
)

func HandleError(out io.Writer, err error) int {
	PrintError(out, err)
	return ExitCode(err)
}

func PrintError(out io.Writer, err error) {
	switch err := err.(type) {
	case errors.Error:
		printErrorWithTrace(out, err)
	default:
		printSep(out)
		printError(out, err)
	}
}

func ExitCode(err error) int {
	var sfe task.ErrScriptFail
	if errors.As(err, &sfe) {
		return sfe.ExitCode
	}
	return 1
}

func printSep(out io.Writer) {
	sepChar := "-"
	sepWidth := 78
	fmt.Fprintln(out, strings.Repeat(sepChar, sepWidth))
}

func printErrorWithTrace(out io.Writer, err errors.Error) {
	if cause := err.Cause(); cause != nil {
		PrintError(out, cause)
	}
	printSep(out)
	err.Print(out)
	fmt.Fprintln(out)
}

func printError(out io.Writer, err error) {
	fmt.Fprintln(out, "Error:", err)
	fmt.Fprintln(out)
}
