package template

import (
	"io"
	"io/ioutil"

	kasia "github.com/ziutek/kasia.go"
)

type Template interface {
	Render(out io.Writer, values interface{}) error
}

type kasiaTemplate struct {
	impl *kasia.Template
}

func Load(path string) (Template, error) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(string(source))
}

func Parse(source string) (Template, error) {
	kt, err := kasia.Parse(source)
	if err != nil {
		return nil, err
	}
	return kasiaTemplate{impl: kt}, nil
}

type emptyValues struct{}

func (t kasiaTemplate) Render(out io.Writer, values interface{}) error {
	if values == nil {
		values = emptyValues{}
	}
	return t.impl.Run(out, values)
}
