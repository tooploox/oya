package oyafile

import (
	"fmt"
	"unicode"

	"github.com/Masterminds/sprig"
	"github.com/tooploox/oya/pkg/template"
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
