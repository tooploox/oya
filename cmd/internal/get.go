package internal

import (
	"io"
	"strings"

	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/project"
	"github.com/bilus/oya/pkg/semver"
	"github.com/pkg/errors"
)

func Get(workDir, uri string, update bool, stdout, stderr io.Writer) error {
	repoUri, versionStr, err := parseUri(uri)
	if err != nil {
		return wrapErr(err, uri)
	}
	library, err := pack.OpenLibrary(repoUri)

	prj, err := project.Detect(workDir)
	if err != nil {
		return wrapErr(err, uri)
	}

	var pack pack.Pack
	if len(versionStr) == 0 {
		pack, err = library.LatestVersion()
		if err != nil {
			return wrapErr(err, uri)
		}
		if !update {
			installed, err := prj.IsInstalled(pack)
			if err != nil {
				return wrapErr(err, uri)
			}
			if installed {
				return nil
			}
		}
	} else {
		if update {
			return errors.Errorf("Cannot request a specific pack version and use the -u (--update) flag at the same time")
		}
		version, err := semver.Parse(versionStr)
		if err != nil {
			return wrapErr(err, uri)
		}

		pack, err = library.Version(version)
		if err != nil {
			return wrapErr(err, uri)
		}
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
