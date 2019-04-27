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
	*ShellConfig
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
	_, err := io.WriteString(w, s.CommandName)
	return errors.WithStack(err)
}

func TestRegisterShell(t *testing.T) {
	t.Run("register", func(t *testing.T) {
		unregisterAllShells()
		for _, name := range []string{"foo", "bar", "baz"} {
			t.Run(name, func(t *testing.T) {
				registerTestShell(name)
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
		defer func() {
			if err := recover(); err == nil {
				t.Error("must returns an error")
			}
		}()
		RegisterShell("nil", nil)
	})
	t.Run("duplicated", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Error("must returns an error")
			}
		}()
		unregisterAllShells()
		registerTestShell("test")
		registerTestShell("test")
	})
}

func TestNewShell(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		unregisterAllShells()
		registerTestShell("test")
		shell, err := NewShell(&ShellConfig{
			CommandName: "test",
			ShellName:   "test",
		})
		if err != nil {
			t.Fatal(err)
		}
		if shell == nil {
			t.Fatal("shell is nil")
		}
		var buf bytes.Buffer
		if err := shell.OutputScript(&buf); err != nil {
			t.Fatal(err)
		}
		if buf.String() != "test" {
			t.Error("not match to the testcase")
		}
	})
	t.Run("unknown shell", func(t *testing.T) {
		unregisterAllShells()
		if _, err := NewShell(&ShellConfig{ShellName: "test"}); err == nil {
			t.Error("must return nil because test shell is not registered")
		}
	})
	t.Run("failed to initialize", func(t *testing.T) {
		unregisterAllShells()
		registerErrorShell("test")
		if _, err := NewShell(&ShellConfig{ShellName: "test"}); err == nil {
			t.Error("must return nil")
		}
	})
}

func registerTestShell(name string) {
	RegisterShell(name, func(config *ShellConfig) (Shell, error) {
		return &testShell{config}, nil
	})
}

func registerErrorShell(name string) {
	RegisterShell(name, func(*ShellConfig) (Shell, error) {
		return nil, errors.New("this constructor always return an error")
	})
}

func unregisterAllShells() {
	shellsMu.Lock()
	shells = make(map[string]func(*ShellConfig) (Shell, error))
	shellsMu.Unlock()
}
