package internal

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/tooploox/oya/pkg/project"
	"github.com/tooploox/oya/pkg/template"
)

type ErrNoScope struct {
	Scope       string
	OyafilePath string
}

func (err ErrNoScope) Error() string {
	return fmt.Sprintf("Scope not found in %v: %q missing or cannot be used as a scope", err.OyafilePath, err.Scope)
}

func Render(oyafilePath, templatePath string, excludedPaths []string, outputPath string,
	autoScope bool, scopePath string, stdout, stderr io.Writer) error {
	installDir, err := project.InstallDir()
	if err != nil {
		return err
	}
	oyafileFullPath, err := filepath.Abs(oyafilePath)
	if err != nil {
		return err
	}

	proj, err := project.Detect(oyafileFullPath, installDir)
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
		if autoScope && scopePath == "" {
			scopePath, _ = project.LookupOyaScope()
		}
		if scopePath != "" {
			values, err = o.Values.GetScopeAt(scopePath)
		} else {
			values = o.Values
		}
		if err != nil {
			// BUG(bilus): Ignoring err.
			return ErrNoScope{Scope: scopePath, OyafilePath: oyafilePath}
		}
	}

	return template.RenderAll(templatePath, excludedPaths, outputPath, values)
}
