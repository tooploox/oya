package project

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/bilus/oya/pkg/changeset"
	"github.com/bilus/oya/pkg/oyafile"
)

// Changeset returns the list of Oyafiles for the changed directories, based on the
// Changeset: directives.
//
// The algorithm:
//
// - Add default Changeset task to root Oyafile if there isn’t one
// - For each project Oyafile:
//   - If Changeset task defined:
//     - Contribute to changeset
// - Deduplicate changeset
// - Exclude changes outside work directory
func (p Project) Changeset(workDir string) ([]*oyafile.Oyafile, error) {
	oyafiles, err := p.Oyafiles()
	if err != nil {
		return nil, err
	}
	if len(oyafiles) == 0 {
		return nil, ErrNoOyafiles{Path: workDir}
	}

	rootOyafile := oyafiles[0]

	if rel, err := filepath.Rel(rootOyafile.Dir, p.RootDir); err != nil || rel != "." {
		panic("Internal error: expecting root Oyafile to be the first item in array returned from project.Project#Oyafiles method")
	}

	_, ok := rootOyafile.Tasks.LookupTask("Changeset")
	if !ok {
		rootOyafile.Tasks.AddTask("Changeset", defaultRootChangesetTask(oyafiles))
	}

	// return changeset.Calculate(p.Root, oyafiles)
	changeset, err := changeset.Calculate(oyafiles)
	if err != nil {
		return nil, err
	}
	return restrictToDir(workDir, changeset)
}

func restrictToDir(dir string, changeset []*oyafile.Oyafile) ([]*oyafile.Oyafile, error) {
	restricted := make([]*oyafile.Oyafile, 0, len(changeset))
	for _, o := range changeset {
		if isInside(dir, o) {
			restricted = append(restricted, o)
		}
	}
	return restricted, nil
}

func isInside(dir string, o *oyafile.Oyafile) bool {
	d := filepath.Clean(dir)
	r, err := filepath.Rel(d, o.Dir)
	return err == nil && !strings.Contains(r, "..")
}

func defaultRootChangesetTask(oyafiles []*oyafile.Oyafile) oyafile.Task {
	return oyafile.BuiltinTask{
		OnExec: func(stdout, stderr io.Writer) error {
			for _, o := range oyafiles {
				relPath := filepath.Dir(o.RelPath())
				_, err := stdout.Write([]byte(fmt.Sprintf("+%v\n", relPath)))
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}
