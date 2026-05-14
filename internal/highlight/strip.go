package highlight

import "regexp"

// ansiRe matches ANSI escape sequences.
var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// Strip removes all ANSI escape sequences from s, returning plain text.
// This is useful when writing output to a file or a non-TTY sink.
func Strip(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}

// StripBytes removes ANSI escape sequences from a byte slice.
func StripBytes(b []byte) []byte {
	return ansiRe.ReplaceAll(b, nil)
}
