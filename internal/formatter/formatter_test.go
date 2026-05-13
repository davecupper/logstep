package formatter_test

import (
	"testing"

	"github.com/yourorg/logstep/internal/formatter"
)

func TestFormatJSON(t *testing.T) {
	opts := formatter.DefaultOptions()
	opts.Format = formatter.FormatJSON
	f := formatter.New(opts)

	line := []byte(`{"time":"2024-01-01T00:00:00Z","level":"info","msg":"hello"}`)
	out, err := f.Format(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != string(line) {
		t.Errorf("expected raw json, got %q", out)
	}
}

func TestFormatText(t *testing.T) {
	opts := formatter.DefaultOptions()
	opts.Format = formatter.FormatText
	f := formatter.New(opts)

	line := []byte(`{"time":"2024-01-01T00:00:00Z","level":"info","msg":"started"}`)
	out, err := f.Format(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormatCompact(t *testing.T) {
	opts := formatter.DefaultOptions()
	opts.Format = formatter.FormatCompact
	f := formatter.New(opts)

	line := []byte(`{"level":"warn","msg":"low memory"}`)
	out, err := f.Format(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "WARN low memory" {
		t.Errorf("unexpected compact output: %q", out)
	}
}

func TestFormatInvalidJSON(t *testing.T) {
	f := formatter.New(formatter.DefaultOptions())
	_, err := f.Format([]byte(`not json`))
	if err == nil {
		t.Error("expected error for invalid json")
	}
}

func TestFormatMissingFields(t *testing.T) {
	opts := formatter.DefaultOptions()
	opts.Format = formatter.FormatText
	f := formatter.New(opts)

	line := []byte(`{"service":"api","status":200}`)
	out, err := f.Format(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected non-empty output for entry with extra fields")
	}
}
