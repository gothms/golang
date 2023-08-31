package test

import (
	"golang/concurrent"
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
