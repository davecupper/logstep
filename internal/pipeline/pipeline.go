// Package pipeline wires together tail, filter, extractor, and formatter
// into a single processing pipeline for JSON log streams.
package pipeline

import (
	"context"
	"io"

	"github.com/user/logstep/internal/extractor"
	"github.com/user/logstep/internal/filter"
	"github.com/user/logstep/internal/formatter"
	"github.com/user/logstep/internal/tail"
)

// Options holds configuration for the pipeline.
type Options struct {
	Tail      tail.Options
	Filter    []filter.Rule
	Fields    []string
	Formatter formatter.Options
}

// Pipeline reads log lines from a source, applies filtering,
// extracts fields, and writes formatted output.
type Pipeline struct {
	opts      Options
	filter    *filter.Filter
	extractor *extractor.Extractor
	fmt       *formatter.Formatter
}

// New creates a new Pipeline with the given options.
func New(opts Options) (*Pipeline, error) {
	f, err := filter.New(opts.Filter)
	if err != nil {
		return nil, err
	}

	ex := extractor.New(opts.Fields)

	fm, err := formatter.New(opts.Formatter)
	if err != nil {
		return nil, err
	}

	return &Pipeline{
		opts:      opts,
		filter:    f,
		extractor: ex,
		fmt:       fm,
	}, nil
}

// Run starts the pipeline, reading from src and writing to dst.
// It blocks until ctx is cancelled or src is exhausted.
func (p *Pipeline) Run(ctx context.Context, src io.Reader, dst io.Writer) error {
	lines := tail.TailReader(ctx, src, p.opts.Tail)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case line, ok := <-lines:
			if !ok {
				return nil
			}
			if err := p.processLine(line, dst); err != nil {
				return err
			}
		}
	}
}

func (p *Pipeline) processLine(line []byte, dst io.Writer) error {
	matched, err := p.filter.Match(line)
	if err != nil || !matched {
		return nil
	}

	fields, err := p.extractor.Extract(line)
	if err != nil {
		return nil
	}

	formatted, err := p.fmt.Format(line, fields)
	if err != nil {
		return nil
	}

	_, err = dst.Write(append(formatted, '\n'))
	return err
}
