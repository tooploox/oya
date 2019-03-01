package task

import (
	"io"
	"io/ioutil"
	corelog "log"
	"os"
	"strings"

	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
	"github.com/tooploox/oya/pkg/template"
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
	corelog.SetOutput(ioutil.Discard) // BUG(bilus): Suppress logging from the library. This prevents using standard logger anywhere else.

	_, err = sh.Exec(env(), stdout, stderr, s.Shell, scriptFile.Name())
	return err
}

func env() map[string]string {
	env := make(map[string]string)
	for _, v := range os.Environ() {
		parts := strings.SplitN(v, "=", 2)
		switch len(parts) {
		case 1:
			env[parts[0]] = ""
		case 2:
			env[parts[0]] = parts[1]
		}
	}
	return env
}
