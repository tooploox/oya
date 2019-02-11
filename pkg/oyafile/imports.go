package oyafile

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/bilus/oya/pkg/template"
	"github.com/bilus/oya/pkg/types"
	"github.com/pkg/errors"
)

type PackLoader interface {
	Load(importPath types.ImportPath) (*Oyafile, bool, error)
}

func (o *Oyafile) Build(loader PackLoader) error {
	// Do not resolve imports when loading Oyafile. Sometimes, we need to load Oyafile before packs are ready to be imported.
	if !o.IsBuilt {
		err := o.resolveImports(loader)
		if err != nil {
			return err
		}
		o.IsBuilt = true
	}
	return nil
}

func (oyafile *Oyafile) resolveImports(loader PackLoader) error {
	for alias, importPath := range oyafile.Imports {
		o, err := oyafile.loadPackOyafile(loader, importPath)
		if err != nil {
			return err
		}

		oyafile.Values[string(alias)] = o.Values
		for key, val := range collectPackValueOverrides(alias, oyafile.Values) {
			o.Values[key] = val
		}
		oyafile.Tasks.ImportTasks(alias, o.Tasks)
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

// collectPackValueOverrides collects all <alias>.xxx values, overriding values
// in the pack imported under the alias. Example: docker.image.
func collectPackValueOverrides(alias types.Alias, values template.Scope) template.Scope {
	// BUG(bilus): Extract aliased key syntax (dot-separation) from here and other places.
	packValues := template.Scope{}
	find := regexp.MustCompile("^" + string(alias) + "\\.(.*)$")
	for key, val := range values {
		if match := find.FindStringSubmatch(key); len(match) == 2 {
			packValues[match[1]] = val
		}
	}
	return packValues
}
