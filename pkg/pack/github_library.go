package pack

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bilus/oya/pkg/semver"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// GithubLibrary represents all versions of an Oya pack stored in a git repository on Github.com.
type GithubLibrary struct {
	repoUri    string
	packPath   string
	importPath string
}

func OpenLibrary(importPath string) (*GithubLibrary, error) {
	// BUG(bilus): New can return an invalid pack (missing a version). Not good.
	if !strings.HasPrefix(importPath, "github.com/") {
		return nil, ErrNotGithub{ImportPath: importPath}
	}
	repoUri, packPath, err := parseImportPath(importPath)
	if err != nil {
		return nil, err
	}
	return &GithubLibrary{
		repoUri:    repoUri,
		packPath:   packPath,
		importPath: importPath,
	}, nil
}

// AvailableVersions returns a sorted list of remotely available pack versions.
func (l *GithubLibrary) AvailableVersions() ([]semver.Version, error) {
	versions := make([]semver.Version, 0)

	fs := memfs.New()
	storer := memory.NewStorage()
	r, err := git.Clone(storer, fs, &git.CloneOptions{
		URL: l.repoUri,
	})
	if err != nil {
		return nil, err
	}
	tags, err := r.Tags()
	if err != nil {
		return nil, err
	}
	err = tags.ForEach(
		func(t *plumbing.Reference) error {
			n := t.Name()
			if n.IsTag() {
				version, ok := l.parseRef(n.Short())
				if ok {
					versions = append(versions, version)
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	semver.Sort(versions)
	return versions, nil
}

// LatestVersion returns the latest available pack version based on tags in the remote Github repo.
func (l *GithubLibrary) LatestVersion() (*GithubPack, error) {
	versions, err := l.AvailableVersions()
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, ErrNoTaggedVersions{ImportPath: l.importPath}
	}
	latestVersion := versions[len(versions)-1]
	log.Debugf("Updating pack %q to version %v", l.importPath, latestVersion)
	return l.Version(latestVersion)
}

func (l *GithubLibrary) Version(version semver.Version) (*GithubPack, error) {
	// BUG(bilus): Check if version exists?
	return &GithubPack{
		library: l,
		version: version,
	}, nil
}

func (l *GithubLibrary) ImportPath() string {
	return l.importPath
}

func parseImportPath(importPath string) (string, string, error) {
	parts := strings.Split(importPath, "/")
	if len(parts) < 3 {
		return "", "", ErrNotGithub{ImportPath: importPath}
	}
	repoUri := fmt.Sprintf("https://%v.git", strings.Join(parts[0:3], "/"))
	packPath := strings.Join(parts[3:], "/")
	return repoUri, packPath, nil
}

func (l *GithubLibrary) parseRef(tag string) (semver.Version, bool) {
	if len(l.packPath) > 0 && strings.HasPrefix(tag, l.packPath) {
		tag = tag[len(l.packPath)+1:] // e.g. "pack1/v1.0.0" => v1.0.0
	}
	version, err := semver.Parse(tag)
	return version, err == nil
}

func (l *GithubLibrary) makeRef(version semver.Version) string {
	if len(l.packPath) > 0 {
		return fmt.Sprintf("%v/%v", l.packPath, version.String())

	} else {
		return fmt.Sprintf("%v", version.String())
	}
}

func (l *GithubLibrary) Install(version semver.Version, path string) error {
	log.Debugf("Getting %q version %v into %q (tag: %v)", l.ImportPath(), version, path, l.makeRef(version))
	fs := memfs.New()
	storer := memory.NewStorage()
	r, err := git.Clone(storer, fs, &git.CloneOptions{
		URL: l.repoUri,
	})
	if err != nil {
		return err
	}
	tree, err := r.Worktree()
	if err != nil {
		return err
	}
	err = tree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewTagReferenceName(l.makeRef(version)),
	})
	if err != nil {
		return err
	}
	ref, err := r.Head()
	if err != nil {
		return err
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return err
	}

	fIter, err := commit.Files()
	if err != nil {
		return err
	}

	return fIter.ForEach(func(f *object.File) error {
		targetPath := filepath.Join(path, f.Name)
		err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
		if err != nil {
			return err
		}
		reader, err := f.Reader()
		if err != nil {
			return err
		}
		// BUG(bilus): Copy permissions.
		writer, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, reader)
		return err
	})
}
