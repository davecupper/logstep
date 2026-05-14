package highlight_test

import (
	"testing"

	"github.com/user/logstep/internal/highlight"
)

func TestStripPlain(t *testing.T) {
	got := highlight.Strip("hello world")
	if got != "hello world" {
		t.Errorf("Strip of plain text changed value: %q", got)
	}
}

func TestStripRemovesColor(t *testing.T) {
	colored := highlight.Colorize(highlight.Red, "error")
	got := highlight.Strip(colored)
	if got != "error" {
		t.Errorf("Strip(%q) = %q, want %q", colored, got, "error")
	}
}

func TestStripMultiple(t *testing.T) {
	s := highlight.Key("level") + ": " + highlight.Level("warn")
	got := highlight.Strip(s)
	if got != "level: warn" {
		t.Errorf("Strip of multi-colored string = %q, want %q", got, "level: warn")
	}
}

func TestStripBytes(t *testing.T) {
	colored := []byte(highlight.Colorize(highlight.Green, "ok"))
	got := highlight.StripBytes(colored)
	if string(got) != "ok" {
		t.Errorf("StripBytes = %q, want %q", got, "ok")
	}
}
