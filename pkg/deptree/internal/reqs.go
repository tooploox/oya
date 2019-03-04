package internal

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	"github.com/tooploox/oya/pkg/oyafile"
	"github.com/tooploox/oya/pkg/pack"
	"github.com/tooploox/oya/pkg/raw"
	"github.com/tooploox/oya/pkg/repo"
)

type Reqs struct {
	rootDir     string
	installDirs []string
	cache       map[string][]pack.Pack
	mtx         sync.Mutex
}

func NewReqs(rootDir string, installDirs []string) *Reqs {
	return &Reqs{
		rootDir:     rootDir,
		installDirs: installDirs,
		cache:       make(map[string][]pack.Pack),
	}
}

func (r *Reqs) Reqs(pack pack.Pack) ([]pack.Pack, error) {
	reqs, found := r.cachedReqs(pack)
	if found {
		return reqs, nil
	}

	reqs, found, err := r.lookupReqs(pack)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.Errorf("Pack not found: %v", pack.ImportPath())
	}
	r.cacheReqs(pack, reqs)
	return reqs, nil
}

func (r *Reqs) lookupReqs(pack pack.Pack) ([]pack.Pack, bool, error) {
	reqs, found, err := r.localReqs(pack)
	if err != nil {
		return nil, false, err
	}
	if found {
		return reqs, true, nil
	}
	reqs, err = r.remoteReqs(pack)
	if err != nil {
		return nil, false, err
	}
	return reqs, true, nil
}

func (r *Reqs) cachedReqs(pack pack.Pack) ([]pack.Pack, bool) {
	r.mtx.Lock()
	reqs, found := r.cache[id(pack)]
	r.mtx.Unlock()
	if found {
		return reqs, true
	}
	return nil, false
}

func (r *Reqs) cacheReqs(pack pack.Pack, reqs []pack.Pack) {
	r.mtx.Lock()
	r.cache[id(pack)] = reqs
	r.mtx.Unlock()
}

func id(pack pack.Pack) string {
	return fmt.Sprintf("%v@%v", pack.ImportPath(), pack.Version())
}

func (r *Reqs) localReqs(pack pack.Pack) ([]pack.Pack, bool, error) {
	o, found, err := r.LoadLocalOyafile(pack)
	if err != nil {
		return nil, false, err
	}
	if found {
		packs, err := toPacks(o.Requires)
		if err != nil {
			return nil, false, err
		}
		return packs, true, nil
	}
	return nil, false, nil
}

func (r *Reqs) LoadLocalOyafile(pack pack.Pack) (*oyafile.Oyafile, bool, error) {
	if path, ok := pack.ReplacementPath(); ok {
		var fullPath string
		if filepath.IsAbs(path) {
			fullPath = path
		} else {
			fullPath = filepath.Join(r.rootDir, path)
		}
		o, found, err := oyafile.LoadFromDir(fullPath, r.rootDir)
		if !found {
			return nil, false, errors.Errorf("no %v found at the replacement path %v for %q", raw.DefaultName, fullPath, pack.ImportPath())
		}
		if err != nil {
			return nil, false, errors.Wrapf(err, "error resolving replacement path %v for %q", fullPath, pack.ImportPath())

		}
		return o, true, nil

	}
	for _, installDir := range r.installDirs {
		o, found, err := oyafile.LoadFromDir(pack.InstallPath(installDir), r.rootDir)
		if err != nil {
			continue
		}
		if !found {
			continue
		}
		return o, true, nil
	}
	return nil, false, nil
}

func (r *Reqs) remoteReqs(p pack.Pack) ([]pack.Pack, error) {
	l, err := repo.Open(p.ImportPath())
	if err != nil {
		return nil, err
	}
	return l.Reqs(p.Version())
}

func toPacks(references []oyafile.PackReference) ([]pack.Pack, error) {
	packs := make([]pack.Pack, len(references))
	for i, reference := range references {
		repo, err := repo.Open(reference.ImportPath)
		if err != nil {
			return nil, err
		}
		if packs[i], err = repo.Version(reference.Version); err != nil {
			return nil, err
		}
	}
	return packs, nil
}
