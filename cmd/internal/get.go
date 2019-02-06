package internal

import (
	"io"
	"strings"

	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/project"
	"github.com/pkg/errors"
)

func Get(workDir, uri string, stdout, stderr io.Writer) error {
	repoUri, ref, err := parseUri(uri)
	if err != nil {
		return wrapErr(err, uri)
	}
	pack, err := pack.NewFromUri(repoUri, ref)
	if err != nil {
		return wrapErr(err, uri)
	}
	prj, err := project.Detect(workDir)
	if err != nil {
		return wrapErr(err, uri)
	}
	err = prj.Vendor(pack)
	if err != nil {
		return wrapErr(err, uri)
	}
	err = prj.Require(pack)
	if err != nil {
		return wrapErr(err, uri)
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
		return "", "", errors.Errorf("unsupported pack uri: %v", uri)
	}
}

func wrapErr(err error, uri string) error {
	return errors.Wrapf(err, "error getting p %v", uri)
}
