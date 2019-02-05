package raw

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

const DefaultName = "Oyafile"

var importKey = "Import:"
var projectKey = "Project:"
var uriVal = "  %s: %s"
var importRegxp = regexp.MustCompile("(?m)^" + importKey + "$")
var projectRegxp = regexp.MustCompile("^" + projectKey)

type Oyafile struct {
	Path    string
	RootDir string
	file    []byte
}

// DecodedOyafile is an Oyafile that has been loaded from YAML but that hasn't been parsed yet.
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

	return &Oyafile{
		RootDir: rootDir,
		Path:    oyafilePath,
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
	reader := bytes.NewReader(raw.file)
	decoder := yaml.NewDecoder(reader)
	var of DecodedOyafile
	err = decoder.Decode(&of)
	if err != nil {
		return nil, err
	}
	return of, nil
}

func (raw *Oyafile) HasKey(key string) (bool, error) {
	of, err := raw.Decode()
	if err != nil {
		return false, err
	}
	_, ok := of[key]
	return ok, nil
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

func (o *Oyafile) AddImport(alias string, uri string) error {
	var output []string
	uriStr := fmt.Sprintf(uriVal, alias, uri)
	fileContent := string(o.file)
	updated := false

	if gotIt := o.isAlreadyImported(uri, fileContent); gotIt {
		return errors.Errorf("Pack already imported: %v", uri)
	}

	output, updated = o.appendAfter(importRegxp, []string{uriStr})
	if !updated {
		output, updated = o.appendAfter(projectRegxp, []string{importKey, uriStr, ""})
		if !updated {
			output = []string{importKey, uriStr}
			output = append(output, strings.Split(fileContent, "\n")...)
		}
	}

	if err := writeToFile(o.Path, output); err != nil {
		return err
	}

	// BUG(bilus): Does not update o.file!

	return nil
}

func (o *Oyafile) isAlreadyImported(uri string, fileContent string) bool {
	find := regexp.MustCompile("(?m)" + uri + "$")
	return find.MatchString(fileContent)
}

func (o *Oyafile) appendAfter(find *regexp.Regexp, data []string) ([]string, bool) {
	var output []string
	updated := false
	fileArr := strings.Split(string(o.file), "\n")
	for _, line := range fileArr {
		output = append(output, line)
		if find.MatchString(line) {
			updated = true
			output = append(output, data...)
		}
	}
	return output, updated
}

func writeToFile(Path string, content []string) error {
	info, err := os.Stat(Path)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(Path, []byte(strings.Join(content, "\n")), info.Mode()); err != nil {
		return err
	}
	return nil
}

func fullPath(projectDir, name string) string {
	if len(name) == 0 {
		name = DefaultName
	}
	return path.Join(projectDir, name)
}
