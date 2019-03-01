package repo

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/semver"
	"github.com/bilus/oya/pkg/types"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-billy.v4/memfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// GithubRepo represents all versions of an Oya pack stored in a git repository on Github.com.
type GithubRepo struct {
	repoUris   []string
	basePath   string
	packPath   string
	importPath types.ImportPath
}

// AvailableVersions returns a sorted list of remotely available pack versions.
func (l *GithubRepo) AvailableVersions() ([]semver.Version, error) {
	versions := make([]semver.Version, 0)

	r, err := l.clone()
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

func (l *GithubRepo) clone() (*git.Repository, error) {
	var lastErr error
	for _, uri := range l.repoUris {
		fs := memfs.New()
		storer := memory.NewStorage()
		repo, err := git.Clone(storer, fs, &git.CloneOptions{
			URL: uri,
		})
		if err == nil {
			return repo, nil
		}
		lastErr = err

	}
	return nil, lastErr
}

// LatestVersion returns the latest available pack version based on tags in the remote Github repo.
func (l *GithubRepo) LatestVersion() (pack.Pack, error) {
	versions, err := l.AvailableVersions()
	if err != nil {
		return pack.Pack{}, err
	}
	if len(versions) == 0 {
		return pack.Pack{}, ErrNoTaggedVersions{ImportPath: l.importPath}
	}
	latestVersion := versions[len(versions)-1]
	return l.Version(latestVersion)
}

// Version returns the specified version of the pack.
// NOTE: It doesn't check if it's available remotely. This may change.
// It is used when loading Oyafiles so we probably shouldn't do it or use a different function there.
func (l *GithubRepo) Version(version semver.Version) (pack.Pack, error) {
	// BUG(bilus): Check if version exists?
	return pack.New(l, version)
}

// ImportPath returns the pack's import path, e.g. github.com/tooploox/oya-packs/docker.
func (l *GithubRepo) ImportPath() types.ImportPath {
	return l.importPath
}

// InstallPath returns the local path for the specific pack version.
func (l *GithubRepo) InstallPath(version semver.Version, installDir string) string {
	path := filepath.Join(installDir, l.basePath, l.packPath)
	return fmt.Sprintf("%v@%v", path, version.String())
}

func (l *GithubRepo) checkout(version semver.Version) (*object.Commit, error) {
	r, err := l.clone()
	if err != nil {
		return nil, ErrCheckout{ImportPath: l.importPath, ImportVersion: version, GitMsg: err.Error()}
	}
	tree, err := r.Worktree()
	if err != nil {
		return nil, ErrCheckout{ImportPath: l.importPath, ImportVersion: version, GitMsg: err.Error()}
	}
	err = tree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewTagReferenceName(l.makeRef(version)),
	})
	if err != nil {
		return nil, ErrCheckout{ImportPath: l.importPath, ImportVersion: version, GitMsg: err.Error()}
	}
	ref, err := r.Head()
	if err != nil {
		return nil, ErrCheckout{ImportPath: l.importPath, ImportVersion: version, GitMsg: err.Error()}
	}
	return r.CommitObject(ref.Hash())
}

// Install downloads & copies the specified version of the path to the output directory,
// preserving its import path.
// For example, for /home/bilus/.oya output directory and import path github.com/bilus/foo,
// the pack will be extracted to /home/bilus/.oya/github.com/bilus/foo.
func (l *GithubRepo) Install(version semver.Version, installDir string) error {
	commit, err := l.checkout(version)
	if err != nil {
		return err
	}

	fIter, err := commit.Files()
	if err != nil {
		return err
	}

	sourceBasePath := l.packPath
	targetPath := l.InstallPath(version, installDir)
	log.Printf("Installing pack %v version %v into %q (git tag: %v)", l.ImportPath(), version, targetPath, l.makeRef(version))

	return fIter.ForEach(func(f *object.File) error {
		if outside, err := l.isOutsidePack(f.Name); outside || err != nil {
			return err // May be nil if outside true.
		}
		relPath, err := filepath.Rel(sourceBasePath, f.Name)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(targetPath, relPath)
		return copyFile(f, targetPath)
	})
}

func (l *GithubRepo) IsInstalled(version semver.Version, installDir string) (bool, error) {
	fullPath := l.InstallPath(version, installDir)
	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (l *GithubRepo) Reqs(version semver.Version) ([]pack.Pack, error) {
	// BUG(bilus): This is a slow way to get requirements for a pack.
	// It involves installing it out to a local directory.
	// But it's also the simplest one. We can optimize by using HTTP
	// access to pull in Oyafile and then parse the Require: section here.
	// It means duplicating the logic including the assumption that the requires
	// will always be stored in Oyafile, rather than a separate file along the lines
	// of go.mod.

	tempDir, err := ioutil.TempDir("", "oya")
	defer os.RemoveAll(tempDir)

	err = l.Install(version, tempDir)
	if err != nil {
		return nil, err
	}

	fullPath := l.InstallPath(version, tempDir)
	o, found, err := oyafile.LoadFromDir(fullPath, fullPath)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrNoRootOyafile{l.importPath, version}
	}

	// BUG(bilus): This doesn't take Oyafile#Replacements into account.
	// This probably doesn't matter because it's likely meaningless for
	// packs accessed remotely but we may want to revisit it.

	packs := make([]pack.Pack, len(o.Requires))
	for i, require := range o.Requires {
		repo, err := Open(require.ImportPath)
		if err != nil {
			return nil, err
		}
		pack, err := repo.Version(require.Version)
		if err != nil {
			return nil, err
		}
		packs[i] = pack
	}

	return packs, nil
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

func parseImportPath(importPath types.ImportPath) (uris []string, basePath string, packPath string, err error) {
	parts := strings.Split(string(importPath), "/")
	if len(parts) < 3 {
		return nil, "", "", ErrNotGithub{ImportPath: importPath}
	}
	basePath = strings.Join(parts[0:3], "/")
	packPath = strings.Join(parts[3:], "/")
	// Prefer https but fall back on ssh if cannot clone via https
	// as would be the case for private repositories.
	uris = []string{
		fmt.Sprintf("https://%v.git", basePath),
		fmt.Sprintf("git@%s:%s/%s.git", parts[0], parts[1], parts[2]),
	}
	return
}

func (l *GithubRepo) parseRef(tag string) (semver.Version, bool) {
	if len(l.packPath) > 0 && strings.HasPrefix(tag, l.packPath) {
		tag = tag[len(l.packPath)+1:] // e.g. "pack1/v1.0.0" => v1.0.0
	}
	version, err := semver.Parse(tag)
	return version, err == nil
}

func (l *GithubRepo) makeRef(version semver.Version) string {
	if len(l.packPath) > 0 {
		return fmt.Sprintf("%v/%v", l.packPath, version.String())

	} else {
		return fmt.Sprintf("%v", version.String())
	}
}

func (l *GithubRepo) isOutsidePack(relPath string) (bool, error) {
	r, err := filepath.Rel(l.packPath, relPath)
	if err != nil {
		return false, err
	}
	return strings.Contains(r, ".."), nil
}

func toErrClone(url string, err error) error {
	if err == transport.ErrAuthenticationRequired {
		return ErrClone{RepoUrl: url, GitMsg: "Repository not found or private"}
	}
	return ErrClone{RepoUrl: url, GitMsg: err.Error()}
}
