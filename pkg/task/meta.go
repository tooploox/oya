package task

type Meta struct {
	Doc string

	// OriginalTaskName, when different from the actual name associated with the
	// task, contains the name of the original task from an imported package
	// after it has been exposed, becoming accessible in the global scope as
	// taskName. For example, task name is "deploy" while originalTaskName =
	// "docker.deploy".
	//
	// For regular tasks it should either be equal to the task's actual name
	// or left empty.
	OriginalTaskName Name
}

func (meta Meta) IsTaskExposed(taskName Name) bool {
	return taskName != meta.OriginalTaskName && len(meta.OriginalTaskName) > 0
}
