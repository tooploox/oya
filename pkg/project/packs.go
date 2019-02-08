package project

import (
	"path/filepath"

	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/types"
)

const VendorDir = ".oya/vendor"

func (p Project) Require(pack pack.Pack) error {
	raw, err := p.rootRawOyafile()
	if err != nil {
		return err
	}

	return raw.AddRequire(pack)
}

func (p Project) Vendor(pack pack.Pack) error {
	return pack.Vendor(p.vendorDir())
}

func (p Project) IsVendored(pack pack.Pack) (bool, error) {
	return pack.IsVendored(p.vendorDir())
}

func (p Project) InstallPacks() error {
	o, err := p.rootOyafile()
	if err != nil {
		return err
	}
	for _, pack := range o.Require {
		installed, err := p.IsInstalled(pack)
		if err != nil {
			return err
		}
		if installed {

			continue
		}
		err = p.Install(pack)
		if err != nil {
			return err
		}
	}

	return nil
}

// Currently vendoring is the only supported installation method but lets have these functions for clarity.

func (p Project) Install(pack pack.Pack) error {
	return p.Vendor(pack)
}

func (p Project) IsInstalled(pack pack.Pack) (bool, error) {
	return p.IsVendored(pack)
}

func (p Project) vendorDir() string {
	return filepath.Join(p.RootDir, VendorDir)
}

func (p Project) FindRequiredPack(importPath types.ImportPath) (pack.Pack, bool, error) {
	o, err := p.rootOyafile()
	if err != nil {
		return nil, false, err
	}
	for _, pack := range o.Require {
		if pack.ImportPath() == importPath {
			return pack, true, nil
		}
	}
	return nil, false, nil
}
