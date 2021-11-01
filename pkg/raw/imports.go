package raw

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
)

var importStmt = "Import:"
var projectStmt = "Project:"
var uriVal = "  %s: %s"
var importRegexp = regexp.MustCompile("(?m)^" + importStmt + "$")
var projectRegexp = regexp.MustCompile("^" + projectStmt)

func (raw *Oyafile) AddImport(alias, uri string) error {
	if gotIt := raw.isAlreadyImported(uri); gotIt {
		return errors.Errorf("Pack already imported: %v", uri)
	}

	return raw.addImport(alias, uri)
}

func (raw *Oyafile) addImport(alias, uri string) error {
	uriStr := fmt.Sprintf(uriVal, alias, uri)
	if updated, err := raw.insertAfter(importRegexp, uriStr); err != nil || updated {
		return err // nil if updated
	}
	if updated, err := raw.insertAfter(projectRegexp, importStmt, uriStr); err != nil || updated {
		return err // nil if updated
	}
	return raw.prepend(
		importStmt,
		uriStr,
	)
}

func (raw *Oyafile) isAlreadyImported(uri string) bool {
	// BUG(bilus): This is slightly brittle because it will match any line that
	// ends with the uri. It might be a good idea to add a
	// raw.(*Oyafile).matchesWithin function along the lines of how
	// insertBeforeWithin is written.
	rx := regexp.MustCompile("(?m)" + uri + "$")
	return raw.matches(rx)
}
