package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// Options configures the tail behavior.
type Options struct {
	Follow   bool
	PollInterval time.Duration
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Follow:       true,
		PollInterval: 200 * time.Millisecond,
	}
}

// Tail reads lines from a file, optionally following new content.
// Lines are sent on the returned channel. The channel is closed when
// the context is cancelled or an unrecoverable error occurs.
func Tail(ctx context.Context, path string, opts Options) (<-chan string, <-chan error) {
	lines := make(chan string, 64)
	errs := make(chan error, 1)

	go func() {
		defer close(lines)
		defer close(errs)

		f, err := os.Open(path)
		if err != nil {
			errs <- err
			return
		}
		defer f.Close()

		reader := bufio.NewReader(f)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					if !opts.Follow {
						return
					}
					select {
					case <-ctx.Done():
						return
					case <-time.After(opts.PollInterval):
						continue
					}
				}
				errs <- err
				return
			}

			if len(line) > 0 {
				select {
				case lines <- line:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return lines, errs
}

// TailReader tails lines from an io.Reader (useful for stdin or testing).
func TailReader(ctx context.Context, r io.Reader) <-chan string {
	lines := make(chan string, 64)

	go func() {
		defer close(lines)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case lines <- scanner.Text():
			}
		}
	}()

	return lines
}
