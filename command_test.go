package gomplete

import "testing"

func TestCommandComplete(t *testing.T) {
	comp := &Command{Name: "foo", Description: "bar"}
	for _, tc := range []testCase{
		{
			description: "no arg",
			args:        NewArgs([]string{""}),
			completion:  comp,
			expect:      Reply{"foo": "bar"},
			errorFormat: "Completion must return the reply.",
		},
		{
			description: "partial arg",
			args:        NewArgs([]string{"fo"}),
			completion:  comp,
			expect:      Reply{"foo": "bar"},
			errorFormat: "The reply must be the all of what contains the prefix.",
		},
		{
			description: "not match",
			args:        NewArgs([]string{"o"}),
			completion:  comp,
			errorFormat: "The reply must be empty.",
		},
		{
			description: "not last arg",
			args:        NewArgs([]string{"foo", "bar"}),
			completion:  comp,
			errorFormat: "The reply must be empty.",
		},
		{
			description: "sub",
			args:        NewArgs([]string{"foo", ""}),
			completion: &Command{
				Name:        "foo",
				Description: "bar",
				Sub:         &Union{comp},
			},
			expect:      Reply{"foo": "bar"},
			errorFormat: "Completion must return the nested reply.",
		},
	} {
		tc.run(t)
	}
}
