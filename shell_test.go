package gomplete

import (
	"bytes"
	"io"
	"reflect"
	"sort"
	"strconv"
	"testing"

	"github.com/pkg/errors"
)

type testShell struct {
	ShellConfig
}

func (s *testShell) Args() *Args {
	return NewArgs([]string{"foo", "bar", "baz"})
}

func (s *testShell) FormatReply(reply Reply, w io.Writer) error {
	str := strconv.Itoa(len(reply))
	_, err := io.WriteString(w, str)
	return errors.WithStack(err)
}

func (s *testShell) OutputScript(w io.Writer) error {
	_, err := io.WriteString(w, s.Name)
	return errors.WithStack(err)
}

func TestRegisterShell(t *testing.T) {
	t.Run("register", func(t *testing.T) {
		unregisterAllShells()
		for _, name := range []string{"foo", "bar", "baz"} {
			t.Run(name, func(t *testing.T) {
				if err := registerTestShell(name); err != nil {
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
			t.Error("must returns an error")
		}
	})
	t.Run("duplicated", func(t *testing.T) {
		unregisterAllShells()
		if err := registerTestShell("test"); err != nil {
			t.Error(err)
		}
		if err := registerTestShell("test"); err == nil {
			t.Error("must returns an error")
		}
	})
}

func TestNewShell(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		unregisterAllShells()
		if err := registerTestShell("test"); err != nil {
			t.Error(err)
		}
		shell, err := NewShell("test", ShellConfig{
			Name: "test",
		})
		if err != nil {
			t.Error(err)
		}
		if shell == nil {
			t.Error("shell is nil")
		}
		var buf bytes.Buffer
		if err := shell.OutputScript(&buf); err != nil {
			t.Error(err)
		}
		if buf.String() != "test" {
			t.Error("not match to the testcase")
		}
	})
	t.Run("unknown shell", func(t *testing.T) {
		unregisterAllShells()
		if _, err := NewShell("test", ShellConfig{}); err == nil {
			t.Error("must return nil because test shell is not registered")
		}
	})
	t.Run("failed to initialize", func(t *testing.T) {
		unregisterAllShells()
		if err := registerErrorShell("test"); err != nil {
			t.Error(err)
		}
		if _, err := NewShell("test", ShellConfig{}); err == nil {
			t.Error("must return nil")
		}
	})
}

func registerTestShell(name string) error {
	return errors.WithStack(RegisterShell(name, func(config ShellConfig) (Shell, error) {
		return &testShell{config}, nil
	}))
}

func registerErrorShell(name string) error {
	return errors.WithStack(RegisterShell(name, func(ShellConfig) (Shell, error) {
		return nil, errors.New("this constructor always return an error")
	}))
}

func unregisterAllShells() {
	shellsMu.Lock()
	shells = make(map[string]func(ShellConfig) (Shell, error))
	shellsMu.Unlock()
}