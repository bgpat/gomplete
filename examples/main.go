package main

import (
	"flag"

	"github.com/bgpat/gomplete"
	gomplete_flag "github.com/bgpat/gomplete/plugins/flag"
	_ "github.com/bgpat/gomplete/shells"
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
		&gomplete_flag.Value{
			Completion: &comp,
			FlagName:   "completion",
		},
		"completion",
		"output completion script code",
	)
	flag.Parse()
}
