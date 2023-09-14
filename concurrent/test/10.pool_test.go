package test

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestPoolDequeue(t *testing.T) {
	testPoolDequeue(t, NewPoolDequeue(16))
}

func TestPoolChain(t *testing.T) {
	//t.Log(int64(1) << 32)
	testPoolDequeue(t, NewPoolChain())
}

func testPoolDequeue(t *testing.T, d PoolDequeue) {
	const P = 10
	//var N int = 2e6
	var N int = 10
	if testing.Short() {
		N = 1e3
	}
	have := make([]int32, N)
	var stop int32
	var wg sync.WaitGroup
	record := func(val int) {
		atomic.AddInt32(&have[val], 1)
		if val == N-1 {
			atomic.StoreInt32(&stop, 1)
		}
	}

	// Start P-1 consumers.
	for i := 1; i < P; i++ {
		wg.Add(1)
		go func() {
			fail := 0
			for atomic.LoadInt32(&stop) == 0 {
				val, ok := d.PopTail()
				if ok {
					fail = 0
					record(val.(int))
				} else {
					// Speed up the test by
					// allowing the pusher to run.
					if fail++; fail%100 == 0 {
						runtime.Gosched()
					}
				}
			}
			wg.Done()
		}()
	}

	// Start 1 producer.
	nPopHead := 0
	wg.Add(1)
	go func() {
		for j := 0; j < N; j++ {
			for !d.PushHead(j) {
				// Allow a popper to run.
				runtime.Gosched()
			}
			if j%10 == 0 {
				val, ok := d.PopHead()
				if ok {
					nPopHead++
					record(val.(int))
				}
			}
		}
		wg.Done()
	}()
	wg.Wait()

	// Check results.
	for i, count := range have {
		if count != 1 {
			t.Errorf("expected have[%d] = 1, got %d", i, count)
		}
	}
	// Check that at least some PopHeads succeeded. We skip this
	// check in short mode because it's common enough that the
	// queue will stay nearly empty all the time and a PopTail
	// will happen during the window between every PushHead and
	// PopHead.
	if !testing.Short() && nPopHead == 0 {
		t.Errorf("popHead never succeeded")
	}
}

//var Runtime_Semacquire = runtime_Semacquire
//var Runtime_Semrelease = runtime_Semrelease
//var Runtime_procPin = runtime_procPin
//var Runtime_procUnpin = runtime_procUnpin

// poolDequeue testing.
type PoolDequeue interface {
	PushHead(val any) bool
	PopHead() (any, bool)
	PopTail() (any, bool)
}

func NewPoolDequeue(n int) PoolDequeue {
	d := &poolDequeue{
		vals: make([]eface, n),
	}
	// For testing purposes, set the head and tail indexes close
	// to wrapping around.
	d.headTail = d.pack(1<<dequeueBits-500, 1<<dequeueBits-500)
	return d
}

func (d *poolDequeue) PushHead(val any) bool {
	return d.pushHead(val)
}

func (d *poolDequeue) PopHead() (any, bool) {
	return d.popHead()
}

func (d *poolDequeue) PopTail() (any, bool) {
	return d.popTail()
}

func NewPoolChain() PoolDequeue {
	return new(poolChain)
}

func (c *poolChain) PushHead(val any) bool {
	c.pushHead(val)
	return true
}

func (c *poolChain) PopHead() (any, bool) {
	return c.popHead()
}

func (c *poolChain) PopTail() (any, bool) {
	return c.popTail()
}
