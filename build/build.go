package build

import (
	"io"

	"github.com/bilus/oya/pkg/changeset"
	"github.com/bilus/oya/pkg/oyafile"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func Build(rootDir, hookName string, stdout, stderr io.Writer) error {
	log.Debugf("Hook %q at %v", hookName, rootDir)
	oyafile, ok, err := oyafile.LoadFromDir(rootDir)
	if err != nil {
		return err
	}
	if !ok {
		// TODO: Need warn.
		return nil
	}

	changes, err := changeset.Calculate(oyafile)
	if err != nil {
		return err
	}
	for _, o := range changes {
		_, err = o.ExecHook(hookName, nil, stdout, stderr)
		if err != nil {
			return errors.Wrapf(err, "error in %v", o.Path)
		}
	}
	return nil
}
