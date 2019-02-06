package raw

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
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
	return raw.write()
}

func (raw *Oyafile) concat(lines ...string) error {
	output := bytes.NewBuffer(raw.file)
	for _, l := range lines {
		output.WriteString(l)
		output.WriteString("\n")
	}

	raw.file = output.Bytes()
	return raw.write()
}

func (raw *Oyafile) write() error {
	info, err := os.Stat(raw.Path)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(raw.Path, raw.file, info.Mode())
}
