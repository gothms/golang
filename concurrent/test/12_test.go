package test

import (
	"fmt"
	"golang/concurrent"
	"sync"
	"testing"
)

func TestLKQueue(t *testing.T) {
	const n = 150
	lk := concurrent.NewLKQueue()
	var wg sync.WaitGroup
	wg.Add(n * 10)
	for i := 1; i <= n; i++ {
		//time.Sleep(time.Millisecond * 10)
		for j := 0; j < 10; j++ {
			go func(i, j int) {
				cnt := lk.Enqueue(i + j)
				if cnt > 0 {
					fmt.Println(cnt)
				}
				//time.Sleep(time.Millisecond * 10)
				wg.Done()
			}(i*10, j)
		}
	}
	const m = 3
	wg.Add(m)
	for i := 0; i < m; i++ {
		go func(i int) {
			lk.EnqueueLone(9999 + i)
			//lk.Range(nil)
			//cnt, n := lk.EnqueueLone(9999 + i)
			//fmt.Println("============")
			//lk.Range(nil)
			//fmt.Println("------------")
			//fmt.Println(cnt)
			//lk.Range(lone)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestAtomicValueConfig(t *testing.T) {
	concurrent.AtomicValueConfig()
}
