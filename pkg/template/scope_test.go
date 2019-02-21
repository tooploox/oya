package template_test

import (
	"encoding/gob"
	"os"
	"testing"

	"github.com/bilus/oya/pkg/template"
	tu "github.com/bilus/oya/testutil"
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
