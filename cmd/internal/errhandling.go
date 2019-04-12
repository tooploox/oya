package internal

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/tooploox/oya/pkg/deptree"
	"github.com/tooploox/oya/pkg/mvs"
	"github.com/tooploox/oya/pkg/oyafile"
	"github.com/tooploox/oya/pkg/raw"
	"github.com/tooploox/oya/pkg/task"
)

var HEADER_FMT = "--- %v ---------------------- %v\n"

func HandleError(err error) {
	switch err := err.(type) {
	case oyafile.ErrTaskFail:
		handleTaskFail(err)
	case ErrRenderFail:
		handleRenderFail(err)
	case deptree.ErrExplodingDeps:
		handleErrExplodingDeps(err)
	default:
		handleUnknownErr(err)
	}
}

func formatHeader(out *bytes.Buffer, title, path string) {
	// BUG(bilus): Auto-size to keep the same width
	fmt.Fprintf(out, HEADER_FMT, title, path)
}

func handleTaskFail(err oyafile.ErrTaskFail) {
	out := bytes.Buffer{}
	formatTaskFail(&out, err)
	switch cause := err.Cause.(type) {
	case task.ErrScriptFail:
		out.WriteString("\n")
		out.WriteString("Cause: \n")
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
	formatHeader(out, "SCRIPT ERROR", err.OyafilePath)
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

func handleRenderFail(err ErrRenderFail) {
	out := bytes.Buffer{}
	formatRenderFail(&out, err)
	os.Stderr.Write(out.Bytes())
	os.Exit(1)
}

func formatRenderFail(out *bytes.Buffer, err ErrRenderFail) {
	formatHeader(out, "RENDER ERROR", err.TemplatePath)
	fmt.Fprint(out, "Error rendering template\n")
	out.WriteString("\n")
	out.WriteString("Cause: ")
	out.WriteString(err.Cause.Error())
	out.WriteString("\n\n")
}

func handleUnknownErr(err error) {
	out := bytes.Buffer{}
	formatHeader(&out, "ERROR", "")
	fmt.Fprintf(&out, "Error: %v\n", err.Error())
	os.Stderr.Write(out.Bytes())
	os.Exit(1)
}

func handleErrExplodingDeps(err deptree.ErrExplodingDeps) {
	out := bytes.Buffer{}
	formatHeader(&out, "DEPS ERROR", err.ProjectRootDir)
	fmt.Fprint(&out, "Error getting dependencies\n")
	switch cause := err.Cause.(type) {
	case mvs.ErrResolvingReqs:
		traceCount := len(cause.Trace)
		if traceCount > 0 {
			for i := range cause.Trace {
				j := traceCount - i - 1
				fmt.Fprintf(&out, "  required by %v\n", cause.Trace[j])
			}
			fmt.Fprintf(&out, "  required by %v\n", filepath.Join(err.ProjectRootDir, raw.DefaultName))
		}
	default:

	}
	out.WriteString("\n")
	out.WriteString("Cause: ")
	out.WriteString(err.Cause.Error())
	out.WriteString("\n\n")
	os.Stderr.Write(out.Bytes())
	os.Exit(1)
}
