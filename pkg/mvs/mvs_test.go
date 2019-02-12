package mvs_test

import (
	"sort"
	"testing"

	"github.com/bilus/oya/pkg/mvs"
	"github.com/bilus/oya/pkg/pack"
	tu "github.com/bilus/oya/testutil"
)

type MockReqs map[string][]pack.Pack

func (r *MockReqs) Reqs(pack pack.Pack) ([]pack.Pack, error) {
	reqs, ok := (*r)[mvs.Hash(pack)]
	if !ok {
		return nil, nil
	} else {
		return reqs, nil
	}
}

var (
	pack1_100 = tu.MustMakeMockPack("pack1", "v1.0.0")
	pack1_110 = tu.MustMakeMockPack("pack1", "v1.1.0")
	pack1_111 = tu.MustMakeMockPack("pack1", "v1.1.1")
	pack2_100 = tu.MustMakeMockPack("pack2", "v1.0.0")
	pack2_110 = tu.MustMakeMockPack("pack2", "v1.1.0")
	pack2_111 = tu.MustMakeMockPack("pack2", "v1.1.1")
	pack3_100 = tu.MustMakeMockPack("pack3", "v1.0.0")
	pack3_110 = tu.MustMakeMockPack("pack3", "v1.1.0")
	pack3_111 = tu.MustMakeMockPack("pack3", "v1.1.1")
)

func TestNoDependencies(t *testing.T) {
	reqs := &MockReqs{}
	list, err := mvs.List(nil, reqs)
	tu.AssertNoErr(t, err, "Error creating dependency list")
	tu.AssertEqualMsg(t, 0, len(list), "Incorrect dependency list length")
}

func TestOneDependency(t *testing.T) {
	reqs := &MockReqs{}
	list, err := mvs.List([]pack.Pack{pack1_100}, reqs)
	tu.AssertNoErr(t, err, "Error creating dependency list")
	tu.AssertEqualMsg(t, 1, len(list), "Incorrect dependency list length")
	assertSamePacks(t, []pack.Pack{pack1_100}, list)
}

func TestTwoLevels(t *testing.T) {
	reqs := &MockReqs{
		mvs.Hash(pack1_100): []pack.Pack{pack2_100, pack3_100},
		mvs.Hash(pack2_100): []pack.Pack{pack3_100},
	}
	list, err := mvs.List([]pack.Pack{pack1_100, pack2_100}, reqs)
	tu.AssertNoErr(t, err, "Error creating dependency list")
	tu.AssertEqualMsg(t, 3, len(list), "Incorrect dependency list length")
	assertSamePacks(t, []pack.Pack{pack1_100, pack2_100, pack3_100}, list)
}

func TestCycles(t *testing.T) {
	reqs := &MockReqs{
		mvs.Hash(pack1_100): []pack.Pack{pack2_100, pack3_100},
		mvs.Hash(pack2_100): []pack.Pack{pack3_100, pack1_100},
	}
	list, err := mvs.List([]pack.Pack{pack1_100, pack2_100}, reqs)
	tu.AssertNoErr(t, err, "Error creating dependency list")
	tu.AssertEqualMsg(t, 3, len(list), "Incorrect dependency list length")
	assertSamePacks(t, []pack.Pack{pack1_100, pack2_100, pack3_100}, list)
}

func TestPreferHigherVersions(t *testing.T) {
	reqs := &MockReqs{
		mvs.Hash(pack1_100): []pack.Pack{pack2_110, pack3_110},
		mvs.Hash(pack2_100): []pack.Pack{pack3_111},
	}
	list, err := mvs.List([]pack.Pack{pack1_100, pack2_100, pack3_100}, reqs)
	tu.AssertNoErr(t, err, "Error creating dependency list")
	tu.AssertEqualMsg(t, 3, len(list), "Incorrect dependency list length")
	assertSamePacks(t, []pack.Pack{pack1_100, pack2_110, pack3_111}, list)
}

func sortPacks(packs []pack.Pack) {
	sort.Slice(packs, func(i, j int) bool {
		return mvs.Hash(packs[i]) < mvs.Hash(packs[j])
	})
}

func assertSamePacks(t *testing.T, lhs, rhs []pack.Pack) {
	sortPacks(lhs)
	sortPacks(rhs)
	tu.AssertObjectsEqual(t, lhs, rhs)
}
