package oyafile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
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
	// YAML parser does not handle files without at least one node.
	empty, err := isEmptyYAML(oyafilePath)
	if err != nil {
		return nil, false, errors.Wrapf(err, "error loading Oyafile %s", oyafilePath)
	}
	if empty {
		return New(oyafilePath), true, nil
	}
	file, err := os.Open(oyafilePath)
	if err != nil {
		return nil, false, errors.Wrapf(err, "error loading Oyafile %s", oyafilePath)
	}
	defer func() { _ = file.Close() }()
	decoder := yaml.NewDecoder(file)
	var of OyafileFormat
	err = decoder.Decode(&of)
	if err != nil {
		return nil, false, errors.Wrapf(err, "error parsing Oyafile %s", oyafilePath)
	}
	oyafile := New(oyafilePath)
	for name, script := range of {
		oyafile.Hooks[name] = ScriptedHook{
			Name:   name,
			Script: script,
			Shell:  oyafile.Shell,
		}
	}
	return oyafile, true, nil
}

func LoadFromDir(dirPath string) (*Oyafile, bool, error) {
	oyafilePath := fullPath(dirPath, "")
	return Load(oyafilePath)
}

func (oyafile Oyafile) ExecHook(hookName string, env map[string]string, stdout, stderr io.Writer) (found bool, err error) {
	hook, ok := oyafile.Hooks[hookName]
	if !ok {
		return false, fmt.Errorf("no such hook: %v", hookName)
	}
	return true, hook.Exec(nil, stdout, stderr)
}

func (oyafile Oyafile) Equals(other Oyafile) bool {
	// TODO: Far from perfect, we should ensure relative vs absolute paths work.
	// The simplest thing is probably to ensure oyafile.Dir is always absolute.
	return filepath.Clean(oyafile.Dir) == filepath.Clean(other.Dir)
}

func fullPath(projectDir, name string) string {
	if len(name) == 0 {
		name = DefaultName
	}
	return path.Join(projectDir, name)
}

// isEmptyYAML returns true if the Oyafile contains only blank characters or YAML comments.
func isEmptyYAML(oyafilePath string) (bool, error) {
	file, err := os.Open(oyafilePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if isNode(scanner.Text()) {
			return false, nil
		}
	}

	return true, scanner.Err()
}

func isNode(line string) bool {
	for _, c := range line {
		switch c {
		case '#':
			return false
		case ' ', '\t', '\n', '\f', '\r':
			continue
		default:
			return true
		}
	}
	return false
}
