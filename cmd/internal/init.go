package internal

import (
	"io"

	"github.com/bilus/oya/pkg/raw"
	"github.com/pkg/errors"
)

func Init(rootDir string, stdout, stderr io.Writer) error {
	err := raw.InitDir(rootDir)
	if err != nil {
		return errors.Wrapf(err, "Error while initializing %v", rootDir)
	}
	return nil
}
