package oyafile

import (
	"fmt"
	"io"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"
)

type OyafileFormat = map[string]Script

type Oyafile struct {
	Dir   string
	Path  string
	Shell string
	Hooks map[string]Hook
}

func New(oyafilePath string) *Oyafile {
	return &Oyafile{
		Dir:   path.Dir(oyafilePath),
		Path:  oyafilePath,
		Shell: "/bin/sh",
		Hooks: make(map[string]Hook),
	}
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
	oyafile := New(oyafilePath)
	for name, script := range of {
		oyafile.Hooks[name] = Hook{
			Name:   name,
			Script: script,
		}
	}
	return oyafile, true, err
}

func (oyafile Oyafile) ExecHook(hookName string, env map[string]string, stdout, stderr io.Writer) (found bool, err error) {
	hook, ok := oyafile.Hooks[hookName]
	if !ok {
		return false, fmt.Errorf("no such hook: %v", hookName)
	}
	return true, hook.Script.Exec(nil, stdout, stderr, oyafile.Shell)
}
