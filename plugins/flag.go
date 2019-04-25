package plugins

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"

	"github.com/bgpat/gomplete"
	_ "github.com/bgpat/gomplete/shells/bash"
)

// Flag is a implementation for flag.Value.
type Flag struct {
	Completion gomplete.Completion
	Config     gomplete.ShellConfig
}

// String returns empty string.
// This is a dummy function.
func (f *Flag) String() string {
	return ""
}

// Set aborts parsing flags and runs the shell completion.
func (f *Flag) Set(name string) error {
	var args *gomplete.Args
	for i, arg := range os.Args {
		if arg == "--" {
			args = gomplete.NewArgs(os.Args[i+1:])
			break
		}
	}
	if args == nil {
		return errors.WithStack(f.outputScript(name))
	}
	fmt.Fprintf(os.Stderr, "%#v\n", args)

	shell, err := gomplete.NewShell(name, f.Config)
	if err != nil {
		return errors.WithStack(err)
	}
	reply := f.Completion.Complete(context.Background(), args)
	return errors.WithStack(shell.FormatReply(reply, os.Stdout))
}

func (f *Flag) outputScript(name string) error {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Printf("usage: source <(%s)\n", strings.Join(os.Args, " "))
		return nil
	}
	shell, err := gomplete.NewShell(name, f.Config)
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(shell.OutputScript(os.Stdout))
}
