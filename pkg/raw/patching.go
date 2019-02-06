package raw

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"regexp"
)

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

func (raw *Oyafile) concat(lines ...string) error {
	output := bytes.NewBuffer(raw.file)
	for _, l := range lines {
		output.WriteString(l)
		output.WriteString("\n")
	}

	raw.file = output.Bytes()
	return nil
}

func (raw *Oyafile) write() error {
	info, err := os.Stat(raw.Path)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(raw.Path, raw.file, info.Mode())
}
