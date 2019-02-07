package oyafile

import "strings"

type TaskName string

func MakeTaskName(n string) (TaskName, error) {
	return TaskName(n), nil
}

func (n TaskName) Split() (Alias, string) {
	name := string(n)
	parts := strings.Split(name, ".")
	switch len(parts) {
	case 0:
		return "", ""
	case 1:
		return "", parts[0]
	default:
		return Alias(parts[0]), strings.Join(parts[1:], ".")
	}
}

func (n TaskName) IsBuiltIn() bool {
	firstChar := string(n)[0:1]
	return firstChar == strings.ToUpper(firstChar)
}
