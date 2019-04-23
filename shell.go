package gomplete

// A Shell is the shell completion interface.
type Shell interface {
	// FormatReply converts the completion reply to the string the script output by Script() can parse.
	FormatReply(reply Reply) string

	// Script outputs the shell script to parse the reply and register the completion.
	Script(cmdline string) string
}
