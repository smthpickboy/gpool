package pool

import "github.com/pkg/errors"

const (
	errClosed    = "ERROR: Work Unit added/run after the pool had been closed or cancelled"
)

var (
	ErrPoolClosed = errors.New(errClosed)
)

// ErrRecovery contains the error when a consumer goroutine needed to be recovers
type ErrRecovery struct {
	s string
}

// Error prints recovery error
func (e *ErrRecovery) Error() string {
	return e.s
}

// ErrCancelled is the error returned to a Work Unit when it has been cancelled.
type ErrCancelled struct {
	s string
}

// Error prints Work Unit Cancellation error
func (e *ErrCancelled) Error() string {
	return e.s
}
