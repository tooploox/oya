package internal

import (
	"io"
	"regexp"
	"strings"

	"github.com/tooploox/oya/pkg/project"
	"github.com/tooploox/oya/pkg/task"
	"github.com/tooploox/oya/pkg/template"
)

type Args struct {
	All        []string
	Positional []string
	Flags      map[string]string
}

func Run(workDir, taskName string, taskArgs Args, recurse, changeset bool, stdout, stderr io.Writer) error {
	installDir, err := InstallDir()
	if err != nil {
		return err
	}
	p, err := project.Detect(workDir, installDir)
	if err != nil {
		return err
	}
	err = p.InstallPacks()
	if err != nil {
		return err
	}
	values, err := p.Values()
	if err != nil {
		return err
	}
	tn := task.Name(taskName)

	alias, _ := tn.Split()
	oldOyaScope, _ := lookupOyaScope()
	if err := setOyaScope(alias.String()); err != nil {
		return err
	}
	defer setOyaScope(oldOyaScope) // Mostly useful in tests, child processes naturally implement stacks.

	return p.Run(workDir, tn, recurse, changeset, taskArgs.All,
		toScope(taskArgs).Merge(values), stdout, stderr)
}

func toScope(taskArgs Args) template.Scope {
	return template.Scope{
		"Args":  taskArgs.Positional,
		"Flags": camelizeFlags(taskArgs.Flags),
	}
}

func camelizeFlags(flags map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range flags {
		result[camelize(k)] = v
	}
	return result
}

var sepRx = regexp.MustCompile("(-|_).")

// camelize turns - or _ separated identifiers into camel case.
// Example: "aa-bb" becomes "aaBb".
func camelize(s string) string {
	return sepRx.ReplaceAllStringFunc(s, func(match string) string {
		return strings.ToUpper(match[1:])
	})

}
