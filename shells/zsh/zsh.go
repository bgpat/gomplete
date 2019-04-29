package zsh

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

const scriptTemplate = `_{{sanitize .CommandName}}() {
	IFS=$'\n'
	local -a completion
	completion=( $( CURRENT=$CURRENT {{join .CompleteCommand}} -- "${words[@]}" ) )
	if [ ${#completion} -gt 0 ]; then
		_values {{.CommandName}} "${completion[@]}"
	fi
}
compdef _{{sanitize .CommandName}} {{.CommandName}}
`

var (
	funcMap    template.FuncMap
	sanitizeRe = regexp.MustCompile("[^a-zA-Z0-9-]+")
)

// Shell is the implementation of gomplete.Shell for zsh.
type Shell struct {
	*gomplete.ShellConfig
	current int
}

func init() {
	gomplete.RegisterShell("zsh", NewShell)

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
	current := len(config.Args)
	if v, ok := config.Env["CURRENT"]; ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		current = i
	}
	return &Shell{
		ShellConfig: config,
		current:     current,
	}, nil
}

// Args returns returns command-line arguments to complete.
func (s *Shell) Args() *gomplete.Args {
	return gomplete.NewArgs(s.ShellConfig.Args[:s.current]).Next()
}

// FormatReply returns reply keys joined by newline.
func (s *Shell) FormatReply(reply gomplete.Reply, w io.Writer) error {
	values := make([]string, 0, len(reply))
	for k, v := range reply {
		k = strings.ReplaceAll(k, "\\", "\\\\")
		k = strings.ReplaceAll(k, "[", "\\]")
		if v == "" {
			values = append(values, k)
			continue
		}
		v = strings.ReplaceAll(v, "\\", "\\\\")
		v = strings.ReplaceAll(v, "]", "\\]")
		values = append(values, fmt.Sprintf("%s[%s]", k, v))
	}
	sort.Strings(values)
	_, err := io.WriteString(w, strings.Join(values, "\n"))
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
