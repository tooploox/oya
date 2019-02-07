package oyafile

import (
	"fmt"
	"sort"
)

type Meta struct {
	Doc string
}

type TaskTable struct {
	tasks map[TaskName]Task
	meta  map[TaskName]Meta
}

func newTaskTable() TaskTable {
	return TaskTable{
		tasks: make(map[TaskName]Task),
		meta:  make(map[TaskName]Meta),
	}
}

func (tt TaskTable) LookupTask(name TaskName) (Task, bool) {
	t, ok := tt.tasks[name]
	return t, ok
}

func (tt TaskTable) AddTask(name TaskName, task Task) {
	tt.tasks[name] = task
}

func (tt TaskTable) AddDoc(taskName TaskName, s string) {
	tt.meta[taskName] = Meta{
		Doc: s,
	}
}

func (tt TaskTable) ImportTasks(alias Alias, other TaskTable) {
	for key, task := range other.tasks {
		// TODO: Detect if task already set.
		tt.AddTask(TaskName(fmt.Sprintf("%v.%v", alias, key)), task)
	}
}

func (tt TaskTable) ForEach(f func(taskName TaskName, task Task, meta Meta) error) error {
	for taskName, task := range tt.tasks {
		meta, _ := tt.meta[taskName]
		if err := f(taskName, task, meta); err != nil {
			return err
		}
	}
	return nil
}

type TaskNames []TaskName

func (names TaskNames) Len() int {
	return len(names)
}

func (names TaskNames) Swap(i, j int) {
	names[i], names[j] = names[j], names[i]
}

func (names TaskNames) Less(i, j int) bool {
	return string(names[i]) < string(names[j])
}

func (tt TaskTable) ForEachSorted(f func(taskName TaskName, task Task, meta Meta) error) error {
	taskNames := make([]TaskName, 0, len(tt.tasks))
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
