package task_test

import (
	"testing"

	"github.com/tooploox/oya/pkg/task"
	"github.com/tooploox/oya/pkg/types"
	tu "github.com/tooploox/oya/testutil"
)

func TestName_Split_NoAlias(t *testing.T) {
	name := task.Name("foo")
	alias, task := name.Split()
	tu.AssertEqual(t, types.Alias(""), alias)
	tu.AssertEqual(t, "foo", task)
}

func TestName_Split_WithAlias(t *testing.T) {
	name := task.Name("foo.bar")
	alias, task := name.Split()
	tu.AssertEqual(t, types.Alias("foo"), alias)
	tu.AssertEqual(t, "bar", task)
}

func TestName_Split_WithNestedAlias(t *testing.T) {
	name := task.Name("foo.bar.zoo")
	alias, task := name.Split()
	tu.AssertEqual(t, types.Alias("foo.bar"), alias)
	tu.AssertEqual(t, "zoo", task)
}

func TestName_IsAliased(t *testing.T) {
	testCases := []struct {
		name      task.Name
		alias     types.Alias
		isAliased bool
	}{
		{task.Name("name"), types.Alias("alias"), false},
		{task.Name("name").Aliased("alias"), types.Alias("alias"), true},
		{task.Name("name").Aliased("otherAlias"), types.Alias("alias"), false},
	}

	for _, tc := range testCases {
		tu.AssertEqual(t, tc.isAliased, tc.name.IsAliased(tc.alias))
	}
}

func TestName_Unaliased(t *testing.T) {
	globalName := task.Name("name")
	aliasedName := globalName.Aliased(types.Alias("alias"))
	tu.AssertEqual(t, globalName, aliasedName.Unaliased())
	tu.AssertEqual(t, globalName, globalName.Unaliased())
}
