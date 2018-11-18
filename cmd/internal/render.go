package internal

import (
	"io"

	"github.com/bilus/oya/pkg/project"
	"github.com/bilus/oya/pkg/template"
)

func Render(oyafilePath, templatePath, outputPath string, stdout, stderr io.Writer) error {
	proj, err := project.Detect(outputPath)
	if err != nil {
		return err
	}

	o, found, err := proj.LoadOyafile(oyafilePath)
	if err != nil {
		return err
	}

	var values template.Scope
	if found {
		values = o.Values
	}

	return template.RenderAll(templatePath, outputPath, values)
}
