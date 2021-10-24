package raw

import (
	"errors"
	"fmt"
	"regexp"
)

var exposeStmt = "Expose: %s"
var exposeRegexp = regexp.MustCompile("(?m)^Expose:.+$")

func (raw *Oyafile) Expose(alias string) error {
	if raw.isAlreadyExposed() {
		return errors.New("an import is already exposed")
	}
	return raw.expose(alias)
}

func (raw *Oyafile) expose(alias string) error {
	stmt := fmt.Sprintf(exposeStmt, alias)
	if updated, err := raw.insertAfter(projectRegexp, stmt); err != nil || updated {
		return err // nil if updated
	}
	return raw.prepend(
		stmt,
	)
}

func (raw *Oyafile) isAlreadyExposed() bool {
	return raw.matches(exposeRegexp)
}
