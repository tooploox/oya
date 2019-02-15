package raw

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/bilus/oya/pkg/secrets"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

const DefaultName = "Oyafile"
const SecretsSuffix = ".secrets"

// Oyafile represents an unparsed Oyafile.
type Oyafile struct {
	Path    string // Path contains normalized absolute path.
	RootDir string // RootDir is the absolute, normalized path to the project root directory.
	file    []byte // file contains Oyafile contents.
	secrets []byte // secrets contains Oyafile.secrets contents.
}

// DecodedOyafile is an Oyafile that has been loaded from YAML
// but hasn't been parsed yet.
type DecodedOyafile map[string]interface{}

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
	secs, err := secrets.DecryptSecrets(oyafilePath + SecretsSuffix)
	if err != nil {
		log.Debug("Secrets could not be loaded at %v. %v", oyafilePath+SecretsSuffix, err)
	}
	return &Oyafile{
		Path:    oyafilePath,
		RootDir: rootDir,
		file:    file,
		secrets: secs,
	}, nil
}

func (raw *Oyafile) Decode() (DecodedOyafile, error) {
	// YAML parser does not handle files without at least one node.
	// TODO: do we need this ? maybe we can chack if raw.file is not empty?
	empty, err := isEmptyYAML(raw.Path)
	if err != nil {
		return nil, err
	}
	if empty {
		return make(DecodedOyafile), nil
	}
	reader := bytes.NewReader(raw.Content())
	decoder := yaml.NewDecoder(reader)
	var of DecodedOyafile
	err = decoder.Decode(&of)
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

func (raw *Oyafile) Content() []byte {
	if len(raw.secrets) > 0 {
		output := append(raw.file, []byte("\n")...)
		return append(output, raw.secrets...)
	} else {
		return raw.file
	}
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
