package extractor

import (
	"testing"
)

func TestExtractFlatFields(t *testing.T) {
	e := New([]string{"level", "msg"})
	result, err := e.Extract(`{"level":"info","msg":"hello","ts":"2024-01-01"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["level"] != "info" {
		t.Errorf("expected level=info, got %q", result["level"])
	}
	if result["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %q", result["msg"])
	}
	if _, ok := result["ts"]; ok {
		t.Error("ts should not be extracted")
	}
}

func TestExtractNestedField(t *testing.T) {
	e := New([]string{"user.id", "user.name"})
	result, err := e.Extract(`{"user":{"id":"42","name":"alice"}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["user.id"] != "42" {
		t.Errorf("expected user.id=42, got %q", result["user.id"])
	}
	if result["user.name"] != "alice" {
		t.Errorf("expected user.name=alice, got %q", result["user.name"])
	}
}

func TestExtractMissingField(t *testing.T) {
	e := New([]string{"missing"})
	result, err := e.Extract(`{"level":"warn"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["missing"]; ok {
		t.Error("missing field should not appear in result")
	}
}

func TestExtractInvalidJSON(t *testing.T) {
	e := New([]string{"level"})
	_, err := e.Extract(`not json`)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestExtractNoFields(t *testing.T) {
	e := New([]string{})
	result, err := e.Extract(`{"level":"info"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}
