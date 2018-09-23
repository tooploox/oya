package oyafile

import (
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"
)

const DefaultName = "Oyafile"

type OyafileFormat = map[string]Script

type Oyafile struct {
	Hooks map[string]Hook
}

func Load(oyafilePath string) (*Oyafile, bool, error) {
	if _, err := os.Stat(oyafilePath); os.IsNotExist(err) {
		return nil, false, nil
	}
	file, err := os.Open(oyafilePath)
	if err != nil {
		return nil, false, err
	}
	defer func() { _ = file.Close() }()
	decoder := yaml.NewDecoder(file)
	var of OyafileFormat
	err = decoder.Decode(&of)
	if err != nil {
		return nil, false, err
	}
	oyafile := Oyafile{
		Hooks: make(map[string]Hook),
	}
	for name, script := range of {
		oyafile.Hooks[name] = Hook{
			Name:   name,
			Script: script,
		}
	}
	return &oyafile, true, err
}

func FullPath(projectDir, name string) string {
	if len(name) == 0 {
		name = DefaultName
	}
	return path.Join(projectDir, name)
}
