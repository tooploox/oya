package raw

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bilus/oya/pkg/pack"
)

type ErrNotRootOyafile struct {
	Path string
}

func (e ErrNotRootOyafile) Error() string {
	return fmt.Sprintf("missing Project: directive in %v; not root Oyafile?", e.Path)
}

var requireKeyRegxp = regexp.MustCompile("^Require:")
var requireEntryRegexp = regexp.MustCompile("^(\\s*)([^:]+)\\:\\s*([^ #]+)")
var topLevelKeyRegexp = regexp.MustCompile("^[\\s]+:")

var defaultIndent = 2

// AddRequire adds a Require: entry for the pack.
func (raw *Oyafile) AddRequire(pack pack.Pack) error {
	if err := raw.addRequire(pack); err != nil {
		return err
	}
	return raw.write()
}

// addRequire adds a require for a pack using the following algorithm:
// 1. Look for and update an existing entry for the path.
// 2. Look for ANY pack under Require:; if found, insert the new entry beneath it.
// 3. Look for Require: key (we know it's empty), insert the new entry inside it.
// 4. Look for Project: key, insert the new entry beneath it (under Require:).
// 5. Fail because Oyafile has no Project: so we shouldn't be trying to add a require to it.
// The method stops if any of the steps succeeds.
func (raw *Oyafile) addRequire(pack pack.Pack) error {
	if found, err := raw.updateExistingEntry(pack); err != nil || found {
		return err // nil if found
	}
	if found, err := raw.insertBeforeExistingEntry(pack); err != nil || found {
		return err // nil if found
	}
	if found, err := raw.insertAfter(requireKeyRegxp, formatRequire(2, pack)); err != nil || found {
		return err // nil if found
	}

	found, err := raw.insertAfter(projectRegexp, "Require:", formatRequire(defaultIndent, pack))
	if err != nil {
		return err
	}

	if !found {
		return ErrNotRootOyafile{Path: raw.Path}
	}
	return raw.write()
}

func (raw *Oyafile) updateExistingEntry(pack pack.Pack) (bool, error) {
	found := false
	err := raw.flatMap(func(line string) []string {
		if found {
			return []string{line}
		}

		if matches := requireEntryRegexp.FindStringSubmatch(line); len(matches) == 4 {
			if matches[2] == pack.ImportUrl() {
				found = true
				indent := len(matches[1])
				return []string{formatRequire(indent, pack)}
			}
		}

		return []string{line}
	})
	return found, err
}

func (raw *Oyafile) insertBeforeExistingEntry(pack pack.Pack) (bool, error) {
	found := false
	insideRequire := false
	err := raw.flatMap(func(line string) []string {
		if found {
			return []string{line}
		}
		if insideRequire {
			if topLevelKeyRegexp.MatchString(line) {
				insideRequire = false
			}

			if matches := requireEntryRegexp.FindStringSubmatch(line); len(matches) == 4 {
				found = true
				indent := len(matches[1])
				return []string{formatRequire(indent, pack), line}
			}

		} else {
			if requireKeyRegxp.MatchString(line) {
				insideRequire = true
			}
		}

		return []string{line}
	})
	return found, err
}

func formatRequire(indent int, pack pack.Pack) string {
	return fmt.Sprintf("%v%v: %v",
		strings.Repeat(" ", indent),
		pack.ImportUrl(),
		pack.Version())
}
