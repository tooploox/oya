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
	filePath string
	file     []byte
}

type RawOyafileFormat = map[string]interface{}

func LoadRawFromDir(dirPath string) (RawModifier, bool, error) {
	oyafilePath := fullPath(dirPath, "")
	fi, err := os.Stat(oyafilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return RawModifier{}, false, nil
		}
		return RawModifier{}, false, err
	}
	if fi.IsDir() {
		return RawModifier{}, false, nil
	}
	raw, err := NewRawModifier(oyafilePath)
	if err != nil {
		return RawModifier{}, false, nil
	}
	return raw, true, nil
}

func NewRawModifier(oyafilePath string) (RawModifier, error) {
	file, err := ioutil.ReadFile(oyafilePath)
	if err != nil {
		return RawModifier{}, err
	}

	return RawModifier{
		filePath: oyafilePath,
		file:     file,
	}, nil
}

func (o *RawModifier) HasKey(key string) (bool, error) {
	// YAML parser does not handle files without at least one node.
	empty, err := isEmptyYAML(o.filePath)
	if err != nil {
		return false, err
	}
	if empty {
		return false, nil
	}
	reader := bytes.NewReader(o.file)
	decoder := yaml.NewDecoder(reader)
	var of RawOyafileFormat
	err = decoder.Decode(&of)
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
