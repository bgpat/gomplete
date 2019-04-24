package gomplete

import (
	"reflect"
	"sort"
	"testing"
)

type testShell struct{}

func (s *testShell) FormatReply(reply Reply) string {
	return ""
}

func (s *testShell) Script(cmdline string) string {
	return cmdline
}

func TestRegisterShell(t *testing.T) {
	t.Run("register", func(t *testing.T) {
		for _, name := range []string{"foo", "bar", "baz"} {
			t.Run(name, func(t *testing.T) {
				if err := RegisterShell(name, &testShell{}); err != nil {
					t.Error(err)
				}
				if _, ok := shells[name]; !ok {
					t.Errorf("%q is not registered", name)
				}
			})
		}
		expect := []string{"bar", "baz", "foo"}
		actual := Shells()
		sort.Strings(actual)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect %#v, but actual %#v", expect, actual)
		}
	})
	t.Run("nil", func(t *testing.T) {
		if err := RegisterShell("nil", nil); err == nil {
			t.Error("should returns an error")
		}
	})
	t.Run("duplicated", func(t *testing.T) {
		if err := RegisterShell("awesome", &testShell{}); err != nil {
			t.Error(err)
		}
		if err := RegisterShell("awesome", &testShell{}); err == nil {
			t.Error("should returns an error")
		}
	})
}
