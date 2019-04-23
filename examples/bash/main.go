package main

import (
	"flag"

	"github.com/bgpat/gomplete"
	"github.com/bgpat/gomplete/examples"
	"github.com/bgpat/gomplete/shells/bash"
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
	shell := bash.Shell{}
	flag.Var(&examples.Flag{
		Shell:      &shell,
		Completion: &comp,
	}, "completion", "output completion script")
	flag.Parse()
}
