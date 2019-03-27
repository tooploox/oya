package internal

import (
	"fmt"
	"io"
	"path/filepath"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/tooploox/oya/pkg/project"
	"github.com/tooploox/oya/pkg/task"
)

func Tasks(workDir string, recurse, changeset bool, stdout, stderr io.Writer) error {

	w := tabwriter.NewWriter(stdout, 0, 0, 2, ' ', 0)

	installDir, err := project.InstallDir()
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
	oyafiles, err := p.RunTargets(workDir, recurse, changeset)
	if err != nil {
		return err
	}

	dependencies, err := p.Deps()
	if err != nil {
		return err
	}

	first := true
	for _, o := range oyafiles {
		relPath, err := filepath.Rel(workDir, o.Path)
		if err != nil {
			return err
		}

		err = o.Build(dependencies)
		if err != nil {
			return errors.Wrapf(err, "error in %v", o.Path)
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
					fmt.Fprintf(w, "oya %s\t# %s\n", taskName, meta.Doc)
				} else {
					fmt.Fprintf(w, "oya %s\t\n", taskName)
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

func println(out io.Writer, s string) {
	out.Write([]byte(fmt.Sprintf("%s\n", s)))
}
