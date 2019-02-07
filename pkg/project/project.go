package project

import (
	"io"
	"path/filepath"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/raw"
	"github.com/bilus/oya/pkg/task"
	"github.com/bilus/oya/pkg/template"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// TODO: Duplicated in oyafile module.
const VendorDir = ".oya/vendor"

type Project struct {
	RootDir string
}

func Load(rootDir string) (Project, error) {
	prj, err := Detect(rootDir)
	if err != nil {
		return prj, err
	}

	rel, err := filepath.Rel(rootDir, prj.RootDir)
	if err != nil {
		return prj, errors.Wrapf(err, "%v is not the Oya project root directory (it's %v)", rootDir, prj.RootDir)
	}
	if rel != "." {
		return prj, errors.Errorf("%v is not an Oya project root directory", rootDir)
	}

	return prj, nil
}

func Detect(workDir string) (Project, error) {
	detectedRootDir, found, err := detectRoot(workDir)
	if err != nil {
		return Project{}, err
	}
	if !found {
		return Project{}, ErrNoProject{Path: workDir}
	}
	return Project{
		RootDir: detectedRootDir,
	}, nil
}

func (p Project) Run(workDir string, taskName task.Name, scope template.Scope, stdout, stderr io.Writer) error {
	log.Debugf("Task %q at %v", taskName, workDir)

	changes, err := p.Changeset(workDir)
	if err != nil {
		return err
	}

	if len(changes) == 0 {
		return nil
	}

	foundAtLeastOneTask := false
	for _, o := range changes {
		found, err := o.RunTask(taskName, scope, stdout, stderr)
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

func (p Project) rootOyafile() (*oyafile.Oyafile, error) {
	o, found, err := oyafile.LoadFromDir(p.RootDir, p.RootDir)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrNoOyafile{Path: p.RootDir}
	}

	return o, nil
}

func (p Project) rootRawOyafile() (*raw.Oyafile, error) {
	o, found, err := raw.LoadFromDir(p.RootDir, p.RootDir)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrNoOyafile{Path: p.RootDir}
	}

	return o, nil
}

func (p Project) Oyafile(oyafilePath string) (*oyafile.Oyafile, bool, error) {
	return oyafile.Load(oyafilePath, p.RootDir)
}

func (p Project) Values() (template.Scope, error) {
	oyafilePath := filepath.Join(p.RootDir, "Oyafile")
	o, found, err := p.Oyafile(oyafilePath)
	if err != nil {
		return template.Scope{}, err
	}
	if !found {
		return template.Scope{}, ErrNoOyafile{Path: oyafilePath}
	}
	return template.Scope{
		"Project": o.Project,
	}, nil
}

// detectRoot attempts to detect the root project directory marked by
// root Oyafile, i.e. one containing Project: directive.
// It walks the directory tree, starting from startDir, going upwards,
// looking for root.
func detectRoot(startDir string) (string, bool, error) {
	path := startDir
	maxParts := 256
	for i := 0; i < maxParts; i++ {
		raw, found, err := raw.LoadFromDir(path, path) // "Guess" path is the root dir.
		if err == nil && found {
			isRoot, err := raw.IsRoot()
			if err != nil {
				return "", false, err
			}
			if isRoot {
				return path, true, nil
			}
		}

		if path == "/" {
			break
		}
		path = filepath.Dir(path)
	}

	return "", false, nil
}
