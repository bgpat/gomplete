package gomplete

import (
	"context"
	"strings"
)

// A Command is the simple completion.
type Command struct {
	Name        string
	Description string
	Sub         Completion
}

// Complete returns the pairs that have the prefix of current arg.
func (c *Command) Complete(ctx context.Context, args *Args) Reply {
	if strings.HasPrefix(c.Name, args.Current()) {
		if args.IsLast() {
			return Reply{c.Name: c.Description}
		}
		if c.Sub != nil {
			return c.Sub.Complete(ctx, args.Next())
		}
	}
	return nil
}
