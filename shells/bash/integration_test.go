// +build integration,!zsh,!fish

package bash_test

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kr/pty"

	"github.com/bgpat/gomplete/test"
)

func TestShell(t *testing.T) {
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	}()

	binfile := filepath.Join(dir, "examples")
	excmd := exec.Command("go", "build", "-o", binfile, "../../examples")
	if err := excmd.Run(); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), test.Timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "--noprofile", "--norc", "-o", "errexit")
	cmd.Env = []string{"PATH=" + dir}
	tty, err := pty.Start(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tty.WriteString("source <(examples -completion=bash)\n"); err != nil {
		t.Fatal(err)
	}
	test.WriteAndWait(ctx, t, tty, "examples ")
	if _, err := tty.WriteString("\t"); err != nil {
		t.Fatal(err)
	}
	test.WaitString(ctx, t, tty, "foo")
	if _, err := tty.WriteString("\t"); err != nil {
		t.Fatal(err)
	}
	test.WaitString(ctx, t, tty, "bar")
	if _, err := tty.WriteString("\t\t"); err != nil {
		t.Fatal(err)
	}
	reply := test.WaitString(ctx, t, tty, "examples foo bar")
	if err != nil {
		t.Fatal(err)
	}
	for _, arg := range []string{"hoge", "fuga", "piyo"} {
		if !strings.Contains(reply, arg) {
			t.Errorf("3rd sub-commands must include %q", arg)
		}
	}
}
