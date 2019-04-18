package gomplete

import (
	"context"
	"strings"
)

// A Command is the set of the completion pairs.
// Each element has the candidate word and the description.
type Command map[string]string

// Complete returns the pairs that have the prefix of current arg.
func (c *Command) Complete(ctx context.Context, args *Args) Reply {
	if !args.IsLast() {
		return nil
	}
	reply := Reply{}
	for k, v := range *c {
		if strings.HasPrefix(k, args.Current()) {
			reply[k] = v
		}
	}
	return reply
}
