package tail_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logstep/internal/tail"
)

func TestTailReader(t *testing.T) {
	input := "line1\nline2\nline3\n"
	reader := strings.NewReader(input)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	lines := tail.TailReader(ctx, reader)

	var got []string
	for line := range lines {
		got = append(got, line)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	if got[0] != "line1" || got[1] != "line2" || got[2] != "line3" {
		t.Errorf("unexpected lines: %v", got)
	}
}

func TestTailReaderContextCancel(t *testing.T) {
	// Use a pipe so reading blocks indefinitely.
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer pr.Close()
	defer pw.Close()

	ctx, cancel := context.WithCancel(context.Background())

	lines := tail.TailReader(ctx, pr)

	// Write one line then cancel.
	pw.WriteString("hello\n")
	time.Sleep(50 * time.Millisecond)
	cancel()

	var count int
	for range lines {
		count++
	}
	// We may receive 0 or 1 lines depending on timing; just ensure no hang.
	if count > 1 {
		t.Errorf("expected at most 1 line, got %d", count)
	}
}

func TestTailFileNoFollow(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "logstep-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	f.WriteString("{\"level\":\"info\"}\n")
	f.WriteString("{\"level\":\"error\"}\n")
	f.Close()

	ctx := context.Background()
	opts := tail.Options{Follow: false}
	lines, errs := tail.Tail(ctx, f.Name(), opts)

	var got []string
	for line := range lines {
		got = append(got, strings.TrimSpace(line))
	}

	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(got), got)
	}
}

func TestTailFileNotFound(t *testing.T) {
	ctx := context.Background()
	opts := tail.DefaultOptions()
	opts.Follow = false
	_, errs := tail.Tail(ctx, "/nonexistent/path/file.log", opts)

	err := <-errs
	if err == nil {
		t.Fatal("expected an error for missing file")
	}
}
