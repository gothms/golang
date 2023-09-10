package test

import (
	"fmt"
	"golang/concurrent"
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
	var wg concurrent.WaitGroup
	wg.Add(1)
	for i := 0; i < 20; i++ {
		go func(i int) {
			wg.Add(1)
			//v, w := wg.GetWaitGroupCount()
			//t.Log(v, w)
			t.Log(i)
			time.Sleep(time.Millisecond * 300)
			wg.Done()
			wg.Wait()
		}(i)
	}
	time.Sleep(time.Second * 1)
	go func() {
		wg.Done()
	}()
	time.Sleep(time.Millisecond * 300)

	v, w := wg.GetWaitGroupCount()
	t.Log(v, w)
	wg.Add(2)
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
	t.Log("over")
}
