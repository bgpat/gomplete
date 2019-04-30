package flag

import (
	"context"
	"flag"
	"strings"

	"github.com/bgpat/gomplete"
)

// Completion is a completion implementation for flag.FlagSet.
type Completion struct {
	*flag.FlagSet
	Parent gomplete.Completion
}

type flagCompletion struct {
	*flag.Flag
	parent *Completion
}

type valueCompletion struct {
	flag.Value
	parent *flagCompletion
}

// Complete returns the completion reply based on flag set.
func (c *Completion) Complete(ctx context.Context, args *gomplete.Args) gomplete.Reply {
	if c.FlagSet == nil {
		return nil
	}
	flags := gomplete.Union{}
	c.FlagSet.VisitAll(func(f *flag.Flag) {
		flags = append(flags, &flagCompletion{
			Flag:   f,
			parent: c,
		})
	})
	return (&gomplete.Union{
		&flags,
		c.Parent,
	}).Complete(ctx, args)
}

func (f *flagCompletion) Complete(ctx context.Context, args *gomplete.Args) gomplete.Reply {
	if f.Flag == nil {
		return nil
	}
	if !args.IsLast() {
		return nil
	}
	s := "-" + f.Flag.Name + "="
	if strings.HasPrefix(s, args.Current()) {
		return gomplete.Reply{s: f.Flag.Usage}
	}
	return nil
}

func (v *valueCompletion) Complete(ctx context.Context, args *gomplete.Args) gomplete.Reply {
	return nil
}
