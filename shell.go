package gomplete

// A Shell is the shell completion interface.
type Shell interface {
	Completion

	// Format converts the completion reply to the string the script output by Script() can parse.
	Format(reply Reply) string

	// Script outputs the shell script to parse the reply and register the completion.
	Script(cmdline string) string
}
