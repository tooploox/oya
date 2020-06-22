package shell

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/tooploox/oya/pkg/template"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

const prompt = "$ "
const lineCont = "> "

func StartREPL(workDir string, values template.Scope, stdin io.Reader, stdout, stderr io.Writer, customPreamble *string) error {
	r, err := interp.New(interp.StdIO(nil, stdout, stderr),
		interp.Dir("."),
		interp.Env(nil))
	if err != nil {
		return err
	}

	var lastErr error

	ctx := context.Background()

	parser := syntax.NewParser()
	fmt.Fprint(stdout, prompt)

	preamble, _ := buildPreamble(values, customPreamble)
	if err != nil {
		return err
	}
	file, err := parser.Parse(strings.NewReader(preamble), "")
	for _, stmt := range file.Stmts {
		if err := r.Run(ctx, stmt); err != nil {
			log.Fatalf("Unexpected error in by script preamble: %v", err)
		}
		if r.Exited() {
			log.Fatalf("Unexpected exit caused by script preamble")
		}
	}

	err = parser.Interactive(stdin, func(stmts []*syntax.Stmt) bool {
		if parser.Incomplete() {
			fmt.Fprint(stdout, lineCont)
			return true
		}
		for _, stmt := range stmts {
			lastErr = r.Run(ctx, stmt)
			if r.Exited() {
				return false
			}
		}
		fmt.Fprint(stdout, prompt)
		return true
	})
	if err != nil {
		switch err := err.(type) {
		case syntax.ParseError:
			return fmt.Errorf("Error: %v", err)
			// TODO: Better error reporting?
			// return s.errScriptFail(err.Pos, preambleLines, err, -1)
		default:
			return err
		}
	}

	return lastErr
}

func buildPreamble(scope template.Scope, customPreamble *string) (string, uint) {
	declarations := declarations(scope)
	preamble := strings.Join(declarations, "; ") + ";"
	if customPreamble != nil {
		preamble = *customPreamble + "; " + preamble
	}
	numPreambleLines := uint(strings.Count(preamble, "\n") + 1)
	return preamble, numPreambleLines
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
