package types

import (
	"strings"
)

type Alias string
type ImportPath string

func (p ImportPath) Host() Host {
	if strings.HasPrefix(string(p), "github.com/") {
		return HostGithub
	}
	return HostUnknown
}

type Host int

const (
	HostUnknown = iota
	HostGithub
)
