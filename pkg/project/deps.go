package project

import (
	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/types"
)

type Deps interface {
	Load(importPath types.ImportPath) (*oyafile.Oyafile, bool, error)
	Find(importPath types.ImportPath) (pack.Pack, bool, error)
	ForEach(f func(pack.Pack) error) error
	BuildList() error
}
