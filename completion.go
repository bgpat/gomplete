package gomplete

import "context"

// Completion is the interface has the completion information.
type Completion interface {
	// Complete returns the command line completion reply from args.
	Complete(ctx context.Context, args *Args) Reply
}

// Reply is the alias for map[string]string.
// The key is the candidate completion word, and the value is the description.
type Reply map[string]string
