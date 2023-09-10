package test

import (
	"golang/concurrent"
	"sync"
	"testing"
	"time"
)

func TestRWMutexCircleWaite(t *testing.T) {
	concurrent.RWMutexCircleWaite()
}

func TestRWMutex(t *testing.T) {
	var counter concurrent.Counter
	for i := 0; i < 10; i++ { // 10个reader
		go func() {
			for {
				counter.Count() // 计数器读操作
				time.Sleep(time.Millisecond)
			}
		}()
	}
	for { // 一个writer
		counter.Incr() // 计数器写操作
		time.Sleep(time.Second)
	}
}

func TestRUnlock(t *testing.T) {
	var rw sync.RWMutex
	var wg sync.WaitGroup
	wg.Wait()
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			rw.RLock()
			time.Sleep(time.Millisecond * 1000)
			rw.RUnlock() // fatal error: sync: RUnlock of unlocked RWMutex
			t.Log(i)
			wg.Done()
		}(i)
	}
	wg.Add(1)
	//time.Sleep(time.Millisecond * 200)
	go func() {
		time.Sleep(time.Millisecond * 200)
		rw.RUnlock()
		wg.Done()
	}()
	wg.Wait()
}
