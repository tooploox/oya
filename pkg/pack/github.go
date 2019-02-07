package pack

import (
	"path/filepath"

	"github.com/bilus/oya/pkg/semver"
	"github.com/pkg/errors"
)

// GithubPack represents a specific version of an Oya pack.
type GithubPack struct {
	library *GithubLibrary
	version semver.Version
}

func (p *GithubPack) Vendor(vendorDir string) error {
	fullPath := filepath.Join(vendorDir, p.ImportPath())
	err := p.library.Install(p.version, fullPath)
	if err != nil {
		return errors.Wrapf(err, "error vendoring pack %v", p.ImportPath())
	}
	return nil
}

func (p *GithubPack) Version() semver.Version {
	return p.version
}

func (p *GithubPack) ImportPath() string {
	return p.library.ImportPath()
}
