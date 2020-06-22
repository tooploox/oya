package internal

import "github.com/tooploox/oya/pkg/project"

func prepareProject(workDir string) (*project.Project, error) {
	installDir, err := installDir()
	if err != nil {
		return nil, err
	}
	p, err := project.Detect(workDir, installDir)
	if err != nil {
		return nil, err
	}
	err = p.InstallPacks()
	if err != nil {
		return nil, err
	}
	return p, nil
}
