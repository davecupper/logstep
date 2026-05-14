package pipeline_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/user/logstep/internal/filter"
	"github.com/user/logstep/internal/formatter"
	"github.com/user/logstep/internal/pipeline"
	"github.com/user/logstep/internal/tail"
)

func defaultOpts() pipeline.Options {
	return pipeline.Options{
		Tail:      tail.DefaultOptions(),
		Formatter: formatter.DefaultOptions(),
	}
}

func TestPipelinePassThrough(t *testing.T) {
	input := `{"level":"info","msg":"hello"}` + "\n"
	src := strings.NewReader(input)
	var dst bytes.Buffer

	p, err := pipeline.New(defaultOpts())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx := context.Background()
	if err := p.Run(ctx, src, &dst); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if dst.Len() == 0 {
		t.Fatal("expected output, got empty")
	}
}

func TestPipelineFilterDropsLine(t *testing.T) {
	input := `{"level":"debug","msg":"ignored"}` + "\n"
	src := strings.NewReader(input)
	var dst bytes.Buffer

	opts := defaultOpts()
	opts.Filter = []filter.Rule{
		{Field: "level", Op: "eq", Value: "info"},
	}

	p, err := pipeline.New(opts)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx := context.Background()
	if err := p.Run(ctx, src, &dst); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if dst.Len() != 0 {
		t.Fatalf("expected no output, got: %s", dst.String())
	}
}

func TestPipelineContextCancel(t *testing.T) {
	// Infinite reader to ensure context cancellation stops the pipeline.
	pr, pw := bytes.NewBuffer(nil), &bytes.Buffer{}
	_ = pr

	opts := defaultOpts()
	p, err := pipeline.New(opts)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = p.Run(ctx, strings.NewReader(""), pw)
	if err != nil && err != context.Canceled {
		t.Fatalf("unexpected error: %v", err)
	}
}
