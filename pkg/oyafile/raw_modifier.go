package oyafile

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

var importKey = "Import:"
var projectKey = "Project:"
var uriVal = "  %s: %s"
var importRegxp = regexp.MustCompile("(?m)^" + importKey + "$")
var projectRegxp = regexp.MustCompile("^" + projectKey)

type RawModifier struct {
	rootDir  string
	filePath string
	file     []byte
}

type RawOyafileFormat = map[string]interface{}

func LoadRaw(oyafilePath, rootDir string) (*RawModifier, bool, error) {
	raw, err := NewRawModifier(oyafilePath)
	if err != nil {
		return nil, false, nil
	}
	raw.rootDir = rootDir // BUG(bilus): NewRawModifier has a broken interface.
	return raw, true, nil
}

func LoadRawFromDir(dirPath, rootDir string) (*RawModifier, bool, error) {
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
	return LoadRaw(oyafilePath, rootDir)
}

func NewRawModifier(oyafilePath string) (*RawModifier, error) {
	file, err := ioutil.ReadFile(oyafilePath)
	if err != nil {
		return nil, err
	}

	return &RawModifier{
		filePath: oyafilePath,
		file:     file,
	}, nil
}

func (raw *RawModifier) Parse() (*Oyafile, error) {
	of, err := raw.decode()
	if err != nil {
		return nil, err
	}

	return parseOyafile(raw.filePath, raw.rootDir, of)
}

func (raw *RawModifier) decode() (RawOyafileFormat, error) {
	// YAML parser does not handle files without at least one node.
	empty, err := isEmptyYAML(raw.filePath)
	if err != nil {
		return nil, err
	}
	if empty {
		return make(RawOyafileFormat), nil
	}
	reader := bytes.NewReader(raw.file)
	decoder := yaml.NewDecoder(reader)
	var of RawOyafileFormat
	err = decoder.Decode(&of)
	if err != nil {
		return nil, err
	}
	return of, nil
}

func (raw *RawModifier) HasKey(key string) (bool, error) {
	of, err := raw.decode()
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

func (o *RawModifier) addImport(name string, uri string) error {
	var output []string
	uriStr := fmt.Sprintf(uriVal, name, uri)
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

	if err := writeToFile(o.filePath, output); err != nil {
		return err
	}

	// BUG(bilus): Does not update o.file!

	return nil
}

func (o *RawModifier) isAlreadyImported(uri string, fileContent string) bool {
	find := regexp.MustCompile("(?m)" + uri + "$")
	return find.MatchString(fileContent)
}

func (o *RawModifier) appendAfter(find *regexp.Regexp, data []string) ([]string, bool) {
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

func writeToFile(filePath string, content []string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filePath, []byte(strings.Join(content, "\n")), info.Mode()); err != nil {
		return err
	}
	return nil
}
