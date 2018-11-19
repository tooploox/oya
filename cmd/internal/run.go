package internal

import (
	"io"

	"github.com/bilus/oya/pkg/project"
)

func Run(workDir, taskName string, stdout, stderr io.Writer) error {
	p, err := project.Detect(workDir)
	if err != nil {
		return err
	}

	return p.Run(workDir, taskName, stdout, stderr)
}
