package build

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
)

func Build(projectDir, jobName string, stdout, stderr io.Writer) error {
	tempDir, err := ioutil.TempDir("", "oya")
	defer os.RemoveAll(tempDir)
	if err != nil {
		return err
	}
	return filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		oyafilePath := oyafile.FullPath(path, "")
		oyafile, buildable, err := oyafile.Load(oyafilePath)
		if err != nil {
			return err
		}
		if !buildable {
			return nil
		}

		job, ok := oyafile.Jobs[jobName]
		if !ok {
			return fmt.Errorf("no such job: %v", jobName)
		}

		script := job.Script

		scriptFile, err := ioutil.TempFile(tempDir, "")
		if err != nil {
			return err
		}
		_, err = scriptFile.WriteString(string(script))
		if err != nil {
			_ = scriptFile.Close()
			return err
		}
		err = scriptFile.Close()
		if err != nil {
			return err
		}
		_, err = sh.Exec(nil, stdout, stderr, "sh", scriptFile.Name())
		if err != nil {
			return errors.Wrapf(err, "error in %v", oyafilePath)
		}

		return nil
	})
}
