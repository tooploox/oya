package project

import (
	"io"

	"github.com/pkg/errors"
	"github.com/tooploox/oya/pkg/oyafile"
	"github.com/tooploox/oya/pkg/task"
	"github.com/tooploox/oya/pkg/template"
)

func (p *Project) Run(workDir string, taskName task.Name, recurse, useChangeset bool,
	args []string, scope template.Scope, stdout, stderr io.Writer) error {

	values, err := p.values()
	if err != nil {
		return err
	}
	scope = scope.Merge(values)

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
