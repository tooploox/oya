package oyafile

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/tooploox/oya/pkg/template"
	"github.com/tooploox/oya/pkg/types"
)

type PackLoader interface {
	Load(importPath types.ImportPath) (*Oyafile, bool, error)
}

// Build resolves Oyafile imports.
func (oyafile *Oyafile) Build(loader PackLoader) error {
	// Do not resolve imports when loading Oyafile. Sometimes, we need to load Oyafile before packs are ready to be imported.
	if !oyafile.IsBuilt {
		err := oyafile.resolveImports(loader)
		if err != nil {
			return err
		}
		oyafile.IsBuilt = true
	}
	return nil
}

func (oyafile *Oyafile) resolveImports(loader PackLoader) error {
	for alias, importPath := range oyafile.Imports {
		packOyafile, err := oyafile.loadPackOyafile(loader, importPath)
		if err != nil {
			return err
		}
		err = packOyafile.Build(loader)
		if err != nil {
			return err
		}

		// TODO(bilus): Extract function.
		err = oyafile.Values.UpdateScopeAt(string(alias),
			func(scope template.Scope) template.Scope {
				// Values in the main Oyafile overwrite values in the pack Oyafile.
				merged := packOyafile.Values.Merge(scope)
				// Task is keeping a pointer to the scope.
				packOyafile.Values.Replace(merged)
				return merged
			}, false)
		if err != nil {
			return errors.Wrapf(err, "error merging values for imported pack %v", alias)
		}

		oyafile.Tasks.ImportTasks(alias, packOyafile.Tasks)
	}

	return oyafile.expose()
}

func (oyafile *Oyafile) expose() error {
	for _, alias := range oyafile.ExposedAliases {
		oyafile.Tasks.Expose(alias)
	}
	return nil
}

func (oyafile *Oyafile) loadPackOyafile(loader PackLoader, importPath types.ImportPath) (*Oyafile, error) {
	o, found, err := loader.Load(importPath)
	if err != nil {
		return nil, err
	}
	if found {
		return o, nil
	}

	// Attempt to load the Oyafile using the local path.
	fullImportPath := filepath.Join(oyafile.RootDir, string(importPath))
	if isValidImportPath(fullImportPath) {
		o, found, err := LoadFromDir(fullImportPath, oyafile.RootDir)
		if err != nil {
			return nil, err
		}
		if found {
			return o, nil
		}
	}
	return nil, errors.Errorf("missing pack %v", importPath)
}

func isValidImportPath(fullImportPath string) bool {
	f, err := os.Stat(fullImportPath)
	return err == nil && f.IsDir()
}
