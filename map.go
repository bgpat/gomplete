package gomplete

import (
	"context"
	"strings"
)

// A Map is the set of the completion pairs.
// Each element has the candidate word and the description.
type Map map[string]string

// Complete returns the pairs that have the prefix of current arg.
func (m *Map) Complete(ctx context.Context, args *Args) Reply {
	if !args.IsLast() {
		return nil
	}
	reply := Reply{}
	for k, v := range *m {
		if strings.HasPrefix(k, args.Current()) {
			reply[k] = v
		}
	}
	return reply
}
