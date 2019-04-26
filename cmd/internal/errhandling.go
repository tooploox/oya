package internal

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/tooploox/oya/pkg/errors"
)

var HEADER_FMT = "---------------------- %v ---------------------- %v\n"

func HandleError(err error) {
	switch err := err.(type) {
	case errors.Error:
		fmt.Fprintln(os.Stderr, "---------------------------------------------------------------")
		handleError(err)
		os.Exit(1) // BUG(bilus): propagate exit code.
	default:
		handleUnknownErr(err)
	}
}

func handleError(err errors.Error) {
	out := os.Stderr
	if cause := err.Cause(); cause != nil {
		cause, ok := cause.(errors.Error)
		if ok {
			handleError(cause)

		} else {
			fmt.Fprintln(out, "Error:", cause)
		}

		fmt.Fprintln(out)
		fmt.Fprintln(out, "---------------------------------------------------------------")
	}
	err.Print(out)
	fmt.Fprintln(out)
}

func formatHeader(out io.Writer, title, path string) {
	// BUG(bilus): Auto-size to keep the same width
	fmt.Fprintf(out, HEADER_FMT, title, path)
}

func handleUnknownErr(err error) {
	out := bytes.Buffer{}
	fmt.Fprintln(&out, "---------------------------------------------------------------")
	fmt.Fprintf(&out, "Error: %v\n", err.Error())
	os.Stderr.Write(out.Bytes())
	os.Exit(1)
}

// func handleTaskFail(err oyafile.ErrTaskFail) {
// 	out := bytes.Buffer{}
// 	formatTaskFail(&out, err)
// 	switch cause := err.Cause.(type) {
// 	case task.ErrScriptFail:
// 		out.WriteString("\n")
// 		out.WriteString("Cause: \n")
// 		out.WriteString("\n")
// 		formatScriptFail(&out, cause)
// 		out.WriteString("\n")
// 		os.Stderr.Write(out.Bytes())
// 		os.Exit(cause.ExitCode)
// 	default:
// 		os.Stderr.Write(out.Bytes())
// 		os.Exit(1)
// 	}
// }

// func formatTaskFail(out *bytes.Buffer, err oyafile.ErrTaskFail) {
// 	formatHeader(out, "SCRIPT ERROR", err.OyafilePath)
// 	var showArgs string
// 	if len(err.Args) > 0 {
// 		showArgs = fmt.Sprintf(" invoked with arguments %q", strings.Join(err.Args, " "))
// 	}
// 	fmt.Fprintf(out, "Error in task %q%v: %v\n", string(err.TaskName), showArgs, err.Cause.Error())
// 	if err.ImportPath != nil {
// 		fmt.Fprintf(out, "  (task imported from %v)\n", *err.ImportPath)
// 	}
// }

// func formatScriptFail(out *bytes.Buffer, err task.ErrScriptFail) {
// 	lines := strings.Split(err.Script, "\n")
// 	var start uint = 1
// 	if err.Line > 1 {
// 		start = err.Line - 1
// 	}
// 	digits := int(math.Log10(float64(err.Line)) + 1)
// 	lineFmt := fmt.Sprintf("%%v %%%vv %%v\n", digits)
// 	for i := start; i <= err.Line; i++ {
// 		var marker string
// 		if i == err.Line {
// 			marker = ">"
// 		} else {
// 			marker = " "
// 		}
// 		fmt.Fprintf(out, lineFmt, marker, i, lines[i-1])
// 	}
// 	fmt.Fprintf(out, lineFmt, " ", " ", colMarker(err.Column))
// }

// func colMarker(col uint) string {
// 	return strings.Repeat(" ", int(col)-1) + "^"
// }

// func handleRenderFail(err template.ErrRenderFail) {
// 	out := bytes.Buffer{}
// 	formatRenderFail(&out, err)
// 	os.Stderr.Write(out.Bytes())
// 	os.Exit(1)
// }

// func formatRenderFail(out *bytes.Buffer, err template.ErrRenderFail) {
// 	cwd, _ := os.Getwd()
// 	relPath, _ := filepath.Rel(cwd, err.TemplatePath)
// 	formatHeader(out, "RENDER ERROR", relPath)
// 	fmt.Fprint(out, "Error rendering template\n")
// 	out.WriteString("\n")
// 	out.WriteString("Cause: ")
// 	out.WriteString(err.Cause.Error())
// 	out.WriteString("\n\n")
// }

// func handleErrExplodingDeps(err deptree.ErrExplodingDeps) {
// 	out := bytes.Buffer{}
// 	formatHeader(&out, "DEPS ERROR", err.ProjectRootDir)
// 	fmt.Fprint(&out, "Error getting dependencies\n")
// 	switch cause := err.Cause.(type) {
// 	case mvs.ErrResolvingReqs:
// 		traceCount := len(cause.Trace)
// 		if traceCount > 0 {
// 			for i := range cause.Trace {
// 				j := traceCount - i - 1
// 				fmt.Fprintf(&out, "  required by %v\n", cause.Trace[j])
// 			}
// 			fmt.Fprintf(&out, "  required by %v\n", filepath.Join(err.ProjectRootDir, raw.DefaultName))
// 		}
// 	default:

// 	}
// 	out.WriteString("\n")
// 	out.WriteString("Cause: ")
// 	out.WriteString(err.Cause.Error())
// 	out.WriteString("\n\n")
// 	os.Stderr.Write(out.Bytes())
// 	os.Exit(1)
// }
