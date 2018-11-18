package project

import "fmt"

type ErrNoOyafile struct {
	Path string
}

func (e ErrNoOyafile) Error() string {
	return fmt.Sprintf("no Oyafile in %v", e.Path)
}

type ErrNoOyafiles struct {
	Path string
}

func (e ErrNoOyafiles) Error() string {
	return fmt.Sprintf("no Oyafile in %v", e.Path)
}

type ErrNoProject struct {
	Path string
}

func (e ErrNoProject) Error() string {
	return fmt.Sprintf("no Oyafile project in %v or any parent directories", e.Path)
}

type ErrNoHook struct {
	Hook string
}

func (e ErrNoHook) Error() string {
	return fmt.Sprintf("missing hook %q", e.Hook)
}
