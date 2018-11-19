package oyafile

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const DefaultName = "Oyafile"

// TODO: Cleanup, should probably be Project.List.
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
		// TODO: Remove once we start verifying that all imported plugins have Project:
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
		if oyafile.Project != "" && path != rootDir {
			// Exclude projects nested under the current project.
			return filepath.SkipDir
		}
		oyafiles = append(oyafiles, oyafile)
		return nil
	})
}
