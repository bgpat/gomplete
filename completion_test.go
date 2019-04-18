package gomplete

import (
	"context"
	"reflect"
	"testing"
)

type testCase struct {
	description string
	args        *Args
	completion  Completion
	expect      Reply
	errorFormat string
}

func (tc *testCase) run(t *testing.T) {
	t.Run(tc.description, func(t *testing.T) {
		actual := tc.completion.Complete(context.Background(), tc.args)
		if tc.expect == nil {
			if len(actual) > 0 {
				t.Errorf(tc.errorFormat+" expect: %#v, actual: %#v", tc.expect, actual)
			}
		} else {
			if !reflect.DeepEqual(actual, tc.expect) {
				t.Errorf(tc.errorFormat+" expect: %#v, actual: %#v", tc.expect, actual)
			}
		}
	})
}
