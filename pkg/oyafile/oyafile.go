package oyafile

import (
	"bufio"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/bilus/oya/pkg/template"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

const DefaultName = "Oyafile"

// OyaCmdOverride is used in tests, to override the path to the current oya executable.
// It is used to invoke other tasks from a task body.
// When tests are run, the current process executable path points to the test runner
// so it has to be overridden (with 'go run oya.go', roughly speaking).
var OyaCmdOverride *string

type OyafileFormat = map[string]interface{}

type Alias string
type ImportPath string

type Oyafile struct {
	Dir     string
	Path    string
	RootDir string
	Shell   string
	Imports map[Alias]ImportPath
	Tasks   TaskTable
	Values  template.Scope
	Project string   // Project is set for root Oyafile.
	Ignore  []string // Ignore contains directory exclusion rules.

	relPath string

	OyaCmd string // OyaCmd contains the path to the current oya executable.
}

func New(oyafilePath string, rootDir string) (*Oyafile, error) {
	var oyaCmd string
	if OyaCmdOverride != nil {
		oyaCmd = *OyaCmdOverride
	} else {
		var err error
		oyaCmd, err = os.Executable()
		if err != nil {
			return nil, err
		}
	}

	relPath, err := filepath.Rel(rootDir, oyafilePath)
	log.Debug("Oyafile at", oyafilePath)
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
		Tasks:   newTaskTable(),
		Values:  template.Scope{},
		relPath: relPath,
		OyaCmd:  oyaCmd,
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
	err = oyafile.addBuiltIns()
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

func (oyafile Oyafile) RunTask(taskName string, scope template.Scope, stdout, stderr io.Writer) (found bool, err error) {
	task, ok := oyafile.Tasks.LookupTask(taskName)
	if !ok {
		return false, nil
	}
	tasks, err := oyafile.bindTasks(task, stdout, stderr)
	if err != nil {
		return true, err
	}
	scope["Tasks"] = tasks
	if err != nil {
		return true, err
	}
	return true, task.Exec(oyafile.Dir, scope, stdout, stderr)
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

func (o *Oyafile) Ignores() string {
	return strings.Join(o.Ignore, "\n")
}

func (o *Oyafile) RelPath() string {
	return o.relPath
}
