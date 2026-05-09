package extractor

import (
	"fmt"
	"strings"
)

// Format controls how extracted fields are rendered.
type Format int

const (
	// FormatText renders fields as key=value pairs separated by spaces.
	FormatText Format = iota
	// FormatCSV renders field values separated by commas in declaration order.
	FormatCSV
)

// Render formats the extracted fields map according to the given format.
// The fields slice determines ordering for CSV output.
func Render(fields []string, extracted map[string]string, format Format) string {
	switch format {
	case FormatCSV:
		return renderCSV(fields, extracted)
	default:
		return renderText(extracted)
	}
}

func renderText(extracted map[string]string) string {
	parts := make([]string, 0, len(extracted))
	for k, v := range extracted {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, " ")
}

func renderCSV(fields []string, extracted map[string]string) string {
	values := make([]string, len(fields))
	for i, f := range fields {
		if v, ok := extracted[f]; ok {
			values[i] = v
		} else {
			values[i] = ""
		}
	}
	return strings.Join(values, ",")
}
