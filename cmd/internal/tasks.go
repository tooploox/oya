package internal

import (
	"fmt"
	"io"
	"path/filepath"
	"text/tabwriter"

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
	oyafiles, err := p.LoadOyafiles(workDir, recurse, changeset)
	if err != nil {
		return err
	}

	first := true
	for _, o := range oyafiles {
		relPath, err := filepath.Rel(workDir, o.Path)
		if err != nil {
			return err
		}
		if !first {
			fmt.Fprintln(w)
		} else {
			first = false
		}

		fmt.Fprintf(w, "# in ./%s\n", relPath)

		err = o.Tasks.ForEachSorted(func(taskName task.Name, task task.Task, meta task.Meta) error {
			if !taskName.IsBuiltIn() {
				if len(meta.Doc) > 0 {
					fmt.Fprintf(w, "oya run %s\t# %s\n", taskName, meta.Doc)
				} else {
					fmt.Fprintf(w, "oya run %s\t\n", taskName)
				}
			}
			return nil
		})

		if err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}
