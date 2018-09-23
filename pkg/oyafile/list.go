package oyafile

import (
	"os"
	"path"
	"path/filepath"
)

const DefaultName = "Oyafile"

func List(rootDir string) ([]*Oyafile, error) {
	var oyafiles []*Oyafile
	return oyafiles, filepath.Walk(rootDir, func(path string, info os.FileInfo, _ error) error {
		if !info.IsDir() {
			return nil
		}
		oyafilePath := fullPath(path, "")
		oyafile, buildable, err := Load(oyafilePath)
		if err != nil {
			return err
		}
		if !buildable {
			return nil
		}
		oyafiles = append(oyafiles, oyafile)
		return nil
	})
}

func fullPath(projectDir, name string) string {
	if len(name) == 0 {
		name = DefaultName
	}
	return path.Join(projectDir, name)
}
