package test

import (
	"sync"
	"testing"
)

func TestTestWaitGroupMisuse(t *testing.T) {
	for i := 0; i < 1000; i++ {
		TestWaitGroupMisuse(t)
	}
}

// TestWaitGroupMisuse 然后并没有 panic 呢？
func TestWaitGroupMisuse(t *testing.T) {
	defer func() {
		err := recover()
		if err != "sync: negative WaitGroup counter" {
			t.Fatalf("Unexpected panic: %#v", err)
		}
	}()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Done()
	wg.Done()
	t.Fatal("Should panic")
}
