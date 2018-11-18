package run

import (
	"io"
	"os"

	"github.com/bilus/oya/pkg/project"
)

func Run(workDir, hookName string, stdout, stderr io.Writer) error {
	if workDir == "" {
		var err error
		workDir, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	p, err := project.Detect(workDir)
	if err != nil {
		return err
	}

	return p.Run(workDir, hookName, stdout, stderr)
}
