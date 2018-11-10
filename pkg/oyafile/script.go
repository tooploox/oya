package oyafile

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
	kasia "github.com/ziutek/kasia.go"
)

type Script string

func (s Script) Exec(values map[string]interface{}, stdout, stderr io.Writer, shell string) error {
	scriptFile, err := ioutil.TempFile("", "oya-script-")
	if err != nil {
		return err
	}
	defer os.Remove(scriptFile.Name())
	scriptTpl, err := kasia.Parse(string(s))
	if err != nil {
		return errors.Wrapf(err, "error running script")
	}
	err = scriptTpl.Run(scriptFile, values)
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
