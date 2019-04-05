package project

import (
	"path/filepath"

	"github.com/tooploox/oya/pkg/oyafile"
	"github.com/tooploox/oya/pkg/raw"
)

func (p *Project) oyafileIn(dir string) (*oyafile.Oyafile, bool, error) {
	normalizedDir := filepath.Clean(dir)
	o, found := p.oyafileCache[normalizedDir]
	if found {
		return o, true, nil
	}
	o, found, err := oyafile.LoadFromDir(dir, p.RootDir)
	if err != nil {
		return nil, false, err
	}
	if found {
		p.oyafileCache[normalizedDir] = o
		return o, true, nil
	}

	return nil, false, nil
}

func (p *Project) rawOyafileIn(dir string) (*raw.Oyafile, bool, error) {
	// IMPORTANT: Call invalidateOyafileCache after patching raw Oyafiles
	// obtained using this method!
	normalizedDir := filepath.Clean(dir)
	o, found := p.rawOyafileCache[normalizedDir]
	if found {
		return o, true, nil
	}
	o, found, err := raw.LoadFromDir(dir, p.RootDir)
	if err != nil {
		return nil, false, err
	}
	if found {
		p.rawOyafileCache[normalizedDir] = o
		return o, true, nil
	}
	return nil, false, nil
}

func (p *Project) invalidateOyafileCache(dir string) {
	delete(p.oyafileCache, dir)
	delete(p.rawOyafileCache, dir)
}

func (p *Project) Oyafile(oyafilePath string) (*oyafile.Oyafile, bool, error) {
	// BUG(bilus): Uncached (but used only by Render).
	return oyafile.Load(oyafilePath, p.RootDir)
}
