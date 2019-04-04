package raw

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tooploox/oya/pkg/secrets"
	"github.com/tooploox/oya/pkg/template"
	yaml "gopkg.in/yaml.v2"
)

const DefaultName = "Oyafile"

// Oyafile represents an unparsed Oyafile.
type Oyafile struct {
	Path    string // Path contains normalized absolute path to the Oyafile.
	Dir     string // Dir contains normalized absolute path to the containing directory.
	RootDir string // RootDir is the absolute, normalized path to the project root directory.
	file    []byte // file contains Oyafile contents.
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
		Path:    normalizedOyafilePath,
		Dir:     filepath.Dir(normalizedOyafilePath),
		RootDir: rootDir,
		file:    file,
	}, nil
}

func (raw *Oyafile) Decode() (DecodedOyafile, error) {
	// YAML parser does not handle files without at least one node.
	empty, err := isEmptyYAML(raw.Path)
	if err != nil {
		return nil, err
	}
	if empty {
		return make(DecodedOyafile), nil
	}
	decodedOyafileI, err := decodeYaml(raw.file)
	if err != nil {
		return nil, err
	}
	decodedOyafile := DecodedOyafile(decodedOyafileI)

	secs, err := secrets.Decrypt(raw.Dir)
	if err != nil {
		if _, ok := err.(secrets.ErrNoSecretsFile); !ok {
			log.Debug(fmt.Sprintf("Secrets could not be loaded at %v: %v", raw.Dir, err))
		}
	} else {
		if len(secs) > 0 {
			decodedSecrets, err := decodeYaml(secs)
			if err != nil {
				log.Warn(fmt.Sprintf("Secrets could not be parsed after loading from %v: %v", raw.Dir, err))
			}
			secrets, ok := template.ParseScope(decodedSecrets)
			if !ok {
				return nil, errors.Errorf("Internal: error parsing scope trying to merge secrets, unexpected type: %T", decodedSecrets)
			}
			if err := mergeSecrets(&decodedOyafile, secrets); err != nil {
				return nil, err
			}
		}
	}

	return decodedOyafile, nil
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

func (raw *Oyafile) LookupKey(key string) (interface{}, bool, error) {
	of, err := raw.Decode()
	if err != nil {
		return nil, false, err
	}
	val, ok := of[key]
	return val, ok, nil
}

func (raw *Oyafile) IsRoot() (bool, error) {
	_, hasProject, err := raw.LookupKey("Project")
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

func mergeSecrets(of *DecodedOyafile, secrets template.Scope) error {
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
