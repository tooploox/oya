package oyafile

import (
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"
)

const DefaultName = "Oyafile"

type OyafileFormat struct {
	Jobs map[string]Script `yaml:"jobs"`
}

type Oyafile struct {
	Jobs map[string]Job
}

func Load(oyafilePath string) (*Oyafile, error) {
	file, err := os.Open(oyafilePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	decoder := yaml.NewDecoder(file)
	var of OyafileFormat
	err = decoder.Decode(&of)
	if err != nil {
		return nil, err
	}
	oyafile := Oyafile{
		Jobs: make(map[string]Job),
	}
	for name, script := range of.Jobs {
		oyafile.Jobs[name] = Job{
			Name:   name,
			Script: script,
		}
	}
	return &oyafile, err
}

func FullPath(projectDir, name string) string {
	if len(name) == 0 {
		name = DefaultName
	}
	return path.Join(projectDir, name)
}
