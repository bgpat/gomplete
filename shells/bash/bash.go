package bash

import (
	"io"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/bgpat/gomplete"
	"github.com/pkg/errors"
)

const scriptTemplate = `_{{sanitize .CommandName}}_completion() {
	IFS=$'\n'
	COMPREPLY=( $({{.CompleteCommand}} {{.ShellName}} -- "${COMP_WORDS[@]}") )
}
complete -o default -F _{{sanitize .CommandName}}_completion {{.CommandName}}
`

var (
	funcMap    template.FuncMap
	sanitizeRe = regexp.MustCompile("[^a-zA-Z0-9]+")
)

// Shell is the implementation of gomplete.Shell for bash.
type Shell struct {
	gomplete.ShellConfig
}

func init() {
	if err := gomplete.RegisterShell("bash", NewShell); err != nil {
		panic(err)
	}

	funcMap = template.FuncMap{
		"sanitize": func(str string) string {
			return sanitizeRe.ReplaceAllString(str, "_")
		},
	}
}

// NewShell returns a shell instance from shell config.
func NewShell(config gomplete.ShellConfig) (gomplete.Shell, error) {
	return newShell(config)
}

func newShell(config gomplete.ShellConfig) (*Shell, error) {
	return &Shell{
		ShellConfig: config,
	}, nil
}

// Args returns returns command-line arguments to complete.
func (s *Shell) Args() *gomplete.Args {
	return gomplete.NewArgs(s.ShellConfig.Args)
}

// FormatReply returns reply keys joined by newline.
func (s *Shell) FormatReply(reply gomplete.Reply, w io.Writer) error {
	keys := make([]string, 0, len(reply))
	for k := range reply {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	_, err := io.WriteString(w, strings.Join(keys, "\n"))
	return errors.WithStack(err)
}

// OutputScript returns the shell script to parse replies.
func (s *Shell) OutputScript(w io.Writer) error {
	t, err := template.New(s.CommandName).Funcs(funcMap).Parse(scriptTemplate)
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(t.Execute(w, s))
}
