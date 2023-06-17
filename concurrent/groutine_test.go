package concurrent

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
			调度效率非常高，但线程切换牵涉到内核对象的切换，消耗很大
		Groutine 是 M:N，多对多
			多个协程可对应一个内核对象，切换的消耗小很多
Go 协程调度机制
	M：System Thread
		系统线程，也就是Kernel Entity
	P：Processor
		并不是真正意义的处理器，而是Go语言实现的协程处理器
		维护协程队列
	G：Goroutine
		协程
注意点
	1.守护线程会计数每个 Processor 完成的 G 的数量
		如果某段时间完成的数量没有发生变化，会往这个协程的任务栈里插入特殊标记
		当这个协程运行到这个标记时，就会中断并插到队列的队尾，把资源让给别的 G
	2.当某个协程被系统中断（比如IO）需要等待时，为了提高整体并发
		Processor 会把自己移动到另一个可使用的系统线程 M 中，继续执行 G 队列
		当这个被中断的协程被唤醒完成之后，会把自己加到某个 Processor 的队列里，或全局等待队列中
	3.某个协程被中断时，它在寄存器里的运行状态也会保存在协程对象里
		再次运行时这些状态会重新写入寄存器，继续运行

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
