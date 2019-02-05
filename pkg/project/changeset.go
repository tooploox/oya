package project

import (
	"github.com/bilus/oya/pkg/changeset"
	"github.com/bilus/oya/pkg/oyafile"
	log "github.com/sirupsen/logrus"
)

func (p Project) Changeset(workDir string) ([]*oyafile.Oyafile, error) {
	oyafiles, err := listOyafiles(workDir, p.Root.RootDir)
	if err != nil {
		return nil, err
	}
	for _, o := range oyafiles {
		log.Println(o.Path)
	}
	if len(oyafiles) == 0 {
		return nil, ErrNoOyafiles{Path: workDir}
	}

	// _, ok := rootOyafile.Tasks.LookupTask("Changeset")
	// if !ok {
	// 	rootOyafile.Tasks.AddTask("Changeset", defaultRootChangesetTask(oyafiles))
	// }

	// return changeset.Calculate(p.Root, oyafiles)
	return changeset.Calculate(oyafiles)
}

// func defaultRootChangesetTask(candidates []*oyafile.Oyafile) oyafile.Task {
// 	return oyafile.BuiltinTask{
// 		Name: "Changeset",
// 		OnExec: func(stdout, stderr io.Writer) error {
// 			for _, o := range candidates {
// 				relPath := filepath.Dir(o.RelPath())
// 				_, err := stdout.Write([]byte(fmt.Sprintf("+%v\n", relPath)))
// 				if err != nil {
// 					return err
// 				}
// 			}

// 			return nil
// 		},
// 	}
// }
