package pack

import "github.com/bilus/oya/pkg/semver"

type Pack interface {
	Vendor(vendorDir string) error
	IsVendored(vendorDir string) (bool, error)
	Version() semver.Version
	ImportPath() string
}
