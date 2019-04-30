package flag

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/bgpat/gomplete"
	test "github.com/bgpat/gomplete/test"
)

func TestMain(m *testing.M) {
	osExit = func(code int) {
		panic(fmt.Sprintf("exit code %v", code))
	}

	stdout = &bytes.Buffer{}

	gomplete.RegisterShell("fake", func(*gomplete.ShellConfig) (gomplete.Shell, error) {
		return &test.FakeShell{
			ShellConfig: &gomplete.ShellConfig{
				CommandName:     "cmd",
				CompleteCommand: []string{"sub1", "sub2"},
				Args:            []string{"arg1", "arg2", "arg3"},
				Env:             map[string]string{},
				ShellName:       "fake",
			},
		}, nil
	})
	gomplete.RegisterShell("fake_noargs", func(*gomplete.ShellConfig) (gomplete.Shell, error) {
		return &test.FakeShell{
			ShellConfig: &gomplete.ShellConfig{
				CommandName:     "cmd",
				CompleteCommand: []string{"sub1", "sub2"},
				Env:             map[string]string{},
				ShellName:       "fake_noargs",
			},
		}, nil
	})
	os.Exit(m.Run())
}
