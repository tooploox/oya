package build

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func Build(projectDir, job string) error {
	oyafilePath := path.Join(projectDir, "Oyafile")
	file, err := os.Open(oyafilePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	decoder := yaml.NewDecoder(file)
	var jobs map[string]string
	err = decoder.Decode(&jobs)
	if err != nil {
		return err
	}
	script, ok := jobs[job]
	if !ok {
		return fmt.Errorf("no such job: %v", job)
	}

	cmds := strings.Split(script, "\n")
	for line, cmd := range cmds {
		fmt.Printf("$ %v\n", cmd)
		err := sh.RunV("sh", "-c", cmd)
		if err != nil {
			return errors.Wrapf(err, "error in %v:%d", oyafilePath, line)
		}

	}
	return nil
}
