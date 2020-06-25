package shell

import (
	"context"
	"fmt"
	"strings"

	"github.com/tooploox/oya/pkg/template"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

var ErrPreambleExit = fmt.Errorf("Unexpected exit caused by script preamble")

type ErrPreambleRun struct {
	Err error
}

func (e ErrPreambleRun) Error() string {
	return fmt.Sprintf("Unexpected error in script preamble: %v", e.Err)
}

func addPreamble(ctx context.Context, runner *interp.Runner, parser *syntax.Parser, values template.Scope, custom *string) error {
	preamble := buildPreamble(values, custom)
	file, err := parser.Parse(strings.NewReader(preamble), "")
	if err != nil {
		return err
	}
	for _, stmt := range file.Stmts {
		if err := runner.Run(ctx, stmt); err != nil {
			return ErrPreambleRun{Err: err}
		}
		if runner.Exited() {
			return ErrPreambleExit
		}
	}
	return nil
}

func buildPreamble(scope template.Scope, customPreamble *string) string {
	declarations := declarations(scope)
	preamble := strings.Join(declarations, "; ") + ";"
	if customPreamble != nil {
		preamble = *customPreamble + "; " + preamble
	}
	return preamble
}

func declarations(scope template.Scope) []string {
	dfs := append([]string{}, "declare -A Oya=()")

	for k, v := range scope.Flat() {
		ks, ok := k.(string)
		if !ok {
			continue
		}
		dfs = append(dfs, declaration(ks, v))
	}
	return dfs
}

func declaration(k, v interface{}) string {
	switch vt := v.(type) {
	case string:
		return fmt.Sprintf("Oya[%v]='%v'", k, escapeQuotes(vt))
	default:
		return fmt.Sprintf("Oya[%v]='%v'", k, v)
	}
}

func escapeQuotes(s string) string {
	s1 := strings.Replace(s, "\\", "\\\\", -1)
	return strings.Replace(s1, "'", "\\'", -1)
}

func substitution(k interface{}) string {
	return fmt.Sprintf("${Oya[%v]}", k)
}
