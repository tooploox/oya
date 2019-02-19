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
		result[k] = v
	}
	return result
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
// as a Scope (see ParseScope), the function signals ErrScopeMergeConflict.
func (scope Scope) UpdateScopeAt(path string, f func(Scope) Scope) error {
	var pathArr []string
	if len(path) > 0 {
		pathArr = strings.Split(path, ".")
	}
	targetScope, err := scope.resolveScope(pathArr, true)
	if err != nil {
		// Ignore the reason.
		return ErrScopeMergeConflict{Path: path}
	}
	targetScope.Replace(f(targetScope))
	return nil
}

// GetScopeAt returns scope at the specified path. If path doesn't exist or points
// to a value that cannot be interpreted as a Scope (see ParseScope), the function
// signals an error.
func (scope Scope) GetScopeAt(path string) (Scope, error) {
	pathArr := strings.Split(path, ".")
	return scope.resolveScope(pathArr, false)
}

func (scope Scope) resolveScope(path []string, create bool) (Scope, error) {
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
		return nil, errors.Errorf("Unsupported value under %q", scopeName)
	}
	return subScope.resolveScope(path[1:], create)
}
