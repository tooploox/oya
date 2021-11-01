package task_test

import (
	"testing"

	"github.com/tooploox/oya/pkg/task"
	"github.com/tooploox/oya/pkg/types"
	tu "github.com/tooploox/oya/testutil"
)

func TestTable_Expose(t *testing.T) {
	tt := task.NewTable()

	globalName := task.Name("task")
	aliasedName := globalName.Aliased("alias")
	tt.AddTask(aliasedName, task.MockTask{})

	tt.Expose(types.Alias("alias"))

	_, ok := tt.LookupTask(aliasedName)
	tu.AssertTrue(t, ok, "Aliased task not found")
	_, ok = tt.LookupTask(globalName)
	tu.AssertTrue(t, ok, "Global task not found")
}
