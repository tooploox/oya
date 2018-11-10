package changeset

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/bilus/oya/pkg/oyafile"
	log "github.com/sirupsen/logrus"
)

func Calculate(candidates []*oyafile.Oyafile) ([]*oyafile.Oyafile, error) {
	if len(candidates) == 0 {
		return candidates, nil
	}

	rootOyafile := candidates[0]
	// Set to default if not present.
	rootOyafile.Hooks["Changeset"] = rootChangesetHook(rootOyafile)

	var changeset []*oyafile.Oyafile
	for _, candidate := range candidates {
		oyafiles, err := calculateChangeset(candidate)
		if err != nil {
			return nil, err
		}
		changeset = append(changeset, oyafiles...)
	}

	return changeset, nil
}

func rootChangesetHook(rootOyafile *oyafile.Oyafile) oyafile.Hook {
	defaultHook := oyafile.BuiltinHook{
		Name: "Changeset",
		OnExec: func(values map[string]interface{}, stdout, stderr io.Writer) error {
			oyafiles, err := oyafile.List(rootOyafile.Dir)
			if err != nil {
				return err
			}
			for _, o := range oyafiles {
				relPath, err := filepath.Rel(rootOyafile.Dir, o.Dir)
				if err != nil {
					return err
				}
				_, err = stdout.Write([]byte(fmt.Sprintf("+%v\n", relPath)))
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	customHook, ok := rootOyafile.Hooks["Changeset"]
	if ok {
		return customHook
	}
	return defaultHook
}

func execChangesetHook(changesetHook oyafile.Hook) ([]string, error) {
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	err := changesetHook.Exec(nil, stdout, stderr)
	if err != nil {
		return nil, err
	}
	// TODO: We shouldn't be ignoring stderr.
	return parseChangeset(stdout.String())
}

func calculateChangeset(currOyafile *oyafile.Oyafile) ([]*oyafile.Oyafile, error) {
	changesetHook, ok := currOyafile.Hooks["Changeset"]
	if !ok {
		return nil, nil
	}

	dirs, err := execChangesetHook(changesetHook)
	if err != nil {
		return nil, err
	}

	oyafiles := make([]*oyafile.Oyafile, 0, len(dirs))
	for _, dir := range dirs {
		fullPath := filepath.Join(currOyafile.Dir, dir)
		o, exists, err := oyafile.LoadFromDir(fullPath, currOyafile.VendorDir)
		if !exists {
			// TODO: Warning that changeset contains paths without Oyafiles?
			continue
		}
		if err != nil {
			return nil, err
		}
		oyafiles = append(oyafiles, o)
	}
	return oyafiles, nil
}

func parseChangeset(changeset string) ([]string, error) {
	dirs := make([]string, 0)

	// TODO: Ignores stderr for the time being.
	changes := strings.Split(changeset, "\n")
	for _, change := range changes {
		if len(change) == 0 {
			continue
		}
		if change[0] != '+' || len(change) < 2 {
			log.Debugf("  Error: %v", change)
			return nil, fmt.Errorf("Unexpected changeset entry %q expected \"+path\"", change)
		}
		path := change[1:]
		log.Debugf("  Addition: %v", path)
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

func unique(oyafiles []*oyafile.Oyafile) []*oyafile.Oyafile {
	result := make([]*oyafile.Oyafile, 0, len(oyafiles))
	unique := make(map[string]struct{})
	for _, o := range oyafiles {
		_, ok := unique[o.Dir]
		if ok {
			continue
		}
		unique[o.Dir] = struct{}{}
		result = append(result, o)
	}

	return result
}
