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

func Run(workDir, taskName string, recurse, changeset bool, positionalArgs []string, flags map[string]string,
	autoScope bool, stdout, stderr io.Writer) error {
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
	oldOyaScope, ok := lookupOyaScope()
	if ok && oldOyaScope != "" {
		tn = tn.Aliased(types.Alias(oldOyaScope))
	}
	if alias != "" && autoScope {
		alias, _ = tn.Split()
		if err := setOyaScope(alias.String()); err != nil {
			return err
		}
	}
	defer setOyaScope(oldOyaScope) // Mostly useful in tests, child processes naturally implement stacks.

	return p.Run(workDir, tn, recurse, changeset, toScope(positionalArgs, flags).Merge(values), stdout, stderr)
}

func toScope(positionalArgs []string, flags map[string]string) template.Scope {
	return template.Scope{
		"Args":  positionalArgs,
		"Flags": camelizeFlags(flags),
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
