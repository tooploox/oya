package oyafile

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/bilus/oya/pkg/template"
	"github.com/bilus/oya/pkg/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const VendorDir = ".oya/vendor"

func (oyafile *Oyafile) resolveImports() error {
	for alias, path := range oyafile.Imports {
		log.Debugf("Importing pack %v as %v", path, alias)
		pack, err := oyafile.loadPack(path)
		if err != nil {
			return err
		}

		oyafile.Values[string(alias)] = pack.Values
		for key, val := range collectPackValueOverrides(alias, oyafile.Values) {
			pack.Values[key] = val
		}

		oyafile.Tasks.ImportTasks(alias, pack.Tasks)
	}
	return nil
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

func (oyafile *Oyafile) loadPack(path types.ImportPath) (*Oyafile, error) {
	for _, importDir := range oyafile.importDirs() {
		fullPath := filepath.Join(importDir, string(path))
		if !isValidImportPath(fullPath) {
			continue
		}
		pack, found, err := LoadFromDir(fullPath, oyafile.RootDir)
		if err != nil {
			continue
		}
		if !found {
			continue
		}
		return pack, nil
	}

	return nil, errors.Errorf("missing pack %v", path)
}

func (oyafile *Oyafile) importDirs() []string {
	return []string{
		oyafile.RootDir,
		filepath.Join(oyafile.RootDir, VendorDir),
	}
}

func isValidImportPath(fullImportPath string) bool {
	f, err := os.Stat(fullImportPath)
	return err == nil && f.IsDir()
}
