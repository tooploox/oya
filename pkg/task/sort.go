package task

import "sort"

type TaskNames []Name

func (names TaskNames) Len() int {
	return len(names)
}

func (names TaskNames) Swap(i, j int) {
	names[i], names[j] = names[j], names[i]
}

func (names TaskNames) Less(i, j int) bool {
	return string(names[i]) < string(names[j])
}

func Sort(taskNames []Name) {
	sort.Sort(TaskNames(taskNames))
}
