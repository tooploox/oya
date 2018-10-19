package fixtures

import (
	"path/filepath"

	"github.com/bilus/oya/pkg/oyafile"
)

func Oyafile(dirPath string, kvs ...string) *oyafile.Oyafile {
	o := oyafile.New(filepath.Join(dirPath, oyafile.DefaultName))
	for i := 0; i < len(kvs); i = i + 2 {
		hook := kvs[i]
		script := kvs[i+1]
		o.Hooks[hook] = oyafile.ScriptedHook{
			Name:   hook,
			Script: oyafile.Script(script),
		}
	}
	return o
}
