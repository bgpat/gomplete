// +build integration,!bash,!zsh

package e2e

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kr/pty"
)

func TestFish(t *testing.T) {
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

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "fish")
	cmd.Dir = dir
	cmd.Env = []string{
		"PATH=" + dir + ":" + os.Getenv("PATH"),
		"TERM=xterm",
	}
	tty, err := pty.Start(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if err := pty.Setsize(tty, &pty.Winsize{Rows: 20, Cols: 80}); err != nil {
		t.Fatal(err)
	}
	if _, err := tty.WriteString("source (examples -completion=fish | psub)\n"); err != nil {
		t.Fatal(err)
	}
	if err := writeAndWait(ctx, tty, "examples "); err != nil {
		t.Fatal(err)
	}
	if _, err := tty.WriteString("\t"); err != nil {
		t.Fatal(err)
	}
	reply, err := waitString(ctx, tty, "foo")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("reply: %q", reply)
	if _, err := tty.WriteString("\t"); err != nil {
		t.Fatal(err)
	}
	reply, err = waitString(ctx, tty, "bar")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("reply: %q", reply)
	if _, err := tty.WriteString("\t\t"); err != nil {
		t.Fatal(err)
	}
	reply, err = waitString(ctx, tty, "hoge")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("reply: %q", reply)
	for _, arg := range []string{"hoge", "fuga", "piyo"} {
		if !strings.Contains(reply, arg) {
			t.Errorf("3rd sub-commands must include %q", arg)
		}
	}
}
