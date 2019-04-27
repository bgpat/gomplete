package bash

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bgpat/gomplete"
)

func TestRegisterShell(t *testing.T) {
	_, err := gomplete.NewShell(&gomplete.ShellConfig{ShellName: "bash"})
	if err != nil {
		t.Error(err)
	}
}

func TestNewShell(t *testing.T) {
	shell, err := NewShell(&gomplete.ShellConfig{})
	if err != nil {
		t.Error(err)
	}
	if shell == nil {
		t.Error("shell is nil")
	}
}

func TestArgs(t *testing.T) {
	testcase := []string{"command", "foo", "bar", "baz"}
	shell := Shell{
		&gomplete.ShellConfig{
			Args: testcase,
		},
	}
	args := shell.Args()
	if args == nil {
		t.Error("args is nil")
	}
	for i, expect := range testcase[1:] {
		expect := expect
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := args.Current()
			if actual != expect {
				t.Errorf("expect %q, but actual %q", expect, actual)
			}
			args = args.Next()
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
		t.Error(err)
	}
	expect := "bar\nbaz\nfoo"
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
				t.Error(shell)
			}
			var buf bytes.Buffer
			if err := shell.OutputScript(&buf); err != nil {
				t.Error(err)
			}
			cupaloy.SnapshotT(t, buf.String())
		})
	}
}
