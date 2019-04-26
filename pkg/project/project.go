package project

import (
	"path/filepath"

	"github.com/tooploox/oya/pkg/oyafile"
	"github.com/tooploox/oya/pkg/raw"
)

// TODO: Duplicated in oyafile module.
type Project struct {
	RootDir      string
	installDir   string
	dependencies Deps

	oyafileCache    map[string]*oyafile.Oyafile
	rawOyafileCache map[string]*raw.Oyafile
}

func Detect(workDir, installDir string) (*Project, error) {
	detectedRootDir, found, err := detectRoot(workDir)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrNoProject{Path: workDir}
	}
	return &Project{
		RootDir:         detectedRootDir,
		installDir:      installDir,
		dependencies:    nil, // lazily-loaded in Deps()
		oyafileCache:    make(map[string]*oyafile.Oyafile),
		rawOyafileCache: make(map[string]*raw.Oyafile),
	}, nil
}

// detectRoot attempts to detect the root project directory marked by
// root Oyafile, i.e. one containing Project: directive.
// It walks the directory tree, starting from startDir, going upwards,
// looking for root.
func detectRoot(startDir string) (string, bool, error) {
	path := startDir
	maxParts := 256
	for i := 0; i < maxParts; i++ {
		raw, found, err := raw.LoadFromDir(path, path) // "Guess" path is the root dir.
		if err == nil && found {
			isRoot, err := raw.IsRoot()
			if err != nil {
				return "", false, err
			}
			if isRoot {
				return path, true, nil
			}
		}

		if path == "/" {
			break
		}
		path = filepath.Dir(path)
	}

	return "", false, nil
}
