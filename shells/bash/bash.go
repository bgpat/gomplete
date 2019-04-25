package bash

import (
	"io"
	"sort"
	"strings"
	"text/template"

	"github.com/bgpat/gomplete"
	"github.com/pkg/errors"
)

const scriptTemplate = `
_{{.Name}}_completion() {
	local words cword
	if type _get_comp_words_by_ref &>/dev/null; then
		_get_comp_words_by_ref -n = -n @ -n : -w words -i cword
	else
		cword="$COMP_CWORD"
		words=("${COMP_WORDS[@]}")
	fi
	local si="$IFS"
	IFS=$'\n' COMPREPLY=($( \
		COMP_CWORD="$cword" \
		COMP_LINE="$COMP_LINE" \
		COMP_POINT="$COMP_POINT" \
		{{.CmdLine}} "${words[@]}" \
		2>/dev/null \
	)) || return $?
	IFS="$si"
	if type __ltrim_colon_completions &>/dev/null; then
		__ltrim_colon_completions "${words[cword]}"
	fi
}
complete -o default -F _{{.Name}}_completion {{.Name}}
`

// Shell is the implementation of gomplete.Shell for bash.
type Shell struct {
	gomplete.ShellConfig
}

func init() {
	gomplete.RegisterShell("bash", NewShell)
}

// NewShell returns a shell instance from shell config.
func NewShell(config gomplete.ShellConfig) (gomplete.Shell, error) {
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
	t := template.Must(template.New(s.Name).Parse(scriptTemplate))
	return errors.WithStack(t.Execute(w, s))
}
