package bash

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bgpat/gomplete"
)

func TestRegisterShell(t *testing.T) {
	_, err := gomplete.NewShell(&gomplete.ShellConfig{ShellName: "bash"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewShell(t *testing.T) {
	os.Args = []string{"/path/to/command", "-completion", "bash", "--", "hoge", "fuga", "piyo"}
	for name, testcase := range map[string]struct {
		want *Shell
		env  map[string]string
	}{
		"no env": {
			want: &Shell{current: 2},
		},
		"cword": {
			want: &Shell{current: 1},
			env:  map[string]string{"COMP_CWORD": "1"},
		},
		"cword error": {
			env: map[string]string{"COMP_CWORD": "NaN"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range testcase.env {
				os.Setenv(k, v)
			}
			cfg := gomplete.NewShellConfig("bash")
			if testcase.want != nil {
				testcase.want.ShellConfig = cfg
			}
			got, err := NewShell(cfg)
			if testcase.want == nil {
				if err == nil {
					t.Error("must return an error")
				}
			} else if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(testcase.want, got) {
				t.Errorf("want %#v, got %#v", testcase.want, got)
			}
		})
	}
}

func TestArgs(t *testing.T) {
	for name, testcase := range map[string]struct {
		want    *gomplete.Args
		args    []string
		current int
	}{
		"right": {
			want:    gomplete.NewArgs([]string{"command", "foo", "bar", "baz"}).Next(),
			args:    []string{"command", "foo", "bar", "baz"},
			current: 3,
		},
		"middle": {
			want:    gomplete.NewArgs([]string{"command", "foo", "bar"}).Next(),
			args:    []string{"command", "foo", "bar", "baz"},
			current: 2,
		},
	} {
		t.Run(name, func(t *testing.T) {
			shell := Shell{
				ShellConfig: &gomplete.ShellConfig{Args: testcase.args},
				current:     testcase.current,
			}
			got := shell.Args()
			if got == nil {
				t.Error("args is nil")
			}
			if !reflect.DeepEqual(testcase.want, got) {
				t.Errorf("want %#v, but got %#v", testcase.want, got)
			}
		})
	}
}

func TestFormatReply(t *testing.T) {
	shell := Shell{}
	reply := gomplete.Reply{
		"foo": "FOO",
		"bar": "BAR",
		"baz": "BAZ",
	}
	var buf bytes.Buffer
	if err := shell.FormatReply(reply, &buf); err != nil {
		t.Fatal(err)
	}
	want := "bar\nbaz\nfoo"
	got := buf.String()
	if got != want {
		t.Errorf("output is wrong. want %q, but got %q", want, got)
	}
}

func TestOutputScript(t *testing.T) {
	for _, cfg := range []gomplete.ShellConfig{
		{
			CommandName:     "simple",
			CompleteCommand: []string{"simple", "-completion"},
		},
		{
			CommandName:     "kebab-case",
			CompleteCommand: []string{"kebab-case", "-completion"},
		},
	} {
		cfg := cfg
		t.Run(cfg.CommandName, func(t *testing.T) {
			shell, err := newShell(&cfg)
			if err != nil {
				t.Fatal(err)
			}
			var buf bytes.Buffer
			if err := shell.OutputScript(&buf); err != nil {
				t.Fatal(err)
			}
			cupaloy.SnapshotT(t, buf.String())
		})
	}
}

func TestUsage(t *testing.T) {
	shell := Shell{}
	want := "source <(cmd -completion=bash)"
	got := shell.Usage("cmd -completion=bash")
	if got != want {
		t.Errorf("want %q, but got %q", want, got)
	}
}
