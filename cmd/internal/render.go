package internal

import (
	"fmt"
	"io"

	"github.com/bilus/oya/pkg/project"
	"github.com/bilus/oya/pkg/template"
)

type ErrNoAlias struct {
	Alias       string
	OyafilePath string
}

func (err ErrNoAlias) Error() string {
	return fmt.Sprintf("Unknown import alias %q in %v", err.Alias, err.OyafilePath)
}

func Render(oyafilePath, templatePath, outputPath, alias string, stdout, stderr io.Writer) error {
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

	err = o.Build(installDir)
	if err != nil {
		return err
	}

	var values template.Scope
	if found {
		if alias != "" {
			av, ok := o.Values[alias]
			if !ok {
				return ErrNoAlias{Alias: alias, OyafilePath: oyafilePath}
			}
			aliasScope, ok := av.(template.Scope)
			if !ok {
				return ErrNoAlias{Alias: alias, OyafilePath: oyafilePath}
			}
			values = aliasScope
		} else {
			values = o.Values
		}
	}

	return template.RenderAll(templatePath, outputPath, values)
}
