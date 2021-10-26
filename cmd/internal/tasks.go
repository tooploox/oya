package internal

import (
	"io"
	"text/tabwriter"

	"github.com/tooploox/oya/cmd/internal/printers"
	"github.com/tooploox/oya/pkg/project"
	"github.com/tooploox/oya/pkg/task"
)

// Tasks loads Oyafiles for the project detected at workDir and prints their
// tasks to stdout, after ensuring that all required packs have been installed.
func Tasks(workDir string, recurse, changeset bool, stdout, stderr io.Writer) error {
	w := tabwriter.NewWriter(stdout, 0, 0, 2, ' ', 0)

	installDir, err := installDir()
	if err != nil {
		return err
	}
	p, err := project.Detect(workDir, installDir)
	if err != nil {
		return err
	}
	err = p.InstallPacks()
	if err != nil {
		return err
	}

	printer := printers.NewTaskList(workDir)
	err = p.ForEachTask(workDir, recurse, changeset, false, // No built-ins.
		func(i int, oyafilePath string, taskName task.Name, task task.Task, meta task.Meta) error {
			return printer.AddTask(taskName, meta, oyafilePath)
		})
	printer.Print(w)
	w.Flush()
	return err
}
