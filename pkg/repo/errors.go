package repo

import (
	"fmt"

	"github.com/tooploox/oya/pkg/semver"
	"github.com/tooploox/oya/pkg/types"
)

// ErrNotGithub indicates that the import path doesn't start with github.com.
type ErrNotGithub struct {
	ImportPath types.ImportPath
}

func (err ErrNotGithub) Error() string {
	return fmt.Sprintf("incorrect Github.com import path %q; expected to start with \"github.com/<user>/<repository>\"", err.ImportPath)
}

// ErrNoTaggedVersions indicates there are no available remote versions of the pack.
type ErrNoTaggedVersions struct {
	ImportPath types.ImportPath
}

func (err ErrNoTaggedVersions) Error() string {
	return fmt.Sprintf("no available remote versions for import path %q", err.ImportPath)
}

// ErrClone indacates problem during repository clone (doesn't exist or private)
type ErrClone struct {
	RepoUrl string
}

func (err ErrClone) Error() string {
	return fmt.Sprintf("failed to clone repository %v", err.RepoUrl)
}

// ErrCheckout indicates problem during repository checkout for ex. unknown version
type ErrCheckout struct {
	ImportPath    types.ImportPath
	ImportVersion semver.Version
}

func (err ErrCheckout) Error() string {
	return fmt.Sprintf("failed to get pack %v@%v", err.ImportPath, err.ImportVersion)
}

// ErrNoRootOyafile indicates that pack's root Oyafile is missing.
type ErrNoRootOyafile struct {
	ImportPath types.ImportPath
	Version    semver.Version
}

func (err ErrNoRootOyafile) Error() string {
	return fmt.Sprintf("missing top-level Oyafile for pack %v version %v", err.ImportPath, err.Version)
}
