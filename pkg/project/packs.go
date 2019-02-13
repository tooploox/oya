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
	err = raw.AddRequire(pack)
	if err != nil {
		return err
	}
	p.dependencies = nil // Force reload.
	return nil
}

func (p Project) Install(pack pack.Pack) error {
	return pack.Install(p.installDir)
}

func (p Project) IsInstalled(pack pack.Pack) (bool, error) {
	return pack.IsInstalled(p.installDir)
}

// InstallPacks installs packs used by the project.
// It works in two steps:
// 1. It goes through all Import: directives and updates the Require: section with missing packs in their latest versions.
// 2. It installs all packs that haven't been installed.
func (p Project) InstallPacks() error {
	err := p.updateDependencies()
	if err != nil {
		return err
	}

	deps, err := p.Dependencies()
	if err != nil {
		return err
	}
	return deps.ForEach(
		func(pack pack.Pack) error {
			installed, err := p.IsInstalled(pack)
			if err != nil {
				return err
			}
			if installed {
				return nil
			}
			return p.Install(pack)
		},
	)
}

func (p Project) FindRequiredPack(importPath types.ImportPath) (pack.Pack, bool, error) {
	deps, err := p.Dependencies()
	if err != nil {
		return nil, false, err
	}
	return deps.Find(importPath)
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

func (p Project) updateDependencies() error {
	files, err := p.List(p.RootDir)
	if err != nil {
		return err
	}

	importPaths := make(map[types.ImportPath]struct{})
	for _, o := range files {
		for _, importPath := range o.Imports {
			importPaths[importPath] = struct{}{}
		}
	}

	deps, err := p.Dependencies()
	if err != nil {
		return err
	}

	for importPath := range importPaths {
		_, found, err := deps.Find(importPath)
		if err != nil {
			return err
		}
		if found {
			continue
		}

		l, err := pack.OpenLibrary(importPath)
		if err != nil {
			// Import paths can also be relative to the root directory.
			// BUG(bilus): I don't particularly like it how tihs logic is split. Plus we may be masking some other errors this way
			if _, ok := err.(pack.ErrNotGithub); ok {
				continue
			}
			return err
		}

		pack, err := l.LatestVersion()
		if err != nil {
			return err
		}
		err = p.Require(pack)
		if err != nil {
			return err
		}
	}

	p.dependencies = nil // Force reload.
	return nil
}
