package project

import (
	"io"
	"path/filepath"

	"github.com/bilus/oya/pkg/changeset"
	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/pack"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// TODO: Duplicated in oyafile module.
const VendorDir = "oya/vendor"

type Project struct {
	Root *oyafile.Oyafile
}

func Load(rootDir string) (Project, error) {
	root, ok, err := oyafile.LoadFromDir(rootDir, rootDir)
	if !ok {
		err = errors.Errorf("Missing Oyafile at %v", rootDir)
	}
	return Project{
		Root: root,
	}, err
}

func Detect(workDir string) (Project, error) {
	o, found, err := detectRoot(workDir)
	if err != nil {
		return Project{}, err
	}
	if !found {
		return Project{}, ErrNoProject{Path: workDir}
	}
	return Project{
		Root: o,
	}, nil
}

func (p Project) Run(workDir, taskName string, stdout, stderr io.Writer) error {
	log.Debugf("Task %q at %v", taskName, workDir)

	oyafiles, err := listOyafiles(workDir)
	if err != nil {
		return err
	}
	for _, o := range oyafiles {
		log.Println(o.Path)
	}
	if len(oyafiles) == 0 {
		return ErrNoOyafiles{Path: workDir}
	}

	changes, err := changeset.Calculate(oyafiles)
	if err != nil {
		return err
	}

	if len(changes) == 0 {
		return nil
	}

	foundAtLeastOneTask := false
	for _, o := range changes {
		found, err := o.RunTask(taskName, stdout, stderr)
		if err != nil {
			return errors.Wrapf(err, "error in %v", o.Path)
		}
		if found {
			foundAtLeastOneTask = found
		}
	}

	if !foundAtLeastOneTask {
		return ErrNoTask{
			Task: taskName,
		}
	}
	return nil
}

func (p Project) Oyafile(oyafilePath string) (*oyafile.Oyafile, bool, error) {
	return oyafile.Load(oyafilePath, p.Root.RootDir)
}

func (p Project) Vendor(pack pack.Pack) error {
	return pack.Vendor(filepath.Join(p.Root.RootDir, VendorDir))
}

func isRoot(o *oyafile.Oyafile) bool {
	return len(o.Project) > 0
}

func detectRoot(startDir string) (*oyafile.Oyafile, bool, error) {
	path := startDir
	maxParts := 256
	for i := 0; i < maxParts; i++ {
		o, found, err := oyafile.LoadFromDir(path, path)
		if err != nil {
			return nil, false, err
		}
		if err == nil && found && isRoot(o) {
			return o, true, nil
		}
		if path == "/" {
			break
		}
		path = filepath.Dir(path)
	}

	return nil, false, nil
}
