package fixtures

import "github.com/bilus/oya/pkg/oyafile"

func Oyafile(path string, kvs ...string) *oyafile.Oyafile {
	o := oyafile.New(path)
	for i := 0; i < len(kvs); i = i + 2 {
		hook := kvs[i]
		script := kvs[i+1]
		o.Hooks[hook] = oyafile.Hook{
			Name:   hook,
			Script: oyafile.Script(script),
		}
	}
	return o
}
