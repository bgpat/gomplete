package bash

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"github.com/bgpat/gomplete"
)

const scriptTemplate = `_{{sanitize .CommandName}}_completion() {
	IFS=$'\n'
	COMPREPLY=( $( COMP_CWORD=$COMP_CWORD {{join .CompleteCommand}} -- "${COMP_WORDS[@]}") )
}
complete -F _{{sanitize .CommandName}}_completion {{.CommandName}}
`

var (
	funcMap    template.FuncMap
	sanitizeRe = regexp.MustCompile("[^a-zA-Z0-9]+")
)

// Shell is the implementation of gomplete.Shell for bash.
type Shell struct {
	*gomplete.ShellConfig
	current int
}

func init() {
	gomplete.RegisterShell("bash", NewShell)

	funcMap = template.FuncMap{
		"sanitize": func(str string) string {
			return sanitizeRe.ReplaceAllString(str, "_")
		},
		"join": func(src []string) string {
			tmp := make([]string, 0, len(src))
			for _, s := range src {
				tmp = append(tmp, strconv.Quote(s))
			}
			return strings.Join(tmp, " ")
		},
	}
}

// NewShell returns a shell instance from shell config.
func NewShell(config *gomplete.ShellConfig) (gomplete.Shell, error) {
	return newShell(config)
}

func newShell(config *gomplete.ShellConfig) (*Shell, error) {
	current := len(config.Args) - 1
	if v, ok := config.Env["COMP_CWORD"]; ok {
		cword, err := strconv.Atoi(v)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		current = cword
	}
	return &Shell{
		ShellConfig: config,
		current:     current,
	}, nil
}

// Args returns returns command-line arguments to complete.
func (s *Shell) Args() *gomplete.Args {
	return gomplete.NewArgs(s.ShellConfig.Args[:s.current+1]).Next()
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
	t := template.Must(template.New(s.CommandName).Funcs(funcMap).Parse(scriptTemplate))
	return errors.WithStack(t.Execute(w, s))
}

// Usage returns the usage of the shell script.
func (s *Shell) Usage(cmdline string) string {
	return fmt.Sprintf("source <(%s)", cmdline)
}
