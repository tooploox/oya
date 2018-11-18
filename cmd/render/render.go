package render

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/project"
	"github.com/bilus/oya/pkg/template"
	log "github.com/sirupsen/logrus"
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

	var values oyafile.Scope
	if found {
		values = o.Values
	}

	return filepath.Walk(templatePath, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(templatePath, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			// templatePath is a path to a file.
			relPath = filepath.Base(templatePath)
		}

		filePath, err := renderString(filepath.Join(outputPath, relPath), values)
		if err != nil {
			return err
		}
		log.Println(outputPath, "+", relPath, "=", filePath)
		return renderFile(path, filePath, values)
	})
}

func renderFile(templatePath, outputPath string, values oyafile.Scope) error {
	t, err := template.Load(templatePath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(outputPath), 0700)
	if err != nil {
		return err
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	return t.Render(out, values)
}

func renderString(templateSource string, values oyafile.Scope) (string, error) {
	t, err := template.Parse(templateSource)
	if err != nil {
		return "", err
	}
	out := new(bytes.Buffer)
	err = t.Render(out, values)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
