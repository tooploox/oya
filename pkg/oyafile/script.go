package oyafile

import (
	"io"
	"io/ioutil"

	"github.com/magefile/mage/sh"
)

type Script string

func (s Script) Exec(env map[string]string, stdout, stderr io.Writer, shell string) error {
	scriptFile, err := ioutil.TempFile("", "oya-script-")
	if err != nil {
		return err
	}
	_, err = scriptFile.WriteString(string(s))
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
