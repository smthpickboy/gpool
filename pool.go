package pool

import "context"

// Pool contains all information for a pool instance.
type Pool interface {

	// Queue queues the work to be run, and starts processing immediately
	Queue(ctx context.Context, work WorkFunc, report ReportFunc) WaitFunc

	// Reset reinitializes a pool that has been closed/cancelled back to a working
	// state. if the pool has not been closed/cancelled, nothing happens as the pool
	// is still in a valid running state
	Reset()

	// Close cleans up pool data and cancels any pending work still not committed
	// to processing. Call Reset() to reinitialize the pool for use.
	Close()
}

// WorkFunc is the function type needed by the pool for execution
type WorkFunc func(ctx context.Context) (interface{}, error)
type ReportFunc func(interface{}, error)
type WaitFunc func() (interface{}, error)
