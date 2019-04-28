package internal

import (
	"fmt"
	"io"
	"strings"

	"github.com/tooploox/oya/pkg/errors"
)

func HandleError(out io.Writer, err error) {
	printSep(out)
	switch err := err.(type) {
	case errors.Error:
		printError(out, err)
	default:
		printUnknownErr(out, err)
	}
}

func printSep(out io.Writer) {
	sepChar := "-"
	sepWidth := 78
	fmt.Fprintln(out, strings.Repeat(sepChar, sepWidth))
}

func printError(out io.Writer, err errors.Error) {
	if cause := err.Cause(); cause != nil {
		if causeError, ok := cause.(errors.Error); ok {
			printError(out, causeError)
		} else {
			fmt.Fprintln(out, "Error:", cause)
			fmt.Fprintln(out)
		}

		printSep(out)
	}
	err.Print(out)
	fmt.Fprintln(out)
}

func printUnknownErr(out io.Writer, err error) {
	fmt.Fprintf(out, "Error: %v\n", err.Error())
}
