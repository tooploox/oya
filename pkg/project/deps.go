package project

import (
	"github.com/tooploox/oya/pkg/oyafile"
	"github.com/tooploox/oya/pkg/pack"
	"github.com/tooploox/oya/pkg/types"
)

type Deps interface {
	Load(importPath types.ImportPath) (*oyafile.Oyafile, bool, error)
	Find(importPath types.ImportPath) (pack.Pack, bool, error)
	ForEach(f func(pack.Pack) error) error
	Explode() error
}
