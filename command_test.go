package gomplete

import (
	"context"
	"reflect"
	"testing"
)

func TestCommandComplete(t *testing.T) {
	ctx := context.Background()
	for desc, testcase := range map[string]struct {
		args   *Args
		comp   Command
		expect Reply
		ef     string
	}{
		"no arg": {
			args: NewArgs([]string{""}),
			comp: Command{
				"foo": "foo",
				"bar": "bar",
				"baz": "baz",
			},
			expect: Reply{
				"foo": "foo",
				"bar": "bar",
				"baz": "baz",
			},
			ef: "What same as the completion candidates must be returned.",
		},
		"partial arg": {
			args: NewArgs([]string{"ba"}),
			comp: Command{
				"foo": "foo",
				"bar": "bar",
				"baz": "baz",
			},
			expect: Reply{
				"bar": "bar",
				"baz": "baz",
			},
			ef: "The reply must be the all of what contains the prefix.",
		},
		"not match": {
			args: NewArgs([]string{"hoge"}),
			comp: Command{
				"foo": "foo",
				"bar": "bar",
				"baz": "baz",
			},
			ef: "The reply must be empty.",
		},
		"not last arg": {
			args: NewArgs([]string{"foo", "bar"}),
			comp: Command{
				"foo": "foo",
				"bar": "bar",
				"baz": "baz",
			},
			ef: "The reply must be empty.",
		},
	} {
		t.Run(desc, func(t *testing.T) {
			actual := testcase.comp.Complete(ctx, testcase.args)
			if testcase.expect == nil {
				if len(actual) > 0 {
					t.Errorf(testcase.ef+" expect: %#v, actual: %#v", testcase.expect, actual)
				}
			} else {
				if !reflect.DeepEqual(actual, testcase.expect) {
					t.Errorf(testcase.ef+" expect: %#v, actual: %#v", testcase.expect, actual)
				}
			}
		})
	}
}
