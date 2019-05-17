package raw

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func InitDir(dirPath, projectName string) error {
	// BUG(bilus): Use raw access.
	_, found, err := LoadFromDir(dirPath, dirPath)
	if err == nil && found {
		return errors.Errorf("already an Oya project")
	}
	f, err := os.Create(fullPath(dirPath, ""))
	if err != nil {
		return err
	}
	_, err = f.WriteString(fmt.Sprintf("Project: %s\n", projectName))
	if err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
