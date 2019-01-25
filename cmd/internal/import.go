package internal

import (
	"io"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/pkg/errors"
)

func Import(rootDir, uri string, stdout, stderr io.Writer) error {
	err := oyafile.AddImport(rootDir, uri)
	if err != nil {
		return errors.Wrapf(err, "Error while importing pack %v", uri)
	}
	return nil
}
