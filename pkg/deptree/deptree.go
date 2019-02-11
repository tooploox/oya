package deptree

import (
	"os"
	"path/filepath"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/types"
)

type DependencyTree struct {
	rootDir      string
	installDirs  []string
	dependencies []pack.Pack
}

func New(rootDir string, installDirs []string, dependencies []pack.Pack) (DependencyTree, error) {
	return DependencyTree{
		rootDir:      rootDir,
		installDirs:  installDirs,
		dependencies: dependencies,
	}, nil
}

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
