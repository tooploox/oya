package internal

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/project"
)

func Tasks(workDir string, stdout, stderr io.Writer) error {
	p, err := project.Detect(workDir)
	if err != nil {
		return err
	}
	tasks, err := p.Tasks(workDir, stdout, stderr)
	if err != nil {
		return err
	}

	relWorkDir, err := filepath.Rel(p.Root.RootDir, workDir)
	if err != nil {
		return err
	}

	first := true

	for path, tt := range tasks {
		relPath, err := filepath.Rel(relWorkDir, path)
		if err != nil {
			return err
		}

		if !first {
			println(stdout, "")
		} else {
			first = false
		}

		println(stdout, fmt.Sprintf("# in ./%s", relPath))

		err = tt.ForEach(func(taskName string, task oyafile.Task) error {
			if !task.IsBuiltIn() {
				println(stdout, fmt.Sprintf("oya run %s", taskName))
			}
			return nil
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func println(out io.Writer, s string) {
	out.Write([]byte(fmt.Sprintf("%s\n", s)))
}
