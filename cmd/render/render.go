package render

import (
	"io"
	"os"
	"path/filepath"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/template"
)

func Render(oyafilePath, templatePath string, stdout, stderr io.Writer) error {
	t, err := template.Load(templatePath)
	if err != nil {
		return err
	}

	rootDir, err := detectRootDir(filepath.Dir(oyafilePath))
	if err != nil {
		return err
	}

	o, found, err := oyafile.Load(oyafilePath, rootDir)
	if err != nil {
		return err
	}

	fname := filepath.Base(templatePath)

	out, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	if found {
		return t.Render(out, o.Values)
	}
	return t.Render(out, nil)
}

func detectRootDir(currentDir string) (string, error) {
	return oyafile.DetectRootDir(currentDir)
}
