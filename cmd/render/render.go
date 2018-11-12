package render

import (
	"io"
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

	if found {
		return t.Render(stdout, o.Values)
	} else {
		return t.Render(stdout, nil)
	}
}

func detectRootDir(currentDir string) (string, error) {
	return oyafile.DetectRootDir(currentDir)
}
