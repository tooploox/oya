package fixtures

import (
	"path/filepath"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/raw"
	"github.com/bilus/oya/pkg/task"
)

func Oyafile(dirPath string, kvs ...string) *oyafile.Oyafile {
	o, err := oyafile.New(filepath.Join(dirPath, raw.DefaultName), filepath.Join(dirPath, "oya/vendor"))
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(kvs); i = i + 2 {
		taskName := kvs[i]
		script := kvs[i+1]
		o.Tasks.AddTask(task.Name(taskName), task.Script{
			Script: script,
		})
	}
	return o
}
