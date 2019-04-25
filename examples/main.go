package main

import (
	"flag"

	"github.com/bgpat/gomplete"
	"github.com/bgpat/gomplete/plugins"
)

func main() {
	comp := gomplete.Command{
		Name:        "foo",
		Description: "first sub command",
		Sub: &gomplete.Command{
			Name:        "bar",
			Description: "second sub command",
		},
	}
	cfg := gomplete.ShellConfig{
		CommandName:     "examples",
		CompleteCommand: "examples -completion",
	}
	flag.Var(&plugins.Flag{
		Completion: &comp,
		Config:     cfg,
	}, "completion", "output completion script")
	flag.Parse()
}
