package oyafile

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type OyafileRawModifier struct {
	filePath string
	file     []byte
}

func NewOyafileRawModifier(oyafilePath string) (OyafileRawModifier, error) {
	file, err := ioutil.ReadFile(oyafilePath)
	if err != nil {
		return OyafileRawModifier{}, err
	}

	return OyafileRawModifier{
		filePath: oyafilePath,
		file:     file,
	}, nil
}

func (o *OyafileRawModifier) addImport(name string, uri string) error {
	importStr := "Import:"
	uriStr := fmt.Sprintf("  %s: %s", name, uri)
	fileContent := string(o.file)
	var output []string
	if strings.Contains(fileContent, "Import:") {
		output = o.appendAfter("Import:", []string{uriStr})
	} else if strings.Contains(fileContent, "Project:") {
		output = o.appendAfter("Project:", []string{importStr, uriStr})
	} else {
		output = []string{importStr, uriStr}
		output = append(output, strings.Split(fileContent, "\n")...)
	}
	if err := writeToFile(o.filePath, output); err != nil {
		return err
	}

	return nil
}

func (o *OyafileRawModifier) appendAfter(find string, data []string) []string {
	var output []string
	fileArr := strings.Split(string(o.file), "\n")
	for _, line := range fileArr {
		output = append(output, line)
		if strings.Contains(line, find) {
			output = append(output, data...)
		}
	}
	return output
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
