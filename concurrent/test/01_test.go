package test

import (
	"golang/concurrent"
	"sync"
	"testing"
)

func TestMutexPractice(t *testing.T) {
	concurrent.MutexPractice()
}

func TestMutexConcurrent(t *testing.T) {
	const GCount = 10
	var (
		cnt int
		mut sync.Mutex     // 互斥锁保护计数器
		wg  sync.WaitGroup // 使用WaitGroup等待10个goroutine完成
	)
	wg.Add(GCount)
	for i := 0; i < GCount; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100_000; j++ { // 对变量count执行10次加1
				mut.Lock()
				cnt++
				mut.Unlock()
			}
		}()
	}
	wg.Wait() // 等待10个goroutine完成
	t.Log(cnt)
}

// for j := 0; j < 100_000; j++ {
// mut.Lock()	// 锁 cnt++
// BenchmarkMutex-8              31          45788819 ns/op
// mut.Lock()	// 锁 for 循环
// for j := 0; j < 100_000; j++ {
// BenchmarkMutex-8             902           1400905 ns/op
func BenchmarkMutex(b *testing.B) {
	const GCount = 10
	var (
		cnt int
		mut sync.Mutex     // 互斥锁保护计数器
		wg  sync.WaitGroup // 使用WaitGroup等待10个goroutine完成
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(GCount)
		for idx := 0; idx < GCount; idx++ {
			go func() {
				defer wg.Done()
				for j := 0; j < 100_000; j++ { // 对变量count执行10次加1
					mut.Lock()
					cnt++
					mut.Unlock()
				}
			}()
		}
		wg.Wait() // 等待10个goroutine完成
	}
	b.StopTimer()
	b.Log(cnt)
}

// TestMutex 创建了 10 个 goroutine，同时不断地对一个变量（count）进行加 1 操作，每个 goroutine 负责执行 10 万次的加 1 操作
// 我们期望的最后计数的结果是 10 * 100000 = 1000000 (一百万)
func TestMutex(t *testing.T) {
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
	t.Log(cnt)
}
