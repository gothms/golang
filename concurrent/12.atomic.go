package concurrent

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

/*
atomic：要保证原子操作，一定要使用这几种方法

原子操作的基础知识
	atomic 包
		Package sync/atomic 实现了同步算法底层的原子的内存操作原语，我们把它叫做原子操作原语
		它提供了一些实现原子操作的方法
	原子操作特性
		一个原子在执行的时候，其它线程不会看到执行一半的操作结果
		在其它线程看来，原子操作要么执行完了，要么还没有执行，就像一个最小的粒子 - 原子一样，不可分割
	CPU 提供了基础的原子操作，不过，不同架构的系统的原子操作是不一样的
		单处理器单核系统
			如果一个操作是由一个 CPU 指令来实现的，那么它就是原子操作，比如它的 XCHG 和 INC 等指令
			如果操作是基于多条指令来实现的，那么，执行的过程中可能会被中断，并执行上下文切换
			这样的话，原子性的保证就被打破了，因为这个时候，操作可能只执行了一半
		多处理器多核系统，原子操作的实现就比较复杂了
			由于 cache 的存在，单个核上的单个指令进行原子操作的时候，你要确保其它处理器或者核不访问此原子操作的地址，或者是确保其它处理器或者核总是访问原子操作之后的最新的值
			x86 架构中提供了指令前缀 LOCK，LOCK 保证了指令（比如 LOCK CMPXCHG op1、op2）不会受其它处理器或 CPU 核的影响，有些指令（比如 XCHG）本身就提供 Lock 的机制
			不同的 CPU 架构提供的原子操作指令的方式也是不同的，比如对于多核的 MIPS 和 ARM，提供了 LL/SC（Load Link/Store Conditional）指令，可以帮助实现原子操作（ARMLL/SC 指令 LDREX 和 STREX）
		因为不同的 CPU 架构甚至不同的版本提供的原子操作的指令是不同的
			所以，要用一种编程语言实现支持不同架构的原子操作是相当有难度的
	Go
		Go 提供了一个通用的原子操作的 API，将更底层的不同的架构下的实现封装成 atomic 包
		提供了修改类型的原子操作（atomic read-modify-write，RMW）和加载存储类型的原子操作（Load 和 Store）的 API
		如果要想保证原子操作，切记一定要使用 atomic 提供的方法
	有的代码也会因为架构的不同而不同
		有时看起来貌似一个操作是原子操作，但实际上，对于不同的架构来说，情况是不一样的
		示例
			const x int64 = 1 + 1<<33
			func atomicDemo() {
				var i = x
				_ = i
			}
		var i = x：将一个 64 位的值赋值给变量 i
		_ = i：
			如果你使用 GOARCH=386 的架构去编译这段代码，那么，'_ = i' 其实是被拆成了两个指令，分别操作低 32 位和高 32 位
				使用 GOARCH=386 go tool compile -N -l test.go；GOARCH=386 go tool objdump -gnu test.o 反编译试试
			如果 GOARCH=amd64 的架构去编译这段代码，那么，'_ = i' 其中的赋值操作其实是一条指令
				12.atomic_compile_goarch_amd64.jpg

atomic 原子操作的应用场景
	使用 atomic 的一些方法，我们可以实现更底层的一些优化
		如果使用 Mutex 等并发原语进行这些优化，虽然可以解决问题，但是这些并发原语的实现逻辑比较复杂，对性能还是有一定的影响的
	应用场景
		简单原子操作
		使用 atomic 实现自己定义的基本并发原语
		atomic 原子操作还是实现 lock-free 数据结构的基石
	举例一
		描述
			假设你想在程序中使用一个标志（flag，比如一个 bool 类型的变量），来标识一个定时任务是否已经启动执行了
		加锁方法
			如果使用 Mutex 和 RWMutex，在读取和设置这个标志的时候加锁，是可以做到互斥的、保证同一时刻只有一个定时任务在执行的
			所以使用 Mutex 或者 RWMutex 是一种解决方案
		原子操作
			这个场景中的问题不涉及到对资源复杂的竞争逻辑，只是会并发地读写这个标志，这类场景就适合使用 atomic 的原子操作
			可以使用一个 uint32 类型的变量
			如果这个变量的值是 0，就标识没有任务在执行，如果它的值是 1，就标识已经有任务在完成了
	举例二
		描述
			在开发应用程序的时候，需要从配置服务器中读取一个节点的配置信息
			而且，在这个节点的配置发生变更的时候，你需要重新从配置服务器中拉取一份新的配置并更新
		加锁
			程序中可能有多个 goroutine 都依赖这份配置，涉及到对这个配置对象的并发读写，你可以使用读写锁实现对配置对象的保护
		原子操作
			在大部分情况下，你也可以利用 atomic 实现配置对象的更新和加载
	使用 atomic 实现自己定义的基本并发原语
		比如 Go issue 有人提议的 CondMutex、Mutex.LockContext、WaitGroup.Go 等，我们可以使用 atomic 或者基于它的更高一级的并发原语去实现
		前面的几种基本并发原语的底层（比如 Mutex），就是基于通过 atomic 的方法实现的
	atomic 原子操作还是实现 lock-free 数据结构的基石
		在实现 lock-free 的数据结构时，我们可以不使用互斥锁，这样就不会让线程因为等待互斥锁而阻塞休眠，而是让线程保持继续处理的状态
		另外，不使用互斥锁的话，lock-free 的数据结构还可以提供并发的性能
	lock-free
		lock-free 的数据结构实现起来比较复杂，需要考虑的东西很多
		一位微软专家写的一篇经验分享：Lockless Programming Considerations for Xbox 360 and Microsoft Windows
		link：

atomic 提供的方法
	sync/atomic/doc.go
		atomic 为了支持 int32、int64、uint32、uint64、uintptr、Pointer（Add 方法不支持）类型
		分别提供了 AddXXX、CompareAndSwapXXX、SwapXXX、LoadXXX、StoreXXX 等方法
	切记
		atomic 操作的对象是一个地址，你需要把可寻址的变量的地址作为参数传递给方法
		而不是把变量的值传递给方法
	Add
		方法签名
			func AddInt32(addr *int32, delta int32) (new int32)
			func AddUint32(addr *uint32, delta uint32) (new uint32)
			func AddInt64(addr *int64, delta int64) (new int64)
			func AddUint64(addr *uint64, delta uint64) (new uint64)
			func AddUintptr(addr *uintptr, delta uintptr) (new uintptr)
		Add 方法就是给第一个参数地址中的值增加一个 delta 值
			对于有符号的整数来说，delta 可以是一个负数，相当于减去一个值
			对于无符号的整数和 uinptr 类型来说，可以利用计算机补码的规则，把减法变成加法
				AddUint32(&x, ^uint32(c-1))
				AddUint32(&x, ^uint32(0))	// 减 1 操作的简化
	CAS （CompareAndSwap）
		在 CAS 的方法签名中，需要提供要操作的地址、原数据值、新值
			func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
		相当于“判断相等才替换”
			比较当前 addr 地址里的值是不是 old，如果不等于 old，就返回 false
			如果等于 old，就把此地址的值替换成 new 值，返回 true
	Swap
		如果不需要比较旧值，只是比较粗暴地替换的话，就可以使用 Swap 方法，它替换后还可以返回旧值
			func SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer)
	Load
		取出 addr 地址中的值，即使在多处理器、多核、有 CPU cache 的情况下，这个操作也能保证 Load 是一个原子操作
			func LoadPointer(addr *unsafe.Pointer) (val unsafe.Pointer)
	Store
		把一个值存入到指定的 addr 地址中，即使在多处理器、多核、有 CPU cache 的情况下，这个操作也能保证 Store 是一个原子操作
		别的 goroutine 通过 Load 读取出来，不会看到存取了一半的值
			func StoreUintptr(addr *uintptr, val uintptr)
Value 类型
	特殊的类型 Value：它可以原子地存取对象类型，常常用在配置变更等场景中
		type Value
			func (v *Value) Load() (val any)
			func (v *Value) Store(val any)
			func (v *Value) Swap(new any) (old any)
			func (v *Value) CompareAndSwap(old, new any)
	示例：配置变更
		AtomicValueConfig & TestAtomicValueConfig
		启动一个 goroutine 等待配置变更的信号，一旦有变更，它就会加载最新的配置

第三方库的扩展
	atomic 的 API 已经算是很简单的了，它提供了包一级的函数，可以对几种类型的数据执行原子操作
	uber-go/atomic：提供了面向对象的使用方式
		定义和封装了几种与常见类型相对应的原子操作类型，这些类型提供了原子操作的方法
		这些类型包括 Bool、Duration、Error、Float64、Int32、Int64、String、Uint32、Uint64 等
		比如 Bool 类型，提供了 CAS、Store、Swap、Toggle 等原子方法
		还提供 String、MarshalJSON、UnmarshalJSON 等辅助方法，确实是一个精心设计的 atomic 扩展库
	示例
		var running atomic.Bool
		running.Store(true)
		running.Toggle()
		fmt.Println(running.Load()) // false

使用 atomic 实现 Lock-Free queue
	Lock-Free queue 论文
		Lock-Free queue 最出名的就是 Maged M. Michael 和 Michael L. Scott 1996 年发表的论文中的算法
		算法比较简单，容易实现，伪代码的每一行都提供了注释
		论文 link：
	Go 实现
		LKQueue ==========lock-free queue==========
		问题：Dequeue 后，为什么不把 q.head 置为 nil？

总结
	对一个地址的赋值是原子操作吗？
		如果是原子操作，还要 atomic 包干什么？
		官方的文档中并没有特意的介绍，不过在一些 issue 或者论坛中，每当有人谈到这个问题时，总是会被建议用 atomic 包
	Dave Cheney 就谈到过这个问题，讲得非常好（总结一下他讲的知识点，这样你就比较容易理解使用 atomic 和直接内存操作的区别了）
		在现在的系统中，write 的地址基本上都是对齐的（aligned）
			比如，32 位的操作系统、CPU 以及编译器，write 的地址总是 4 的倍数，64 位的系统总是 8 的倍数（如 WaitGroup 针对 64 位系统和 32 位系统对 state1 的字段不同的处理）
			对齐地址的写，不会导致其他人看到只写了一半的数据，因为它通过一个指令就可以实现对地址的操作
			如果地址不是对齐的话，那么，处理器就需要分成两个指令去处理，如果执行了一个指令，其它人就会看到更新了一半的错误的数据，这被称做撕裂写（torn write）
			所以，你可以认为赋值操作是一个原子操作，这个“原子操作”可以认为是保证数据的完整性
		多处理多核系统
			但是，对于现代的多处理多核的系统来说，由于 cache、指令重排，可见性等问题，我们对原子操作的意义有了更多的追求
			在多核系统中，一个核对地址的值的更改，在更新到主内存中之前，是在多级缓存中存放的
			这时，多个核看到的数据可能是不一样的，其它的核可能还没有看到更新的数据，还在使用旧的数据
		内存屏障（memory fence 或 memory barrier）
			多处理器多核心系统为了处理这类问题，使用了一种叫做内存屏障（memory fence 或 memory barrier）的方式
			一个写内存屏障会告诉处理器，必须要等到它管道中的未完成的操作（特别是写操作）都被刷新到内存中，再进行操作
			此操作还会让相关的处理器的 CPU 缓存失效，以便让它们从主存中拉取最新的值
		atomic 包
			atomic 包提供的方法会提供内存屏障的功能
			所以，atomic 不仅仅可以保证赋值的数据完整性，还能保证数据的可见性，一旦一个核更新了该地址的值，其它处理器总是能读取到它的最新值
			但是，需要注意的是，因为需要处理器之间保证数据的一致性，atomic 的操作也是会降低性能的

思考
	atomic.Value 只有 Load/Store 方法，你是不是感觉意犹未尽？你可以尝试为 Value 类型增加 Swap 和 CompareAndSwap 方法
	link：
*/

// LKQueue ==========lock-free queue==========
type LKQueue struct {
	head unsafe.Pointer // 辅助头指针，头指针不包含有意义的数据，只是一个辅助的节点
	tail unsafe.Pointer
}
type node struct { // 通过链表实现，这个数据结构代表链表中的节点
	value interface{}
	next  unsafe.Pointer
}

func NewLKQueue() *LKQueue {
	n := unsafe.Pointer(&node{})
	return &LKQueue{head: n, tail: n}
}

// Enqueue 入队
// 入队的时候，通过 CAS 操作将一个元素添加到队尾，并且移动尾指针
func (q *LKQueue) Enqueue(v interface{}) {
	n := &node{value: v}
	for {
		tail := load(&q.tail)
		next := load(&tail.next)
		if tail == load(&q.tail) { // 尾还是尾
			if next == nil { // 还没有新数据入队
				if cas(&tail.next, next, n) { //增加到队尾
					cas(&q.tail, tail, n) //入队成功，移动尾巴指针
					return
				}
			} else { // 已有新数据加到队列后面，需要移动尾指针
				cas(&q.tail, tail, next)
			}
		}
	}
}

// Dequeue 出队，没有元素则返回nil
// 出队的时候移除一个节点，并通过 CAS 操作移动 head 指针，同时在必要的时候移动尾指针
func (q *LKQueue) Dequeue() interface{} {
	for {
		head := load(&q.head)
		tail := load(&q.tail)
		next := load(&head.next)
		if head == load(&q.head) { // head还是那个head
			if head == tail { // head和tail一样
				if next == nil { // 说明是空队列
					return nil
				}
				// 只是尾指针还没有调整，尝试调整它指向下一个
				cas(&q.tail, tail, next)
			} else {
				// 读取出队的数据
				v := next.value
				// 既然要出队了，头指针移动到下一个
				if cas(&q.head, head, next) {
					return v // Dequeue is done. return
				}
			}
		}
	}
}
func load(p *unsafe.Pointer) (n *node) { // 将unsafe.Pointer原子加载转换成node
	return (*node)(atomic.LoadPointer(p))
}
func cas(p *unsafe.Pointer, old, new *node) (ok bool) { // 封装CAS,避免直接将*node转换成unsafe.Pointer
	return atomic.CompareAndSwapPointer(
		p, unsafe.Pointer(old), unsafe.Pointer(new))
}

// AtomicValueConfig ==========atomic.Value 配置变更示例==========
func AtomicValueConfig() {
	var config atomic.Value
	config.Store(loadNewConfig())
	var cond = sync.NewCond(&sync.Mutex{}) // 不要一直读
	go func() {                            // 设置新的config
		for {
			time.Sleep(time.Duration(5+rand.Int63n(5)) * time.Second)
			config.Store(loadNewConfig())
			cond.Broadcast() // 通知等待着配置已变更
		}
	}()
	go func() {
		//var last Config
		for {
			cond.L.Lock()
			cond.Wait()                 // 等待变更信号：不要一直读
			c := config.Load().(Config) // 读取新的配置
			fmt.Printf("new config: %+v\n %[1]T\n", c)
			cond.L.Unlock()

			//if c := config.Load().(Config); c != last {
			//	last = c
			//	fmt.Printf("new config: %+v\n", last)
			//}
			//time.Sleep(time.Duration(3+rand.Int63n(3)) * time.Second)
		}
	}()
	select {}
}

type Config struct {
	NodeName string
	Addr     string
	Count    int32
}

func loadNewConfig() Config {
	return Config{
		NodeName: "北京",
		Addr:     "10.77.95.27",
		Count:    rand.Int31(),
	}
}
