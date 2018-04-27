package pool

import (
	//"sync/atomic"
	"context"
)

// WorkUnit contains a single uint of works values
type WorkUnit interface {
	Wait() (interface{}, error)
}

var _ WorkUnit = new(workUnit)

// workUnit contains a single unit of works values
type workUnit struct {
	ctx    context.Context
	done   chan struct{}
	work   WorkFunc
	report ReportFunc
	v      interface{}
	e      error
}

// Wait blocks until WorkUnit has been processed or cancelled
func (wu *workUnit) Wait() (interface{}, error) {
	<-wu.done
	return nil, nil
}
