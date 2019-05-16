package internal

import (
	"io"

	"github.com/pkg/errors"
	"github.com/tooploox/oya/pkg/raw"
)

func Init(rootDir, projectName string, stdout, stderr io.Writer) error {
	err := raw.InitDir(rootDir, projectName)
	if err != nil {
		return errors.Wrapf(err, "Error while initializing %v", rootDir)
	}
	return nil
}
