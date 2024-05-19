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
		操作系统调度，一般等同于 CPU 核数
		M必须持有一个P才可以运行Go代码
	P：Processor
		并不是真正意义的处理器，而是Go语言实现的协程处理器
		维护协程队列
		包含 Go 的必要资源，也有调度goroutine的能力
		P的个数在程序启动时决定，默认情况下等同于CPU的核数（一般情况下M的个数会略大于P的个数，多出来的M将会在G产生系统调用时发挥作用）
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

Go专家编程
	系统调用.jpg
		G0即将进入系统调用时，M0将释放P，进而某个空闲的M1获取P，继续执行P队列中剩下的G
		而M0由于陷入系统调用而进被阻塞，M1接替M0的工作，只要P不空闲，就可以保证充分利用CPU
		M1的来源有可能是M的缓存池，也可能是新建的。当G0系统调用结束后，根据M0是否能获取到P，将会将G0做不同的处理
			1.如果有空闲的P，则获取一个P，继续执行G0
			2.如果没有空闲的P，则将G0放入全局队列，等待被其他的P调度。然后M0将进入缓存池睡眠
		GOMAXPROCS设置对性能的影响
			一般来讲，程序运行时就将GOMAXPROCS大小设置为CPU核数，可让Go程序充分利用CPU
			在某些IO密集型的应用里，这个值可能并不意味着性能最好。理论上当某个Goroutine进入系统调用时，会有一个新的M被启用或创建，继续占满CPU
			但由于Go调度器检测到M被阻塞是有一定延迟的，也即旧的M被阻塞和新的M得到运行之间是有一定间隔的
			所以在IO密集型应用中不妨把GOMAXPROCS设置的大一些，或许会有好的效果
	工作量窃取.jpg
		竖线左侧中右边的P已经将G全部执行完，然后去查询全局队列，全局队列中也没有G，而另一个M中除了正在运行的G外，队列中还有3个G待运行
		此时，空闲的P会将其他P中的G偷取一部分过来，一般每次偷取一半
		偷取完如右图所示

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
