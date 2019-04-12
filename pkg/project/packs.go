package project

import (
	"sort"

	"github.com/tooploox/oya/pkg/deptree"
	"github.com/tooploox/oya/pkg/oyafile"
	"github.com/tooploox/oya/pkg/pack"
	"github.com/tooploox/oya/pkg/repo"
	"github.com/tooploox/oya/pkg/types"
)

func (p *Project) Require(pack pack.Pack) error {
	raw, found, err := p.rawOyafileIn(p.RootDir)
	if err != nil {
		return err
	}
	if !found {
		return ErrNoOyafile{Path: p.RootDir}
	}

	err = raw.AddRequire(pack)
	if err != nil {
		return err
	}
	p.invalidateOyafileCache(p.RootDir)

	p.dependencies = nil // Force reload.
	return nil
}

func (p *Project) Install(pack pack.Pack) error {
	return pack.Install(p.installDir)
}

func (p *Project) IsInstalled(pack pack.Pack) (bool, error) {
	return pack.IsInstalled(p.installDir)
}

// InstallPacks installs packs used by the project.
// It works in two steps:
// 1. It goes through all Import: directives and updates the Require: section with missing packs in their latest versions.
// 2. It installs all packs that haven't been installed.
func (p *Project) InstallPacks() error {
	err := p.updateDependencies()
	if err != nil {
		return err
	}

	deps, err := p.Deps()
	if err != nil {
		return err
	}
	return deps.ForEach(
		func(pack pack.Pack) error {
			_, ok := pack.ReplacementPath()
			if ok {
				return nil
			}
			installed, err := p.IsInstalled(pack)
			if err != nil {
				return err
			}
			if installed {
				return nil
			}
			err = p.Install(pack)
			return err
		},
	)
}

func (p *Project) wrapInstallErr(err error) error {
	return ErrInstallingPacks{
		Cause:          err,
		ProjectRootDir: p.RootDir,
	}
}

func (p *Project) FindRequiredPack(importPath types.ImportPath) (pack.Pack, bool, error) {
	deps, err := p.Deps()
	if err != nil {
		return pack.Pack{}, false, err
	}
	return deps.Find(importPath)
}

func (p *Project) Deps() (Deps, error) {
	if p.dependencies != nil {
		return p.dependencies, nil
	}

	o, found, err := p.oyafileIn(p.RootDir)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrNoOyafile{Path: p.RootDir}
	}

	installDirs := []string{
		p.installDir,
	}
	requires, err := resolvePackReferences(o.Requires)
	if err != nil {
		return nil, err
	}
	ldr, err := deptree.New(p.RootDir, installDirs, requires)
	if err != nil {
		return nil, err
	}
	err = ldr.Explode()

	if err != nil {
		return nil, err
	}
	p.dependencies = ldr
	return ldr, nil
}

func (p *Project) updateDependencies() error {
	files, err := p.List(p.RootDir)
	if err != nil {
		return err
	}

	deps, err := p.Deps()
	if err != nil {
		return err
	}

	// Collect all import paths from all Oyafiles. Make them unique.
	// Also, sort them to make writing reliable tests easier.
	importPaths := uniqueSortedImportPaths(files)
	for _, importPath := range importPaths {
		_, found, err := deps.Find(importPath)
		if err != nil {
			return err
		}
		if found {
			continue
		}

		l, err := repo.Open(importPath)
		if err != nil {
			// Import paths can also be relative to the root directory.
			// BUG(bilus): I don't particularly like it how tihs logic is split. Plus we may be masking some other errors this way
			if _, ok := err.(repo.ErrNotGithub); ok {
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
	return nil
}

func uniqueSortedImportPaths(oyafiles []*oyafile.Oyafile) []types.ImportPath {
	importPathSet := make(map[types.ImportPath]struct{})
	importPaths := make([]types.ImportPath, 0)
	for _, o := range oyafiles {
		for _, importPath := range o.Imports {
			if _, exists := importPathSet[importPath]; !exists {
				importPaths = append(importPaths, importPath)
			}
			importPathSet[importPath] = struct{}{}
		}
	}

	sort.Slice(importPaths, func(i, j int) bool {
		return importPaths[i] < importPaths[j]
	})

	return importPaths
}

func resolvePackReferences(references []oyafile.PackReference) ([]pack.Pack, error) {
	packs := make([]pack.Pack, len(references))
	for i, reference := range references {
		l, err := repo.Open(reference.ImportPath)
		if err != nil {
			return nil, err
		}
		pack, err := l.Version(reference.Version)
		if err != nil {
			return nil, err
		}
		if len(reference.ReplacementPath) > 0 {
			pack = pack.LocalReplacement(reference.ReplacementPath)
		}
		packs[i] = pack
	}
	return packs, nil
}
