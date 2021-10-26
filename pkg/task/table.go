package task

import (
	"github.com/tooploox/oya/pkg/types"
)

type Table struct {
	tasks map[Name]Task
	meta  map[Name]Meta
}

func NewTable() Table {
	return Table{
		tasks: make(map[Name]Task),
		meta:  make(map[Name]Meta),
	}
}

func (tt Table) LookupTask(name Name) (Task, bool) {
	t, ok := tt.tasks[name]
	return t, ok
}

func (tt Table) AddTask(name Name, task Task) {
	tt.tasks[name] = task
}

func (tt Table) AddDoc(taskName Name, s string) {
	tt.meta[taskName] = Meta{
		Doc:              s,
		OriginalTaskName: taskName,
	}
}

func (tt Table) ImportTasks(alias types.Alias, other Table) {
	for name, t := range other.tasks {
		// TODO: Detect if task already set.
		tt.AddTask(Name(name).Aliased(alias), t)
	}
}

// Expose copies tasks under an alias to global scope (without the alias)
// never overriding the existing global tasks.
func (tt Table) Expose(alias types.Alias) {
	for name, task := range tt.tasks {
		if name.IsAliased(alias) {
			globalName := name.Unaliased()
			_, ok := tt.LookupTask(globalName)
			if !ok {
				meta := tt.meta[name]
				meta.OriginalTaskName = name
				tt.addTaskWithMeta(globalName, task, meta)
			}
		}
	}
}

func (tt Table) addTaskWithMeta(name Name, task Task, meta Meta) {
	tt.tasks[name] = task
	tt.meta[name] = meta
}

// ForEach invokes the callback for each task.
func (tt Table) ForEach(f func(taskName Name, task Task, meta Meta) error) error {
	for taskName, task := range tt.tasks {
		meta := tt.meta[taskName]
		if err := f(taskName, task, meta); err != nil {
			return err
		}
	}
	return nil
}

// ForEachSorted invokes the callback for each task, the tasks sorted alphabetically.
func (tt Table) ForEachSorted(f func(taskName Name, task Task, meta Meta) error) error {
	taskNames := make([]Name, 0, len(tt.tasks))
	for taskName := range tt.tasks {
		taskNames = append(taskNames, taskName)
	}
	Sort(taskNames)

	for _, taskName := range taskNames {
		task := tt.tasks[taskName]
		meta := tt.meta[taskName]
		if err := f(taskName, task, meta); err != nil {
			return err
		}
	}
	return nil
}
