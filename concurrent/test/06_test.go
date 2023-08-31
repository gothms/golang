package test

import (
	"fmt"
	"golang/concurrent"
	"sync"
	"testing"
	"time"
)

func TestWaitGroupBug(t *testing.T) {
	concurrent.WaitGroupBug()
}
func TestWGCopy(t *testing.T) {
	concurrent.WGCopy()
}
func TestWGWrongAdd(t *testing.T) {
	concurrent.WGWrongAdd()
}
func TestWGCounterMain(t *testing.T) {
	concurrent.WGCounterMain()
}

// 试图复现：panic("sync: WaitGroup misuse: Add called concurrently with Wait")
func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	for i := 0; i < 3; i++ {
		go func() {
			wg.Wait()
		}()
	}
	time.Sleep(time.Millisecond * 300)
	wg.Add(1)
	go func() {
		fmt.Println(".")
		wg.Done()
	}()
	go func() {
		fmt.Println("。")
		wg.Done()
	}()
	//time.Sleep(time.Millisecond * 300)
	wg.Wait()
}
