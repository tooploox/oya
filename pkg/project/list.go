package project

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/raw"
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
	skip := makeSkipFunc(startDir, rootDir)
	ignore, err := makeIgnoreFunc(rootDir)
	if err != nil {
		return nil, errors.Wrapf(err, "error setting up ignores in %v", startDir)
	}
	var oyafiles []*oyafile.Oyafile
	return oyafiles, filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		doSkip, err := skip(path)
		if err != nil {
			return errors.Wrapf(err, "error trying to determine if %v should be skipped", path)
		}
		if doSkip {
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

// makeSkipFunc returns a function that given a path, returns
// true if the entire subdirectory should be ignored.
// Similar to makeIgnoreFunc but does not parse Oyafile, thus allowing
// for broken Oyafile projects nested under the current project.
func makeSkipFunc(startDir, rootDir string) func(path string) (bool, error) {
	vendorDir := filepath.Join(startDir, VendorDir)
	return func(path string) (bool, error) {
		// Exclude anything under .oya/vendor

		if path == vendorDir {
			return true, nil
		}

		// Exclude projects nested under the current project.

		raw, ok, err := raw.LoadFromDir(path, rootDir)
		if !ok {
			return false, nil
		}
		if err != nil {
			return false, err
		}

		isRoot, err := raw.IsRoot()
		if err != nil {
			return false, err
		}

		// BUG(bilus): Clean up this magic string & logic duplication everywhere.
		_, isProject, err := raw.LookupKey("Project")
		if err != nil {
			return false, err
		}

		return isProject && !isRoot, nil
	}
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
		fi, err := os.Stat(o.Path)
		if err != nil {
			return true, err
		}

		return ignore.Ignore(o.RelPath(), fi), nil
	}, nil
}
