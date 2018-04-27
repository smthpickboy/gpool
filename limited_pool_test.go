package pool

import (
	"testing"
	"time"

	"context"
	"fmt"
	"sync/atomic"

	. "gopkg.in/go-playground/assert.v1"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called
//
// or
//
// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html
//

func TestPool(t *testing.T) {

	var res []WaitFunc

	pool := NewLimited(4)
	defer pool.Close()

	newFunc := func(d time.Duration) WorkFunc {
		return func(context.Context) (interface{}, error) {
			time.Sleep(d)
			return nil, nil
		}
	}

	reportCount := int64(0)
	report := func(v interface{}, err error) {
		atomic.AddInt64(&reportCount, 1)
	}

	for i := 0; i < 4; i++ {
		wu := pool.Queue(context.Background(), newFunc(time.Second*1), report)
		res = append(res, wu)
	}

	var count int

	for i, wu := range res {
		fmt.Println(i)
		v, e := wu()
		Equal(t, e, nil)
		Equal(t, v, nil)
		count++
	}

	Equal(t, count, 4)
	Equal(t, reportCount, int64(4))

	pool.Close() // testing no error occurs as Close will be called twice once defer pool.Close() fires
}
