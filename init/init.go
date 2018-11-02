package init

import (
	"io"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/pkg/errors"
)

func Init(rootDir string, stdout, stderr io.Writer) error {
	err := oyafile.InitDir(rootDir)
	if err != nil {
		return errors.Wrapf(err, "Error while initializing %v", rootDir)
	}
	return nil
}
