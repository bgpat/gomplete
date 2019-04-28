// +build !bash,!fish

package e2e

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kr/pty"
)

func TestZsh(t *testing.T) {
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	binfile := filepath.Join(dir, "examples")
	excmd := exec.Command("go", "build", "-o", binfile, "../../examples")
	if err := excmd.Run(); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	cmd := exec.CommandContext(ctx, "zsh", "--no-rcs", "-o", "errexit")
	cmd.Env = []string{"PATH=" + dir}
	tty, err := pty.Start(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tty.WriteString("autoload -U compinit && compinit\n"); err != nil {
		t.Fatal(err)
	}
	if _, err := tty.WriteString("source <(examples -completion=zsh)\n"); err != nil {
		t.Fatal(err)
	}
	if err := writeAndWait(ctx, tty, "examples "); err != nil {
		t.Fatal(err)
	}
	if _, err := tty.WriteString("\t"); err != nil {
		t.Fatal(err)
	}
	if _, err := waitString(ctx, tty, "foo"); err != nil {
		t.Fatal(err)
	}
	if _, err := tty.WriteString("\t"); err != nil {
		t.Fatal(err)
	}
	if _, err := waitString(ctx, tty, "bar"); err != nil {
		t.Fatal(err)
	}
	if _, err := tty.WriteString("\t\t"); err != nil {
		t.Fatal(err)
	}
	reply, err := waitString(ctx, tty, "examples foo bar")
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
