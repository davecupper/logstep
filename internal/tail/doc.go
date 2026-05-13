// Package tail provides utilities for reading log lines from files or
// io.Reader sources, with optional follow (tail -f) behaviour.
//
// Basic usage with a file:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	lines, errs := tail.Tail(ctx, "/var/log/app.log", tail.DefaultOptions())
//	for line := range lines {
//		fmt.Println(line)
//	}
//	if err := <-errs; err != nil {
//		log.Fatal(err)
//	}
//
// Reading from stdin:
//
//	lines := tail.TailReader(ctx, os.Stdin)
//	for line := range lines {
//		fmt.Println(line)
//	}
//
// Following a file (tail -f behaviour):
//
//	opts := tail.DefaultOptions()
//	opts.Follow = true
//	lines, errs := tail.Tail(ctx, "/var/log/app.log", opts)
//	for line := range lines {
//		fmt.Println(line)
//	}
//
// The context can be cancelled at any time to stop tailing. When the context
// is cancelled, the lines channel is closed and any pending error is sent on
// the errs channel.
package tail
