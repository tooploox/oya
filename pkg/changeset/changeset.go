package changeset

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bilus/oya/pkg/oyafile"
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
	included := uniqueDirs(dirs)

	changeset := make([]*oyafile.Oyafile, 0, len(oyafiles))
	for _, oyafile := range oyafiles {
		_, ok := included[filepath.Clean(oyafile.Dir)]
		if ok {
			changeset = append(changeset, oyafile)
		}
	}
	return changeset, nil
}

func calculateChangeset(oyafile *oyafile.Oyafile) ([]string, error) {
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
			continue
		}
		if change[0] != '+' || len(change) < 2 {
			return nil, fmt.Errorf("Unexpected changeset entry %q expected \"+path\"", change)
		}
		// TODO: Check if path is valid.
		dirs = append(dirs, change[1:])
	}

	return dirs, nil
}

func uniqueDirs(dirs []string) map[string]struct{} {
	unique := make(map[string]struct{})
	for _, dir := range dirs {
		cleanDir := filepath.Clean(dir)
		unique[cleanDir] = struct{}{}
	}
	return unique
}
