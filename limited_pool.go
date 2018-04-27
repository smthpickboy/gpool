package pool

import (
	"context"
	"sync"
)

var _ Pool = new(limitedPool)

// limitedPool contains all information for a limited pool instance.
type limitedPool struct {
	workers uint
	work    chan *workUnit
	cancel  chan struct{}
	closed  bool
	m       sync.RWMutex
}

// NewLimited returns a new limited pool instance
func NewLimited(workers uint) Pool {

	if workers == 0 {
		panic("invalid workers '0'")
	}

	p := &limitedPool{
		workers: workers,
	}

	p.initialize()

	return p
}

func (p *limitedPool) initialize() {

	p.work = make(chan *workUnit, p.workers*2)
	p.cancel = make(chan struct{})
	p.closed = false

	// fire up workers here
	for i := 0; i < int(p.workers); i++ {
		p.newWorker()
	}
}

// passing work and cancel channels to newWorker() to avoid any potential race condition
// betweeen p.work read & write
func (p *limitedPool) newWorker() {
	go func(p *limitedPool) {
		for {
			select {
			case wu := <-p.work:
				// in case work and cancel are closed at the same time
				if wu == nil {
					continue
				}
				wu.v, wu.e = wu.work(wu.ctx)
				if wu.report != nil {
					wu.report(wu.v, wu.e)
				}
				close(wu.done)
			case <-p.cancel:
				return
			}
		}

	}(p)
}

// Queue queues the work to be run, and starts processing immediately
func (p *limitedPool) Queue(ctx context.Context, work WorkFunc, report ReportFunc) WaitFunc {
	w := &workUnit{
		done:   make(chan struct{}),
		work:   work,
		report: report,
		ctx:    ctx,
	}
	p.m.RLock()
	defer p.m.RUnlock()
	if p.closed {
		w.e = ErrPoolClosed
		report(nil, ErrPoolClosed)
	} else {
		select {
		case p.work <- w:
		case <-ctx.Done():
			w.e = ctx.Err()
		}
	}
	return w.Wait
}

// Reset reinitializes a pool that has been closed/cancelled back to a working state.
// if the pool has not been closed/cancelled, nothing happens as the pool is still in
// a valid running state
func (p *limitedPool) Reset() {

	p.m.Lock()

	if !p.closed {
		p.m.Unlock()
		return
	}

	// cancelled the pool, not closed it, pool will be usable after calling initialize().
	p.initialize()
	p.m.Unlock()
}

func (p *limitedPool) closeWithError(err error) {

	p.m.Lock()

	if !p.closed {
		p.closed = true
		close(p.cancel)
		close(p.work)
	}

	for wu := range p.work {
		wu.report(nil, ErrPoolClosed)
	}

	p.m.Unlock()
}

// Close cleans up the pool workers and channels and cancels any pending
// work still yet to be processed.
// call Reset() to reinitialize the pool for use.
func (p *limitedPool) Close() {
	p.closeWithError(ErrPoolClosed)
}
