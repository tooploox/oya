package oyafile

import (
	"fmt"
	"io"

	"github.com/bilus/oya/pkg/template"
)

func (o *Oyafile) addBuiltIns() error {
	o.Values = o.Values.Merge(o.defaultValues())
	return nil
}

func (o *Oyafile) defaultValues() template.Scope {
	return template.Scope{
		"BasePath": o.Dir,
	}
}

// bindTasks returns a map of functions allowing invoking other tasks via $Tasks.xyz().
// It makes invokable only tasks defined in the same Oyafile, stripping away any aliases, so the tasks are accessible names exactly as they appear in a given Oyafile.
func (o *Oyafile) bindTasks(task Task, stdout, stderr io.Writer) (map[string]func() string, error) {
	tasks := make(map[string]func() string)

	importAlias, _ := task.SplitName()

	o.Tasks.ForEach(func(taskName string, task Task, _ Meta) error {
		alias, baseName := task.SplitName()
		if alias == importAlias {
			tasks[baseName] = func() string {
				return fmt.Sprintf("%s run %s\n", o.OyaCmd, taskName)
			}
		}
		return nil
	})

	return tasks, nil
}
