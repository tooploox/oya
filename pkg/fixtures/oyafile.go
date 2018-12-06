package fixtures

import (
	"path/filepath"

	"github.com/bilus/oya/pkg/oyafile"
)

func Oyafile(dirPath string, kvs ...string) *oyafile.Oyafile {
	o, err := oyafile.New(filepath.Join(dirPath, oyafile.DefaultName), filepath.Join(dirPath, "oya/vendor"))
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(kvs); i = i + 2 {
		task := kvs[i]
		script := kvs[i+1]
		o.Tasks[task] = oyafile.ScriptedTask{
			Name:   task,
			Script: oyafile.Script(script),
		}
	}
	return o
}
