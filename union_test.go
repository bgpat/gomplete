package gomplete

import "testing"

func TestUnionComplete(t *testing.T) {
	for _, tc := range []testCase{
		{
			description: "empty",
			args:        NewArgs([]string{""}),
			completion:  &Union{},
			errorFormat: "The reply must be empty.",
		},
		{
			description: "match all",
			args:        NewArgs([]string{""}),
			completion: &Union{
				&Command{Name: "foo", Description: "foo"},
				&Command{Name: "bar", Description: "bar"},
				&Command{Name: "baz", Description: "baz"},
			},
			expect: Reply{
				"foo": "foo",
				"bar": "bar",
				"baz": "baz",
			},
			errorFormat: "The reply must be the union of all replies.",
		},
		{
			description: "match partial",
			args:        NewArgs([]string{"ba"}),
			completion: &Union{
				&Command{Name: "foo", Description: "foo"},
				&Command{Name: "bar", Description: "bar"},
				&Command{Name: "baz", Description: "baz"},
			},
			expect: Reply{
				"bar": "bar",
				"baz": "baz",
			},
			errorFormat: "The reply must be the union of matched replies.",
		},
		{
			description: "not match",
			args:        NewArgs([]string{"hoge"}),
			completion: &Union{
				&Command{Name: "foo", Description: "foo"},
				&Command{Name: "bar", Description: "bar"},
				&Command{Name: "baz", Description: "baz"},
			},
			errorFormat: "The reply must be empty.",
		},
		{
			description: "nil",
			args:        NewArgs([]string{""}),
			completion:  (*Union)(nil),
			errorFormat: "The reply must be empty.",
		},
	} {
		tc.run(t)
	}
}
