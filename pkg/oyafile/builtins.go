package oyafile

import (
	"fmt"
	"io"

	"github.com/bilus/oya/pkg/task"
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
func (o *Oyafile) bindTasks(taskName task.Name, t task.Task, stdout, stderr io.Writer) (map[string]func() string, error) {
	tasks := make(map[string]func() string)

	importAlias, _ := taskName.Split()

	o.Tasks.ForEach(func(tn task.Name, _ task.Task, _ task.Meta) error {
		alias, baseName := tn.Split()
		if alias == importAlias {
			tasks[baseName] = func() string {
				return fmt.Sprintf("%s run %s\n", o.OyaCmd, tn)
			}
		}
		return nil
	})

	return tasks, nil
}
