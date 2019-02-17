package internal

import (
	"fmt"
	"io"

	"github.com/bilus/oya/pkg/project"
	"github.com/bilus/oya/pkg/template"
)

type ErrNoScope struct {
	Scope       string
	OyafilePath string
}

func (err ErrNoScope) Error() string {
	return fmt.Sprintf("Scope %q not found in %v", err.Scope, err.OyafilePath)
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
		if scopeSelector != "" {
			av, ok := o.Values[scopeSelector]
			if !ok {
				return ErrNoScope{Scope: scopeSelector, OyafilePath: oyafilePath}
			}
			selectedScope, ok := av.(template.Scope)
			if !ok {
				return ErrNoScope{Scope: scopeSelector, OyafilePath: oyafilePath}
			}
			values = selectedScope
		} else {
			values = o.Values
		}
	}

	return template.RenderAll(templatePath, outputPath, values)
}
