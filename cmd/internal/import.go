package internal

import (
	"io"
	"strings"

	"github.com/bilus/oya/pkg/project"
	"github.com/bilus/oya/pkg/raw"
	"github.com/pkg/errors"
)

func Import(workDir, uri string, stdout, stderr io.Writer) error {
	uriArr := strings.Split(uri, "/")
	alias := uriArr[len(uriArr)-1]

	proj, err := project.Detect(workDir)
	if err != nil {
		return err
	}

	raw, found, err := raw.LoadFromDir(workDir, proj.RootDir)
	if err != nil {
		return err
	}
	if !found {
		return errors.Errorf("No Oyafile found in %v", workDir)
	}

	if err := raw.AddImport(alias, uri); err != nil {
		return err
	}

	return nil
}
