package semver

import (
	"fmt"
	"regexp"

	"github.com/blang/semver"
	"github.com/pkg/errors"
)

type Version semver.Version

var versionRx = regexp.MustCompile("^v([\\d\\.]+)$")

// Parse parses version string in vX.Y.Z format and returns a validated Version or error.
func Parse(v string) (Version, error) {
	if matches := versionRx.FindStringSubmatch(v); len(matches) == 2 {
		ver, err := semver.Parse(matches[1])
		return Version(ver), err
	} else {
		return Version{}, errors.Errorf("Unrecognized version syntax: %v, expected: vX.Y.Z", v)
	}
}

// MustParse is like Parse but panics if the version cannot be parsed.
func MustParse(v string) Version {
	ver, err := Parse(v)
	if err != nil {
		panic(err)
	}
	return ver
}

// String returns display representation of the version in the vX.Y.Z format.
func (ver Version) String() string {
	return fmt.Sprintf("v%v", semver.Version(ver))
}

func (ver Version) LessThan(other Version) bool {
	return semver.Version(ver).LT(semver.Version(other))
}
