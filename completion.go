package gomplete

// Completion is the interface has the completion infromation.
type Completion interface {
	Complete(args *Args) Reply
}

// Reply is the alias for map[string]string
type Reply map[string]string
