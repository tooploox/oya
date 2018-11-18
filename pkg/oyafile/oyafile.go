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

type OyafileFormat = map[string]interface{}

type Alias string
type ImportPath string

type Oyafile struct {
	Dir     string
	Path    string
	RootDir string
	Shell   string
	Imports map[Alias]ImportPath
	Hooks   map[string]Hook
	Values  Scope
	Module  string // Set for root Oyafile
}

type Scope map[string]interface{}

func New(oyafilePath string, rootDir string) *Oyafile {

	dir := path.Dir(oyafilePath)
	return &Oyafile{
		Dir:     dir,
		Path:    oyafilePath,
		RootDir: rootDir,
		Shell:   "/bin/sh",
		Imports: make(map[Alias]ImportPath),
		Hooks:   make(map[string]Hook),
		Values:  defaultValues(dir),
	}
}

func Load(oyafilePath string, rootDir string) (*Oyafile, bool, error) {
	if _, err := os.Stat(oyafilePath); os.IsNotExist(err) {
		return nil, false, nil
	}
	// YAML parser does not handle files without at least one node.
	empty, err := isEmptyYAML(oyafilePath)
	if err != nil {
		return nil, false, wrapLoadErr(err, oyafilePath)
	}
	if empty {
		return New(oyafilePath, rootDir), true, nil
	}
	file, err := os.Open(oyafilePath)
	if err != nil {
		return nil, false, wrapLoadErr(err, oyafilePath)
	}
	defer func() { _ = file.Close() }()
	decoder := yaml.NewDecoder(file)
	var of OyafileFormat
	err = decoder.Decode(&of)
	if err != nil {
		return nil, false, wrapLoadErr(err, oyafilePath)
	}
	oyafile, err := parseOyafile(oyafilePath, rootDir, of)
	if err != nil {
		return nil, false, wrapLoadErr(err, oyafilePath)
	}
	err = oyafile.resolveImports()
	if err != nil {
		return nil, false, wrapLoadErr(err, oyafilePath)
	}

	return oyafile, true, nil
}

func LoadFromDir(dirPath, rootDir string) (*Oyafile, bool, error) {
	oyafilePath := fullPath(dirPath, "")
	return Load(oyafilePath, rootDir)
}

func InitDir(dirPath string) error {
	_, found, err := LoadFromDir(dirPath, dirPath)
	if err == nil && found {
		return errors.Errorf("already an Oya project")
	}
	f, err := os.Create(fullPath(dirPath, ""))
	if err != nil {
		return err
	}
	_, err = f.WriteString("Module: project\n")
	if err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

func (oyafile Oyafile) ExecHook(hookName string, stdout, stderr io.Writer) (found bool, err error) {
	hook, ok := oyafile.Hooks[hookName]
	if !ok {
		return false, nil
	}
	return true, hook.Exec(stdout, stderr)
}

func (oyafile Oyafile) Equals(other Oyafile) bool {
	// TODO: Far from perfect, we should ensure relative vs absolute paths work.
	// The simplest thing is probably to ensure oyafile.Dir is always absolute.
	return filepath.Clean(oyafile.Dir) == filepath.Clean(other.Dir)
}

func wrapLoadErr(err error, oyafilePath string) error {
	return errors.Wrapf(err, "error loading Oyafile %v", oyafilePath)
}

func fullPath(projectDir, name string) string {
	if len(name) == 0 {
		name = DefaultName
	}
	return path.Join(projectDir, name)
}

func defaultValues(dirPath string) Scope {
	return Scope{
		"BasePath": dirPath,
	}
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

func parseOyafile(path, rootDir string, of OyafileFormat) (*Oyafile, error) {
	oyafile := New(path, rootDir)
	for name, value := range of {
		switch name {
		case "Import":
			imports, ok := value.(map[interface{}]interface{})
			if !ok {
				return nil, fmt.Errorf("map of aliases to paths expected for key %q", name)
			}
			for alias, path := range imports {
				alias, ok := alias.(string)
				if !ok {
					return nil, fmt.Errorf("expected string alias for key %q", name)
				}
				path, ok := path.(string)
				if !ok {
					return nil, fmt.Errorf("expected path for key %q", name)
				}
				oyafile.Imports[Alias(alias)] = ImportPath(path)
			}
		case "Values":
			values, ok := value.(map[interface{}]interface{})
			if !ok {
				return nil, fmt.Errorf("map of keys to values expected for key %q", name)
			}
			for k, v := range values {
				valueName, ok := k.(string)
				if !ok {
					return nil, fmt.Errorf("map of keys to values expected for key %q", name)
				}
				oyafile.Values[valueName] = v
			}
		case "Module":
			moduleName, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("expected Module to point to a string, actual: name")
			}
			oyafile.Module = moduleName
		default:
			script, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("script expected for key %q", name)
			}
			oyafile.Hooks[name] = ScriptedHook{
				Name:   name,
				Script: Script(script),
				Shell:  oyafile.Shell,
				Scope:  &oyafile.Values,
			}
		}
	}

	return oyafile, nil
}
