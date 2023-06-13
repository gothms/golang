package advanced

import (
	"fmt"
	"testing"
	"time"
)

// cancel_01 共享内存并发机制
func cancel_01(ch chan struct{}) {
	ch <- struct{}{}
}

// close(chan) 广播
func cancel_02(ch chan struct{}) {
	close(ch)
}
func isCancelled(ch chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}
func TestCancle(t *testing.T) {
	ch := make(chan struct{}, 0)
	for i := 0; i < 5; i++ {
		go func(i int, ch chan struct{}) {
			//fmt.Printf("%p\n", ch) // 同一个 ch
			for {
				if isCancelled(ch) {
					break
				}
				time.Sleep(time.Millisecond * 50)
			}
			fmt.Println(i, "cancelled")
		}(i, ch)
	}
	//cancel_01(ch) // 共享内存并发机制
	cancel_02(ch) // close(chan) 广播
	time.Sleep(time.Second * 1)
}
