package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println(0 &^ 2)

	const GCount = 10
	var (
		cnt int
		wg  sync.WaitGroup // 使用WaitGroup等待10个goroutine完成
	)
	wg.Add(GCount)
	for i := 0; i < GCount; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100_000; j++ { // 对变量count执行10次加1
				cnt++
			}
		}()
	}
	wg.Wait() // 等待10个goroutine完成
	fmt.Println(cnt)
}
