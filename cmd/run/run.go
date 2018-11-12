package run

import (
	"fmt"
	"io"
	"os"

	"github.com/bilus/oya/pkg/changeset"
	"github.com/bilus/oya/pkg/oyafile"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var ErrNoOyafiles = fmt.Errorf("missing Oyafile")

type ErrNoHook struct {
	Hook string
}

func (e ErrNoHook) Error() string {
	return fmt.Sprintf("missing hook %q", e.Hook)
}

func Run(rootDir, hookName string, stdout, stderr io.Writer) error {
	if len(rootDir) == 0 {
		var err error
		rootDir, err = detectRootDir()
		if err != nil {
			return errors.Wrapf(err, "error detecting Oya root directory")
		}
	}

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

	if len(changes) == 0 {
		return nil
	}

	foundAtLeastOnHook := false
	for _, o := range changes {
		found, err := o.ExecHook(hookName, stdout, stderr)
		if err != nil {
			return errors.Wrapf(err, "error in %v", o.Path)
		}
		if found {
			foundAtLeastOnHook = found
		}
	}

	if !foundAtLeastOnHook {
		return ErrNoHook{
			Hook: hookName,
		}
	}
	return nil
}

func detectRootDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return oyafile.DetectRootDir(cwd)
}
