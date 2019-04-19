package bash

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/kr/pty"

	"github.com/bgpat/gomplete"
)

func TestComplete(t *testing.T) {
	for _, tc := range []struct {
		description string
		args        *gomplete.Args
		completion  Completion
		expect      gomplete.Reply
		errorFormat string
	}{
		{
			description: "no sub",
			args:        gomplete.NewArgs([]string{"foo"}),
			completion:  Completion{},
			errorFormat: "The reply must be emptry.",
		},
		{
			description: "no args",
			args:        gomplete.NewArgs([]string{}),
			completion:  Completion{Sub: &gomplete.Union{}},
			errorFormat: "The reply must be emptry.",
		},
		{
			description: "next",
			args:        gomplete.NewArgs([]string{"foo", "bar"}),
			completion: Completion{Sub: &gomplete.Command{
				Name:        "bar",
				Description: "bar",
			}},
			expect:      gomplete.Reply{"bar": "bar"},
			errorFormat: "The completion must return the reply of sub.",
		},
	} {
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
}

func TestFormat(t *testing.T) {
	comp := Completion{}
	reply := gomplete.Reply{
		"foo": "foo",
		"bar": "bar",
		"baz": "baz",
	}
	actual := comp.Format(reply)
	count := 0
	for _, expect := range []string{
		"foo\nbar\nbaz",
		"foo\nbaz\nbar",
		"bar\nfoo\nbaz",
		"bar\nbaz\nfoo",
		"baz\nfoo\nbar",
		"baz\nbar\nfoo",
	} {
		if actual == expect {
			count++
		}
	}
	if count != 1 {
		t.Errorf("The output is wrong. %q", actual)
	}
}

func TestScript(t *testing.T) {
	comp := Completion{Name: "awesome"}
	script := comp.Script("awesome -completion -- ")
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	compfile := filepath.Join(dir, "awesome")
	if err := ioutil.WriteFile(compfile, []byte(script), 0644); err != nil {
		t.Error(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "--noprofile", "--norc", "-o", "errexit")
	cmd.Env = []string{}
	tty, err := pty.Start(cmd)
	fmt.Fprintf(tty, "source %q\n", compfile)
	if err := writeAndWait(ctx, tty, "awesome "); err != nil {
		t.Error(err)
	}
	tty.WriteString("\t\t")
	if err := waitString(ctx, tty, "aaa abb bbb"); err != nil {
		t.Error(err)
	}
}

func writeAndWait(ctx context.Context, tty io.ReadWriter, s string) error {
	io.WriteString(tty, s)
	return waitString(ctx, tty, s)
}

func waitString(ctx context.Context, tty io.Reader, s string) error {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	suffix := ""
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			buf := make([]byte, 1024)
			n, err := tty.Read(buf)
			println(n, err)
			if err != nil {
				return err
			}
			suffix += string(buf[:n])
			if len(suffix) > len(s) {
				suffix = suffix[len(suffix)-len(s) : len(suffix)]
			}
			if suffix == s {
				return nil
			}
			fmt.Printf("%q", suffix)
		}
	}
}
