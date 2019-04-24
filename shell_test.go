package gomplete

import (
	"io"
	"reflect"
	"sort"
	"strconv"
	"testing"

	"github.com/pkg/errors"
)

type testShell struct{}

func (s *testShell) FormatReply(reply Reply, w io.Writer) error {
	str := strconv.Itoa(len(reply))
	_, err := io.WriteString(w, str)
	return errors.WithStack(err)
}

func (s *testShell) Script(cmdline string) string {
	return cmdline
}

func TestRegisterShell(t *testing.T) {
	t.Run("register", func(t *testing.T) {
		unregisterAllShells()
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
		unregisterAllShells()
		if err := RegisterShell("awesome", &testShell{}); err != nil {
			t.Error(err)
		}
		if err := RegisterShell("awesome", &testShell{}); err == nil {
			t.Error("should returns an error")
		}
	})
}

func TestFormatReply(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		unregisterAllShells()
		if err := RegisterShell("test", &testShell{}); err != nil {
			t.Error(err)
		}
		for i, reply := range []Reply{
			nil,
			{"": ""},
			{"": "", "a": ""},
			{"": "", "a": "", "b": ""},
		} {
			expect := strconv.Itoa(i)
			t.Run(expect, func(t *testing.T) {
				actual, err := FormatReply("test", reply)
				if err != nil {
					t.Error(err)
				}
				if expect != actual {
					t.Errorf("expect %v, but actual %v", expect, actual)
				}
			})
		}
	})
	t.Run("unknown shell", func(t *testing.T) {
		unregisterAllShells()
		if _, err := FormatReply("test", Reply{}); err == nil {
			t.Error("should returns an error")
		}
	})
}

func unregisterAllShells() {
	shellsMu.Lock()
	shells = make(map[string]Shell)
	shellsMu.Unlock()
}
