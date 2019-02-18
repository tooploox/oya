package template

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type ErrScopeMergeConflict struct {
	Path string
}

func (e ErrScopeMergeConflict) Error() string {
	return fmt.Sprintf("path %v already exists", e.Path)
}

type Scope map[string]interface{}

func ParseScope(v interface{}) (Scope, bool) {
	if Scope, ok := v.(Scope); ok {
		return Scope, true
	}
	if aMap, ok := v.(map[interface{}]interface{}); ok {
		scope := make(Scope)
		for k, v := range aMap {
			name, ok := k.(string)
			if !ok {
				return nil, false
			}
			scope[name] = v
		}
		return scope, true
	}
	return nil, false
}

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

func (scope Scope) Replace(other Scope) {
	for k := range scope {
		delete(scope, k)
	}
	for k, v := range other {
		scope[k] = v
	}
}

func (scope Scope) UpdateScopeAt(path string, f func(Scope) Scope) error {
	pathArr := strings.Split(path, ".")
	targetScope, err := scope.resolveScope(pathArr, true)
	if err != nil {
		// Ignore the reason.
		return ErrScopeMergeConflict{Path: path}
	}
	targetScope.Replace(f(targetScope))
	return nil
}

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
