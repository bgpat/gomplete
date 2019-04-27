package gomplete

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

var (
	shellsMu sync.RWMutex
	shells   = make(map[string]func(*ShellConfig) (Shell, error))
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
	CommandName     string            // The name of completion function.
	CompleteCommand []string          // The prefix of completion command.
	Args            []string          // The command line arguments.
	Env             map[string]string // The map of environment variables.
	ShellName       string            //The name of the shell
}

// RegisterShell makes a shell implementation available by the provided name.
// If RegisterShell is called twice with the same name or if constructor is nil, it panics.
func RegisterShell(name string, constructor func(config *ShellConfig) (Shell, error)) {
	shellsMu.Lock()
	defer shellsMu.Unlock()
	if constructor == nil {
		panic(errors.New("shell is nil"))
	}
	if _, dup := shells[name]; dup {
		panic(fmt.Errorf("%q is already registered", name))
	}
	shells[name] = constructor
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
func NewShell(config *ShellConfig) (Shell, error) {
	shellsMu.RLock()
	defer shellsMu.RUnlock()
	constructor := shells[config.ShellName]
	if constructor == nil {
		return nil, fmt.Errorf("unknown shell %q (forgotten import?)", config.ShellName)
	}
	shell, err := constructor(config)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to initialize %q", config.ShellName)
	}
	return shell, nil
}

// NewShellConfig returns the default config from command-line arguments and environment variables.
func NewShellConfig(shell string) *ShellConfig {
	cfg := ShellConfig{
		ShellName: shell,
	}

	if len(os.Args) > 0 {
		cfg.CommandName = filepath.Base(os.Args[0])
	}

	for i, arg := range os.Args {
		if arg == "--" {
			cfg.CompleteCommand = os.Args[:i]
			cfg.Args = os.Args[i+1:]
			break
		}
	}

	if cfg.CompleteCommand == nil {
		cfg.CompleteCommand = os.Args
	}

	cfg.Env = parseEnv()

	return &cfg
}

func parseEnv() map[string]string {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		v := strings.SplitN(env, "=", 2)
		if len(v) == 2 {
			envs[v[0]] = v[1]
		}
	}
	return envs
}
