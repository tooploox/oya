package oyafile

import (
	"fmt"
	"io"
	"unicode"

	"github.com/Masterminds/sprig"
	"github.com/bilus/oya/pkg/task"
	"github.com/bilus/oya/pkg/template"
)

var sprigFunctions = upcaseFuncNames(sprig.GenericFuncMap())

func (o *Oyafile) addBuiltIns() error {
	o.Values = o.Values.Merge(o.defaultValues())
	return nil
}

func (o *Oyafile) defaultValues() template.Scope {
	scope := template.Scope{
		"BasePath": o.Dir,
		"OyaCmd":   o.OyaCmd,
	}

	// Import sprig functions (http://masterminds.github.io/sprig/).
	for name, f := range sprigFunctions {
		_, exists := scope[name]
		if exists {
			panic(fmt.Sprintf("INTERNAL: Conflicting sprig function name: %q", name))
		}
		scope[name] = f
	}

	return scope
}

func upcaseFuncNames(funcs map[string]interface{}) map[string]interface{} {
	upcased := make(map[string]interface{})
	for name, f := range funcs {
		// We know we can cast the first byte to rune, these are function names.
		upcasedName := string(unicode.ToUpper(rune(name[0]))) + name[1:]
		upcased[upcasedName] = f
	}
	return upcased
}

// bindTasks returns a map of functions allowing invoking other tasks via $Tasks.xyz().
// It makes invokable only tasks defined in the same Oyafile, stripping away any aliases, so the tasks are accessible names exactly as they appear in a given Oyafile.
func (o *Oyafile) bindTasks(taskName task.Name, t task.Task, stdout, stderr io.Writer) (map[string]func() string, error) {
	tasks := make(map[string]func() string)

	importAlias, _ := taskName.Split()

	o.Tasks.ForEach(func(tn task.Name, _ task.Task, _ task.Meta) error {
		alias, baseName := tn.Split()
		if alias == importAlias {
			tasks[baseName] = func() string {
				return fmt.Sprintf("%s run %s\n", o.OyaCmd, tn)
			}
		}
		return nil
	})

	return tasks, nil
}

func (o *Oyafile) bindRender(taskName task.Name, stdout, stderr io.Writer) (func(string) string, error) {
	alias, _ := taskName.Split()
	return func(templatePath string) string {
		return fmt.Sprintf("%s render -a %q %s\n", o.OyaCmd, alias, templatePath)
	}, nil
}
