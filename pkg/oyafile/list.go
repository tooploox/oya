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
		oyafile, ok, err := LoadFromDir(path, rootDir)
		if err != nil {
			return errors.Wrapf(err, "Error loading Oyafile from %v", path)
		}
		if !ok {
			return nil
		}
		oyafiles = append(oyafiles, oyafile)
		return nil
	})
}
