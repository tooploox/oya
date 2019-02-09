package internal

import (
	"io"
	"regexp"
	"strings"

	"github.com/bilus/oya/pkg/project"
	"github.com/bilus/oya/pkg/task"
	"github.com/bilus/oya/pkg/template"
)

func Run(workDir, taskName string, positionalArgs []string, flags map[string]string, stdout, stderr io.Writer) error {
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
	return p.Run(workDir, task.Name(taskName), toScope(positionalArgs, flags).Merge(values), stdout, stderr)
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
