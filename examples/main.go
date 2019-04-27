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
			Sub: &gomplete.Union{
				&gomplete.Command{
					Name:        "hoge",
					Description: "description of hoge",
				},
				&gomplete.Command{
					Name:        "fuga",
					Description: "description of fuga",
				},
				&gomplete.Command{
					Name:        "piyo",
					Description: "description of piyo",
				},
			},
		},
	}
	flag.Var(
		&plugins.Flag{
			Completion: &comp,
			FlagName:   "completion",
		},
		"completion",
		"output completion script code",
	)
	flag.Parse()
}
