package concurrent

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"log"
	"runtime"
	"sync"
	"time"
)

/*
Semaphore：一篇文章搞懂信号量

80%
	标准库的并发原语、原子操作和 Channel，能解决 80% 的并发问题
	进一步提升你的并发编程能力，需要学习一些第三方库
信号量（Semaphore）
	用来控制多个 goroutine 同时访问多个资源的并发原语

信号量是什么？都有什么操作？
	简介
		信号量的概念是荷兰计算机科学家 Edsger Dijkstra 在 1963 年左右提出来的，广泛应用在不同的操作系统中
		在系统中，会给每一个进程一个信号量，代表每个进程目前的状态
		未得到控制权的进程，会在特定的地方被迫停下来，等待可以继续进行的信号到
	最简单的信号量就是一个变量加一些并发控制的能力，这个变量是 0 到 n 之间的一个数值
		当 goroutine 完成对此信号量的等待（wait）时，该计数值就减 1，当 goroutine 完成对此信号量的释放（release）时，该计数值就加 1
		当计数值为 0 的时候，goroutine 调用 wait 等待该信号量是不会成功的，除非计数器又大于 0，等待的 goroutine 才有可能成功返回
	更复杂的信号量类型，就是使用抽象数据类型代替变量，用来代表复杂的资源类型
		实际上，大部分的信号量都使用一个整型变量来表示一组资源，并没有实现太复杂的抽象数据类型
	举例
		图书馆新购买了 10 本《Go 并发编程的独家秘籍》，有 1 万个学生都想读这本书
		所以，图书馆管理员先会让这 1 万个同学进行登记，按照登记的顺序，借阅此书
		如果书全部被借走，那么，其他想看此书的同学就需要等待，如果有人还书了，图书馆管理员就会通知下一位同学来借阅这本书
		这里的资源是《Go 并发编程的独家秘籍》这十本书，想读此书的同学就是 goroutine，图书管理员就是信号量
P/V 操作
	Dijkstra 在他的论文中为信号量定义了两个操作 P 和 V
		P 操作（descrease、wait、acquire）是减少信号量的计数值
		而 V 操作（increase、signal、release）是增加信号量的计数值
	伪代码
		中括号代表原子操作
			function V(semaphore S, integer I):
				[S ← S + I]
			function P(semaphore S, integer I):
				repeat:
					[if S ≥ I:
					S ← S − I
					break]
		信号量的值除了初始化的操作以外，只能由 P/V 操作改变
			初始化信号量 S 有一个指定数量（n）的资源，它就像是一个有 n 个资源的池子
			P 操作相当于请求资源，如果资源可用，就立即返回；如果没有资源或者不够，那么，它可以不断尝试或者阻塞等待
			V 操作会释放自己持有的资源，把资源返还给信号量
	信号量的实现
		初始化信号量：设定初始的资源的数量
		P 操作：将信号量的计数值减去 1，如果新值已经为负，那么调用者会被阻塞并加入到等待队列中
			否则，调用者会继续执行，并且获得一个资源
		V 操作：将信号量的计数值加 1，如果先前的计数值为负，就说明有等待的 P 操作的调用者
			它会从等待队列中取出一个等待的调用者，唤醒它，让它继续执行
	饥饿
		如同 Mutex 饥饿问题，在高并发的极端场景下，会有些 goroutine 始终抢不到锁
		为了处理饥饿的问题，你可以在等待队列中做一些“文章”
		比如实现一个优先级的队列，或者先入先出的队列，等等，保持公平性，并且照顾到优先级
	信号量 vs 互斥锁
		其实，信号量可以分为计数信号量（counting semaphre）和二进位信号量（binary semaphore）
			图书馆借书的例子就是一个计数信号量，它的计数可以是任意一个整数
			在特殊的情况下，如果计数值只能是 0 或者 1，那么，这个信号量就是二进位信号量，提供了互斥的功能（要么是 0，要么是 1）
			所以，有时候互斥锁也会使用二进位信号量来实现
		二进制信号量
			我们一般用信号量保护一组资源，比如数据库连接池、一组客户端的连接、几个打印机资源，等
			如果信号量蜕变成二进位信号量，那么，它的 P/V 就和互斥锁的 Lock/Unlock 一样了
		“区分”二进位信号量和互斥锁
			在 Windows 系统中，互斥锁只能由持有锁的线程释放锁，而二进位信号量则没有这个限制（Stack Overflow 上也有相关的讨论）
			实际上，虽然在 Windows 系统中，它们的确有些区别
			但是对 Go 语言来说，互斥锁也可以由非持有的 goroutine 来释放，所以，从行为上来说，它们并没有严格的区别
		建议
			没必要进行细致的区分，因为互斥锁并不是一个很严格的定义
			实际在遇到互斥并发的问题时，我们一般选用互斥锁

Go 官方扩展库的实现
	Mutex
		在运行时，Go 内部使用信号量来控制 goroutine 的阻塞和唤醒
		比如互斥锁的第二个字段：
			type Mutex struct {
				state int32
				sema  uint32
			}
		信号量的 P/V 操作是通过函数实现的
			func runtime_Semacquire(s *uint32)
			func runtime_SemacquireMutex(s *uint32, lifo bool, skipframes int)
			func runtime_Semrelease(s *uint32, handoff bool, skipframes int)
		它是 Go 运行时内部使用的
			并没有封装暴露成一个对外的信号量并发原语，原则上我们没有办法使用
	semaphore
		cmd/vendor/golang.org/x/sync/semaphore/semaphore.go
		Go 在它的扩展包中提供了信号量 semaphore，不过这个信号量的类型名并不叫 Semaphore，而是叫 Weighted
		它可以在初始化创建这个信号量的时候设置权重（初始化的资源数）
	Weighted
		type Weighted
			func NewWeighted(n int64) *Weighted
			func (s *Weighted) Acquire(ctx context.Context, n int64) error
			func (s *Weighted) Release(n int64)
			func (s *Weighted) TryAcquire(n int64) bool
	API
		1. Acquire 方法：相当于 P 操作，你可以一次获取多个资源，如果没有足够多的资源，调用者就会被阻塞
			它的第一个参数是 Context，这就意味着，你可以通过 Context 增加超时或者 cancel 的机制
			如果是正常获取了资源，就返回 nil；否则，就返回 ctx.Err()，信号量不改变
		2. Release 方法：相当于 V 操作，可以将 n 个资源释放，返还给信号量
		3. TryAcquire 方法：尝试获取 n 个资源，但是它不会阻塞
			要么成功获取 n 个资源，返回 true，要么一个也不获取，返回 false
	示例：Worker Pool
		go get golang.org/x/sync
			创建和 CPU 核数一样多的 Worker，让它们去处理一个 4 倍数量的整数 slice
			每个Worker 一次只能处理一个整数，处理完之后，才能处理下一个
		SemaWorkerPool
			main goroutine 相当于一个 dispacher，负责任务的分发
			它先请求信号量，如果获取成功，就会启动一个 goroutine 去处理计算，然后，这个 goroutine 会释放这个信号量
			有意思的是，信号量的获取是在 main goroutine，信号量的释放是在 worker goroutine 中
			如果获取不成功，就等到有信号量可以使用的时候，再去获取
		技巧
			如果在实际应用中，你想等所有的 Worker 都执行完，就可以获取最大计数值的信号量
	Go 扩展库中的信号量是使用互斥锁 + List 实现的
		互斥锁实现其它字段的保护，而 List 实现了一个等待队列，等待者的通知是通过 Channel 的通知机制实现的
	信号量 Weighted 的数据结构
		type Weighted struct {
			size    int64      // 最大资源数
			cur     int64      // 当前已被使用的资源
			mu      sync.Mutex // 互斥锁，对字段的保护
			waiters list.List  // 等待队列
		}
	Acquire 源码
		它不仅仅要监控资源是否可用，而且还要检测 Context 的 Done 是否已关闭
		其实，为了提高性能，这个方法中的 fast path 之外的代码，可以抽取成 acquireSlow 方法，以便其它 Acquire 被内联
	Release 方法
		将当前计数值减去释放的资源数 n，并唤醒等待队列中的调用者，看是否有足够的资源被获取
	notifyWaiters 方法
		逐个检查等待的调用者，如果资源不够，或者是没有等待者了，就返回
		notifyWaiters 方法是按照先入先出的方式唤醒调用者
		避免饥饿：
		当释放 100 个资源的时候，如果第一个等待者需要 101 个资源，那么，队列中的所有等待者都会继续等待，即使有的等待者只需要 1 个资源
		这样做的目的是避免饥饿，否则的话，资源可能总是被那些请求资源数小的调用者获取
		这样一来，请求资源数巨大的调用者，就没有机会获得资源了

使用信号量的常见错误
	保证信号量不出错的前提是正确地使用它，否则，公平性和安全性就会受到损害，导致程序 panic
	最常见的几个错误
		请求了资源，但是忘记释放它
		释放了从未请求的资源
		长时间持有一个资源，即使不需要它
		不持有一个资源，却直接使用它
	死锁
		即使你规避了这些坑，在同时使用多种资源，不同的信号量控制不同的资源的时候，也可能会出现死锁现象
		比如哲学家就餐问题：https://zh.wikipedia.org/wiki/%E5%93%B2%E5%AD%A6%E5%AE%B6%E5%B0%B1%E9%A4%90%E9%97%AE%E9%A2%98
	参数错误
		Go 扩展库实现的信号量，在调用 Release 方法的时候，你可以传递任意的整数
			但是，如果你传递一个比请求到的数量大的错误的数值，程序就会 panic
			如果传递一个负数，会导致资源永久被持有
		如果你请求的资源数比最大的资源数还大，那么，调用者可能永远被阻塞
	使用信号量遵循的原则
		请求多少资源，就释放多少资源
		必须使用正确的方法传递整数，不要“耍小聪明”，而且，请求的资源数一定不要超过最大资源数

其它信号量的实现
	除了官方扩展库的实现，实际上，我们还有很多方法实现信号量，比较典型的就是使用 Channel 来实现
	示例：ChannelSemaphore 数据结构，并且还实现了Locker接口
		使用一个 buffer 为 n 的 Channel 很容易实现信号量
			在初始化这个信号量的时候，我们设置它的初始容量，代表有多少个资源可以使用
			它使用 Lock 和 Unlock 方法实现请求资源和释放资源，正好实现了 Locker 接口
		扩展方法
			在请求资源的时候使用 Context 参数（Acquire(ctx)）、实现 TryLock 等功能
	疑问
		这个信号量的实现看起来非常简单，而且也能应对大部分的信号量的场景，为什么官方扩展库的信号量的实现不采用这种方法呢？
		官方实现的特色功能：它可以一次请求多个资源，这是通过 Channel 实现的信号量所不具备的
	marusama/semaphore
		实现了一个可以动态更改资源容量的信号量，也是一个非常有特色的实现
		如果你的资源数量并不是固定的，而是动态变化的，建议你考虑一下这个信号量库

总结
	奇怪现象
		标准库中实现基本并发原语（比如 Mutex）的时候，强烈依赖信号量实现等待队列和通知唤醒
		但是，标准库中却没有把这个实现直接暴露出来放到标准库，而是通过第三库提供
	信号量这个并发原语在多资源共享的并发控制的场景中被广泛使用
		有时候也会被 Channel 类型所取代，因为一个 buffered chan 也可以代表 n 个资源
	官方扩展的信号量的优势
		可以一次获取多个资源
		在批量获取资源的场景中，建议尝试使用官方扩展的信号量

思考
	1. 你能用 Channel 实现信号量并发原语吗？你能想到几种实现方式？
	2. 为什么信号量的资源数设计成 int64 而不是 uint64 呢？
*/

// ChannelSemaphore 数据结构，并且还实现了Locker接口 =========Channel 实现 Semaphore=========
type ChannelSemaphore struct {
	//sync.Locker
	ch chan struct{}
}

func NewSemaphore(capacity int) sync.Locker { // 创建一个新的信号量
	if capacity <= 0 {
		capacity = 1 // 容量为1就变成了一个互斥锁
	}
	return &ChannelSemaphore{ch: make(chan struct{}, capacity)}
}
func (s *ChannelSemaphore) Lock() { // 请求一个资源
	s.ch <- struct{}{}
}
func (s *ChannelSemaphore) Unlock() { // 释放资源
	<-s.ch
}

// SemaWorkerPool ==========Worker Pool 示例==========
func SemaWorkerPool() {
	var (
		maxWorkers = runtime.GOMAXPROCS(0)                    // worker数量
		sema       = semaphore.NewWeighted(int64(maxWorkers)) // 信号量
		task       = make([]int, maxWorkers*4)                // 任务数，是worker的四倍
	)
	ctx := context.Background()
	for i := range task {
		// 如果没有worker可用，会阻塞在这里，直到某个worker被释放
		if err := sema.Acquire(ctx, 1); err != nil {
			break
		}
		// 启动worker goroutine
		go func(i int) {
			defer sema.Release(1)
			time.Sleep(100 * time.Millisecond) // 模拟一个耗时操作
			task[i] = i + 1
		}(i)
	}
	// 请求所有的worker,这样能确保前面的worker都执行完
	if err := sema.Acquire(ctx, int64(maxWorkers)); err != nil {
		log.Printf("获取所有的worker失败: %v", err)
	}
	fmt.Println(task)
}
