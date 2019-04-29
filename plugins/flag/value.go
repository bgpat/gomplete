package flag

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"

	"github.com/bgpat/gomplete"
)

var (
	isTerminal = isatty.IsTerminal(os.Stdout.Fd())
	osExit     = os.Exit
	stdout     = io.Writer(os.Stdout)
)

// Value is a implementation for flag.Value.
type Value struct {
	Completion gomplete.Completion
	FlagName   string

	// Calls to be exit.
	// The default is os.Exit(1).
	Exit func(code int)
}

// String returns empty string.
// This is a dummy function.
func (v *Value) String() string {
	return ""
}

// Set aborts parsing flags and runs the shell completion.
func (v *Value) Set(name string) error {
	if name == "true" {
		fmt.Fprintf(
			stdout,
			"specify shell: %v\nexample: -completion=bash\n",
			gomplete.Shells(),
		)
		v.exit(1)
		return nil
	}

	cfg := gomplete.NewShellConfig(name)
	shell, err := gomplete.NewShell(cfg)
	if err != nil {
		return errors.WithStack(err)
	}

	if len(cfg.Args) == 0 {
		if isTerminal {
			fmt.Fprintln(stdout, "usage:", shell.Usage(strings.Join(os.Args, " ")))
			v.exit(1)
			return nil
		}
		return errors.WithStack(shell.OutputScript(stdout))
	}

	reply := v.Completion.Complete(context.Background(), shell.Args())
	return errors.WithStack(shell.FormatReply(reply, stdout))
}

// IsBoolFlag is a method to meet flag.boolValue.
func (v *Value) IsBoolFlag() bool {
	return true
}

func (v *Value) exit(code int) {
	if v.Exit == nil {
		osExit(code)
		return
	}
	v.Exit(code)
}
