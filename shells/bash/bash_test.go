package bash

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kr/pty"

	"github.com/bgpat/gomplete"
)

func TestFormat(t *testing.T) {
	shell := Shell{}
	reply := gomplete.Reply{
		"foo": "foo",
		"bar": "bar",
		"baz": "baz",
	}
	actual := shell.FormatReply(reply)
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
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	binfile := filepath.Join(dir, "example")
	excmd := exec.Command("go", "build", "-o", binfile, "../../examples/bash")
	if err := excmd.Run(); err != nil {
		t.Error(err)
	}

	compfile := filepath.Join(dir, "example.completion")
	shell := Shell{Name: "example"}
	script := shell.Script(binfile + " -completion --")
	if err := ioutil.WriteFile(compfile, []byte(script), 0644); err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "--noprofile", "--norc", "-o", "errexit")
	cmd.Env = []string{"PATH=" + dir}
	tty, err := pty.Start(cmd)
	fmt.Fprintf(tty, "source %q\n", compfile)
	if err := writeAndWait(ctx, tty, "example "); err != nil {
		t.Error(err)
	}
	tty.WriteString("\t")
	if err := waitString(ctx, tty, "foo"); err != nil {
		t.Error(err)
	}
	tty.WriteString("\t")
	if err := waitString(ctx, tty, "bar"); err != nil {
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
	text := ""
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			buf := make([]byte, 1024)
			n, err := tty.Read(buf)
			if err != nil {
				return err
			}
			text += string(buf[:n])
			if strings.LastIndex(text, s) >= 0 {
				return nil
			}
		}
	}
}
