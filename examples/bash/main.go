package main

import (
	"flag"

	"github.com/bgpat/gomplete"
	"github.com/bgpat/gomplete/examples"
	"github.com/bgpat/gomplete/shells/bash"
)

func main() {
	comp := bash.Completion{
		Name: "example",
		Sub: &gomplete.Command{
			Name:        "foo",
			Description: "first sub command",
			Sub: &gomplete.Command{
				Name:        "bar",
				Description: "second sub command",
			},
		},
	}
	flag.Var(&examples.Flag{Shell: &comp}, "completion", "output completion script")
	flag.Parse()
}
