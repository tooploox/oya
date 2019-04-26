package internal

import (
	"io"
	"regexp"
	"strings"

	"github.com/tooploox/oya/pkg/project"
	"github.com/tooploox/oya/pkg/task"
	"github.com/tooploox/oya/pkg/template"
	"github.com/tooploox/oya/pkg/types"
)

type Args struct {
	All        []string
	Positional []string
	Flags      map[string]string
}

func Run(workDir, taskName string, taskArgs Args, recurse, changeset bool, stdout, stderr io.Writer) error {
	installDir, err := installDir()
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
	passedOyaScope, found := lookupOyaScope()
	// BUG(bilus): Refactor using tn.Aliased
	var newOyaScope string
	if len(passedOyaScope) > 0 {
		newOyaScope = passedOyaScope + "." + alias.String()
	} else {
		newOyaScope = alias.String()
	}
	if err := setOyaScope(newOyaScope); err != nil {
		return err
	}
	defer setOyaScope(passedOyaScope) // Mostly useful in tests, child processes naturally implement stacks.

	oyaCmd, found := lookupOyaCmd()
	if found {
		// Tests only.
		task.OyaCmdOverride = &oyaCmd
	}

	if len(passedOyaScope) > 0 {
		tn = tn.Aliased(types.Alias(passedOyaScope))
	}

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
