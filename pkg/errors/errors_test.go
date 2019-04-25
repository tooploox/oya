package errors_test

import (
	"bytes"
	"testing"

	"github.com/tooploox/oya/pkg/errors"
	tu "github.com/tooploox/oya/testutil"
)

func TestSnippet_Format(t *testing.T) {
	cases := []struct {
		desc     string
		snippet  errors.Snippet
		line     uint
		col      uint
		expected string
	}{
		{
			desc:     "empty snippet",
			snippet:  errors.Snippet{},
			line:     1,
			col:      1,
			expected: "",
		},
		{
			desc: "one line snippet",
			snippet: errors.Snippet{
				LineOffset: 0,
				Lines: []string{
					"line1",
				},
			},
			line:     1,
			col:      1,
			expected: "> 1| line1\n     ^\n",
		},
		{
			desc: "multi-line snippet",
			snippet: errors.Snippet{
				LineOffset: 0,
				Lines: []string{
					"line1",
					"line2",
					"line3",
				},
			},
			line:     3,
			col:      1,
			expected: "  2| line2\n> 3| line3\n     ^\n",
		},
		{
			desc: "two digit line number",
			snippet: errors.Snippet{
				LineOffset: 8,
				Lines: []string{
					"line9",
					"line10",
					"line11",
				},
			},
			line:     10,
			col:      2,
			expected: "   9| line9\n> 10| line10\n       ^\n",
		},
	}
	for _, tc := range cases {
		var out bytes.Buffer
		tc.snippet.Print(&out, tc.line, tc.col)
		result := out.String()
		tu.AssertEqual(t, tc.expected, result)
	}
}
