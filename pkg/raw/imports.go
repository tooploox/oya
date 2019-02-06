package raw

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var importKey = "Import:"
var projectKey = "Project:"
var uriVal = "  %s: %s"
var importRegexp = regexp.MustCompile("(?m)^" + importKey + "$")
var projectRegexp = regexp.MustCompile("^" + projectKey)

func (o *Oyafile) AddImport(alias string, uri string) error {
	var output []string
	uriStr := fmt.Sprintf(uriVal, alias, uri)
	fileContent := string(o.file)
	updated := false

	if gotIt := o.isAlreadyImported(uri, fileContent); gotIt {
		return errors.Errorf("Pack already imported: %v", uri)
	}

	output, updated = o.appendAfter(importRegexp, []string{uriStr})
	if !updated {
		output, updated = o.appendAfter(projectRegexp, []string{importKey, uriStr})
		if !updated {
			output = []string{importKey, uriStr}
			output = append(output, strings.Split(fileContent, "\n")...)
		}
	}

	if err := writeToFile(o.Path, output); err != nil {
		return err
	}

	file, err := ioutil.ReadFile(o.Path)
	if err != nil {
		return err
	}

	o.file = file

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
