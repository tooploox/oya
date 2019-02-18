package internal

import (
	"fmt"
	"io"
	"strings"

	"github.com/bilus/oya/pkg/project"
	"github.com/bilus/oya/pkg/template"
	"github.com/pkg/errors"
)

type ErrNoScope struct {
	Scope       string
	OyafilePath string
}

func (err ErrNoScope) Error() string {
	return fmt.Sprintf("Scope not found in %v: %q missing or cannot be used as a scope", err.OyafilePath, err.Scope)
}

func Render(oyafilePath, templatePath, outputPath string, autoScope bool, scopeSelector string, stdout, stderr io.Writer) error {
	installDir, err := installDir()
	if err != nil {
		return err
	}
	proj, err := project.Detect(outputPath, installDir)
	if err != nil {
		return err
	}

	o, found, err := proj.Oyafile(oyafilePath)
	if err != nil {
		return err
	}

	dt, err := proj.Deps()
	if err != nil {
		return err
	}

	err = o.Build(dt)
	if err != nil {
		return err
	}

	var values template.Scope
	if found {
		if autoScope && scopeSelector == "" {
			scopeSelector, _ = lookupOyaScope()
		}
		// BUG(bilus): Breaking encapsulation here (see task.Name#Split)
		if scopeSelector != "" {
			scopeSelectorParts := strings.Split(scopeSelector, ".")
			values, err = resolveScope(scopeSelectorParts, o.Values)
		} else {
			values, err = resolveScope(nil, o.Values)
		}
		if err != nil {
			// BUG(bilus): Ignoring err.
			return ErrNoScope{Scope: scopeSelector, OyafilePath: oyafilePath}
		}
	}

	return template.RenderAll(templatePath, outputPath, values)
}

func resolveScope(scopeSelector []string, scope template.Scope) (template.Scope, error) {
	if len(scopeSelector) == 0 {
		return scope, nil
	}

	scopeName := scopeSelector[0]
	potentialScope, ok := scope[scopeName]
	if !ok {
		return nil, errors.Errorf("Missing key %q", scopeName)
	}
	subScope, ok := template.ParseScope(potentialScope)
	if !ok {
		return nil, errors.Errorf("Unsupported scope under %q", scopeName)
	}
	return resolveScope(scopeSelector[1:], subScope)
}
