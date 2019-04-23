package bash

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/bgpat/gomplete"
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
	Name string
	Sub  gomplete.Completion
}

// FormatReply returns reply keys joined by newline.
func (s *Shell) FormatReply(reply gomplete.Reply) string {
	keys := make([]string, 0, len(reply))
	for k := range reply {
		keys = append(keys, k)
	}
	return strings.Join(keys, "\n")
}

// Script returns the shell script to parse replies.
func (s *Shell) Script(cmdline string) string {
	buf := bytes.Buffer{}
	t := template.Must(template.New(s.Name).Parse(scriptTemplate))
	t.Execute(&buf, struct {
		*Shell
		CmdLine string
	}{
		Shell:   s,
		CmdLine: cmdline,
	})
	return buf.String()
}
