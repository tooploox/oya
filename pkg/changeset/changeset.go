package changeset

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tooploox/oya/pkg/oyafile"
	"github.com/tooploox/oya/pkg/task"
	"github.com/tooploox/oya/pkg/template"
)

// Calculate gets changeset from each Oyafile by invoking their Changeset: tasks
// and parsing the output.
func Calculate(candidates []*oyafile.Oyafile) ([]*oyafile.Oyafile, error) {
	if len(candidates) == 0 {
		return candidates, nil
	}

	var changeset []*oyafile.Oyafile
	for _, candidate := range candidates {
		oyafiles, err := calculateChangeset(candidate)
		if err != nil {
			return nil, err
		}
		changeset = append(changeset, oyafiles...)
	}

	return unique(changeset), nil
}

func execChangesetTask(workDir string, changesetTask task.Task) ([]string, error) {
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	err := changesetTask.Exec(workDir, nil, template.Scope{}, stdout, stderr)
	if err != nil {
		return nil, err
	}
	// TODO: We shouldn't be ignoring stderr.
	return parseChangeset(stdout.String())
}

func calculateChangeset(currOyafile *oyafile.Oyafile) ([]*oyafile.Oyafile, error) {
	changesetTask, ok := currOyafile.Tasks.LookupTask("Changeset")
	if !ok {
		return nil, nil
	}

	dirs, err := execChangesetTask(currOyafile.Dir, changesetTask)
	if err != nil {
		return nil, err
	}

	oyafiles := make([]*oyafile.Oyafile, 0, len(dirs))
	for _, dir := range dirs {
		fullPath := filepath.Join(currOyafile.Dir, dir)
		o, exists, err := oyafile.LoadFromDir(fullPath, currOyafile.RootDir)
		if !exists {
			// TODO: Warning that changeset contains paths without Oyafiles?
			log.Printf("Path %v in changeset does not exist", fullPath)
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
