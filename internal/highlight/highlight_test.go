package highlight_test

import (
	"strings"
	"testing"

	"github.com/user/logstep/internal/highlight"
)

func TestColorizeEmpty(t *testing.T) {
	got := highlight.Colorize("", "hello")
	if got != "hello" {
		t.Errorf("expected plain text, got %q", got)
	}
}

func TestColorizeWraps(t *testing.T) {
	got := highlight.Colorize(highlight.Red, "err")
	if !strings.Contains(got, "err") {
		t.Error("colorized string should contain original text")
	}
	if !strings.HasPrefix(got, highlight.Red) {
		t.Error("colorized string should start with color code")
	}
	if !strings.HasSuffix(got, highlight.Reset) {
		t.Error("colorized string should end with reset code")
	}
}

func TestLevelColorKnown(t *testing.T) {
	cases := []struct {
		level string
		want  string
	}{
		{"error", highlight.Red},
		{"ERROR", highlight.Red},
		{"fatal", highlight.Red},
		{"warn", highlight.Yellow},
		{"WARNING", highlight.Yellow},
		{"info", highlight.Green},
		{"INFO", highlight.Green},
		{"debug", highlight.Cyan},
		{"TRACE", highlight.Cyan},
	}
	for _, tc := range cases {
		got := highlight.LevelColor(tc.level)
		if got != tc.want {
			t.Errorf("LevelColor(%q) = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestLevelColorUnknown(t *testing.T) {
	got := highlight.LevelColor("verbose")
	if got != "" {
		t.Errorf("expected empty color for unknown level, got %q", got)
	}
}

func TestKeyAndField(t *testing.T) {
	key := highlight.Key("ts")
	if !strings.Contains(key, "ts") {
		t.Error("Key output should contain the key name")
	}
	field := highlight.Field("value")
	if !strings.Contains(field, "value") {
		t.Error("Field output should contain the value")
	}
}
