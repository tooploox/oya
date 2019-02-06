package raw

import (
	"fmt"
	"regexp"

	"github.com/bilus/oya/pkg/pack"
)

var requireKeyRegxp = regexp.MustCompile("^Require:\\s*$")
var requireEntryRegexp = regexp.MustCompile("^\\s*([^:]+)\\:\\s*(.+)$")

func (raw *Oyafile) AddRequire(pack pack.Pack) error {
	found := false

	err := raw.flatMap(func(line string) []string {
		if found {
			return []string{line}
		}

		if matches := requireEntryRegexp.FindStringSubmatch(line); len(matches) == 3 {
			if matches[1] == pack.ImportUrl() {
				found = true
				return []string{formatRequire(pack)}
			}
		}

		return []string{line}
	})

	if err != nil {
		return err
	}

	if found {
		return nil
	}

	err = raw.flatMap(func(line string) []string {
		if !found && requireKeyRegxp.MatchString(line) {
			found = true
			return []string{line, formatRequire(pack)}
		} else {
			return []string{line}
		}
	})

	if err != nil {
		return err
	}

	if found {
		return nil
	}

	return raw.concat(
		"Require:",
		formatRequire(pack))
}

func formatRequire(pack pack.Pack) string {
	return fmt.Sprintf("  %v: %v", pack.ImportUrl(), pack.Version())
}
