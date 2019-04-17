package gomplete

// Args is the words of the command line.
type Args struct {
	words []string
	index int
}

// NewArgs returns an Args from the slice of string.
func NewArgs(words []string) *Args {
	buf := make([]string, len(words))
	copy(buf, words)
	return &Args{
		words: buf,
	}
}

// Current returns the current argumrnt.
func (a *Args) Current() string {
	return a.words[a.index]
}

// IsLast returns true if the index is pointed the last argument.
func (a *Args) IsLast() bool {
	return len(a.words)-1 == a.index
}

// Next returns an immutable Args pointed the next argument.
// If the index is the last, returns nil.
func (a *Args) Next() *Args {
	if a.IsLast() {
		return nil
	}
	return &Args{
		words: a.words,
		index: a.index + 1,
	}
}
