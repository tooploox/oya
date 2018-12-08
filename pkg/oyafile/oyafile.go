package oyafile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bilus/oya/pkg/template"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

const DefaultName = "Oyafile"

type OyafileFormat = map[string]interface{}

type Alias string
type ImportPath string

type Oyafile struct {
	Dir     string
	Path    string
	RootDir string
	Shell   string
	Imports map[Alias]ImportPath
	Tasks   map[string]Task
	Values  template.Scope
	Project string   // Set for root Oyafile
	Ignore  []string // Directory exclusion rules

	relPath string
}

func New(oyafilePath string, rootDir string) (*Oyafile, error) {
	relPath, err := filepath.Rel(rootDir, oyafilePath)
	if err != nil {
		return nil, err
	}
	dir := path.Dir(oyafilePath)
	return &Oyafile{
		Dir:     filepath.Clean(dir),
		Path:    filepath.Clean(oyafilePath),
		RootDir: filepath.Clean(rootDir),
		Shell:   "/bin/sh",
		Imports: make(map[Alias]ImportPath),
		Tasks:   make(map[string]Task),
		Values:  defaultValues(dir),
		relPath: relPath,
	}, nil
}

func Load(oyafilePath string, rootDir string) (*Oyafile, bool, error) {
	// YAML parser does not handle files without at least one node.
	empty, err := isEmptyYAML(oyafilePath)
	if err != nil {
		return nil, false, wrapLoadErr(err, oyafilePath)
	}
	if empty {
		o, err := New(oyafilePath, rootDir)
		if err != nil {
			return nil, false, err
		}
		return o, true, nil
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
	fi, err := os.Stat(oyafilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	if fi.IsDir() {
		return nil, false, nil
	}
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
	_, err = f.WriteString("Project: project\n")
	if err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

func (oyafile Oyafile) RunTask(taskName string, stdout, stderr io.Writer) (found bool, err error) {
	task, ok := oyafile.Tasks[taskName]
	if !ok {
		return false, nil
	}
	return true, task.Exec(stdout, stderr)
}

func (oyafile Oyafile) Equals(other Oyafile) bool {
	// TODO: Far from perfect, we should ensure relative vs absolute paths work.
	// The simplest thing is probably to ensure oyafile.Dir is always absolute.
	return filepath.Clean(oyafile.Dir) == filepath.Clean(other.Dir)
}

func (oyafile Oyafile) IsRoot() bool {
	return oyafile.Project != "" && filepath.Clean(oyafile.Dir) == filepath.Clean(oyafile.RootDir)
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

func defaultValues(dirPath string) template.Scope {
	return template.Scope{
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
	oyafile, err := New(path, rootDir)
	if err != nil {
		return nil, err
	}
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
		case "Project":
			projectName, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("expected Project: defining the project name, actual: %v", value)
			}
			oyafile.Project = projectName
		case "Ignore":
			rulesI, ok := value.([]interface{})
			if !ok {
				return nil, fmt.Errorf("expected Ignore: containing an array of ignore rules, actual: %v", value)
			}
			rules := make([]string, len(rulesI))
			for i, ri := range rulesI {
				rule, ok := ri.(string)
				if !ok {
					return nil, fmt.Errorf("expected Ignore: containing an array of ignore rules, actual: %v", ri)
				}
				rules[i] = rule
			}
			oyafile.Ignore = rules
		default:
			script, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("script expected for key %q", name)
			}
			oyafile.Tasks[name] = ScriptedTask{
				Name:   name,
				Script: Script(script),
				Shell:  oyafile.Shell,
				Scope:  &oyafile.Values,
			}
		}
	}

	return oyafile, nil
}

func (o *Oyafile) Ignores() string {
	return strings.Join(o.Ignore, "\n")
}

func (o *Oyafile) RelPath() string {
	return o.relPath
}
