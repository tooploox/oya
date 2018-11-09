package get

import (
	"io"
	"strings"

	"github.com/bilus/oya/pkg/pack"
	"github.com/pkg/errors"
)

func Get(vendorDir, uri string, stdout, stderr io.Writer) error {
	repoUri, ref, err := parseUri(uri)
	if err != nil {
		return errors.Wrapf(err, "error getting p %v", uri)
	}
	p, err := pack.NewFromUri(repoUri, ref)
	if err != nil {
		return errors.Wrapf(err, "error getting p %v", uri)
	}
	err = p.Vendor(vendorDir)
	if err != nil {
		return errors.Wrapf(err, "error getting p %v", uri)
	}
	return nil
}

func parseUri(uri string) (string, string, error) {
	parts := strings.Split(uri, "@")
	switch len(parts) {
	case 1:
		return parts[0], "", nil
	case 2:
		return parts[0], parts[1], nil
	default:
		return "", "", errors.Errorf("unsupported package uri: %v", uri)
	}
}
