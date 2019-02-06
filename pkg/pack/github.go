package pack

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/blang/semver"
	getter "github.com/hashicorp/go-getter"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// GithubPack represents an Oya pack stored in a git repository on Github.com.
type GithubPack struct {
	repoUri    string
	ref        string
	importPath string
}

func New(importPath, ref string) (Pack, error) {
	if !strings.HasPrefix(importPath, "github.com/") {
		return nil, ErrNotGithub{ImportPath: importPath}
	}
	repoUri, err := githubRepoUri(importPath)
	if err != nil {
		return nil, err
	}
	return &GithubPack{
		repoUri:    repoUri,
		ref:        ref,
		importPath: importPath,
	}, nil
}

func (p *GithubPack) Vendor(vendorDir string) error {
	fullPath := filepath.Join(vendorDir, p.importPath)
	log.Debugf("Getting %q into %q", p.src(), fullPath)
	err := getter.GetAny(fullPath, p.src())
	if err != nil {
		return errors.Wrapf(err, "error vendoring pack %v", p.importPath)
	}
	return nil
}

func (p *GithubPack) Version() string {
	return p.ref
}

func (p *GithubPack) ImportUrl() string {
	return p.importPath
}

// Update upgrades pack to the latest available version based on tags in the remote Github repo.
func (p *GithubPack) Update() error {
	versions, err := p.AvailableVersions()
	if err != nil {
		return err
	}
	if len(versions) == 0 {
		return ErrNoTaggedVersions{ImportPath: p.importPath}
	}
	p.ref = makeVersionTag(versions[len(versions)-1])
	log.Debugf("Updating pack %q to version %v", p.importPath, p.ref)
	return nil
}

// AvailableVersions returns a sorted list of remotely available pack versions.
func (p *GithubPack) AvailableVersions() ([]semver.Version, error) {
	versions := make([]semver.Version, 0)

	fs := memfs.New()
	storer := memory.NewStorage()
	r, err := git.Clone(storer, fs, &git.CloneOptions{
		URL: p.repoUri,
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
				version, ok := parseVersionTag(n.Short())
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

func parseVersionTag(tag string) (semver.Version, bool) {
	if tag[0:1] == "v" {
		version, err := semver.Make(tag[1:])
		return version, err == nil
	} else {
		return semver.Version{}, false
	}
}

func makeVersionTag(version semver.Version) string {
	return fmt.Sprintf("v%v", version.String())
}

func githubRepoUri(importPath string) (string, error) {
	return fmt.Sprintf("https://%v.git", importPath), nil
}

func (p *GithubPack) src() string {
	if len(p.ref) > 0 {
		return fmt.Sprintf("%v?ref=%v", p.importPath, p.ref)

	}
	return p.repoUri
}
