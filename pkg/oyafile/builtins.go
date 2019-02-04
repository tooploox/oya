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

func (o *Oyafile) bindTasks(params template.Scope, stdout, stderr io.Writer) (map[string]func() string, error) {
	tasks := make(map[string]func() string)

	o.Tasks.ForEach(func(taskName string, task Task, _ Meta) error {
		tasks[taskName] = func() string {
			// found, err := o.RunTask(taskName, params, stdout, stderr)
			// TODO: Use exe path
			return fmt.Sprintf("%s run %s\n", o.OyaCmd, taskName)
		}
		return nil
	})

	return tasks, nil
}
