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

type mockRepo struct {
	importPath types.ImportPath
}

func (p mockRepo) Install(version semver.Version, installDir string) error {
	return errors.Errorf("mockRepo#Install is not implemented")
}

func (p mockRepo) IsInstalled(version semver.Version, installDir string) (bool, error) {
	return false, errors.Errorf("mockRepo#IsInstalled is not implemented")
}

func (p mockRepo) InstallPath(version semver.Version, installDir string) string {
	panic(errors.Errorf("mockRepo#InstallPath is not implemented"))
}

func (p mockRepo) ImportPath() types.ImportPath {
	return p.importPath
}

func mustMakeMockRepo(importPath string) pack.Repo {
	return mockRepo{
		importPath: types.ImportPath(importPath),
	}
}

func MustMakeMockPack(importPath string, version string) pack.Pack {
	pack, err := pack.New(mustMakeMockRepo(importPath), semver.MustParse(version))
	if err != nil {
		panic(err)
	}
	return pack
}
