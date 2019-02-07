package internal

import (
	"io"
	"strings"

	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/project"
	"github.com/bilus/oya/pkg/semver"
	"github.com/pkg/errors"
)

func Get(workDir, uri string, stdout, stderr io.Writer) error {
	repoUri, versionStr, err := parseUri(uri)
	if err != nil {
		return wrapErr(err, uri)
	}
	library, err := pack.OpenLibrary(repoUri)

	var pack pack.Pack
	if len(versionStr) == 0 {
		pack, err = library.LatestVersion()
		if err != nil {
			return wrapErr(err, uri)
		}
	} else {
		version, err := semver.Parse(versionStr)
		if err != nil {
			return wrapErr(err, uri)
		}

		pack, err = library.Version(version)
		if err != nil {
			return wrapErr(err, uri)
		}
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
	return errors.Wrapf(err, "error getting pack %v", uri)
}
