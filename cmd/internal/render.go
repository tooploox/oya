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
	autoScope bool, scopePath string, overrides map[string]interface{},
	stdout, stderr io.Writer) error {
	installDir, err := installDir()
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
			scopePath, _ = lookupOyaScope()
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

	err = overrideValues(values, overrides)
	if err != nil {
		return err
	}

	return template.RenderAll(templatePath, excludedPaths, outputPath, values)
}

func overrideValues(values template.Scope, overrides map[string]interface{}) error {
	for path, val := range overrides {
		// Force overriding existing paths that have different "shapes".
		// What AssocAt does is it will create scopes along the path if they don't exist.
		// If an intermediate path does exist and is a simple value (not a Scope), it will
		// force it's conversion to a scope, thus losing the original value.
		// For example:
		// values: {"foo": "xxx"}
		// overrides: foo.bar="yyy"
		// result: {"foo": {"bar": "yyy"}}
		// With force set to false, the function would fail for the input ^.
		force := true
		err := values.AssocAt(path, val, force)
		if err != nil {
			return err
		}
	}
	return nil
}
