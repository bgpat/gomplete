package fish

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bgpat/gomplete"
)

func TestRegisterShell(t *testing.T) {
	_, err := gomplete.NewShell(&gomplete.ShellConfig{ShellName: "fish"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewShell(t *testing.T) {
	os.Args = []string{"/path/to/command", "-completion", "fish", "--", "hoge", "fuga", "piyo"}
	for name, testcase := range map[string]struct {
		expect *Shell
		env    map[string]string
	}{
		"no env": {
			expect: &Shell{current: 3},
		},
		"current": {
			expect: &Shell{current: 2},
			env:    map[string]string{"cword": "2"},
		},
		"current error": {
			env: map[string]string{"cword": "NaN"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range testcase.env {
				os.Setenv(k, v)
			}
			cfg := gomplete.NewShellConfig("fish")
			if testcase.expect != nil {
				testcase.expect.ShellConfig = cfg
			}
			actual, err := NewShell(cfg)
			if testcase.expect == nil {
				if err == nil {
					t.Error("must return an error")
				}
			} else if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(testcase.expect, actual) {
				t.Errorf("expect %#v, actual %#v", testcase.expect, actual)
			}
		})
	}
}

func TestArgs(t *testing.T) {
	for name, testcase := range map[string]struct {
		expect  *gomplete.Args
		args    []string
		current int
	}{
		"right": {
			expect:  gomplete.NewArgs([]string{"command", "foo", "bar", "baz"}).Next(),
			args:    []string{"command", "foo", "bar", "baz"},
			current: 4,
		},
		"middle": {
			expect:  gomplete.NewArgs([]string{"command", "foo", "bar"}).Next(),
			args:    []string{"command", "foo", "bar", "baz"},
			current: 3,
		},
	} {
		t.Run(name, func(t *testing.T) {
			shell := Shell{
				ShellConfig: &gomplete.ShellConfig{Args: testcase.args},
				current:     testcase.current,
			}
			actual := shell.Args()
			if actual == nil {
				t.Error("args is nil")
			}
			if !reflect.DeepEqual(testcase.expect, actual) {
				t.Errorf("expect %#v, but actual %#v", testcase.expect, actual)
			}
		})
	}
}

func TestFormatReply(t *testing.T) {
	shell := Shell{}
	reply := gomplete.Reply{
		"foo": "",
		"bar": "BAR",
		"baz": "BAZ",
	}
	var buf bytes.Buffer
	if err := shell.FormatReply(reply, &buf); err != nil {
		t.Fatal(err)
	}
	expect := "bar\tBAR\nbaz\tBAZ\nfoo"
	actual := buf.String()
	if actual != expect {
		t.Errorf("output is wrong. expect %q, but actual %q", expect, actual)
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
