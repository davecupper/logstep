package extractor

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Extractor holds a list of field paths to extract from JSON log lines.
type Extractor struct {
	fields []string
}

// New creates a new Extractor for the given field paths.
// Field paths support dot notation (e.g. "user.id").
func New(fields []string) *Extractor {
	return &Extractor{fields: fields}
}

// Extract parses a JSON line and returns a map of requested fields to their values.
// Fields not present in the JSON are omitted from the result.
func (e *Extractor) Extract(line string) (map[string]string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(line), &data); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	result := make(map[string]string, len(e.fields))
	for _, field := range e.fields {
		val, ok := resolvePath(data, strings.Split(field, "."))
		if ok {
			result[field] = fmt.Sprintf("%v", val)
		}
	}
	return result, nil
}

// resolvePath traverses nested maps using the provided key segments.
func resolvePath(data map[string]interface{}, keys []string) (interface{}, bool) {
	if len(keys) == 0 {
		return nil, false
	}
	val, ok := data[keys[0]]
	if !ok {
		return nil, false
	}
	if len(keys) == 1 {
		return val, true
	}
	nested, ok := val.(map[string]interface{})
	if !ok {
		return nil, false
	}
	return resolvePath(nested, keys[1:])
}
