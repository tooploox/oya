package changeset

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bilus/oya/pkg/debug"
	"github.com/bilus/oya/pkg/oyafile"
	log "github.com/sirupsen/logrus"
)

func Calculate(oyafiles []*oyafile.Oyafile) ([]*oyafile.Oyafile, error) {
	var dirs []string
	for _, oyafile := range oyafiles {
		changes, err := calculateChangeset(oyafile)
		if err != nil {
			return nil, err
		}
		dirs = append(dirs, changes...)
	}
	included := unique(dirs)

	changeset := make([]*oyafile.Oyafile, 0, len(oyafiles))
	for _, oyafile := range oyafiles {
		cleanDir := filepath.Clean(oyafile.Dir)
		_, ok := included[cleanDir]
		if ok {
			log.Debug("Included: %v", cleanDir)
			changeset = append(changeset, oyafile)
		} else {
			log.Debug("Excluded: %v", cleanDir)
		}
	}

	debug.LogOyafiles("Resulting changeset", changeset)
	return changeset, nil
}

func calculateChangeset(oyafile *oyafile.Oyafile) ([]string, error) {
	log.Debugf("Calculating changeset at %v", oyafile.Dir)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	hook, err := oyafile.ExecHook("Changeset", nil, stdout, stderr)
	if !hook {
		return []string{oyafile.Dir}, nil
	}

	if err != nil {
		return nil, err
	}
	// TODO: Ignores stderr for the time being.
	changes := strings.Split(stdout.String(), "\n")
	dirs := make([]string, 0)
	for _, change := range changes {
		if len(change) == 0 {
			log.Debug("Rejected: <empty line>")
			continue
		}
		if change[0] != '+' || len(change) < 2 {
			log.Debug("Error: %v", change)
			return nil, fmt.Errorf("Unexpected changeset entry %q expected \"+path\"", change)
		}
		path := normalizePath(oyafile.Dir, change[1:])
		log.Debugf("Accepted: %v", path)
		// TODO: Check if path is valid.
		dirs = append(dirs, path)
	}

	return dirs, nil
}

func normalizePath(oyafileDir, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Clean(filepath.Join(oyafileDir, path))
}

func unique(strings []string) map[string]struct{} {
	unique := make(map[string]struct{})
	for _, s := range strings {
		unique[s] = struct{}{}
	}
	return unique
}
