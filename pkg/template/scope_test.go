package template_test

import (
	"encoding/gob"
	"fmt"
	"os"
	"testing"

	"github.com/tooploox/oya/pkg/template"
	tu "github.com/tooploox/oya/testutil"
)

func identity(scope template.Scope) template.Scope {
	return scope
}

func mutation(scope template.Scope) template.Scope {
	scope["bar"] = "baz"
	return scope
}

func merge(scope template.Scope) template.Scope {
	return scope.Merge(template.Scope{"bar": "baz"})
}

func deepcopy(dst, src interface{}) error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(w)
	err = enc.Encode(src)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(r)
	return dec.Decode(dst)
}

func TestScope_UpdateScopeAt(t *testing.T) {
	testCases := []struct {
		desc          string
		scope         template.Scope
		path          string
		f             func(template.Scope) template.Scope
		expectedScope template.Scope
	}{
		{
			desc:          "empty path, return original",
			scope:         template.Scope{},
			path:          "",
			f:             identity,
			expectedScope: template.Scope{},
		},
		{
			desc:  "empty path, make an in-place update",
			scope: template.Scope{},
			path:  "",
			f:     mutation,
			expectedScope: template.Scope{
				"bar": "baz",
			},
		},
		{
			desc:  "empty path, return updated copy",
			scope: template.Scope{},
			path:  "",
			f:     merge,
			expectedScope: template.Scope{
				"bar": "baz",
			},
		},
		{
			desc:  "non-existent path, return original",
			scope: template.Scope{},
			path:  "foo",
			f:     identity,
			expectedScope: template.Scope{
				"foo": template.Scope{},
			},
		},
		{
			desc:  "non-existent path, make an in-place update",
			scope: template.Scope{},
			path:  "foo",
			f:     mutation,
			expectedScope: template.Scope{
				"foo": template.Scope{
					"bar": "baz",
				},
			},
		},
		{
			desc:  "non-existent path, return updated copy",
			scope: template.Scope{},
			path:  "foo",
			f:     merge,
			expectedScope: template.Scope{
				"foo": template.Scope{
					"bar": "baz",
				},
			},
		},
		{
			desc:  "non-existent deep path, return original",
			scope: template.Scope{},
			path:  "foo.xxx.yyy",
			f:     identity,
			expectedScope: template.Scope{
				"foo": template.Scope{
					"xxx": template.Scope{
						"yyy": template.Scope{},
					},
				},
			},
		},
		{
			desc:  "non-existent path, make an in-place update",
			scope: template.Scope{},
			path:  "foo.xxx.yyy",
			f:     mutation,
			expectedScope: template.Scope{
				"foo": template.Scope{
					"xxx": template.Scope{
						"yyy": template.Scope{
							"bar": "baz",
						},
					},
				},
			},
		},
		{
			desc:  "non-existent path, return updated copy",
			scope: template.Scope{},
			path:  "foo.xxx.yyy",
			f:     merge,
			expectedScope: template.Scope{
				"foo": template.Scope{
					"xxx": template.Scope{
						"yyy": template.Scope{
							"bar": "baz",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		var scope template.Scope
		err := deepcopy(&scope, tc.scope) // Because f can mutate it.
		tu.AssertNoErr(t, err, "deepcopy failed")
		err = scope.UpdateScopeAt(tc.path, tc.f)
		tu.AssertNoErr(t, err, "UpdateScopeAt failed in test case %q", tc.desc)
		tu.AssertObjectsEqualMsg(t, tc.expectedScope, scope, "In test case %q", tc.desc)
	}
}

func ExampleScope_Flat() {
	scope := template.Scope{
		"foo": map[interface{}]interface{}{
			"bar": "baz",
			"qux": []interface{}{
				"1", "2", 3,
			},
			"abc": map[string]interface{}{
				"123": true,
			},
		},
	}
	flattened := scope.Flat()
	fmt.Println("foo.bar:", flattened["foo.bar"])
	_, ok := flattened["foo"]
	fmt.Println("foo exists?", ok)
	fmt.Println("foo.qux.0:", flattened["foo.qux.0"])
	fmt.Println("foo.qux.1:", flattened["foo.qux.1"])
	fmt.Println("foo.qux.2:", flattened["foo.qux.2"])
	fmt.Println("foo.abc.123:", flattened["foo.abc.123"])
	// Output:
	// foo.bar: baz
	// foo exists? false
	// foo.qux.0: 1
	// foo.qux.1: 2
	// foo.qux.2: 3
	// foo.abc.123: true
}

func TestMerge(t *testing.T) {
	testCases := []struct {
		desc               string
		lhs, rhs, expected template.Scope
	}{
		{
			desc:     "empty scopes",
			lhs:      template.Scope{},
			rhs:      template.Scope{},
			expected: template.Scope{},
		},
		{
			desc:     "no overlapping keys",
			lhs:      template.Scope{"foo": "bar"},
			rhs:      template.Scope{"baz": "qux"},
			expected: template.Scope{"foo": "bar", "baz": "qux"},
		},
		{
			desc:     "overlapping keys",
			lhs:      template.Scope{"foo": "xxx"},
			rhs:      template.Scope{"baz": "qux", "foo": "bar"},
			expected: template.Scope{"foo": "bar", "baz": "qux"},
		},
		{
			desc:     "deep merge",
			lhs:      template.Scope{"foo": map[interface{}]interface{}{"bar": "xxxxx", "baz": "orange"}},
			rhs:      template.Scope{"foo": map[interface{}]interface{}{"bar": "apple", "qux": "peach"}},
			expected: template.Scope{"foo": map[interface{}]interface{}{"bar": "apple", "baz": "orange", "qux": "peach"}},
		},
	}

	for _, tc := range testCases {
		actual := tc.lhs.Merge(tc.rhs)
		tu.AssertObjectsEqual(t, tc.expected, actual)
	}
}
