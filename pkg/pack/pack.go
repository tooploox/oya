package pack

import (
	"github.com/bilus/oya/pkg/semver"
	"github.com/bilus/oya/pkg/types"
)

type Pack interface {
	Vendor(vendorDir string) error
	IsVendored(vendorDir string) (bool, error)
	Version() semver.Version
	ImportPath() types.ImportPath
}
