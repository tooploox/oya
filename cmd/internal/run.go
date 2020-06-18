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
	tn := task.Name(taskName)

	// If OYA_SCOPE alias is present, prefix task with the alias.
	// Then update OYA_SCOPE to contain the newly built task alias,
	// so when oya run are recursively called in imported tasks,
	// the scope is correctly resolved.
	// Examples:
	// | OYA_SCOPE | task    | new OYA_SCOPE | aliased task |
	// |           | xxx     |               | xxx          |
	// | foo       | xxx     | foo           | foo.xxx      |
	// | foo       | bar.xxx | foo.bar       | foo.bar.xxx  |
	passedOyaScope, _ := lookupOyaScope()
	if len(passedOyaScope) > 0 {
		tn = tn.Aliased(types.Alias(passedOyaScope))
	}
	alias, _ := tn.Split()
	if err := setOyaScope(alias.String()); err != nil {
		return err
	}
	defer mustSetOyaScope(passedOyaScope) // Mostly useful in tests, child processes naturally implement stacks.

	oyaCmd, found := lookupOyaCmd()
	if found {
		// Tests only.
		task.OyaCmdOverride = &oyaCmd
	}

	return p.Run(workDir, tn, recurse, changeset, taskArgs.All,
		toScope(taskArgs), stdout, stderr)
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
