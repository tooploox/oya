package pack

import (
	"fmt"
	"path/filepath"

	"github.com/go-distributed/gog/log"
	getter "github.com/hashicorp/go-getter"
	"github.com/pkg/errors"
)

type Pack interface {
	Vendor(vendorDir string) error
}

type GitPack struct {
	repoUri string
	ref     string
	relPath string
}

func NewFromUri(uri, ref string) (Pack, error) {
	return &GitPack{
		repoUri: uri,
		ref:     ref,
		relPath: uri,
	}, nil
}

func (p *GitPack) Vendor(vendorDir string) error {
	fullPath := filepath.Join(vendorDir, p.relPath)
	log.Debugf("Getting %q into %q", p.src(), fullPath)
	err := getter.GetAny(fullPath, p.src())
	if err != nil {
		return errors.Wrapf(err, "error vendoring pack %v", p.repoUri)
	}
	return nil
}

func (p *GitPack) src() string {
	if len(p.ref) > 0 {
		return fmt.Sprintf("%v?ref=%v", p.repoUri, p.ref)

	}
	return p.repoUri
}
