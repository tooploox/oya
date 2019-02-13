package pack

import (
	"github.com/bilus/oya/pkg/semver"
	"github.com/bilus/oya/pkg/types"
)

type Pack interface {
	Version() semver.Version
	ImportPath() types.ImportPath
	Install(installDir string) error
	IsInstalled(installDir string) (bool, error)
	InstallPath(installDir string) string
}
