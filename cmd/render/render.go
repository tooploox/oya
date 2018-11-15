package render

import (
	"io"
	"os"
	"path/filepath"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/template"
)

func Render(oyafilePath, templatePath, outputPath string, stdout, stderr io.Writer) error {
	rootDir, err := detectRootDir(outputPath)
	if err != nil {
		return err
	}

	o, found, err := oyafile.Load(oyafilePath, rootDir)
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

		return renderFile(path, filepath.Join(outputPath, relPath), values)
	})
}

func detectRootDir(currentDir string) (string, error) {
	return oyafile.DetectRootDir(currentDir)
}

func renderFile(templatePath, outputPath string, values oyafile.Scope) error {
	t, err := template.Load(templatePath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(outputPath), 0700)

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	return t.Render(out, values)
}
