package raw

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"regexp"
)

// flatMap maps each line of Oyafile to possibly multiple lines flattening the result. Does not write to the file.
func (raw *Oyafile) flatMap(f func(line string) []string) error {
	scanner := bufio.NewScanner(bytes.NewReader(raw.file))

	output := bytes.NewBuffer(nil)
	for scanner.Scan() {
		line := scanner.Text()
		lines := f(line)
		for _, l := range lines {
			output.WriteString(l)
			output.WriteString("\n")
		}
	}

	raw.file = output.Bytes()
	return nil
}

// insertAfter inserts one or more lines after the first line matching the regular expression. Does not write to the file.
func (raw *Oyafile) insertAfter(rx *regexp.Regexp, lines ...string) (bool, error) {
	found := false
	err := raw.flatMap(func(line string) []string {
		if !found && rx.MatchString(line) {
			found = true
			return append([]string{line}, lines...)
		} else {
			return []string{line}
		}
	})
	return found, err
}

// concat appends one or more lines to the Oyafile. Does not write to the file.
func (raw *Oyafile) concat(lines ...string) error {
	output := bytes.NewBuffer(raw.file)
	for _, l := range lines {
		output.WriteString(l)
		output.WriteString("\n")
	}

	raw.file = output.Bytes()
	return nil
}

// write flushes in-memory Oyafile contents to disk.
func (raw *Oyafile) write() error {
	info, err := os.Stat(raw.Path)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(raw.Path, raw.file, info.Mode())
}
