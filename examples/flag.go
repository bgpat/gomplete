package examples

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bgpat/gomplete"
	"github.com/mattn/go-isatty"
)

// Flag is a implementation for flag.Value.
type Flag struct {
	gomplete.Shell
}

// String returns empty string.
// This is a dummy function.
func (f *Flag) String() string {
	return ""
}

// Set aborts parsing flags and runs the shell completion.
func (f *Flag) Set(s string) error {
	if s != "true" {
		return nil
	}

	for i, arg := range os.Args {
		if arg == "--" {
			args := gomplete.NewArgs(os.Args[i+1:])
			reply := f.Shell.Complete(context.Background(), args)
			fmt.Print(f.Shell.Format(reply))
			return nil
		}
	}

	if isatty.IsTerminal(os.Stdout.Fd()) {
		args := strings.Join(os.Args, " ")
		fmt.Printf(`Usage:
  source <(%s)
	%s > /etc/bash_completion.d/%s
`, args, args, filepath.Base(os.Args[0]))
		return nil
	}

	args := make([]string, 0, len(os.Args)+1)
	args = append(args, os.Args...)
	args = append(args, "--")
	arg0, err := filepath.Abs(os.Args[0])
	if err != nil {
		return nil
	}
	script := f.Shell.Script(arg0 + strings.Join(args[1:], " "))
	os.Stdout.WriteString(script)
	return nil
}

// IsBoolFlag returns always true
func (f *Flag) IsBoolFlag() bool {
	return true
}
