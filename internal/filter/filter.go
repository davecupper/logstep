package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Rule represents a single filter condition applied to a parsed JSON log line.
type Rule struct {
	Field    string
	Operator string // "eq", "contains", "exists"
	Value    string
}

// Filter holds a set of rules and evaluates log lines against them.
type Filter struct {
	Rules []Rule
}

// New creates a Filter from a slice of Rule definitions.
func New(rules []Rule) *Filter {
	return &Filter{Rules: rules}
}

// Match returns true if the raw JSON line satisfies all filter rules.
func (f *Filter) Match(line []byte) (bool, error) {
	if len(f.Rules) == 0 {
		return true, nil
	}

	var record map[string]interface{}
	if err := json.Unmarshal(line, &record); err != nil {
		return false, err
	}

	for _, rule := range f.Rules {
		if !applyRule(rule, record) {
			return false, nil
		}
	}
	return true, nil
}

func applyRule(rule Rule, record map[string]interface{}) bool {
	val, exists := record[rule.Field]

	switch rule.Operator {
	case "exists":
		return exists
	case "eq":
		if !exists {
			return false
		}
		return toString(val) == rule.Value
	case "contains":
		if !exists {
			return false
		}
		return strings.Contains(toString(val), rule.Value)
	default:
		return false
	}
}

func toString(v interface{}) string {
	switch s := v.(type) {
	case string:
		return s
	case float64:
		return fmt.Sprintf("%g", s)
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}
