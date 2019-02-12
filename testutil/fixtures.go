package testutil

import (
	"path/filepath"
	"testing"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/project"
	"github.com/bilus/oya/pkg/semver"
	"github.com/bilus/oya/pkg/types"
	"github.com/pkg/errors"
)

func MustListOyafiles(t *testing.T, rootDir string) []*oyafile.Oyafile {
	t.Helper()
	project, err := project.Detect(rootDir, filepath.Join(rootDir, ".packs"))
	AssertNoErr(t, err, "Error detecting project")
	oyafiles, err := project.Oyafiles()
	AssertNoErr(t, err, "Error listing Oyafiles")
	AssertTrue(t, len(oyafiles) > 0, "No Oyafiles found")
	return oyafiles
}

func MustListOyafilesSubdir(t *testing.T, rootDir, subDir string) []*oyafile.Oyafile {
	t.Helper()
	project, err := project.Detect(rootDir, filepath.Join(rootDir, ".packs"))
	AssertNoErr(t, err, "Error detecting project")
	oyafiles, err := project.List(subDir)
	AssertNoErr(t, err, "Error listing Oyafiles")
	AssertTrue(t, len(oyafiles) > 0, "No Oyafiles found")
	return oyafiles
}

func MustLoadOyafile(t *testing.T, dir, rootDir string) *oyafile.Oyafile {
	t.Helper()
	o, found, err := oyafile.LoadFromDir(dir, rootDir)
	AssertNoErr(t, err, "Error loading root Oyafile")
	AssertTrue(t, found, "Root Oyafile not found")
	return o
}

type mockPack struct {
	importPath types.ImportPath
	version    semver.Version
}

func (p mockPack) Version() semver.Version {
	return p.version
}

func (p mockPack) ImportPath() types.ImportPath {
	return p.importPath
}

func (p mockPack) Install(installPath string) error {
	return errors.Errorf("mockPack#Install is not implemented")
}

func (p mockPack) IsInstalled(installPath string) (bool, error) {
	return false, errors.Errorf("mockPack#IsInstalled is not implemented")
}

func (p mockPack) InstallPath(installPath string) string {
	panic(errors.Errorf("mockPack#InstallPath is not implemented"))
}

func MustMakeMockPack(importPath string, version string) pack.Pack {
	return mockPack{
		importPath: types.ImportPath(importPath),
		version:    semver.MustParse(version),
	}
}
