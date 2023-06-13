package advanced

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

/*
Thread vs. Groutine
	1.创建时默认的 stack 的大小
		JDK5 以后 Java Thread stack 默认为 1M
		Groutine 的 stack 初始化大小为 2k
	2.和 KSE-Kernel Space Entity（内核对象/系统线程）的对应关系
		Java Thread 是 1:1
		Groutine 是 M:N，多对多

1.共享内存并发机制
	sync.Mutex
	sync.WaitGroup
	sync.RWMutex：RLock() RUnlock()
2.CSP vs Actor
	和Actor的直接通讯不同，CSP模式则是通过Channel进行通讯的，更松耦合一些
	Go中Channel有容量限制并且独立于处理Groutine
		而如Erlang，Actor模式中的mailbox容量是无限的
		接收进程也总是被动的处理消息
*/
// 共享内存并发机制
func TestCounterWaitGroup(t *testing.T) {
	var mut sync.Mutex
	var wg sync.WaitGroup
	counter := 0
	wg.Add(5000)
	for i := 0; i < 5000; i++ {
		//wg.Add(1)
		go func() {
			defer mut.Unlock()
			mut.Lock()
			counter++
			wg.Done()
		}()
	}
	wg.Wait() // 代替 time.Sleep(time.Second)
	fmt.Println(counter)
}

// sync.Mutex
func TestCounterThreadSafe(t *testing.T) {
	var mut sync.Mutex
	counter := 0
	for i := 0; i < 5000; i++ {
		go func() {
			defer mut.Unlock() // 释放锁
			mut.Lock()
			counter++
		}()
	}
	time.Sleep(time.Second)
	t.Log(counter)
}

func TestGroutine(t *testing.T) {
	for i := 0; i < 10; i++ {
		//go func() {
		//	fmt.Println(i)	// i 为共享数据，被竞争
		//}()
		go func(i int) {
			fmt.Println(i)
		}(i)
	}
	time.Sleep(time.Millisecond * 50)
}
