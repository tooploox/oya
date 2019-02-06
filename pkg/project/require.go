package project

import (
	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/raw"
)

func (p Project) Require(pack pack.Pack) error {
	raw, found, err := raw.LoadFromDir(p.RootDir, p.RootDir)
	if err != nil {
		return err
	}
	if !found {
		return ErrNoOyafile{Path: p.RootDir}
	}

	return raw.AddRequire(pack)
}
