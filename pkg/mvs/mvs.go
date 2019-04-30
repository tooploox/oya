package mvs

import (
	"fmt"
	"sync"

	"github.com/tooploox/oya/pkg/errors"
	"github.com/tooploox/oya/pkg/mvs/internal"
	"github.com/tooploox/oya/pkg/pack"
	"github.com/tooploox/oya/pkg/types"
)

type ErrResolvingReqs struct {
}

func (e ErrResolvingReqs) Error() string {
	return "problem getting requirements"
}

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
	Trace []errors.Location
	Pack  pack.Pack
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
		queue.Add(Job{nil, r})
	}
	var firstErr error
	queue.Do(10,
		func(j internal.Job) {
			if firstErr != nil {
				return
			}
			job, ok := j.(Job)
			if !ok {
				panic("Internal error: expected pack.Pack passed to work queue")
			}
			crnt := job.Pack
			trace := duplicate(job.Trace)
			mtx.Lock()
			if l, ok := latest[crnt.ImportPath()]; !ok || Hash(max(l, crnt)) != Hash(l) {
				latest[crnt.ImportPath()] = crnt
			}
			mtx.Unlock()

			reqs, err := reqs.Reqs(crnt)
			if err != nil {
				mtx.Lock()
				firstErr = errors.Wrap(
					err,
					ErrResolvingReqs{},
					trace...,
				)
				mtx.Unlock()
				return
			}

			for _, req := range reqs {
				queue.Add(Job{append(trace, toLocation(crnt)), req})
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

func duplicate(trace []errors.Location) []errors.Location {
	dup := make([]errors.Location, len(trace))
	for i, t := range trace {
		dup[i] = t
	}
	return dup
}

func toLocation(pack pack.Pack) errors.Location {
	importPath := pack.ImportPath()
	return errors.Location{
		Name:        importPath.String(),
		VerboseName: fmt.Sprintf("required by %v", importPath),
	}
}
