package project

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/pkg/errors"
	"k8s.io/helm/pkg/ignore"
)

func (p Project) Oyafiles() ([]*oyafile.Oyafile, error) {
	return listOyafiles(p.RootDir, p.RootDir)
}

func (p Project) List(startDir string) ([]*oyafile.Oyafile, error) {
	return listOyafiles(startDir, p.RootDir)
}

// TODO: Cleanup, should probably be Project.List.
func listOyafiles(startDir, rootDir string) ([]*oyafile.Oyafile, error) {
	ignore, err := makeIgnoreFunc(rootDir)
	if err != nil {
		return nil, errors.Wrapf(err, "error setting up ignores in %v", startDir)
	}

	vendorDir := filepath.Join(startDir, VendorDir)
	var oyafiles []*oyafile.Oyafile
	return oyafiles, filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
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
		oyafile, ok, err := oyafile.LoadFromDir(path, rootDir)
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

// makeIgnoreFunc returns a function that given an oyafile returns true if its containing directory tree should be recursively ignored.
// It uses an array of relative paths under "Ignore:" key in the project's root Oyafile.
// BUG(bilus): We should probably make it more intuitive by supporting Ignore: directives in nested dirs as well as the root dir.
func makeIgnoreFunc(rootDir string) (func(*oyafile.Oyafile) (bool, error), error) {
	o, ok, err := oyafile.LoadFromDir(rootDir, rootDir)
	if err != nil {
		return nil, errors.Wrapf(err, "error looking for Ignore: directive")
	}
	if !ok {
		return nil, errors.Errorf("No oyafile found at %v", rootDir)
	}
	ignore, err := ignore.Parse(strings.NewReader(o.Ignores()))
	if err != nil {
		return nil, errors.Wrapf(err, "Ignore: in %v contains invalid entries", o.Path)
	}
	return func(o *oyafile.Oyafile) (bool, error) {
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
