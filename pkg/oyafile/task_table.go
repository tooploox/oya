package oyafile

import (
	"fmt"
	"sort"
)

type Meta struct {
	Doc string
}

type TaskTable struct {
	tasks map[string]Task
	meta  map[string]Meta
}

func newTaskTable() TaskTable {
	return TaskTable{
		tasks: make(map[string]Task),
		meta:  make(map[string]Meta),
	}
}

func (tt TaskTable) LookupTask(name string) (Task, bool) {
	t, ok := tt.tasks[name]
	return t, ok
}

func (tt TaskTable) AddTask(name string, task Task) {
	tt.tasks[name] = task
}

func (tt TaskTable) AddDoc(taskName string, s string) {
	tt.meta[taskName] = Meta{
		Doc: s,
	}
}

func (tt TaskTable) ImportTasks(alias Alias, other TaskTable) {
	for key, task := range other.tasks {
		// TODO: Detect if task already set.
		tt.AddTask(fmt.Sprintf("%v.%v", alias, key), task)
	}
}

func (tt TaskTable) ForEach(f func(taskName string, task Task, meta Meta) error) error {
	for taskName, task := range tt.tasks {
		meta, _ := tt.meta[taskName]
		if err := f(taskName, task, meta); err != nil {
			return err
		}
	}
	return nil
}

func (tt TaskTable) ForEachSorted(f func(taskName string, task Task, meta Meta) error) error {
	taskNames := make([]string, 0, len(tt.tasks))
	for taskName := range tt.tasks {
		taskNames = append(taskNames, taskName)
	}

	sort.Strings(taskNames)
	for _, taskName := range taskNames {
		task := tt.tasks[taskName]
		meta, _ := tt.meta[taskName]
		if err := f(taskName, task, meta); err != nil {
			return err
		}
	}
	return nil
}
