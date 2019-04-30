package flag

import (
	"os"
	"testing"

	"github.com/bgpat/gomplete"
)

func TestValueString(t *testing.T) {
	v := Value{}
	s := v.String()
	if s != "" {
		t.Errorf("want empty string, but got %s", s)
	}
}

func TestValueSet(t *testing.T) {
	t.Run("list shells", func(t *testing.T) {
		var exitCode int
		v := Value{
			Exit: func(code int) { exitCode = code },
		}
		if err := v.Set("true"); err != nil {
			t.Fatal(err)
		}
		if exitCode != 1 {
			t.Errorf("eixt code must 1, but got %v", exitCode)
		}
	})

	t.Run("unknown shell", func(t *testing.T) {
		v := Value{}
		if err := v.Set("unknown"); err == nil {
			t.Error("must return an error because of unknown shell")
		}
	})

	t.Run("output script", func(t *testing.T) {
		args := os.Args
		os.Args = []string{"a", "b", "c"}
		defer func() { os.Args = args }()
		t.Run("terminal", func(t *testing.T) {
			defaultIsTerminal := isTerminal
			isTerminal = true
			defer func() { isTerminal = defaultIsTerminal }()

			var exitCode int
			v := Value{
				Exit: func(code int) { exitCode = code },
			}
			if err := v.Set("fake_noargs"); err != nil {
				t.Fatal(err)
			}
			if exitCode != 1 {
				t.Errorf("eixt code must 1, but got %v", exitCode)
			}
		})
		t.Run("not terminal", func(t *testing.T) {
			defaultIsTerminal := isTerminal
			isTerminal = false
			defer func() { isTerminal = defaultIsTerminal }()

			v := Value{}
			if err := v.Set("fake_noargs"); err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("complete", func(t *testing.T) {
		args := os.Args
		os.Args = []string{"a", "--", "b", "c"}
		defer func() { os.Args = args }()
		v := Value{Completion: &gomplete.Command{}}
		if err := v.Set("fake"); err != nil {
			t.Fatal(err)
		}
	})
}

func TestValueIsBoolFlag(t *testing.T) {
	v := Value{}
	if !v.IsBoolFlag() {
		t.Error("must be true")
	}
}

func TestValue(t *testing.T) {
	var exitCode int
	osExit = func(code int) {
		exitCode = code
	}
	v := Value{}
	v.exit(1)
	if exitCode != 1 {
		t.Errorf("eixt code must 1, but got %v", exitCode)
	}
}
