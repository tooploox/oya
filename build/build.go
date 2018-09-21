package build

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type Script = string

type Oyafile struct {
	Jobs map[string]Script `yaml:"jobs"`
}

func Build(projectDir, job string) error {
	tempDir, err := ioutil.TempDir("", "oya")
	defer os.RemoveAll(tempDir)
	if err != nil {
		return err
	}
	oyafilePath := path.Join(projectDir, "Oyafile")
	file, err := os.Open(oyafilePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	decoder := yaml.NewDecoder(file)
	var oyafile Oyafile
	err = decoder.Decode(&oyafile)
	if err != nil {
		return err
	}
	script, ok := oyafile.Jobs[job]
	if !ok {
		return fmt.Errorf("no such job: %v", job)
	}

	scriptFile, err := ioutil.TempFile(tempDir, "")
	if err != nil {
		return err
	}
	defer scriptFile.Close()
	_, err = scriptFile.WriteString(string(script))
	if err != nil {
		return err
	}

	err = sh.RunV("sh", scriptFile.Name())
	if err != nil {
		return errors.Wrapf(err, "error in %v", oyafilePath)
	}
	return nil
}
