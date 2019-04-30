package gomplete

import "context"

// Union is the alias of []Completion.
// It merges replies of each elements.
type Union []Completion

// Complete returns the union of replies of all elements.
func (u *Union) Complete(ctx context.Context, args *Args) Reply {
	if u == nil {
		return nil
	}
	reply := Reply{}
	for _, comp := range *u {
		if comp == nil {
			continue
		}
		for k, v := range comp.Complete(ctx, args) {
			reply[k] = v
		}
	}
	return reply
}
