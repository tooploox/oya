package project

import (
	"github.com/tooploox/oya/pkg/task"
)

func (p *Project) ForEachTask(workDir string, recurse, useChangeset, includeBuiltins bool,
	callback func(i int, oyafilePath string, taskName task.Name, task task.Task, meta task.Meta) error) error {

	oyafiles, err := p.loadOyafiles(workDir, recurse, useChangeset)
	if err != nil {
		return err
	}

	for _, o := range oyafiles {
		i := 0
		err = o.Tasks.ForEachSorted(func(taskName task.Name, task task.Task, meta task.Meta) error {
			if !taskName.IsBuiltIn() || includeBuiltins {
				if err := callback(i, o.Path, taskName, task, meta); err != nil {
					return err
				}
				i = i + 1
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
