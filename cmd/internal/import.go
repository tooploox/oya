package internal

import (
	"io"
	"strings"

	"github.com/pkg/errors"
	"github.com/tooploox/oya/pkg/project"
	"github.com/tooploox/oya/pkg/raw"
)

func Import(workDir, uri string, stdout, stderr io.Writer) error {
	uriArr := strings.Split(uri, "/")
	alias := uriArr[len(uriArr)-1]

	installDir, err := installDir()
	if err != nil {
		return wrapErr(err, uri)
	}
	proj, err := project.Detect(workDir, installDir)
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

	return proj.InstallPacks()
}
