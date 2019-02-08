package task

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/bilus/oya/pkg/template"
	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
)

type Script struct {
	Script string
	Shell  string
	Scope  *template.Scope
}

func (s Script) Exec(workDir string, values template.Scope, stdout, stderr io.Writer) error {
	scope := values.Merge(*s.Scope)

	scriptFile, err := ioutil.TempFile("", "oya-script-")
	if err != nil {
		return err
	}
	defer os.Remove(scriptFile.Name())
	scriptTpl, err := template.Parse(s.Script)
	if err != nil {
		return errors.Wrapf(err, "error running script")
	}
	err = scriptTpl.Render(scriptFile, scope)
	if err != nil {
		_ = scriptFile.Close()
		return err
	}
	err = scriptFile.Close()
	if err != nil {
		return err
	}

	oldCwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(oldCwd)
	err = os.Chdir(workDir)
	if err != nil {
		return err
	}
	log.SetOutput(ioutil.Discard) // BUG(bilus): Suppress logging from the library. This prevents using standard logger anywhere else.
	_, err = sh.Exec(nil, stdout, stderr, s.Shell, scriptFile.Name())
	return err
}
