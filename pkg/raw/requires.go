package raw

import (
	"fmt"
	"regexp"

	"github.com/bilus/oya/pkg/pack"
)

type ErrNotRootOyafile struct {
	Path string
}

func (e ErrNotRootOyafile) Error() string {
	return fmt.Sprintf("missing Project: directive in %v; not root Oyafile?", e.Path)
}

var requireKeyRegxp = regexp.MustCompile("^Require:")
var requireEntryRegexp = regexp.MustCompile("^\\s*([^:]+)\\:\\s*([^ #]+)")

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
	if found, err := raw.insertAfter(requireKeyRegxp, formatRequire(pack)); err != nil || found {
		return err // nil if found
	}

	found, err := raw.insertAfter(projectRegexp, "Require:", formatRequire(pack))
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

func formatRequire(pack pack.Pack) string {
	return fmt.Sprintf("  %v: %v", pack.ImportUrl(), pack.Version())
}
