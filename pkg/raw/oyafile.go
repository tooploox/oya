package raw

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/tooploox/oya/pkg/secrets"
	"github.com/tooploox/oya/pkg/template"
	yaml "gopkg.in/yaml.v2"
)

const DefaultName = "Oyafile"
const ValueFileExt = ".oya"

// Oyafile represents an unparsed Oyafile.
type Oyafile struct {
	Path            string // Path contains normalized absolute path to the Oyafile.
	Dir             string // Dir contains normalized absolute path to the containing directory.
	RootDir         string // RootDir is the absolute, normalized path to the project root directory.
	oyafileContents []byte // file contains the main Oyafile contents.
}

// DecodedOyafile is an Oyafile that has been loaded from YAML
// but hasn't been parsed yet.
type DecodedOyafile template.Scope

func (lhs DecodedOyafile) Merge(rhs DecodedOyafile) DecodedOyafile {
	return DecodedOyafile(template.Scope(lhs).Merge(template.Scope(rhs)))
}

func Load(oyafilePath, rootDir string) (*Oyafile, bool, error) {
	raw, err := New(oyafilePath, rootDir)
	if err != nil {
		return nil, false, nil
	}
	return raw, true, nil
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

func New(oyafilePath, rootDir string) (*Oyafile, error) {
	file, err := ioutil.ReadFile(oyafilePath)
	if err != nil {
		return nil, err
	}
	normalizedOyafilePath := filepath.Clean(oyafilePath)
	return &Oyafile{
		Path:            normalizedOyafilePath,
		Dir:             filepath.Dir(normalizedOyafilePath),
		RootDir:         rootDir,
		oyafileContents: file,
	}, nil
}

func (raw *Oyafile) Decode() (DecodedOyafile, error) {
	mainOyafile, err := decodeOyafile(raw)
	if err != nil {
		return nil, err
	}

	paths, err := listFiles(raw.Dir, ValueFileExt)
	if err != nil {
		return nil, err
	}
	for _, path := range paths {
		rawValueFile, found, err := Load(path, raw.RootDir)
		if err != nil {
			return nil, err
		}
		if !found {
			return nil, errors.Errorf("Internal error: %s file not found while loading", path)
		}
		valueFile, err := decodeOyafile(rawValueFile)
		if err != nil {
			return nil, err
		}
		values, ok := template.ParseScope(map[interface{}]interface{}(valueFile))
		if !ok {
			return nil, errors.Errorf("Internal: error parsing scope trying to merge values, unexpected type: %T", valueFile)
		}
		if err := mergeValues(&mainOyafile, values); err != nil {
			return nil, err
		}
	}
	return mainOyafile, nil
}

func listFiles(path, ext string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		path := file.Name()
		if !file.IsDir() && filepath.Ext(path) == ext {
			files = append(files, path)
		}
	}
	return files, nil
}

func decodeOyafile(raw *Oyafile) (DecodedOyafile, error) {
	decrypted, found, err := secrets.Decrypt(raw.Path)
	if err != nil {
		return nil, err
	}
	if found {
		decodedSecrets, err := decodeYaml(decrypted)
		if err != nil {
			return nil, errors.Wrapf(err, "error parsing secret file %q", raw.Path)
		}
		return decodedSecrets, nil
	}

	// YAML parser does not handle files without at least one node.
	empty, err := isEmptyYAML(raw.Path)
	if err != nil {
		return nil, err
	}
	if empty {
		return make(DecodedOyafile), nil
	}
	decodedOyafileI, err := decodeYaml(raw.oyafileContents)
	if err != nil {
		return nil, err
	}
	return DecodedOyafile(decodedOyafileI), nil
}

func decodeYaml(content []byte) (map[interface{}]interface{}, error) {
	reader := bytes.NewReader(content)
	decoder := yaml.NewDecoder(reader)
	var of map[interface{}]interface{}
	err := decoder.Decode(&of)
	if err != nil {
		return nil, err
	}
	return of, nil
}

func (raw *Oyafile) Project() (interface{}, bool, error) {
	of, err := decodeOyafile(raw)
	if err != nil {
		return nil, false, err
	}
	val, ok := of["Project"]
	return val, ok, nil
}

func (raw *Oyafile) IsRoot() (bool, error) {
	_, hasProject, err := raw.Project()
	if err != nil {
		return false, err
	}

	rel, err := filepath.Rel(raw.RootDir, raw.Path)
	if err != nil {
		return false, err
	}
	return hasProject && rel == DefaultName, nil
}

// isEmptyYAML returns true if the Oyafile contains only blank characters
// or YAML comments.
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

func fullPath(projectDir, name string) string {
	if len(name) == 0 {
		name = DefaultName
	}
	return path.Join(projectDir, name)
}

func mergeValues(of *DecodedOyafile, secrets template.Scope) error {
	var values template.Scope
	valuesI, ok := (*of)["Values"]
	if ok {
		values, ok = template.ParseScope(valuesI)
		if !ok {
			return errors.Errorf("Internal: error parsing scope")
		}
	}
	(*of)["Values"] = map[interface{}]interface{}(values.Merge(secrets))
	return nil
}
