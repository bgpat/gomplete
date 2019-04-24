package gomplete

import (
	"errors"
	"fmt"
	"sync"
)

var (
	shellsMu sync.RWMutex
	shells   = make(map[string]Shell)
)

// A Shell is the shell completion interface.
type Shell interface {
	// FormatReply converts the completion reply to the string the script output by Script() can parse.
	FormatReply(reply Reply) string

	// Script outputs the shell script to parse the reply and register the completion.
	Script(cmdline string) string
}

// RegisterShell makes a shell implementation available by the provided name.
// If RegisterShell is called twice with the same name or if shell is nil, it returns an error.
func RegisterShell(name string, shell Shell) error {
	shellsMu.Lock()
	defer shellsMu.Unlock()
	if shell == nil {
		return errors.New("shell is nil")
	}
	if _, dup := shells[name]; dup {
		return fmt.Errorf("%q is already registered", name)
	}
	shells[name] = shell
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

func FormatReply(name string, reply Reply) (string, error) {
	shellsMu.RLock()
	shellsMu.RUnlock()
	return "", nil
}
