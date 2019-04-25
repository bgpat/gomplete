package gomplete

import (
	"fmt"
	"io"
	"sync"

	"github.com/pkg/errors"
)

var (
	shellsMu sync.RWMutex
	shells   = make(map[string]func(ShellConfig) (Shell, error))
)

// Shell is the shell completion interface.
type Shell interface {
	// Args returns the command line arguments.
	Args() *Args

	// FormatReply converts the completion reply to the string the script output by Script() can parse.
	FormatReply(reply Reply, w io.Writer) error

	// OutputScript outputs the shell script to parse the reply and register the completion.
	OutputScript(w io.Writer) error
}

// ShellConfig is the configuration for shell.
type ShellConfig struct {
	CommandName     string // The name of completion function.
	CompleteCommand string // The prefix of completion command.

	Args      []string          // The command line arguments.
	Env       map[string]string // The map of environment variables.
	ShellName string            //The name of the shell
}

// RegisterShell makes a shell implementation available by the provided name.
// If RegisterShell is called twice with the same name or if constructor is nil, it returns an error.
func RegisterShell(name string, constructor func(config ShellConfig) (Shell, error)) error {
	shellsMu.Lock()
	defer shellsMu.Unlock()
	if constructor == nil {
		return errors.New("shell is nil")
	}
	if _, dup := shells[name]; dup {
		return fmt.Errorf("%q is already registered", name)
	}
	shells[name] = constructor
	return nil
}

// Shells returns a list of the names of the registered shells.
func Shells() []string {
	shellsMu.RLock()
	defer shellsMu.RUnlock()
	list := make([]string, 0, len(shells))
	for name := range shells {
		list = append(list, name)
	}
	return list
}

// NewShell creates a new Shell instances by the provided name.
func NewShell(name string, config ShellConfig) (Shell, error) {
	shellsMu.RLock()
	defer shellsMu.RUnlock()
	constructor := shells[name]
	if constructor == nil {
		return nil, fmt.Errorf("unknown shell %q (forgotten import?)", name)
	}
	cfg := config
	cfg.ShellName = name
	shell, err := constructor(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to initialize %q", name)
	}
	return shell, nil
}
