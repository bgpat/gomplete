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
	want        Reply
	errorFormat string
}

func (tc *testCase) run(t *testing.T) {
	t.Run(tc.description, func(t *testing.T) {
		got := tc.completion.Complete(context.Background(), tc.args)
		if tc.want == nil {
			if len(got) > 0 {
				t.Errorf(tc.errorFormat+" want: %#v, got: %#v", tc.want, got)
			}
		} else {
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf(tc.errorFormat+" want: %#v, got: %#v", tc.want, got)
			}
		}
	})
}
