package task

import (
	"fmt"
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
		lastIdx := len(parts) - 1
		return types.Alias(strings.Join(parts[0:lastIdx], ".")), parts[lastIdx]
	}
}

func (n Name) IsBuiltIn() bool {
	firstChar := string(n)[0:1]
	return firstChar == strings.ToUpper(firstChar)
}

func (n Name) Aliased(alias types.Alias) Name {
	return Name(fmt.Sprintf("%v.%v", alias, n))
}
