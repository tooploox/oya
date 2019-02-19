package mvs

import (
	"fmt"
	"sync"

	"github.com/bilus/oya/pkg/mvs/internal"
	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/types"
)

type Reqs interface {
	Reqs(pack pack.Pack) ([]pack.Pack, error)
}

func max(p1, p2 pack.Pack) pack.Pack {
	if p1.Version().LessThan(p2.Version()) {
		return p2
	} else {
		return p1
	}
}

func Hash(pack pack.Pack) string {
	return fmt.Sprintf("%v@%v", pack.ImportPath(), pack.Version())
}

type Job struct {
	pack.Pack
}

func (j Job) Payload() interface{} {
	return j.Pack
}

func (j Job) ID() interface{} {
	return Hash(j.Pack)
}

// List creates a list of requirements based on initial list of required packs, taking inter-pack requirements into account.
func List(required []pack.Pack, reqs Reqs) ([]pack.Pack, error) {
	mtx := sync.Mutex{}
	latest := make(map[types.ImportPath]pack.Pack)
	queue := internal.Work{}
	for _, r := range required {
		queue.Add(Job{r})
	}
	var firstErr error
	queue.Do(10,
		func(job internal.Job) {
			if firstErr != nil {
				return
			}
			mtx.Lock()
			crnt, ok := job.Payload().(pack.Pack)
			if !ok {
				mtx.Unlock()
				panic("Internal error: expected pack.Pack passed to work queue")
			}
			if l, ok := latest[crnt.ImportPath()]; !ok || Hash(max(l, crnt)) != Hash(l) {
				latest[crnt.ImportPath()] = crnt
			}
			mtx.Unlock()

			reqs, err := reqs.Reqs(crnt)
			if err != nil {
				firstErr = err
				return
			}

			for _, req := range reqs {
				queue.Add(Job{req})
			}
		})

	if firstErr != nil {
		return nil, firstErr
	}

	packs := make([]pack.Pack, 0)
	for _, pack := range latest {
		packs = append(packs, pack)
	}
	return packs, nil
}
