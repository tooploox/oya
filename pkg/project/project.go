package project

import (
	"io"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tooploox/oya/pkg/oyafile"
	"github.com/tooploox/oya/pkg/raw"
	"github.com/tooploox/oya/pkg/task"
	"github.com/tooploox/oya/pkg/template"
)

// TODO: Duplicated in oyafile module.
type Project struct {
	RootDir      string
	installDir   string
	dependencies Deps

	oyafileCache    map[string]*oyafile.Oyafile
	rawOyafileCache map[string]*raw.Oyafile
}

func Detect(workDir, installDir string) (*Project, error) {
	detectedRootDir, found, err := detectRoot(workDir)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrNoProject{Path: workDir}
	}
	return &Project{
		RootDir:         detectedRootDir,
		installDir:      installDir,
		dependencies:    nil, // lazily-loaded in Deps()
		oyafileCache:    make(map[string]*oyafile.Oyafile),
		rawOyafileCache: make(map[string]*raw.Oyafile),
	}, nil
}

func (p *Project) Run(workDir string, taskName task.Name, recurse, useChangeset bool,
	args []string, scope template.Scope, stdout, stderr io.Writer) error {
	log.Debugf("Task %q at %v", taskName, workDir)

	targets, err := p.RunTargets(workDir, recurse, useChangeset)
	if err != nil {
		return err
	}

	if len(targets) == 0 {
		return nil
	}

	dependencies, err := p.Deps()
	if err != nil {
		return err
	}

	foundAtLeastOneTask := false
	for _, o := range targets {
		err = o.Build(dependencies)
		if err != nil {
			return errors.Wrapf(err, "error in %v", o.Path)
		}
		found, err := o.RunTask(taskName, args, scope, stdout, stderr)
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

func (p *Project) RunTargets(workDir string, recurse, useChangeset bool) ([]*oyafile.Oyafile, error) {
	if useChangeset {
		changes, err := p.Changeset(workDir)
		if err != nil {
			return nil, err
		}

		if len(changes) == 0 {
			return nil, nil
		}

		if !recurse {
			return p.oneTargetIn(workDir)
		} else {
			return changes, nil
		}
	} else {
		if !recurse {
			return p.oneTargetIn(workDir)
		} else {
			return p.List(workDir)
		}
	}
}

func (p *Project) oneTargetIn(dir string) ([]*oyafile.Oyafile, error) {
	o, found, err := p.oyafileIn(dir)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrNoOyafile{Path: dir}
	}
	return []*oyafile.Oyafile{o}, nil
}

func (p *Project) oyafileIn(dir string) (*oyafile.Oyafile, bool, error) {
	normalizedDir := filepath.Clean(dir)
	o, found := p.oyafileCache[normalizedDir]
	if found {
		return o, true, nil
	}
	o, found, err := oyafile.LoadFromDir(dir, p.RootDir)
	if err != nil {
		return nil, false, err
	}
	if found {
		p.oyafileCache[normalizedDir] = o
		return o, true, nil
	}

	return nil, false, nil
}

func (p *Project) rawOyafileIn(dir string) (*raw.Oyafile, bool, error) {
	// IMPORTANT: Call InvalidateOyafileCache after patching raw Oyafiles
	// obtained using this method!
	normalizedDir := filepath.Clean(dir)
	o, found := p.rawOyafileCache[normalizedDir]
	if found {
		return o, true, nil
	}
	o, found, err := raw.LoadFromDir(dir, p.RootDir)
	if err != nil {
		return nil, false, err
	}
	if found {
		p.rawOyafileCache[normalizedDir] = o
		return o, true, nil
	}
	return nil, false, nil
}

func (p *Project) InvalidateOyafileCache(dir string) {
	delete(p.oyafileCache, dir)
	delete(p.rawOyafileCache, dir)
}

func (p *Project) Oyafile(oyafilePath string) (*oyafile.Oyafile, bool, error) {
	// BUG(bilus): Uncached (but used only by Render).
	return oyafile.Load(oyafilePath, p.RootDir)
}

func (p *Project) Values() (template.Scope, error) {
	o, found, err := p.rawOyafileIn(p.RootDir)
	if err != nil {
		return template.Scope{}, err
	}
	if !found {
		return template.Scope{}, ErrNoOyafile{Path: p.RootDir}
	}
	project, _, err := o.Project()
	if err != nil {
		return template.Scope{}, err
	}
	return template.Scope{
		"Project": project,
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
