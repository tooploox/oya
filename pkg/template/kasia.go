package template

import (
	"io"

	kasia "github.com/ziutek/kasia.go"
)

type kasiaTemplate struct {
	impl *kasia.Template
}

// parseKasia parses a kasia template in the source string.
func parseKasia(source string) (Template, error) {
	kt, err := kasia.Parse(source)
	if err != nil {
		return nil, err
	}
	kt.Strict = true
	kt.EscapeFunc = nil
	return kasiaTemplate{impl: kt}, nil
}

// Render writes the result of rendering the kasia template using the provided values to the IO writer.
func (t kasiaTemplate) Render(out io.Writer, values ...interface{}) error {
	return t.impl.Run(out, values...)
}

// Render returns the result of rendering the kasia template using the provided values.
func (t kasiaTemplate) RenderString(values ...interface{}) (string, error) {
	return t.impl.RenderString(values...)
}
