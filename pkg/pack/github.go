package pack

import (
	"github.com/bilus/oya/pkg/semver"
	"github.com/bilus/oya/pkg/types"
	"github.com/pkg/errors"
)

// GithubPack represents a specific version of an Oya pack.
type GithubPack struct {
	library *GithubLibrary
	version semver.Version
}

func (p *GithubPack) Install(installDir string) error {
	err := p.library.Install(p.version, installDir)
	if err != nil {
		return errors.Wrapf(err, "error installing pack %v", p.ImportPath())
	}
	return nil
}

func (p *GithubPack) IsInstalled(installDir string) (bool, error) {
	return p.library.IsInstalled(p.version, installDir)
}

func (p *GithubPack) Version() semver.Version {
	return p.version
}

func (p *GithubPack) ImportPath() types.ImportPath {
	return p.library.ImportPath()
}

func (p *GithubPack) InstallPath(installPath string) string {
	return p.library.InstallPath(p.version, installPath)
}
