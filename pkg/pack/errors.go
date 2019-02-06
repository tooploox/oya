package pack

import "fmt"

// ErrNotGithub indicates that the import path doesn't start with github.com.
type ErrNotGithub struct {
	ImportPath string
}

func (err ErrNotGithub) Error() string {
	return fmt.Sprintf("incorrect Github.com import path %q; expected to start with \"github.com/<user>/<repository>\"", err.ImportPath)
}

// ErrNoTaggedVersions indicates there are no available remote versions of the pack.
type ErrNoTaggedVersions struct {
	ImportPath string
}

func (err ErrNoTaggedVersions) Error() string {
	return fmt.Sprintf("no available remote versions for import path %q", err.ImportPath)
}
