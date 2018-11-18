package oyafile

import (
	"io"
	"io/ioutil"

	"github.com/bilus/oya/pkg/template"
	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
)

type Script string

func (s Script) Exec(values template.Scope, stdout, stderr io.Writer, shell string) error {
	scriptFile, err := ioutil.TempFile("", "oya-script-")
	if err != nil {
		return err
	}
	// defer os.Remove(scriptFile.Name())
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
	_, err = sh.Exec(nil, stdout, stderr, "sh", scriptFile.Name())
	return err
}
