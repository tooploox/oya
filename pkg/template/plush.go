package template

import (
	"html"
	"io"
	"strings"
	"sync"

	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plush/token"
	"github.com/pkg/errors"
)

type plushTemplate struct {
	impl *plush.Template
}

var once sync.Once

// parsePlush parses a plush template in the source string.
func parsePlush(source string, delimiters Delimiters) (Template, error) {
	once.Do(prepHelpers)

	if err := token.SetTemplatingDelimiters(delimiters.Start, delimiters.End); err != nil {
		return plushTemplate{}, err
	}

	kt, err := plush.Parse(source)
	if err != nil {
		return nil, err
	}
	return plushTemplate{impl: kt}, nil
}

// Render writes the result of rendering the plush template using the provided scope to the IO writer.
func (t plushTemplate) Render(out io.Writer, scope Scope) error {
	result, err := t.RenderString(scope)
	if err != nil {
		return err
	}
	result = html.UnescapeString(result)
	_, err = out.Write([]byte(result)) // Makes copy of result.
	return err
}

// Render returns the result of rendering the plush template using the provided scope.
func (t plushTemplate) RenderString(scope Scope) (string, error) {
	context := plush.NewContext()
	for k, v := range scope {
		ks, ok := k.(string)
		if !ok {
			return "", errors.Errorf("Internal error: Expected all scope keys to be strings, unexpected: %v", k)
		}
		context.Set(ks, v)
	}
	return t.impl.Exec(context)
}

func prepHelpers() {
	// Support the following plush helpers:
	whitelist := []string{
		"markdown", // Markdown
		"len",      // Len
		"range",    // Range
		"between",  // Between
		"until",    // Until
		"inspect",  // Inspect
	}
	helpers, _ := plush.NewHelperMap()
	for k, v := range plush.Helpers.Helpers() {
		if contains(whitelist, k) {
			helpers.Helpers()[strings.Title(k)] = v
		}
	}
	plush.Helpers = helpers
}

func contains(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}
