package project

import (
	"io"
	"log"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/tooploox/oya/pkg/oyafile"
	"github.com/tooploox/oya/pkg/raw"
	"github.com/tooploox/oya/pkg/shell"
	"github.com/tooploox/oya/pkg/task"
	"github.com/tooploox/oya/pkg/template"
)

// Run runs a task within a project's context, the task looked up in an Oyafile
// in the current directory, unless recurse or useChangeset arguments are true.
// See LoadOyafiles for details regarding these arguments.
func (p *Project) Run(workDir string, taskName task.Name, recurse, useChangeset bool,
	args []string, scope template.Scope, stdout, stderr io.Writer) error {

	values, err := p.values()
	if err != nil {
		return err
	}
	scope = scope.Merge(values)

	targets, err := p.LoadOyafiles(workDir, recurse, useChangeset)
	if err != nil {
		return err
	}

	if len(targets) == 0 {
		return nil
	}

	foundAtLeastOneTask := false
	for _, o := range targets {
		found, err := o.RunTask(taskName, args, scope, stdout, stderr)
		if err != nil {
			return err
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

// StartREPL starts an interactivee shell. If terminal is available, it automatically upgrades,
// turning on support for auto-completion and command history.
func (p *Project) StartREPL(workDir string, stdin io.Reader, stdout, stderr io.Writer) error {
	builtins, err := p.values()
	if err != nil {
		return err
	}
	o, found, err := p.Oyafile(filepath.Join(workDir, raw.DefaultName))
	if err != nil {
		return err
	}

	var scope template.Scope
	if found {
		scope = builtins.Merge(o.Values)
	} else {
		log.Println("WARNING: No Oyafile in the current directory")

	}

	// TODO: Pass oya-cmd-override.
	return shell.StartREPL(scope, workDir, stdin, stdout, stderr, nil)

}

func (p *Project) LoadOyafiles(workDir string, recurse, useChangeset bool) ([]*oyafile.Oyafile, error) {
	var oyafiles []*oyafile.Oyafile
	var err error
	if useChangeset {
		changes, err := p.Changeset(workDir)
		if err != nil {
			return nil, err
		}

		if len(changes) == 0 {
			return nil, nil
		}

		if !recurse {
			oyafiles, err = p.oneTargetIn(workDir)
			if err != nil {
				return nil, err
			}
		} else {
			oyafiles = changes
		}
	} else {
		if !recurse {
			oyafiles, err = p.oneTargetIn(workDir)
			if err != nil {
				return nil, err
			}
		} else {
			oyafiles, err = p.List(workDir)
			if err != nil {
				return nil, err
			}
		}
	}

	dependencies, err := p.Deps()
	if err != nil {
		return nil, err
	}

	for _, o := range oyafiles {
		err = o.Build(dependencies)
		if err != nil {
			return nil, errors.Wrapf(err, "error in %v", o.Path)
		}
	}

	return oyafiles, nil
}

func (p *Project) values() (template.Scope, error) {
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
