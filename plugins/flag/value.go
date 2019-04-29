package flag

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"

	"github.com/bgpat/gomplete"
)

// Value is a implementation for flag.Value.
type Value struct {
	Completion gomplete.Completion
	FlagName   string
}

// String returns empty string.
// This is a dummy function.
func (v *Value) String() string {
	return ""
}

// Set aborts parsing flags and runs the shell completion.
func (v *Value) Set(name string) error {
	if name == "true" {
		fmt.Printf("specify shell: %v\n", gomplete.Shells())
		fmt.Println("example: -completion=bash")
		os.Exit(1)
	}

	cfg := gomplete.NewShellConfig(name)
	shell, err := gomplete.NewShell(cfg)
	if err != nil {
		return errors.WithStack(err)
	}

	if len(cfg.Args) == 0 {
		if isatty.IsTerminal(os.Stdout.Fd()) {
			fmt.Println("usage:", shell.Usage(strings.Join(os.Args, " ")))
			os.Exit(1)
			return nil
		}
		return errors.WithStack(shell.OutputScript(os.Stdout))
	}

	reply := v.Completion.Complete(context.Background(), shell.Args())
	return errors.WithStack(shell.FormatReply(reply, os.Stdout))
}

// IsBoolFlag is a method to meet flag.boolValue.
func (v *Value) IsBoolFlag() bool {
	return true
}
