package template

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// ErrScopeMergeConflict indicates that a conflict occurred when merging two scopes
// at the specified path.
type ErrScopeMergeConflict struct {
	Path string
}

func (e ErrScopeMergeConflict) Error() string {
	return fmt.Sprintf("path %v already exists", e.Path)
}

type ErrPathEmpty struct{}

func (e ErrPathEmpty) Error() string {
	return "Path must not be empty"
}

// Scope represents a single lexical value scope. Scopes can be nested.
type Scope map[interface{}]interface{}

// ParseScope casts the value to Scope.
func ParseScope(v interface{}) (Scope, bool) {
	if scope, ok := v.(Scope); ok {
		return scope, true
	}
	if scope, ok := v.(map[interface{}]interface{}); ok {
		return scope, true
	}
	return nil, false
}

// Merge combines two scopes together. If a key appears in both scopes, the other scope wins.
func (scope Scope) Merge(other Scope) Scope {
	result := Scope{}
	for k, v := range scope {
		result[k] = v
	}
	for k, v := range other {
		existing, ok := result[k]
		if !ok {
			result[k] = v
		} else {
			result[k] = merge(existing, v)
		}
		result[k] = merge(existing, v)
	}
	return result
}

func merge(e, n interface{}) interface{} {
	es, ok := ParseScope(e)
	if !ok {
		return n // Overwrite.
	}
	ns, ok := ParseScope(n)
	if !ok {
		return n // Overwrite.
	}
	// Merge:
	return map[interface{}]interface{}(es.Merge(ns))
}

// Replace replaces contents of this scope with keys and values of the other one.
func (scope Scope) Replace(other Scope) {
	if reflect.ValueOf(scope).Pointer() == reflect.ValueOf(other).Pointer() {
		return
	}
	for k := range scope {
		delete(scope, k)
	}
	for k, v := range other {
		scope[k] = v
	}
}

// UpdateScopeAt transforms a scope pointed to by the path (e.g. "foo.bar.baz").
// It will create scopes along the way if they don't exist.
// In case the value pointed by the path already exists and cannot be interpreted
// as a Scope (see ParseScope), the function signals ErrScopeMergeConflict unless
// force is true. If that's the case, the value is overridden.
func (scope Scope) UpdateScopeAt(path string, f func(Scope) Scope, force bool) error {
	var pathArr []string
	if len(path) > 0 {
		pathArr = strings.Split(path, ".")
	}
	targetScope, err := scope.resolveScope(pathArr, true, force)
	if err != nil {
		// Ignore the reason.
		return ErrScopeMergeConflict{Path: path}
	}
	targetScope.Replace(f(targetScope))
	return nil
}

// AssocAt modifies scope by setting value at the provided path.
// It will create scopes along the way if they don't exist.
// In case the value pointed by the path already exists and cannot be interpreted
// as a Scope (see ParseScope), the function signals ErrScopeMergeConflict unless
// force is true. If that's the case, the value is overridden.
func (scope Scope) AssocAt(path string, value interface{}, force bool) error {
	if len(path) == 0 {
		return ErrPathEmpty{}
	}
	pathArr := strings.Split(path, ".")
	scopePath := strings.Join(pathArr[0:len(pathArr)-1], ".")
	key := pathArr[len(pathArr)-1]
	return scope.UpdateScopeAt(scopePath, func(s Scope) Scope {
		s[key] = value
		return s
	}, force)
}

// GetScopeAt returns scope at the specified path. If path doesn't exist or points
// to a value that cannot be interpreted as a Scope (see ParseScope), the function
// signals an error.
func (scope Scope) GetScopeAt(path string) (Scope, error) {
	pathArr := strings.Split(path, ".")
	return scope.resolveScope(pathArr, false, false)
}

// Flat returns a flattened scope.

func (scope Scope) resolveScope(path []string, create, force bool) (Scope, error) {
	if len(path) == 0 {
		return scope, nil
	}

	scopeName := path[0]
	potentialScope, ok := scope[scopeName]
	if !ok {
		if !create {
			return nil, errors.Errorf("Missing key %q", scopeName)

		}
		// Create scope along the path.
		potentialScope = Scope{}
		scope[scopeName] = potentialScope
	}
	subScope, ok := ParseScope(potentialScope)
	if !ok {
		if !force {
			return nil, errors.Errorf("Unsupported value under %q", scopeName)
		}
		subScope = Scope{}
		scope[scopeName] = subScope
	}
	return subScope.resolveScope(path[1:], create, force)
}

func (scope Scope) Flat() Scope {
	result := make(Scope)
	flatten(scope, &result, "")
	return result
}

func flatten(value interface{}, result *Scope, parent string) {
	if scope, ok := ParseScope(value); ok {
		for k, v := range scope {
			ks, ok := k.(string)
			if !ok {
				panic("Scope keys are expected to be strings!")
			}
			key := ks
			if len(parent) > 0 {
				// E.g. parent.key
				key = fmt.Sprintf("%s.%s", parent, ks)
			}
			flatten(v, result, key)
		}
	} else if xs := reflect.ValueOf(value); xs.Kind() == reflect.Slice || xs.Kind() == reflect.Array {
		for i := 0; i < xs.Len(); i++ {
			x := xs.Index(i)
			var key string
			if len(parent) > 0 {
				// E.g. parent.0
				key = fmt.Sprintf("%s.%v", parent, i)
			} else {
				key = fmt.Sprintf("%v", i)
			}
			flatten(x, result, key)
		}
	} else if m := reflect.ValueOf(value); m.Kind() == reflect.Map {
		for _, k := range m.MapKeys() {
			v := m.MapIndex(k)
			var key string
			if len(parent) > 0 {
				// E.g. parent.0
				key = fmt.Sprintf("%s.%v", parent, k)
			} else {
				key = fmt.Sprintf("%v", k)
			}
			flatten(v, result, key)
		}
	} else {
		(*result)[parent] = value
	}
}
