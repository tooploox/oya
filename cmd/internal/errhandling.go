package internal

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tooploox/oya/pkg/oyafile"
	"github.com/tooploox/oya/pkg/task"
)

func HandleError(err error) {
	switch err := err.(type) {
	case oyafile.ErrTaskFail:
		handleTaskFail(err)
	default:
		logrus.Println(err)
		os.Stderr.WriteString("Error: ")
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
}

func handleTaskFail(err oyafile.ErrTaskFail) {
	out := bytes.Buffer{}
	formatTaskFail(&out, err)
	switch cause := err.Cause.(type) {
	case task.ErrScriptFail:
		out.WriteString("\n")
		formatScriptFail(&out, cause)
		out.WriteString("\n")
		os.Stderr.Write(out.Bytes())
		os.Exit(cause.ExitCode)
	default:
		os.Stderr.Write(out.Bytes())
		os.Exit(1)
	}
}

func formatTaskFail(out *bytes.Buffer, err oyafile.ErrTaskFail) {
	fmt.Fprintf(out, "--- RUN ERROR ---------------------- %v\n", err.OyafilePath)
	var showArgs string
	if len(err.Args) > 0 {
		showArgs = fmt.Sprintf(" invoked with arguments %q", strings.Join(err.Args, " "))
	}
	fmt.Fprintf(out, "Error in task %q%v: %v\n", string(err.TaskName), showArgs, err.Cause.Error())
	if err.ImportPath != nil {
		fmt.Fprintf(out, "  (task imported from %v)\n", *err.ImportPath)
	}
}

func formatScriptFail(out *bytes.Buffer, err task.ErrScriptFail) {
	lines := strings.Split(err.Script, "\n")
	var start uint = 1
	if err.Line > 1 {
		start = err.Line - 1
	}
	digits := int(math.Log10(float64(err.Line)) + 1)
	lineFmt := fmt.Sprintf("%%v %%%vv %%v\n", digits)
	for i := start; i <= err.Line; i++ {
		var marker string
		if i == err.Line {
			marker = ">"
		} else {
			marker = " "
		}
		fmt.Fprintf(out, lineFmt, marker, i, lines[i-1])
	}
	fmt.Fprintf(out, lineFmt, " ", " ", colMarker(err.Column))
}

func colMarker(col uint) string {
	return strings.Repeat(" ", int(col)-1) + "^"
}
