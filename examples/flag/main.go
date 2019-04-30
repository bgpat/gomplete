package main

import (
	"flag"

	gomplete_flag "github.com/bgpat/gomplete/plugins/flag"
	_ "github.com/bgpat/gomplete/shells"
)

func main() {
	comp := gomplete_flag.Completion{
		FlagSet: flag.CommandLine,
	}
	flag.String("string", "foo", "string flag")
	flag.Int("int", 100, "int flag")
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
