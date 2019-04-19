package bash

import (
	"bytes"
	"context"
	"strings"
	"text/template"

	"github.com/bgpat/gomplete"
)

const scriptTemplate = `
_{{.Name}}_completion () {
  COMPREPLY=( hoge )
}
complete -o default -F _{{.Name}}_completion {{.Name}}
`

// A Completion is the implementation of Shell for bash.
type Completion struct {
	Name string
	Sub  gomplete.Completion
}

// Complete seeks the args by 1.
func (c *Completion) Complete(ctx context.Context, args *gomplete.Args) gomplete.Reply {
	if c.Sub == nil {
		return nil
	}
	a := args.Next()
	if a == nil {
		return nil
	}
	return c.Sub.Complete(ctx, a)
}

// Format returns reply keys joined by newline.
func (c *Completion) Format(reply gomplete.Reply) string {
	keys := make([]string, 0, len(reply))
	for k := range reply {
		keys = append(keys, k)
	}
	return strings.Join(keys, "\n")
}

// Script returns the shell script to parse replies.
func (c *Completion) Script(cmdline string) string {
	buf := bytes.Buffer{}
	t := template.Must(template.New(c.Name).Parse(scriptTemplate))
	t.Execute(&buf, c)
	return buf.String()
}
