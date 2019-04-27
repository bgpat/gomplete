package plugins

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"

	"github.com/bgpat/gomplete"
	_ "github.com/bgpat/gomplete/shells/bash" // support bash
)

// Flag is a implementation for flag.Value.
type Flag struct {
	Completion gomplete.Completion
}

// String returns empty string.
// This is a dummy function.
func (f *Flag) String() string {
	return ""
}

// Set aborts parsing flags and runs the shell completion.
func (f *Flag) Set(name string) error {
	cfg := gomplete.NewShellConfig(name)

	if len(cfg.Args) == 0 {
		if isatty.IsTerminal(os.Stdout.Fd()) {
			fmt.Printf("usage: source <(%s)\n", strings.Join(os.Args, " "))
			os.Exit(1)
			return nil
		}
		shell, err := gomplete.NewShell(cfg)
		if err != nil {
			return errors.WithStack(err)
		}
		return errors.WithStack(shell.OutputScript(os.Stdout))
	}

	shell, err := gomplete.NewShell(cfg)
	if err != nil {
		return errors.WithStack(err)
	}
	reply := f.Completion.Complete(context.Background(), shell.Args())
	return errors.WithStack(shell.FormatReply(reply, os.Stdout))
}
