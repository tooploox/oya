package task

import (
	"strings"

	"github.com/tooploox/oya/pkg/types"
)

type Name string

func (n Name) Split() (types.Alias, string) {
	name := string(n)
	parts := strings.Split(name, ".")
	switch len(parts) {
	case 0:
		return "", ""
	case 1:
		return "", parts[0]
	default:
		return types.Alias(parts[0]), strings.Join(parts[1:], ".")
	}
}

func (n Name) IsBuiltIn() bool {
	firstChar := string(n)[0:1]
	return firstChar == strings.ToUpper(firstChar)
}
