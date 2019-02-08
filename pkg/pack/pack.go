package pack

import "github.com/bilus/oya/pkg/semver"

type Pack interface {
	Vendor(vendorDir string) error
	Version() semver.Version
	ImportPath() string
}
