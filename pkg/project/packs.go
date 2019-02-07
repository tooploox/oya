package project

import (
	"path/filepath"

	"github.com/bilus/oya/pkg/pack"
)

func (p Project) Require(pack pack.Pack) error {
	raw, err := p.rootRawOyafile()
	if err != nil {
		return err
	}

	return raw.AddRequire(pack)
}

func (p Project) Vendor(pack pack.Pack) error {
	return pack.Vendor(filepath.Join(p.RootDir, VendorDir))
}

func (p Project) InstallPacks() error {
	o, err := p.rootOyafile()
	if err != nil {
		return err
	}
	for _, pack := range o.Require {
		err := p.Vendor(pack)
		if err != nil {
			return err
		}
	}

	return nil
}
