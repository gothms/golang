package concurrent

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

/*
WaitGroup：协同等待，任务编排利器

简介
	WaitGroup 很简单，就是 package sync 用来做任务编排的一个并发原语
它要解决的就是并发 - 等待的问题
	现在有一个 goroutine A 在检查点（checkpoint）等待一组 goroutine 全部完成
	如果在执行任务的这些 goroutine 还没全部完成，那么 goroutine A 就会阻塞在检查点
	直到所有 goroutine 都完成后才能继续执行
需求场景
	我们要完成一个大的任务，需要使用并行的 goroutine 执行三个小任务，只有这三个小任务都完成，我们才能去执行后面的任务
	如果通过轮询的方式定时询问三个小任务是否完成，会存在两个问题：
	一是，性能比较低，因为三个小任务可能早就完成了，却要等很长时间才被轮询到
	二是，会有很多无谓的轮询，空耗 CPU 资源
解决方案：WaitGroup
	这个时候使用 WaitGroup 并发原语就比较有效了，它可以阻塞等待的 goroutine
	等到三个小任务都完成了，再即时唤醒它们
类似的并发原语
	很多操作系统和编程语言都提供了类似的并发原语
	比如，Linux 中的 barrier、Pthread（POSIX 线程）中的 barrier、C++ 中的 std::barrier、Java 中的 CyclicBarrier 和 CountDownLatch 等
	可见，这个并发原语还是一个非常基础的并发类型

WaitGroup 的基本用法
	Go 标准库中的 WaitGroup 提供了三个方法，保持了 Go 简洁的风格
		func (wg *WaitGroup) Add(delta int)
		func (wg *WaitGroup) Done()
		func (wg *WaitGroup) Wait()
		Add
			用来设置 WaitGroup 的计数值
		Done
			用来将 WaitGroup 的计数值减 1，其实就是调用了 Add(-1)
		Wait
			调用这个方法的 goroutine 会一直阻塞，直到 WaitGroup 的计数值变为 0
	示例
		WGCounterMain & TestWGCounterMain
		使用 WaitGroup 编排这类任务的常用方式
		需要启动多个 goroutine 执行任务，主 goroutine 则需要等待子 goroutine 都完成后才继续执行

WaitGroup 的实现
	数据结构
		它包括了一个 noCopy 的辅助字段，一个 state1 记录 WaitGroup 状态的数组
		noCopy 辅助字段：主要就是辅助 vet 工具检查是否通过 copy 赋值这个 WaitGroup 实例
		state1：一个具有复合意义的字段，包含 WaitGroup 的计数、阻塞在检查点的 waiter 数
		sema  uint32：信号量
	WaitGroup 的数据结构定义以及 state 信息的获取方法
		因为对 64 位整数的原子操作要求整数的地址是 64 位对齐的，所以针对 64 位和 32 位环境的 state 字段的组成是不一样的
		在 64 位环境下，state1 的第一个元素是 waiter 数，第二个元素是 WaitGroup 的计数值，第三个元素是信号量
			06.waitgroup_state_64.jpg
		在 32 位环境下，如果 state1 不是 64 位对齐的地址，那么 state1 的第一个元素是信号量，后两个元素分别是 waiter 数和计数值
			06.waitgroup_state_32.jpg
	源码分析
		除了这些方法本身的实现外，还会有一些额外的代码，主要是 race 检查和异常检查的代码
		其中，有几个检查非常关键，如果检查不通过，会出现 panic
	Add 方法逻辑
		Add 方法主要操作的是 state 的计数部分
		你可以为计数值增加一个 delta 值，内部通过原子操作把这个值加到计数值上
		需要注意的是，这个 delta 也可以是个负数，相当于为计数值减去一个值，Done 方法内部其实就是通过 Add(-1) 实现的
	Done 方法
		实际就是计数器减 1
	Wait 方法逻辑
		不断检查 state 的值
		如果其中的计数值变为了 0，那么说明所有的任务已完成，调用者不必再等待，直接返回
		如果计数值大于 0，说明此时还有任务没完成，那么调用者就变成了等待者，需要加入 waiter 队列，并且阻塞住自己

使用 WaitGroup 时的常见错误
常见问题一：计数器设置为负值
	WaitGroup 的计数器的值必须大于等于 0
		我们在更改这个计数值的时候，WaitGroup 会先做检查，如果计数值被设置为负数，就会导致 panic
		一般情况下，有两种方法会导致计数器设置为负数
	第一种方法：调用 Add 的时候传递一个负数
		如果你能保证当前的计数器加上这个负数后还是大于等于 0 的话，也没有问题，否则就会导致 panic
		示例：panic
			var wg sync.WaitGroup
			wg.Add(10)
			wg.Add(-10) //将-10作为参数调用Add，计数值被设置为0
			wg.Add(-1)  //将-1作为参数调用Add，如果加上-1计数值就会变为负数。这是不对的，所以会触发
	第二种方法：调用 Done 方法的次数过多，超过了 WaitGroup 的计数值
		使用 WaitGroup 的正确姿势
			预先确定好 WaitGroup 的计数值，然后调用相同次数的 Done 完成相应的任务
			比如，在 WaitGroup 变量声明之后，就立即设置它的计数值，或者在 goroutine 启动之前增加 1
			然后在 goroutine 中调用 Done
		如果你没有遵循这些规则，就很可能会导致 Done 方法调用的次数和计数值不一致
			进而造成死锁（Done 调用次数比计数值少）或者 panic（Done 调用次数比计数值多）
		示例：panic
			var wg sync.WaitGroup
			wg.Add(1)
			wg.Done()
			wg.Done()
常见问题二：不期望的 Add 时机
	原则
		等所有的 Add 方法调用之后再调用 Wait，否则就可能导致 panic 或者不期望的结果
	示例：只有部分的 Add/Done 执行完后，Wait 就返回
		WGWrongAdd & TestWGWrongAdd
	分析
		错误之处在于，将 WaitGroup.Add 方法的调用放在了子 gorotuine 中
		等主 goorutine 调用 Wait 的时候，因为四个任务 goroutine 一开始都休眠
		所以可能 WaitGroup 的 Add 方法还没有被调用，WaitGroup 的计数还是 0
		所以它并没有等待四个子 goroutine 执行完毕才继续执行，而是立刻执行了下一步
	原因
		没有遵循先完成所有的 Add 之后才 Wait
	两种解决方式
		一个方法是，预先设置计数值
		另一种方法是在启动子 goroutine 之前才调用 Add
		无论是怎么修复，都要保证所有的 Add 方法是在 Wait 方法之前被调用的
常见问题三：前一个 Wait 还没结束就重用 WaitGroup
	WaitGroup 是可以重用的
		只要 WaitGroup 的计数值恢复到零值的状态，那么它就可以被看作是新创建的 WaitGroup，被重复使用
		但是，如果我们在 WaitGroup 的计数值还没有恢复到零值的时候就重用，就会导致程序 panic
	示例
		WGCopy & TestWGCopy
	panic
		panic: sync: WaitGroup is reused before previous Wait has returned
	小结
		WaitGroup 虽然可以重用，但是是有一个前提的，那就是必须等到上一轮的 Wait 完成之后，才能重用 WaitGroup 执行下一轮的 Add/Wait
		如果你在 Wait 还没执行完的时候就调用下一轮 Add 方法，就有可能出现 panic

noCopy：辅助 vet 检查
	作用
		指示 vet 工具在做检查的时候，这个数据结构不能做值复制使用
		更严谨地说，是不能在第一次使用之后复制使用(must not be copied after first use)
		重要的是，noCopy 是一个通用的计数技术，其他并发原语中也会用到
	vet
		vet 会对实现 Locker 接口的数据类型做静态检查，一旦代码中有复制使用这种数据类型的情况，就会发出警告
	原理
		通过给 WaitGroup 添加一个 noCopy 字段，我们就可以为 WaitGroup 实现 Locker 接口
		这样 vet 工具就可以做复制检查了
		而且因为 noCopy 字段是未输出类型，所以 WaitGroup 不会暴露 Lock/Unlock 方法
	noCopy 字段的类型是 noCopy，它只是一个辅助的、用来帮助 vet 检查用的类型
		type noCopy struct{}
		// Lock is a no-op used by -copylocks checker from `go vet`.
		func (*noCopy) Lock()   {}
		func (*noCopy) Unlock() {}
	使用方式
		如果你想要自己定义的数据结构不被复制使用，或者说，不能通过 vet 工具检查出复制使用的报警
		就可以通过嵌入 noCopy 这个数据类型来实现

流行的 Go 开发项目中的坑
	有网友在 Go 的 issue 28123中提了例子
		示例 WaitGroupBug & TestWaitGroupBug
		分析
			代码最大的一个问题，就是 copy 了 WaitGroup 的实例 w
			虽然这段代码能执行成功，但确实是违反了 WaitGroup 使用之后不要复制的规则
			在项目中，我们可以通过 vet 工具检查出这样的错误
	Docker：issue 28161 和 issue 27011
		都是因为在重用 WaitGroup 的时候，没等前一次的 Wait 结束就 Add 导致的错误
	Etcd：issue 6534
		是重用 WaitGroup 的 Bug，没有等前一个 Wait 结束就 Add
	Kubernetes：issue 59574
		忘记 Wait 之前增加计数了，这就属于我们通常认为几乎不可能出现的 Bug
		图示
			06.waitgroup_demo_kubernetes.jpg
	Go 语言：issue 12813
		因为 defer 的使用，Add 方法可能在 Done 之后才执行，导致计数负值的 panic
		图示
			06.waitgroup_demo_go.jpg

避免错误使用 WaitGroup
	不重用 WaitGroup。新建一个 WaitGroup 不会带来多大的资源开销，重用反而更容易出错
	保证所有的 Add 方法调用都在 Wait 之前
	不传递负数给 Add 方法，只通过 Done 来给计数值减 1
	不做多余的 Done 方法调用，保证 Add 的计数值和 Done 方法调用的数量是一样的
	不遗漏 Done 方法的调用，否则会导致 Wait hang 住无法返回

思考
	通常我们可以把 WaitGroup 的计数值，理解为等待要完成的 waiter 的数量
	你可以试着扩展下 WaitGroup，来查询 WaitGroup 的当前的计数值吗？
*/

// WaitGroup ==========查询 WaitGroup 计数值==========
type WaitGroup struct {
	sync.WaitGroup
}

func (wg *WaitGroup) GetWaitGroupCount() (int32, uint32) {
	state := atomic.LoadUint64((*uint64)(unsafe.Pointer(&wg.WaitGroup)))
	v := int32(state >> 32)
	w := uint32(state)
	return v, w
}

// TestWGStruct ==========copy bug==========
type TestWGStruct struct {
	Wait sync.WaitGroup
}

func WaitGroupBug() {
	w := sync.WaitGroup{}
	w.Add(1)
	t := &TestWGStruct{
		Wait: w,
	}
	t.Wait.Done()
	fmt.Println("Finished")
	t.Wait.Wait()
}

// WGCopy ==========错误重用 WaitGroup==========
func WGCopy() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(time.Millisecond)
		wg.Done() // 计数器减1
		wg.Add(1) // 计数值加1
	}()
	wg.Wait() // 主goroutine等待，有可能和第7行并发执行
}

// WGWrongAdd ==========线程安全的计数器==========
func WGWrongAdd() {
	var wg sync.WaitGroup
	// 方式一
	//wg.Add(4)                // 预先设定WaitGroup的计数值
	//go dosomething(100, &wg) // 启动第一个goroutine
	//go dosomething(110, &wg) // 启动第二个goroutine
	//go dosomething(120, &wg) // 启动第三个goroutine
	//go dosomething(130, &wg) // 启动第四个goroutine

	// 方式二
	dothing(100, &wg) // 调用方法，把计数值加1，并启动任务goroutine
	dothing(110, &wg) // 调用方法，把计数值加1，并启动任务goroutine
	dothing(120, &wg) // 调用方法，把计数值加1，并启动任务goroutine
	dothing(130, &wg) // 调用方法，把计数值加1，并启动任务goroutine

	wg.Wait() // 主goroutine等待完成
	fmt.Println("Done")
}
func dosomething(millisecs time.Duration, wg *sync.WaitGroup) {
	duration := millisecs * time.Millisecond
	time.Sleep(duration) // 故意sleep一段时间
	wg.Add(1)
	fmt.Println("后台执行, duration:", duration)
	wg.Done()
}
func dothing(millisecs time.Duration, wg *sync.WaitGroup) {
	wg.Add(1) // 计数值加1，再启动goroutine
	go func() {
		duration := millisecs * time.Millisecond
		time.Sleep(duration)
		fmt.Println("后台执行, duration:", duration)
		wg.Done()
	}()
}

// WGCounter ==========线程安全的计数器==========
type WGCounter struct {
	mu    sync.Mutex
	count uint64
}

// Incr 对计数值加一
func (c *WGCounter) Incr() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

// WGCounter 获取当前的计数值
func (c *WGCounter) WGCounter() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// sleep 1秒，然后计数值加1
func worker(c *WGCounter, wg *sync.WaitGroup) {
	defer wg.Done() // 调用 Done 方法，把 WaitGroup 的计数值减 1
	time.Sleep(time.Second)
	c.Incr()
}
func WGCounterMain() {
	var wgCounter WGCounter
	var wg sync.WaitGroup     // 声明了一个 WaitGroup 变量，初始值为零
	wg.Add(10)                // WaitGroup的值设置为10
	for i := 0; i < 10; i++ { // 启动10个goroutine执行加1任务
		go worker(&wgCounter, &wg) // WaitGroup 指针当作参数传递进去
	}
	wg.Wait()                          // 阻塞：检查点，等待goroutine都完成任务
	fmt.Println(wgCounter.WGCounter()) // 输出当前计数器的值
}
