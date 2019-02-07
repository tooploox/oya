package oyafile

import (
	"fmt"
	"sort"

	"github.com/bilus/oya/pkg/task"
	"github.com/bilus/oya/pkg/types"
)

type TaskTable struct {
	tasks map[task.Name]task.Task
	meta  map[task.Name]task.Meta
}

func newTaskTable() TaskTable {
	return TaskTable{
		tasks: make(map[task.Name]task.Task),
		meta:  make(map[task.Name]task.Meta),
	}
}

func (tt TaskTable) LookupTask(name task.Name) (task.Task, bool) {
	t, ok := tt.tasks[name]
	return t, ok
}

func (tt TaskTable) AddTask(name task.Name, task task.Task) {
	tt.tasks[name] = task
}

func (tt TaskTable) AddDoc(taskName task.Name, s string) {
	tt.meta[taskName] = task.Meta{
		Doc: s,
	}
}

func (tt TaskTable) ImportTasks(alias types.Alias, other TaskTable) {
	for key, t := range other.tasks {
		// TODO: Detect if task already set.
		tt.AddTask(task.Name(fmt.Sprintf("%v.%v", alias, key)), t)
	}
}

func (tt TaskTable) ForEach(f func(taskName task.Name, task task.Task, meta task.Meta) error) error {
	for taskName, task := range tt.tasks {
		meta, _ := tt.meta[taskName]
		if err := f(taskName, task, meta); err != nil {
			return err
		}
	}
	return nil
}

type TaskNames []task.Name

func (names TaskNames) Len() int {
	return len(names)
}

func (names TaskNames) Swap(i, j int) {
	names[i], names[j] = names[j], names[i]
}

func (names TaskNames) Less(i, j int) bool {
	return string(names[i]) < string(names[j])
}

func (tt TaskTable) ForEachSorted(f func(taskName task.Name, task task.Task, meta task.Meta) error) error {
	taskNames := make([]task.Name, 0, len(tt.tasks))
	for taskName := range tt.tasks {
		taskNames = append(taskNames, taskName)
	}

	sort.Sort(TaskNames(taskNames))
	for _, taskName := range taskNames {
		task := tt.tasks[taskName]
		meta, _ := tt.meta[taskName]
		if err := f(taskName, task, meta); err != nil {
			return err
		}
	}
	return nil
}
