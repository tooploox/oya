package template

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobwas/glob"
	"github.com/tooploox/oya/pkg/errors"
)

type ErrRenderFail struct {
}

func (e ErrRenderFail) Error() string {
	return "render error"
}

// Template represents a template that can be rendered using provided values.
type Template interface {
	Render(out io.Writer, values Scope) error
	RenderString(values Scope) (string, error)
}

// Load loads template from the path.
func Load(path string, delimiters Delimiters) (Template, error) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(string(source), delimiters)
}

// Parse parses template in the source string.
func Parse(source string, delimiters Delimiters) (Template, error) {
	return parsePlush(source, delimiters)
}

// RenderAll renders all templates in the path (directory or a single file) to an output path (directory or file) using the provided value scope.
func RenderAll(templatePath string, excludedPaths []string, outputPath string, values Scope, delimiters Delimiters) error {
	return filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
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
		if ok, err := pathMatches(excludedPaths, relPath); ok || err != nil {
			return err // err is nil if ok
		}

		filePath, err := renderString(filepath.Join(outputPath, relPath), values, delimiters)
		if err != nil {
			return errors.Wrap(
				err,
				ErrRenderFail{},
				errors.Location{
					Name:        path,
					VerboseName: fmt.Sprintf("in template %v", path),
				},
			)
		}
		err = renderFile(path, filePath, values, delimiters)
		if err != nil {
			return errors.Wrap(err,
				ErrRenderFail{},
				errors.Location{
					Name:        path,
					VerboseName: fmt.Sprintf("in template %v", path),
				},
			)
		}
		return nil
	})
}

func renderFile(templatePath, outputPath string, values Scope, delimiters Delimiters) error {
	t, err := Load(templatePath, delimiters)
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

func renderString(templateSource string, values Scope, delimiters Delimiters) (string, error) {
	t, err := Parse(templateSource, delimiters)
	if err != nil {
		return "", err
	}
	str, err := t.RenderString(values)
	if err != nil {
		return "", err
	}
	return str, nil
}

func pathMatches(patterns []string, path string) (bool, error) {
	for _, pattern := range patterns {
		p, err := glob.Compile(pattern)
		if err != nil {
			return false, err
		}
		if p.Match(path) {
			return true, nil
		}
	}
	return false, nil
}
