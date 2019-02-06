package oyafile

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/bilus/oya/pkg/template"
	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
)

type Script string

func (s Script) Exec(workDir string, values template.Scope, stdout, stderr io.Writer, shell string) error {
	scriptFile, err := ioutil.TempFile("", "oya-script-")
	if err != nil {
		return err
	}
	defer os.Remove(scriptFile.Name())
	scriptTpl, err := template.Parse(string(s))
	if err != nil {
		return errors.Wrapf(err, "error running script")
	}
	err = scriptTpl.Render(scriptFile, values)
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
	_, err = sh.Exec(nil, stdout, stderr, shell, scriptFile.Name())
	return err
}
