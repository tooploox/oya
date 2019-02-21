package deptree

import (
	"path/filepath"

	"github.com/bilus/oya/pkg/deptree/internal"
	"github.com/bilus/oya/pkg/mvs"
	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/raw"
	"github.com/bilus/oya/pkg/types"
	"github.com/pkg/errors"
)

// DependencyTree defines a project's dependencies, allowing for loading them.
type DependencyTree struct {
	rootDir      string
	installDirs  []string
	dependencies []pack.Pack
	reqs         *internal.Reqs
}

// New returns a new dependency tree.
// BUG(bilus): It's called a 'tree' but it currently does not take into account inter-pack
// dependencies. This will likely change and then the name will fit like a glove. ;)
func New(rootDir string, installDirs []string, dependencies []pack.Pack) (*DependencyTree, error) {
	return &DependencyTree{
		rootDir:      rootDir,
		installDirs:  installDirs,
		dependencies: dependencies,
		reqs:         internal.NewReqs(rootDir, installDirs),
	}, nil
}

// Explode takes the initial list of dependencies and builds the full list,
// taking into account packs' dependencies and using Minimal Version Selection.
func (dt *DependencyTree) Explode() error {
	list, err := mvs.List(dt.dependencies, dt.reqs)
	if err != nil {
		return err
	}
	dt.dependencies = list
	return nil
}

// Load loads an pack's Oyafile based on its import path.
// It supports two types of import paths:
// - referring to the project's Require: section (e.g. github.com/tooploox/oya-packs/docker), in this case it will load, the required version;
// - path relative to the project's root (e.g. /) -- does not support versioning, loads Oyafile directly from the path (<root dir>/<import path>).
func (dt *DependencyTree) Load(importPath types.ImportPath) (*oyafile.Oyafile, bool, error) {
	pack, found, err := dt.findRequiredPack(importPath)
	if err != nil {
		return nil, false, err
	}
	if found {
		return dt.loadOyafile(pack)
	}
	return nil, false, nil
}

// Find lookups pack by its import path.
func (dt *DependencyTree) Find(importPath types.ImportPath) (pack.Pack, bool, error) {
	for _, pack := range dt.dependencies {
		if pack.ImportPath() == importPath {
			return pack, true, nil
		}
	}
	return pack.Pack{}, false, nil
}

// ForEach iterates through the packs.
func (dt *DependencyTree) ForEach(f func(pack.Pack) error) error {
	for _, pack := range dt.dependencies {
		if err := f(pack); err != nil {
			return err
		}
	}
	return nil
}

func (dt *DependencyTree) loadOyafile(pack pack.Pack) (*oyafile.Oyafile, bool, error) {
	if path, ok := pack.ReplacementPath(); ok {
		fullPath := filepath.Join(dt.rootDir, path)
		o, found, err := oyafile.LoadFromDir(fullPath, dt.rootDir)
		if !found {
			return nil, false, errors.Errorf("no %v found at the replacement path %v for %q", raw.DefaultName, fullPath, pack.ImportPath())
		}
		if err != nil {
			return nil, false, errors.Wrapf(err, "error resolving replacement path %v for %q", fullPath, pack.ImportPath())

		}
		return o, true, nil

	}
	for _, installDir := range dt.installDirs {
		o, found, err := oyafile.LoadFromDir(pack.InstallPath(installDir), dt.rootDir)
		if err != nil {
			continue
		}
		if !found {
			continue
		}
		return o, true, nil
	}
	return nil, false, nil
}

func (dt *DependencyTree) findRequiredPack(importPath types.ImportPath) (pack.Pack, bool, error) {
	for _, pack := range dt.dependencies {
		if pack.ImportPath() == importPath {
			return pack, true, nil
		}
	}
	return pack.Pack{}, false, nil
}
