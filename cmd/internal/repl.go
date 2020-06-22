package internal

import (
	"io"
)

func REPL(workDir string, stdin io.Reader, stdout, stderr io.Writer) error {
	project, err := prepareProject(workDir)
	if err != nil {
		return err
	}
	return project.StartREPL(workDir, stdin, stdout, stderr)
}
