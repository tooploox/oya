package printers_test

import (
	"strings"
	"testing"

	"github.com/tooploox/oya/cmd/internal/printers"
	"github.com/tooploox/oya/pkg/task"
	"github.com/tooploox/oya/pkg/types"
	tu "github.com/tooploox/oya/testutil"
)

type mockWriter struct {
	bs []byte
}

func (w *mockWriter) Write(p []byte) (n int, err error) {
	w.bs = append(w.bs, p...)
	return len(p), nil
}

func (w *mockWriter) Lines() []string {
	if len(w.bs) == 0 {
		return []string{}
	}
	return strings.Split(string(w.bs), "\n")
}

func TestTaskList(t *testing.T) {
	type taskDef struct {
		name        task.Name
		meta        task.Meta
		oyafilePath string
	}
	testCases := []struct {
		desc           string
		workDir        string
		tasks          []taskDef
		expectedOutput []string
	}{
		{
			desc:           "no tasks",
			workDir:        "/project/",
			tasks:          nil,
			expectedOutput: []string{},
		},
		{
			desc:    "one global task",
			workDir: "/project/",
			tasks: []taskDef{
				{task.Name("task1"), task.Meta{}, "/project/Oyafile"},
			},
			expectedOutput: []string{
				"# in ./Oyafile",
				"oya run task1",
				"",
			},
		},
		{
			desc:    "global tasks before imported tasks",
			workDir: "/project/",
			tasks: []taskDef{
				{task.Name("task1"), task.Meta{}, "/project/Oyafile"},
				{task.Name("othertask"), task.Meta{}, "/project/Oyafile"},
				{task.Name("task2").Aliased(types.Alias("pack")), task.Meta{}, "/project/Oyafile"},
				{task.Name("task3").Aliased(types.Alias("pack")), task.Meta{}, "/project/Oyafile"},
			},
			expectedOutput: []string{
				"# in ./Oyafile",
				"oya run othertask",
				"oya run task1",
				"oya run pack.task2",
				"oya run pack.task3",
				"",
			},
		},
		{
			desc:    "top-level tasks before tasks in subdirectories",
			workDir: "/project/",
			tasks: []taskDef{
				{task.Name("task1"), task.Meta{}, "/project/Oyafile"},
				{task.Name("otherTask"), task.Meta{}, "/project/Oyafile"},
				{task.Name("aTask"), task.Meta{}, "/project/subdir/Oyafile"},
				{task.Name("task3"), task.Meta{}, "/project/subdir/Oyafile"},
			},
			expectedOutput: []string{
				"# in ./Oyafile",
				"oya run otherTask",
				"oya run task1",
				"",
				"# in ./subdir/Oyafile",
				"oya run aTask",
				"oya run task3",
				"",
			},
		},
		{
			desc:    "sort aliases and task names alphabetically",
			workDir: "/project/",
			tasks: []taskDef{
				{task.Name("yyy"), task.Meta{}, "/project/Oyafile"},
				{task.Name("zzz"), task.Meta{}, "/project/Oyafile"},
				{task.Name("aaa").Aliased(types.Alias("BBB")), task.Meta{}, "/project/Oyafile"},
				{task.Name("ddd").Aliased(types.Alias("AAA")), task.Meta{}, "/project/Oyafile"},
				{task.Name("ccc").Aliased(types.Alias("AAA")), task.Meta{}, "/project/Oyafile"},
			},
			expectedOutput: []string{
				"# in ./Oyafile",
				"oya run yyy",
				"oya run zzz",
				"oya run AAA.ccc",
				"oya run AAA.ddd",
				"oya run BBB.aaa",
				"",
			},
		},
		{
			desc:    "task descriptions are aligned",
			workDir: "/project/",
			tasks: []taskDef{
				{task.Name("y"), task.Meta{Doc: "a description"}, "/project/Oyafile"},
				{task.Name("zzz"), task.Meta{Doc: "another description"}, "/project/Oyafile"},
				{task.Name("aaa").Aliased(types.Alias("BBB")), task.Meta{Doc: "task aa"}, "/project/Oyafile"},
				{task.Name("ddddd").Aliased(types.Alias("AAA")), task.Meta{}, "/project/Oyafile"},
				{task.Name("ccc").Aliased(types.Alias("AAA")), task.Meta{}, "/project/Oyafile"},
			},
			expectedOutput: []string{
				"# in ./Oyafile",
				"oya run y         # a description",
				"oya run zzz       # another description",
				"oya run AAA.ccc",
				"oya run AAA.ddddd",
				"oya run BBB.aaa   # task aa",
				"",
			},
		},
		{
			desc:    "exposed tasks",
			workDir: "/project/",
			tasks: []taskDef{
				{task.Name("y"), task.Meta{Doc: "a description", OriginalTaskName: task.Name("y").Aliased("somepack")}, "/project/Oyafile"},
				{task.Name("z"), task.Meta{OriginalTaskName: task.Name("z").Aliased("somepack")}, "/project/Oyafile"},
			},
			expectedOutput: []string{
				"# in ./Oyafile",
				"oya run y # a description (somepack.y)",
				"oya run z # (somepack.z)",
				"",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tl := printers.NewTaskList(tc.workDir)
			for _, task := range tc.tasks {
				tu.AssertNoErr(t, tl.AddTask(task.name, task.meta, task.oyafilePath), "Adding task failed")

			}
			w := mockWriter{}
			tl.Print(&w)
			tu.AssertObjectsEqual(t, tc.expectedOutput, w.Lines())

		})
	}
}
