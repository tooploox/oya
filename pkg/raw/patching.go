package raw

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// ApplyChanges flushes in-memory Oyafile contents to disk.
func (raw *Oyafile) ApplyChanges() error {
	info, err := os.Stat(raw.Path)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(raw.Path, raw.oyafileContents, info.Mode())
}

// flatMap maps each line of Oyafile to possibly multiple lines flattening the result. Does not write to the file.
func (raw *Oyafile) flatMap(f func(line string) []string) error {
	scanner := bufio.NewScanner(bytes.NewReader(raw.oyafileContents))

	output := bytes.NewBuffer(nil)
	for scanner.Scan() {
		line := scanner.Text()
		lines := f(line)
		for _, l := range lines {
			output.WriteString(l)
			output.WriteString("\n")
		}
	}

	raw.oyafileContents = output.Bytes()
	return nil
}

// insertAfter inserts one or more lines after the first line matching the regular expression.
// Does not write to the file. Preserves indentation; the new lines will have the same indentation
// as the preceding line.
func (raw *Oyafile) insertAfter(rx *regexp.Regexp, lines ...string) (bool, error) {
	found := false
	err := raw.flatMap(func(line string) []string {
		if !found && rx.MatchString(line) {
			found = true
			return append([]string{line}, indent(lines, detectIndent(line))...)
		} else {
			return []string{line}
		}
	})
	return found, err
}

// insertBefore inserts one or more lines before the first line matching the regular expression.
// Does not write to the file. Preserves indentation; the new lines will have the same indentation
// as the preceding line.
func (raw *Oyafile) insertBefore(rx *regexp.Regexp, lines ...string) (bool, error) {
	found := false
	err := raw.flatMap(func(line string) []string {
		if !found && rx.MatchString(line) {
			found = true
			return append(indent(lines, detectIndent(line)), line)
		} else {
			return []string{line}
		}
	})
	return found, err
}

// replaceAllWhen replaces all lines matching the predicate with one or more lines.
// Does not write to the file. Preserves indentation; the new lines will have the same indentation
// as the line being replaced.
func (raw *Oyafile) replaceAllWhen(pred func(string) bool, lines ...string) (bool, error) {
	found := false
	err := raw.flatMap(func(line string) []string {
		if pred(line) {
			found = true
			return indent(lines, detectIndent(line))
		} else {
			return []string{line}
		}
	})
	return found, err
}

// insertBeforeWithin inserts one or more lines before the first line matching the regular expression
// but - unlike insertBefore - only for lines nested under the specified top-level YAML key.
// Does not write to the file. Preserves indentation; the new lines will have the same indentation
// as the preceding line.
func (raw *Oyafile) insertBeforeWithin(key string, rx *regexp.Regexp, lines ...string) (bool, error) {
	keyRegexp := regexp.MustCompile(fmt.Sprintf("^%v:", key))
	topLevelKeyRegexp := regexp.MustCompile("^[\\s]+:")
	withinKey := false
	found := false
	err := raw.flatMap(func(line string) []string {
		if found {
			return []string{line}
		}
		if withinKey {
			if topLevelKeyRegexp.MatchString(line) {
				withinKey = false
			}

			if rx.MatchString(line) {
				found = true
				return append(indent(lines, detectIndent(line)), line)
			} else {
				return []string{line}
			}
		} else {
			if keyRegexp.MatchString(line) {
				withinKey = true
			}
		}
		return []string{line}
	})
	return found, err
}

// concat appends one or more lines to the Oyafile. Does not write to the file.
func (raw *Oyafile) concat(lines ...string) error {
	output := bytes.NewBuffer(raw.oyafileContents)
	if err := writeLines(lines, output); err != nil {
		return err
	}
	raw.oyafileContents = output.Bytes()
	return nil
}

// prepend takes an Oyafile and prepends its content with one or more lines.
// Does not write to the file.
func (raw *Oyafile) prepend(lines ...string) error {
	output := bytes.NewBuffer(nil)
	if err := writeLines(lines, output); err != nil {
		return err
	}
	if _, err := output.Write(raw.oyafileContents); err != nil {
		return err
	}
	raw.oyafileContents = output.Bytes()
	return nil
}

func writeLines(lines []string, output *bytes.Buffer) error {
	for _, l := range lines {
		if _, err := output.WriteString(l); err != nil {
			return err
		}
		if _, err := output.WriteString("\n"); err != nil {
			return err
		}
	}
	return nil
}

// containsLineMatching returns true if Oyafile matches the regular expression.
// (?m) directive may be used to match file line by line
func (raw *Oyafile) matches(rx *regexp.Regexp) bool {
	return rx.MatchString(string(raw.oyafileContents))
}

func indent(lines []string, level int) []string {
	indented := make([]string, len(lines))
	for i, line := range lines {
		indented[i] = fmt.Sprintf("%v%v", strings.Repeat(" ", level), line)
	}
	return indented
}

func detectIndent(line string) int {
	i := 0
	for _, runeValue := range line {
		if runeValue == ' ' {
			i++
		} else {
			break
		}
	}
	return i
}
