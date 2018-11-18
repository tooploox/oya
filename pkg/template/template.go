package template

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	kasia "github.com/ziutek/kasia.go"
)

type Scope map[string]interface{}

type Template interface {
	Render(out io.Writer, values interface{}) error
}

type kasiaTemplate struct {
	impl *kasia.Template
}

func Load(path string) (Template, error) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(string(source))
}

func Parse(source string) (Template, error) {
	kt, err := kasia.Parse(source)
	if err != nil {
		return nil, err
	}
	return kasiaTemplate{impl: kt}, nil
}

func RenderAll(templatePath, outputPath string, values Scope) error {
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
		return renderFile(path, filePath, values)
	})
}

func renderFile(templatePath, outputPath string, values Scope) error {
	t, err := Load(templatePath)
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

func renderString(templateSource string, values Scope) (string, error) {
	t, err := Parse(templateSource)
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

type emptyValues struct{}

func (t kasiaTemplate) Render(out io.Writer, values interface{}) error {
	if values == nil {
		values = emptyValues{}
	}
	return t.impl.Run(out, values)
}
