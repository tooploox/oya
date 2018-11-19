package oyafile

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const DefaultName = "Oyafile"

func List(rootDir string) ([]*Oyafile, error) {
	vendorDir := filepath.Join(rootDir, VendorDir)
	var oyafiles []*Oyafile
	return oyafiles, filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		if path == vendorDir {
			return filepath.SkipDir
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
