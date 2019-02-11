package project

import (
	"github.com/bilus/oya/pkg/deptree"
	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/types"
)

func (p Project) Require(pack pack.Pack) error {
	raw, err := p.rootRawOyafile()
	if err != nil {
		return err
	}

	return raw.AddRequire(pack)
}

func (p Project) Install(pack pack.Pack) error {
	return pack.Install(p.installDir)
}

func (p Project) IsInstalled(pack pack.Pack) (bool, error) {
	return pack.IsInstalled(p.installDir)
}

func (p Project) InstallPacks() error {
	o, err := p.rootOyafile()
	if err != nil {
		return err
	}
	for _, pack := range o.Require {
		installed, err := p.IsInstalled(pack)
		if err != nil {
			return err
		}
		if installed {

			continue
		}
		err = p.Install(pack)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p Project) FindRequiredPack(importPath types.ImportPath) (pack.Pack, bool, error) {
	o, err := p.rootOyafile()
	if err != nil {
		return nil, false, err
	}
	for _, pack := range o.Require {
		if pack.ImportPath() == importPath {
			return pack, true, nil
		}
	}
	return nil, false, nil
}

func (p Project) Dependencies() (deptree.DependencyTree, error) {
	if p.dependencies != nil {
		return *p.dependencies, nil
	}

	o, err := p.rootOyafile()
	if err != nil {
		return deptree.DependencyTree{}, err
	}
	installDirs := []string{
		p.installDir,
	}
	ldr, err := deptree.New(p.RootDir, installDirs, o.Require)
	if err != nil {
		return deptree.DependencyTree{}, err
	}
	p.dependencies = &ldr
	return ldr, nil
}
