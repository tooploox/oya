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
	basePath   string
	packPath   string
	importPath string
}

// OpenLibrary opens a library containing all versions of a single Oya pack.
func OpenLibrary(importPath string) (*GithubLibrary, error) {
	if !strings.HasPrefix(importPath, "github.com/") {
		return nil, ErrNotGithub{ImportPath: importPath}
	}
	repoUri, basePath, packPath, err := parseImportPath(importPath)
	if err != nil {
		return nil, err
	}
	return &GithubLibrary{
		repoUri:    repoUri,
		packPath:   packPath,
		basePath:   basePath,
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

// Version returns the specified version of the pack.
// NOTE: It doesn't check if it's available remotely. This may change.
// It is used when loading Oyafiles so we probably shouldn't do it or use a different function there.
func (l *GithubLibrary) Version(version semver.Version) (*GithubPack, error) {
	// BUG(bilus): Check if version exists?
	return &GithubPack{
		library: l,
		version: version,
	}, nil
}

// ImportPath returns the pack's import path, e.g. github.com/tooploox/oya-packs/docker.
func (l *GithubLibrary) ImportPath() string {
	return l.importPath
}

// Install downloads & copies the specified version of the path to the output directory,
// preserving its import path.
// For example, for /home/bilus/.oya output directory and import path github.com/bilus/foo,
// the pack will be extracted to /home/bilus/.oya/github.com/bilus/foo.
func (l *GithubLibrary) Install(version semver.Version, outputDir string) error {
	path := filepath.Join(outputDir, l.basePath)
	log.Printf("Getting %q version %v into %q (git tag: %v)", l.ImportPath(), version, path, l.makeRef(version))
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
		if outside, err := l.isOutsidePack(f.Name); outside || err != nil {
			return err // May be nil if outside true.
		}
		targetPath := filepath.Join(path, f.Name)
		return copyFile(f, targetPath)
	})
}

func (l *GithubLibrary) IsInstalled(version semver.Version, outputDir string) (bool, error) {
	fullPath := filepath.Join(outputDir, l.basePath, l.packPath)
	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func copyFile(f *object.File, targetPath string) error {
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
	if err != nil {
		return err
	}
	err = writer.Sync()
	if err != nil {
		return err
	}
	mode, err := f.Mode.ToOSFileMode()
	if err != nil {
		return err
	}
	err = os.Chmod(targetPath, mode)
	if err != nil {
		return err
	}
	return err
}

func parseImportPath(importPath string) (string, string, string, error) {
	parts := strings.Split(importPath, "/")
	if len(parts) < 3 {
		return "", "", "", ErrNotGithub{ImportPath: importPath}
	}
	basePath := strings.Join(parts[0:3], "/")
	repoUri := fmt.Sprintf("https://%v.git", basePath)
	packPath := strings.Join(parts[3:], "/")
	return repoUri, basePath, packPath, nil
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

func (l *GithubLibrary) isOutsidePack(relPath string) (bool, error) {
	r, err := filepath.Rel(l.packPath, relPath)
	if err != nil {
		return false, err
	}
	return strings.Contains(r, ".."), nil
}
