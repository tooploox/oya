package loader

import (
	"os"
	"path/filepath"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/types"
)

type Loader struct {
	rootDir      string
	installDirs  []string
	dependencies []pack.Pack
}

func New(rootDir string, installDirs []string, dependencies []pack.Pack) (Loader, error) {
	return Loader{
		rootDir:      rootDir,
		installDirs:  installDirs,
		dependencies: dependencies,
	}, nil
}

func (l Loader) Load(importPath types.ImportPath) (*oyafile.Oyafile, bool, error) {
	pack, found, err := l.findRequiredPack(importPath)
	if err != nil {
		return nil, false, err
	}
	if found {
		return l.loadOyafile(pack)
	}

	// Attempt to load the Oyafile using the local path.
	fullImportPath := filepath.Join(l.rootDir, string(importPath))
	if isValidImportPath(fullImportPath) {
		o, found, err := oyafile.LoadFromDir(fullImportPath, l.rootDir)
		if err != nil {
			return nil, false, err
		}
		if found {
			return o, true, nil
		}
	}
	return nil, false, nil
}

func (l Loader) loadOyafile(pack pack.Pack) (*oyafile.Oyafile, bool, error) {
	for _, installDir := range l.installDirs {
		o, found, err := oyafile.LoadFromDir(pack.InstallPath(installDir), l.rootDir)
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

func (l Loader) findRequiredPack(importPath types.ImportPath) (pack.Pack, bool, error) {
	for _, pack := range l.dependencies {
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
