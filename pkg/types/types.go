package types

import (
	"strings"
)

type Alias string

func (a Alias) String() string {
	return string(a)
}

func (a Alias) IsEmpty() bool {
	return len(string(a)) == 0
}

type ImportPath string

func (p ImportPath) Host() Host {
	if strings.HasPrefix(string(p), "github.com/") {
		return HostGithub
	}
	return HostUnknown
}

func (p ImportPath) String() string {
	return string(p)
}

type Host int

const (
	HostUnknown = iota
	HostGithub
)
