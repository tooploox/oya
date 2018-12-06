package oyafile

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/helm/pkg/ignore"
)

const DefaultName = "Oyafile"

// TODO: Cleanup, should probably be Project.List.
func List(rootDir string) ([]*Oyafile, error) {
	ignore, err := makeIgnoreFunc(rootDir)
	if err != nil {
		return nil, errors.Wrapf(err, "error setting up ignores in %v", rootDir)
	}

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
			return errors.Wrapf(err, "error loading Oyafile from %v", path)
		}
		if !ok {
			return nil
		}
		doIgnore, err := ignore(oyafile)
		if err != nil {
			return errors.Wrapf(err, "error trying to determine if %v should be ignored", oyafile.Path)
		}
		if doIgnore {
			return filepath.SkipDir
		}
		oyafiles = append(oyafiles, oyafile)
		return nil
	})
}
func makeIgnoreFunc(rootDir string) (func(*Oyafile) (bool, error), error) {
	oyafile, ok, err := LoadFromDir(rootDir, rootDir)
	if !ok {
		return nil, errors.Errorf("%v not found at %v", DefaultName, rootDir)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "error loading Oyafile from %v", rootDir)
	}
	ignore, err := ignore.Parse(strings.NewReader(oyafile.Ignores()))
	if err != nil {
		return nil, errors.Wrapf(err, "Ignore: in %v contains invalid entries", oyafile.Path)
	}
	return func(o *Oyafile) (bool, error) {
		if o.Project != "" && !o.IsRoot() {
			// Exclude projects nested under the current project.
			return true, nil
		}
		fi, err := os.Stat(o.Path)
		if err != nil {
			return true, err
		}

		return ignore.Ignore(o.RelPath(), fi), nil
	}, nil
}
