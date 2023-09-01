package concurrent

import (
	"fmt"
	"sync"
)

/*
Mutex：如何解决资源并发访问问题？

需求场景
	并发访问问题
		多个 goroutine 并发更新同一个资源，像计数器
		同时更新用户的账户信息
		秒杀系统
		往同一个 buffer 中并发写入数据
	如果没有互斥控制，就会出现一些异常情况
		比如计数器的计数不准确、用户的账户可能出现透支、秒杀系统出现超卖、buffer 中的数据混乱

互斥锁的实现机制
	互斥锁/排它锁
		并发控制的一个基本手段，是为了避免竞争而建立的一种并发控制机制
		Go 标准库中，它提供了 Mutex 来实现互斥锁这个功能
	临界区
		在并发编程中，如果程序中的一部分会被并发访问或修改
		为了避免并发访问导致的意想不到的结果，这部分程序需要被保护起来，这部分被保护起来的程序，就叫做临界区
		可以说，临界区就是一个被共享的资源，或者说是一个整体的一组共享资源
		比如对数据库的访问、对某一个共享数据结构的操作、对一个 I/O 设备的使用、对一个连接池中的连接的调用，等
	可以使用互斥锁，限定临界区只能同时由一个线程持有
		当临界区由一个线程持有的时候，其它线程如果想进入这个临界区，就会返回失败，或者是等待
		直到持有的线程退出临界区，这些等待线程中的某一个才有机会接着持有这个临界区
		示例：01.mutex_01.jpg
	同步原语 & 并发原语
		根据 2019 年第一篇全面分析 Go 并发 Bug 的论文 Understanding Real-World Concurrency Bugs in Go
		Mutex 是使用最广泛的同步原语（Synchronization primitives，有人也叫做并发原语）
		直译优先用同步原语，但是并发原语的指代范围更大，还可以包括任务编排的类型，Channel 或者扩展类型时也会用并发原语
		关于同步原语，并没有一个严格的定义，你可以把它看作解决并发问题的一个基础的数据结构
		论文：https://songlh.github.io/paper/go-study.pdf
	同步原语适用场景
		共享资源
			并发地读写共享资源，会出现数据竞争（data race）的问题，所以需要Mutex、RWMutex 这样的并发原语来保护
		任务编排
			需要 goroutine 按照一定的规律执行，而 goroutine 之间有相互等待或者依赖的顺序关系，我们常常使用 WaitGroup 或者 Channel 来实现
		消息传递
			信息交流以及不同的 goroutine 之间的线程安全的数据交流，常常使用 Channel 来实现

Mutex 的基本使用方法
	Locker 接口
		在 Go 的标准库中，package sync 提供了锁相关的一系列同步原语
		这个 package 还定义了一个 Locker 的接口，Mutex、RWMutex 就实现了这个接口
		Go 定义的锁接口的方法集很简单，就是请求锁（Lock）和释放锁（Unlock）这两个方法
		这个接口在实际项目应用得不多，因为我们一般会直接使用具体的同步原语，而不是通过接口
	Mutex
		互斥锁 Mutex 就提供两个方法 Lock 和 Unlock
		进入临界区之前调用 Lock 方法，退出临界区的时候调用 Unlock 方法
		当一个 goroutine 通过调用 Lock 方法获得了这个锁的拥有权后， 其它请求锁的 goroutine 就会阻塞在 Lock 方法的调用上
		直到锁被释放并且自己获取到了这个锁的拥有权
	示例
		问题分析
			count++ 不是一个原子操作，它至少包含几个步骤
			比如读取变量 count 的当前值，对这个值加 1，把结果再保存到 count 中。因为不是原子操作，就可能有并发的问题
			比如，10 个 goroutine 同时读取到 count 的值为 9527，接着各自按照自己的逻辑加 1，值变成了 9528，然后把这个结果再写回到 count 变量
			但是，实际上，此时我们增加的总数应该是 10 才对，这里却只增加了 1，好多计数都被“吞”掉了
			这是并发访问共享数据的常见错误
		count++ 操作
			// count++操作的汇编代码
			MOVQ "".count(SB), AX
			LEAQ 1(AX), CX
			MOVQ CX, "".count(SB)
		Mutex 基本用法
			共享资源是 count 变量，临界区是 count++
			只要在临界区前面获取锁，在离开临界区的时候释放锁，就能完美地解决 data race 的问题了
			$ go test -v -run TestMutexConcurrent golang/concurrent/test
	注意
		Mutex 的零值是还没有 goroutine 等待的未加锁的状态，所以你不需要额外的初始化
		直接声明变量（如 var mu sync.Mutex）即可
	Mutex 用法
		Mutex 嵌入到其它 struct 中使用
			在初始化嵌入的 struct 时，也不必初始化这个 Mutex 字段，不会因为没有初始化出现空指针或者是无法获取到锁的情况
			type Counter struct {
				mu sync.Mutex
				Count uint64
			}
		采用嵌入字段的方式
			通过嵌入字段，你可以在这个 struct 上直接调用 Lock/Unlock 方法
			func main() {
				var counter Counter
				var wg sync.WaitGroup
				wg.Add(10)
				for i := 0; i < 10; i++ {
					go func() {
						defer wg.Done()
						for j := 0; j < 100000; j++ {
							counter.Lock()
							counter.Count++
							counter.Unlock()
						}
					}()
				}
				wg.Wait()
				fmt.Println(counter.Count)
			}
		如果嵌入的 struct 有多个字段，我们一般会把 Mutex 放在要控制的字段上面，然后使用空格把字段分隔开来
			即使你不这样做，代码也可以正常编译，只不过，用这种风格去写的话，逻辑会更清晰，也更易于维护
			甚至，你还可以把获取锁、释放锁、计数加一的逻辑封装成一个方法，对外不需要暴露锁等逻辑

race detector
	Go 提供的一个检测并发访问共享资源是否有问题的工具
		它可以帮助我们自动发现程序有没有 data race 的问题
	简介
		Go race detector 是基于 Google 的 C/C++ sanitizers 技术实现的
		编译器通过探测所有的内存访问，加入代码能监视对这些内存地址的访问（读还是写）
		在代码运行的时候，race detector 就能监控到对共享变量的非同步访问，出现 race 的时候，就会打印出警告信息
	应用
		这个技术在 Google 内部帮了大忙，探测出了 Chromium 等代码的大量并发问题
		Go 1.1 中就引入了这种技术，并且一下子就发现了标准库中的 42 个并发问题
		现在，race detector 已经成了 Go 持续集成过程中的一部分
	使用：在编译（compile）、测试（test）或者运行（run）Go 代码的时候，加上 race 参数，就有可能发现并发问题
		管理员权限开启 cmd：set CGO_ENABLED=1
		main
			$ go run -race mutex.go
			E:\gothmslee\golang\main>go run -race mutex.go
		go test
			E:\gothmslee\golang\concurrent\test>go test -v -race 01_test.go
			=== RUN   TestMutex
			==================
			WARNING: DATA RACE
			Read at 0x00c0000182a8 by goroutine 9:
			  command-line-arguments.TestMutex.func1()
				  E:/gothmslee/golang/concurrent/test/01_test.go:21 +0xa8

			Previous write at 0x00c0000182a8 by goroutine 8:
			  command-line-arguments.TestMutex.func1()
				  E:/gothmslee/golang/concurrent/test/01_test.go:21 +0xba

			Goroutine 9 (running) created at:
			  command-line-arguments.TestMutex()
				  E:/gothmslee/golang/concurrent/test/01_test.go:18 +0x8d
			  testing.tRunner()
				  E:/Go/src/testing/testing.go:1576 +0x216
			  testing.(*T).Run.func1()
				  E:/Go/src/testing/testing.go:1629 +0x47

			Goroutine 8 (running) created at:
			  command-line-arguments.TestMutex()
				  E:/gothmslee/golang/concurrent/test/01_test.go:18 +0x8d
			  testing.tRunner()
				  E:/Go/src/testing/testing.go:1576 +0x216
			  testing.(*T).Run.func1()
				  E:/Go/src/testing/testing.go:1629 +0x47
			==================
				01_test.go:26: 329539
				testing.go:1446: race detected during execution of test
			--- FAIL: TestMutex (0.07s)
			=== NAME
				testing.go:1446: race detected during execution of test
			FAIL
			FAIL    command-line-arguments  0.138s
			FAIL
		分析
			WARNING: DATA RACE
				输出警告信息
			Read at 0x00c0000182a8 by goroutine 9:
			  command-line-arguments.TestMutex.func1()
				  E:/gothmslee/golang/concurrent/test/01_test.go:21 +0xa8
			...
				哪个 goroutine 在哪一行对哪个变量有读/写操作
				就是这些并发的读写访问，引起了 data race
		局限性
			虽然这个工具使用起来很方便，但是，因为它的实现方式，只能通过真正对实际地址进行读写访问的时候才能探测
			所以它并不能在编译的时候发现 data race 的问题
			而且，在运行的时候，只有在触发了 data race 之后，才能检测到
			如果碰巧没有触发（比如一个 data race 问题只能在 2 月 14 号零点或者 11 月 11 号零点才出现），是检测不出来的
			而且，把开启了 race 的程序部署在线上，还是比较影响性能的
		go tool compile -race -S mutex.go
			报错：
				mutex.go:5:2: could not import sync (file not found)
			解决：https://github.com/golang/go/issues/58629
				go 工具编译 的命令行标志和应用程序接口非常复杂，而且每个版本都会有变化
				使用 "go build"，而不是直接运行 "go tool compile"
			go build -race -gcflags=-S mutex.go
				...
				0x009f 00159 (E:\gothmslee\golang\main\mutex.go:19)     MOVQ    main.&cnt+104(SP), AX
				0x00a4 00164 (E:\gothmslee\golang\main\mutex.go:19)     PCDATA  $1, $2
				0x00a4 00164 (E:\gothmslee\golang\main\mutex.go:19)     CALL    runtime.raceread(SB)
				0x00a9 00169 (E:\gothmslee\golang\main\mutex.go:19)     MOVQ    main.&cnt+104(SP), AX
				0x00ae 00174 (E:\gothmslee\golang\main\mutex.go:19)     MOVQ    (AX), CX
				0x00b1 00177 (E:\gothmslee\golang\main\mutex.go:19)     MOVQ    CX, main..autotmp_9+16(SP)
				0x00b6 00182 (E:\gothmslee\golang\main\mutex.go:19)     CALL    runtime.racewrite(SB)
				0x00bb 00187 (E:\gothmslee\golang\main\mutex.go:19)     MOVQ    main..autotmp_9+16(SP), CX
				0x00c0 00192 (E:\gothmslee\golang\main\mutex.go:19)     INCQ    CX
				0x00c3 00195 (E:\gothmslee\golang\main\mutex.go:19)     MOVQ    main.&cnt+104(SP), AX
				0x00c8 00200 (E:\gothmslee\golang\main\mutex.go:19)     MOVQ    CX, (AX)
				0x00cb 00203 (E:\gothmslee\golang\main\mutex.go:18)     MOVQ    main.j+8(SP), CX
				0x00d0 00208 (E:\gothmslee\golang\main\mutex.go:18)     INCQ    CX
				0x00d3 00211 (E:\gothmslee\golang\main\mutex.go:18)     MOVQ    CX, AX
				0x00d6 00214 (E:\gothmslee\golang\main\mutex.go:18)     CMPQ    AX, $100000
				0x00dc 00220 (E:\gothmslee\golang\main\mutex.go:18)     JLT     154
				0x00de 00222 (E:\gothmslee\golang\main\mutex.go:21)     PCDATA  $1, $0
				0x00de 00222 (E:\gothmslee\golang\main\mutex.go:21)     NOP
				0x00e0 00224 (E:\gothmslee\golang\main\mutex.go:21)     CALL    runtime.deferreturn(SB)
				0x00e5 00229 (E:\gothmslee\golang\main\mutex.go:21)     CALL    runtime.racefuncexit(SB)
				...
			分析
				查看计数器例子的代码，重点关注一下 count++ 前后的编译后的代码
				在编译的代码中，增加了 runtime.racefuncenter、runtime.raceread、runtime.racewrite、runtime.racefuncexit 等检测 data race 的方法
				通过这些插入的指令，Go race detector 工具就能够成功地检测出 data race 问题了
			参见
				如 runtime/race.go 中，runtime.racefuncexit
		data race 工具的实现机制
			通过在编译的时候插入一些指令，在运行时通过这些插入的指令检测并发读写从而发现 data race 问题

总结
	意识到共享资源的并发访问
		在项目开发的初始阶段，我们可能并没有仔细地考虑资源的并发问题，因为在初始阶段，我们还不确定这个资源是否被共享
		经过更加深入的设计，或者新功能的增加、代码的完善，这个时候，我们就需要考虑共享资源的并发问题了
		当然，如果你能在初始阶段预见到资源会被共享并发访问就更好了
	意识到共享资源的并发访问的早晚不重要
		重要的是，一旦你意识到这个问题，你就要及时通过互斥锁等手段去解决
		比如 Docker issue 37583、35517、32826、30696等，kubernetes issue 72361、 71617等，都是后来发现的 data race
		而采用互斥锁 Mutex 进行修复的

思考
	如果 Mutex 已经被一个 goroutine 获取了锁，其它等待中的 goroutine 们只能一直等待
	那么，等这个锁释放后，等待中的 goroutine 中哪一个会优先获取 Mutex 呢？
*/

func MutexPractice() {
	// 封装好的计数器
	var counter counter
	var wg sync.WaitGroup
	wg.Add(10)
	// 启动10个goroutine
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			// 执行10万次累加
			for j := 0; j < 100000; j++ {
				counter.incr() // 受到锁保护的方法
			}
		}()
	}
	wg.Wait()
	fmt.Println(counter.getCount())
}

// Counter 线程安全的计数器类型
type counter struct {
	CounterType int
	Name        string
	mu          sync.Mutex
	count       uint64
}

// Incr 加1的方法，内部使用互斥锁保护
func (c *counter) incr() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

// Count 得到计数器的值，也需要锁保护
func (c *counter) getCount() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}
