// Package highlight provides ANSI color highlighting for log output fields and levels.
package highlight

import "fmt"

// ANSI escape codes.
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"
	White  = "\033[97m"
)

// LevelColor returns an ANSI color code for a known log level string.
// Unrecognised levels return the empty string (no color).
func LevelColor(level string) string {
	switch level {
	case "error", "ERROR", "fatal", "FATAL":
		return Red
	case "warn", "WARN", "warning", "WARNING":
		return Yellow
	case "info", "INFO":
		return Green
	case "debug", "DEBUG", "trace", "TRACE":
		return Cyan
	default:
		return ""
	}
}

// Colorize wraps text with the given ANSI color code and resets afterward.
// If color is empty the original text is returned unchanged.
func Colorize(color, text string) string {
	if color == "" {
		return text
	}
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}

// Field wraps a field value in Cyan for visual distinction.
func Field(value string) string {
	return Colorize(Cyan, value)
}

// Key wraps a field key in Bold+Gray.
func Key(key string) string {
	return fmt.Sprintf("%s%s%s%s", Bold, Gray, key, Reset)
}

// Level colorises a log-level string according to its severity.
func Level(level string) string {
	return Colorize(LevelColor(level), level)
}
