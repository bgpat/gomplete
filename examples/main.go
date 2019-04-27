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
					Description: "hoge",
				},
				&gomplete.Command{
					Name:        "fuga",
					Description: "fuga",
				},
				&gomplete.Command{
					Name:        "piyo",
					Description: "piyo",
				},
			},
		},
	}
	flag.Var(&plugins.Flag{
		Completion: &comp,
	}, "completion", "output completion script")
	flag.Parse()
}
