package raw

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tooploox/oya/pkg/semver"
	"github.com/tooploox/oya/pkg/types"
)

type Pack interface {
	ImportPath() types.ImportPath
	Version() semver.Version
}

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
func (raw *Oyafile) AddRequire(pack Pack) error {
	return raw.addRequire(pack)
}

// addRequire adds a Require: entry for a pack using the following algorithm:
// 1. Look for and update an existing entry for the path.
// 2. Look for ANY pack under Require:; if found, insert the new entry beneath it.
// 3. Look for Require: key (we know it's empty), insert the new entry inside it.
// 4. Look for Project: key, insert the new entry beneath it (under Require:).
// 5. Fail because Oyafile has no Project: so we shouldn't be trying to add a require to it.
// The method stops if any of the steps succeeds.
// NOTE: It does not modify the Oyafile on disk.
func (raw *Oyafile) addRequire(pack Pack) error {
	if found, err := raw.updateExistingEntry(pack); err != nil || found {
		return err // nil if found
	}
	if found, err := raw.insertBeforeExistingEntry(pack); err != nil || found {
		return err // nil if found
	}
	if found, err := raw.insertAfter(requireKeyRegxp, formatRequire(defaultIndent, pack)); err != nil || found {
		return err // nil if found
	}

	found, err := raw.insertAfter(projectRegexp, "Require:", formatRequire(defaultIndent, pack))
	if err != nil {
		return err
	}

	if !found {
		return ErrNotRootOyafile{Path: raw.Path}
	}
	return nil
}

func (raw *Oyafile) updateExistingEntry(pack Pack) (bool, error) {
	return raw.replaceAllWhen(
		func(line string) bool {
			if matches := requireEntryRegexp.FindStringSubmatch(line); len(matches) == 4 {
				return types.ImportPath(matches[2]) == pack.ImportPath()
			}
			return false

		}, []string{formatRequire(0, pack)}...)
}

func (raw *Oyafile) insertBeforeExistingEntry(pack Pack) (bool, error) {
	return raw.insertBeforeWithin("Require", requireEntryRegexp, formatRequire(0, pack))
}

func formatRequire(indent int, pack Pack) string {
	return fmt.Sprintf("%v%v: %v",
		strings.Repeat(" ", indent),
		pack.ImportPath(),
		pack.Version())
}
