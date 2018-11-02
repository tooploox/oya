package oyafile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type OyafileFormat = map[string]interface{}

type Alias string
type ImportPath string

type Oyafile struct {
	Dir       string
	Path      string
	VendorDir string
	Shell     string
	Imports   map[Alias]ImportPath
	Hooks     map[string]Hook
}

func New(oyafilePath string, vendorDir string) *Oyafile {
	return &Oyafile{
		Dir:       path.Dir(oyafilePath),
		Path:      oyafilePath,
		VendorDir: vendorDir,
		Shell:     "/bin/sh",
		Imports:   make(map[Alias]ImportPath),
		Hooks:     make(map[string]Hook),
	}
}

func Load(oyafilePath string, vendorDir string) (*Oyafile, bool, error) {
	if _, err := os.Stat(oyafilePath); os.IsNotExist(err) {
		return nil, false, nil
	}
	// YAML parser does not handle files without at least one node.
	empty, err := isEmptyYAML(oyafilePath)
	if err != nil {
		return nil, false, errors.Wrapf(err, "error loading Oyafile %s", oyafilePath)
	}
	if empty {
		return New(oyafilePath, vendorDir), true, nil
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
	oyafile, err := parseOyafile(oyafilePath, vendorDir, of)
	if err != nil {
		return nil, false, errors.Wrapf(err, "error parsing Oyafile %s", oyafilePath)
	}

	return oyafile, true, nil
}

func LoadFromDir(dirPath, vendorDir string) (*Oyafile, bool, error) {
	oyafilePath := fullPath(dirPath, "")
	return Load(oyafilePath, vendorDir)
}

func InitDir(dirPath string) error {
	f, err := os.Create(fullPath(dirPath, ""))
	if err != nil {
		return err
	}
	return f.Close()
}

func (oyafile Oyafile) ExecHook(hookName string, env map[string]string, stdout, stderr io.Writer) (found bool, err error) {
	hook, ok := oyafile.Hooks[hookName]
	if !ok {
		return false, nil
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

func parseOyafile(path, vendorDir string, of OyafileFormat) (*Oyafile, error) {
	oyafile := New(path, vendorDir)
	for name, value := range of {
		switch name {
		case "Import":
			imports, ok := value.(map[interface{}]interface{})
			if !ok {
				return nil, fmt.Errorf("Map of aliases to paths expected for key %q", name)
			}
			for alias, path := range imports {
				alias, ok := alias.(string)
				if !ok {
					return nil, fmt.Errorf("Expected string alias for key %q", name)
				}
				path, ok := path.(string)
				if !ok {
					return nil, fmt.Errorf("Expected path for key %q", name)
				}
				oyafile.Imports[Alias(alias)] = ImportPath(path)
			}
		default:
			script, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("Script expected for key %q", name)
			}
			oyafile.Hooks[name] = ScriptedHook{
				Name:   name,
				Script: Script(script),
				Shell:  oyafile.Shell,
			}
		}
	}

	return oyafile, oyafile.resolveImports()
}

func (oyafile *Oyafile) resolveImports() error {
	for alias, path := range oyafile.Imports {
		fullPath := filepath.Join(oyafile.VendorDir, string(path))
		log.Debugf("Importing Oyafile in %v as %v", fullPath, alias)
		imported, found, err := LoadFromDir(fullPath, oyafile.VendorDir)
		if err != nil {
			return errors.Wrap(err, "error resolving imports")
		}
		if !found {
			log.Debugf("Import %v has no Oyafile", path)
			// TODO: Warn?
			return nil
		}
		for key, hook := range imported.Hooks {
			// TODO: Detect if hook already set.
			log.Printf("Importing hook %v/%v", alias, key)
			oyafile.Hooks[fmt.Sprintf("%v/%v", alias, key)] = hook
		}
	}
	return nil
}
