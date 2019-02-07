package raw

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
)

var importKey = "Import:"
var projectKey = "Project:"
var uriVal = "  %s: %s"
var importRegexp = regexp.MustCompile("(?m)^" + importKey + "$")
var projectRegexp = regexp.MustCompile("^" + projectKey)

func (raw *Oyafile) AddImport(alias, uri string) error {
	if gotIt := raw.isAlreadyImported(uri); gotIt {
		return errors.Errorf("Pack already imported: %v", uri)
	}

	err := raw.addImport(alias, uri)
	if err != nil {
		return err
	}

	return raw.write()
}

func (raw *Oyafile) addImport(alias, uri string) error {
	uriStr := fmt.Sprintf(uriVal, alias, uri)
	if updated, err := raw.insertAfter(importRegexp, uriStr); err != nil || updated {
		return err // nil if updated
	}
	if updated, err := raw.insertAfter(projectRegexp, importKey, uriStr); err != nil || updated {
		return err // nil if updated
	}

	return raw.prepend(
		importKey,
		uriStr,
	)
}

func (raw *Oyafile) isAlreadyImported(uri string) bool {
	// BUG(bilus): This is slightly brittle because it will match any line that ends with the uri. It might be a good idea to add a raw.(*Oyafile).matchesWithin function along the lines of how insertBeforeWithin is written.
	rx := regexp.MustCompile("(?m)" + uri + "$")
	return raw.matches(rx)
}
