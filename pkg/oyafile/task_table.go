package oyafile

import "fmt"

type TaskTable struct {
	impl map[string]Task
}

func newTaskTable() TaskTable {
	return TaskTable{impl: make(map[string]Task)}
}

func (tt TaskTable) LookupTask(name string) (Task, bool) {
	t, ok := tt.impl[name]
	return t, ok
}

func (tt TaskTable) AddTask(name string, task Task) {
	tt.impl[name] = task
}

func (tt TaskTable) AliasTasks(alias Alias, other TaskTable) {
	for key, task := range other.impl {
		// TODO: Detect if task already set.
		tt.AddTask(fmt.Sprintf("%v.%v", alias, key), task)
	}
}

func (tt TaskTable) ForEach(f func(taskName string, task Task) error) error {
	for taskName, task := range tt.impl {
		if err := f(taskName, task); err != nil {
			return err
		}
	}
	return nil
}
