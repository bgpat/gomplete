package test

import (
	"fmt"
	"io"
	"strings"

	"github.com/bgpat/gomplete"
	"github.com/pkg/errors"
)

// FakeShell is a Shell implementation for testing.
type FakeShell struct {
	*gomplete.ShellConfig
}

// Args returns args of the config.
func (s *FakeShell) Args() *gomplete.Args {
	return gomplete.NewArgs(s.ShellConfig.Args)
}

// FormatReply writes formatted reply to the buffer.
func (s *FakeShell) FormatReply(reply gomplete.Reply, w io.Writer) error {
	lines := make([]string, len(reply))
	for k, v := range reply {
		lines = append(lines, fmt.Sprintf("%q %q\n", k, v))
	}
	_, err := io.WriteString(w, strings.Join(lines, ""))
	return errors.WithStack(err)
}

// OutputScript returns just the command line of the config.
func (s *FakeShell) OutputScript(w io.Writer) error {
	_, err := io.WriteString(w, s.CommandName)
	return errors.WithStack(err)
}

// Usage returns specified command line.
func (s *FakeShell) Usage(cmdline string) string {
	return cmdline
}
