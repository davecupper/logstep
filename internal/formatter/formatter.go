// Package formatter provides output formatting for log entries.
package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Format represents the output format type.
type Format string

const (
	FormatJSON    Format = "json"
	FormatText    Format = "text"
	FormatCompact Format = "compact"
)

// Options configures the formatter behavior.
type Options struct {
	Format        Format
	TimeField     string
	TimeFormat    string
	LevelField    string
	MessageField  string
	Colorize      bool
}

// DefaultOptions returns sensible defaults for formatting.
func DefaultOptions() Options {
	return Options{
		Format:       FormatText,
		TimeField:    "time",
		TimeFormat:   time.RFC3339,
		LevelField:   "level",
		MessageField: "msg",
		Colorize:     false,
	}
}

// Formatter formats a parsed log entry for output.
type Formatter struct {
	opts Options
}

// New creates a new Formatter with the given options.
func New(opts Options) *Formatter {
	return &Formatter{opts: opts}
}

// Format formats a raw JSON log line according to the configured output format.
func (f *Formatter) Format(line []byte) (string, error) {
	var entry map[string]interface{}
	if err := json.Unmarshal(line, &entry); err != nil {
		return "", fmt.Errorf("invalid json: %w", err)
	}

	switch f.opts.Format {
	case FormatJSON:
		return string(line), nil
	case FormatCompact:
		return f.formatCompact(entry), nil
	default:
		return f.formatText(entry), nil
	}
}

func (f *Formatter) formatText(entry map[string]interface{}) string {
	var parts []string

	if ts, ok := entry[f.opts.TimeField]; ok {
		parts = append(parts, fmt.Sprintf("%v", ts))
	}
	if lvl, ok := entry[f.opts.LevelField]; ok {
		parts = append(parts, fmt.Sprintf("[%v]", strings.ToUpper(fmt.Sprintf("%v", lvl))))
	}
	if msg, ok := entry[f.opts.MessageField]; ok {
		parts = append(parts, fmt.Sprintf("%v", msg))
	}

	for k, v := range entry {
		if k == f.opts.TimeField || k == f.opts.LevelField || k == f.opts.MessageField {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}

	return strings.Join(parts, " ")
}

func (f *Formatter) formatCompact(entry map[string]interface{}) string {
	lvl, _ := entry[f.opts.LevelField]
	msg, _ := entry[f.opts.MessageField]
	return fmt.Sprintf("%v %v", strings.ToUpper(fmt.Sprintf("%v", lvl)), msg)
}
