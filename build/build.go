package build

import (
	"fmt"
	"io"

	"github.com/bilus/oya/pkg/changeset"
	"github.com/bilus/oya/pkg/oyafile"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var ErrNoOyafiles = fmt.Errorf("No Oyafiles found")

func Build(rootDir, hookName string, stdout, stderr io.Writer) error {
	log.Debugf("Hook %q at %v", hookName, rootDir)

	oyafiles, err := oyafile.List(rootDir)
	if err != nil {
		return err
	}
	if len(oyafiles) == 0 {
		return ErrNoOyafiles
	}

	if oyafiles[0].Dir != rootDir {
		panic("oyafile.List post-condition failed: expected first oyafile to be root Oyafile")
	}

	changes, err := changeset.Calculate(oyafiles)
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
