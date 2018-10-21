package oyafile

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const DefaultName = "Oyafile"

func List(rootDir string) ([]*Oyafile, error) {
	var oyafiles []*Oyafile
	return oyafiles, filepath.Walk(rootDir, func(path string, info os.FileInfo, _ error) error {
		if !info.IsDir() {
			return nil
		}
		oyafile, ok, err := LoadFromDir(path)
		if err != nil {
			return errors.Wrapf(err, "error listing Oyafiles in %s", rootDir)
		}
		if !ok {
			return nil
		}
		oyafiles = append(oyafiles, oyafile)
		return nil
	})
}
