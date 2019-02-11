package deptree

import (
	"os"
	"path/filepath"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/types"
)

// DependencyTree defines a project's dependencies, allowing for loading them.
type DependencyTree struct {
	rootDir      string
	installDirs  []string
	dependencies []pack.Pack
}

// New returns a new dependency tree.
// BUG(bilus): It's called a 'tree' but it currently does not take into account inter-pack
// dependencies. This will likely change and then the name will fit like a glove. ;)
func New(rootDir string, installDirs []string, dependencies []pack.Pack) (DependencyTree, error) {
	return DependencyTree{
		rootDir:      rootDir,
		installDirs:  installDirs,
		dependencies: dependencies,
	}, nil
}

func (dt DependencyTree) WithDependencies(dependencies []pack.Pack) (DependencyTree, error) {
	return DependencyTree{
		rootDir:      dt.rootDir,
		installDirs:  dt.installDirs,
		dependencies: dt.dependencies,
	}, nil
}

// Load loads an pack's Oyafile based on its import path.
// It supports two types of import paths:
// - referring to the project's Require: section (e.g. github.com/tooploox/oya-packs/docker), in this case it will load, the required version;
// - path relative to the project's root (e.g. /) -- does not support versioning, loads Oyafile directly from the path (<root dir>/<import path>).
func (dt DependencyTree) Load(importPath types.ImportPath) (*oyafile.Oyafile, bool, error) {
	pack, found, err := dt.findRequiredPack(importPath)
	if err != nil {
		return nil, false, err
	}
	if found {
		return dt.loadOyafile(pack)
	}

	// Attempt to load the Oyafile using the local path.
	fullImportPath := filepath.Join(dt.rootDir, string(importPath))
	if isValidImportPath(fullImportPath) {
		o, found, err := oyafile.LoadFromDir(fullImportPath, dt.rootDir)
		if err != nil {
			return nil, false, err
		}
		if found {
			return o, true, nil
		}
	}
	return nil, false, nil
}

// Find lookups pack by its import path.
func (dt DependencyTree) Find(importPath types.ImportPath) (pack.Pack, bool, error) {
	for _, pack := range dt.dependencies {
		if pack.ImportPath() == importPath {
			return pack, true, nil
		}
	}
	return nil, false, nil
}

// ForEach iterates through the packs.
func (dt DependencyTree) ForEach(f func(pack.Pack) error) error {
	for _, pack := range dt.dependencies {
		if err := f(pack); err != nil {
			return err
		}
	}
	return nil
}

func (dt DependencyTree) loadOyafile(pack pack.Pack) (*oyafile.Oyafile, bool, error) {
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

func (dt DependencyTree) findRequiredPack(importPath types.ImportPath) (pack.Pack, bool, error) {
	for _, pack := range dt.dependencies {
		if pack.ImportPath() == importPath {
			return pack, true, nil
		}
	}
	return nil, false, nil
}

func isValidImportPath(fullImportPath string) bool {
	f, err := os.Stat(fullImportPath)
	return err == nil && f.IsDir()
}
