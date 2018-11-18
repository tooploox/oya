package project

import (
	"io"
	"path/filepath"

	"github.com/bilus/oya/pkg/changeset"
	"github.com/bilus/oya/pkg/oyafile"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Project struct {
	RootDir string
}

func Detect(workDir string) (Project, error) {
	rootDir, found, err := detectRootDir(workDir)
	if err != nil {
		return Project{}, err
	}
	if !found {
		return Project{}, ErrNoProject{Path: workDir}
	}
	return Project{
		RootDir: rootDir,
	}, nil
}

func (p Project) Run(workDir, hookName string, stdout, stderr io.Writer) error {
	_, found, _ := oyafile.LoadFromDir(workDir, p.RootDir)
	if !found {
		return ErrNoOyafile{Path: workDir}
	}

	log.Debugf("Hook %q at %v", hookName, workDir)

	oyafiles, err := oyafile.List(workDir)
	if err != nil {
		return err
	}
	if len(oyafiles) == 0 {
		return ErrNoOyafiles{Path: workDir}
	}

	if oyafiles[0].Dir != p.RootDir {
		panic("oyafile.List post-condition failed: expected first oyafile to be root Oyafile")
	}

	changes, err := changeset.Calculate(oyafiles)
	if err != nil {
		return err
	}

	if len(changes) == 0 {
		return nil
	}

	foundAtLeastOnHook := false
	for _, o := range changes {
		found, err := o.ExecHook(hookName, stdout, stderr)
		if err != nil {
			return errors.Wrapf(err, "error in %v", o.Path)
		}
		if found {
			foundAtLeastOnHook = found
		}
	}

	if !foundAtLeastOnHook {
		return ErrNoHook{
			Hook: hookName,
		}
	}
	return nil
}

func (p Project) LoadOyafile(oyafilePath string) (*oyafile.Oyafile, bool, error) {
	return oyafile.Load(oyafilePath, p.RootDir)
}

func isRoot(o *oyafile.Oyafile) bool {
	return len(o.Module) > 0
}

func detectRootDir(startDir string) (string, bool, error) {
	path := startDir
	maxParts := 256
	for i := 0; i < maxParts; i++ {
		o, found, err := oyafile.LoadFromDir(path, path)
		if err != nil {
			return "", false, err
		}
		if err == nil && found && isRoot(o) {
			return filepath.Clean(path), true, nil
		}
		if path == "/" {
			break
		}
		path = filepath.Dir(path)
	}

	return "", false, nil
}
