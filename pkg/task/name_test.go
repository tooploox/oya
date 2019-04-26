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
