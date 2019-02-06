package raw

import (
	"fmt"
	"regexp"

	"github.com/bilus/oya/pkg/pack"
)

var requireKeyRegxp = regexp.MustCompile("^Require:\\s*$")
var requireEntryRegexp = regexp.MustCompile("^\\s*([^:]+)\\:\\s*(.+)$")

func (raw *Oyafile) AddRequire(pack pack.Pack) error {
	if err := raw.addRequire(pack); err != nil {
		return err
	}
	return raw.write()
}

func (raw *Oyafile) addRequire(pack pack.Pack) error {
	if found, err := raw.updateExistingEntry(pack); err != nil || found {
		return err // nil if found
	}
	if found, err := raw.addToRequire(pack); err != nil || found {
		return err // nil if found
	}
	if err := raw.concat("Require:", formatRequire(pack)); err != nil {
		return err // never nil
	}
	return raw.write()
}

func (raw *Oyafile) updateExistingEntry(pack pack.Pack) (bool, error) {
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
	return found, err
}

func (raw *Oyafile) addToRequire(pack pack.Pack) (bool, error) {
	return raw.insertAfter(requireKeyRegxp, formatRequire(pack))
}

func formatRequire(pack pack.Pack) string {
	return fmt.Sprintf("  %v: %v", pack.ImportUrl(), pack.Version())
}
