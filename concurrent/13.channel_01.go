package concurrent

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

/*
Channel：另辟蹊径，解决并发问题

Channel 是 Go 语言内建的 first-class 类型，也是 Go 语言与众不同的特性之一
	Go 语言的 Channel 设计精巧简单，以至于也有人用其它语言编写了类似 Go 风格的 Channel 库
	比如 docker/libchan、tylertreat/chan，但是并不像 Go 语言一样把 Channel 内置到了语言规范中
	从这一点，你也可以看出来，Channel 的地位在编程语言中的地位之高，比较罕见

Channel 的发展
	CSP
		Communicating Sequential Process 的简称，中文直译为通信顺序进程，或者叫做交换信息的循序进程
		是用来描述并发系统中进行交互的一种模式
	CSP 历史
		CSP 最早出现于计算机科学家 Tony Hoare 在 1978 年发表的论文中（Tony Hoare 也是 Quicksort 排序算法的作者，图灵奖的获得者）
		最初，论文中提出的 CSP 版本在本质上不是一种进程演算，而是一种并发编程语言，但之后又经过了一系列的改进，最终发展并精炼出 CSP 的理论
		CSP 允许使用进程组件来描述系统，它们独立运行，并且只通过消息传递的方式通信
	CSP 模型对 Go 创始人设计 Channel 类型的影响
		就像 Go 的创始人之一 Rob Pike 所说的：“每一个计算机程序员都应该读一读 Tony Hoare 1978 年的关于 CSP 的论文。”
		他和 Ken Thompson 在设计 Go 语言的时候也深受此论文的影响
		并将 CSP 理论真正应用于语言本身（Russ Cox 专门写了一篇文章记录这个历史），通过引入 Channel 这个新的类型，来实现 CSP 的思想
		link：历史
		https://zhuanlan.zhihu.com/p/143846211
	内置类型 Channel
		Channel 类型是 Go 语言内置的类型，你无需引入某个包，就能使用它
		虽然 Go 也提供了传统的并发原语，但是它们都是通过库的方式提供的，你必须要引入 sync 包或者 atomic 包才能使用它们
		而 Channel 就不一样了，它是内置类型，使用起来非常方便
		Channel 和 Go 的另一个独特的特性 goroutine 一起为并发编程提供了优雅的、便利的、与传统并发控制不同的方案，并演化出很多并发模式

Channel 的应用场景
	Don’t communicate by sharing memory, share memory by communicating.
		Go Proverbs by Rob Pike
		这是 Rob Pike 在 2015 年的一次 Gopher 会议中提到的一句话，虽然有一点绕，但也指出了使用 Go 语言的哲学
		翻译：“执行业务处理的 goroutine 不要通过共享内存的方式通信，而是要通过 Channel 通信的方式分享数据。”
	“communicate by sharing memory”和“share memory by communicating”是两种不同的并发处理模式
		“communicate by sharing memory”
			是传统的并发编程处理方式，就是指，共享的数据需要用锁进行保护，goroutine 需要获取到锁，才能并发访问数据
		“share memory by communicating”
			则是类似于 CSP 模型的方式，通过通信的方式，一个 goroutine 可以把数据的“所有权”交给另外一个 goroutine
			虽然 Go 中没有“所有权”的概念，但是从逻辑上说，你可以把它理解为是所有权的转移
	竞争关系
		从 Channel 的历史和设计哲学上，我们就可以了解到，Channel 类型和基本并发原语是有竞争关系的
		它应用于并发场景，涉及到 goroutine 之间的通讯，可以提供并发的保护，等
	Channel 的应用场景分为五种类型
		1. 数据交流：当作并发的 buffer 或者 queue，解决生产者 - 消费者问题
			多个 goroutine 可以并发当作生产者（Producer）和消费者（Consumer）
		2. 数据传递：一个 goroutine 将数据交给另一个 goroutine，相当于把数据的拥有权 (引用) 托付出去
		3. 信号通知：一个 goroutine 可以将信号 (closing、closed、data ready 等) 传递给另一个或者另一组 goroutine
		4. 任务编排：可以让一组 goroutine 按照一定的顺序并发或者串行的执行，这就是编排的功能
		5. 锁：利用 Channel 也可以实现互斥锁的机制

Channel 基本用法
	Channel 类型
		只能接收、只能发送、既可以接收又可以发送
			官方文档：ChannelType = ( "chan" | "chan" "<-" | "<-" "chan" ) ElementType
			既能接收又能发送的 chan 叫做双向的 chan，把只能发送和只能接收的 chan 叫做单向的 chan
		箭头总是射向左边的，元素类型总在最右边
			如果箭头指向 chan，就表示可以往 chan 中塞数据
			如果箭头远离 chan，就表示 chan 会往外吐数据
		“<-”规则
			总是尽量和左边的 chan 结合（The <- operator associates with the leftmost chan possible:）
		chan 中的元素是任意的类型，所以也可能是 chan 类型
			a := make(chan<- chan int)
			b := make(chan<- <-chan int)
			c := make(<-chan <-chan int)
			d := make(chan (<-chan int))
	初始化
		通过 make，我们可以初始化一个 chan
		未初始化的 chan 的零值是 nil
		nil 是 chan 的零值，是一种特殊的 chan，对值是 nil 的 chan 的发送接收调用者总是会阻塞
	分类
		buffered chan
			如果 chan 中还有数据，那么从这个 chan 接收数据的时候就不会阻塞
			如果 chan 还未满（“满”指达到其容量），给它发送数据也不会阻塞，否则就会阻
		unbuffered chan
			只有读写都准备好之后才不会阻塞，这也是很多使用 unbuffered chan 时的常见 Bug
	1. 发送数据
		ch <- 2000
	2. 接收数据
		示例
			x := <-ch // 把接收的一条数据赋值给变量x
			foo(<-ch) // 把接收的一个的数据作为参数传给函数
			<-ch      // 丢弃接收的一条数据
		返回值：可以返回两个值
			第一个值是返回的 chan 中的元素
			第二个值是 bool 类型，代表是否成功地从 chan 中读取到一个值
				如果第二个参数是 false，chan 已经被 close 而且 chan 中没有缓存的数据，这个时候，第一个值是零值
				所以，如果从 chan 读取到一个零值，可能是 sender 真正发送的零值，也可能是 closed 的并且没有缓存元素产生的零值
	3. 其它操作
		Go 内建的函数 close、cap、len 都可以操作 chan 类型
			close 会把 chan 关闭掉，cap 返回 chan 的容量，len 返回 chan 中缓存的还未被取走的元素数量
		select 语句
			send 和 recv 都可以作为 select 语句的 case clause
		for-range 语句
			for v := range ch {
				fmt.Println(v)
			}
		清空 chan
			for range ch {
			}

Channel 的实现原理
	重点：chan 的数据结构、初始化的方法以及三个重要的操作方法，分别是 send、recv 和 close
chan 数据结构
	数据类型 runtime.hchan
		源码 runtime/chan.go
	图示
		13.channel_01_structure.jpg
	字段
		qcount：代表 chan 中已经接收但还没被取走的元素的个数
			内建函数 len 可以返回这个字段的值
		dataqsiz：队列的大小
			chan 使用一个循环队列来存放元素，循环队列很适合这种生产者 - 消费者的场景
			好奇为什么这个字段省略 size 中的 e
		buf：存放元素的循环队列的 buffer
		elemtype 和 elemsize：chan 中元素的类型和 size
			因为 chan 一旦声明，它的元素类型是固定的，即普通类型或者指针类型，所以元素大小也是固定的
		sendx：处理发送数据的指针在 buf 中的位置
			一旦接收了新的数据，指针就会加上
		elemsize，移向下一个位置
			buf 的总大小是 elemsize 的整数倍，而且 buf 是一个循环列表
		recvx：处理接收请求时的指针在 buf 中的位置
			一旦取出数据，此指针会移动到下一个位置
		recvq：chan 是多生产者多消费者的模式，如果消费者因为没有数据可读而被阻塞了，就会被加入到 recvq 队列中
		sendq：如果生产者因为 buf 满了而阻塞，会被加入到 sendq 队列中
初始化
	Go 在编译的时候，会根据容量的大小选择调用 makechan64，还是 makechan
		处理 make chan 的逻辑，它会决定是使用 makechan 还是 makechan64 来实现 chan 的初始化
		makechan64 只是做了 size 检查，底层还是调用 makechan 实现的
		makechan 的目标就是生成 hchan 对象
	图示 13.channel_01_make_chan.jpg
		参考：https://www.bookstack.cn/read/draveness-golang/c666731d5f1a2820.md
	makechan 的主要逻辑
		它会根据 chan 的容量的大小和元素的类型不同，初始化不同的存储空间
		针对不同的容量和元素类型，这段代码分配了不同的对象来初始化 hchan 对象的字段，返回 hchan 对象
send
	Go 在编译发送数据给 chan 的时候，会把 send 语句转换成 chansend1 函数，chansend1 函数会调用 chansend，
	第一部分
		进行判断：如果 chan 是 nil 的话，就把调用者 goroutine park（阻塞休眠），调用者就永远被阻塞住了
	第二部分
		当你往一个已经满了的 chan 实例发送数据时，并且想不阻塞当前调用，那么这里的逻辑是直接返回
		chansend1 方法在调用 chansend 的时候设置了阻塞参数，所以不会执行到第二部分的分支里
	第三部分
		如果 chan 已经被 close 了，再往里面发送数据的话会 panic
	第四部分
		如果等待队列中有等待的 receiver，那么这段代码就把它从队列中弹出
		然后直接把数据交给它（通过 memmove(dst, src, t.size)），而不需要放入到 buf 中，速度可以更快一些
	第五部分
		当前没有 receiver，需要把数据放入到 buf 中，放入之后，就成功返回了
	第六部分
		处理 buf 满的情况
		如果 buf 满了，发送者的 goroutine 就会加入到发送者的等待队列中，直到被唤醒
		这个时候，数据或者被取走了，或者 chan 被 close 了
recv
	在处理从 chan 中接收数据时，Go 会把代码转换成 chanrecv1 函数
		如果要返回两个返回值，会转换成 chanrecv2
		chanrecv1 函数和 chanrecv2 会调用 chanrecv
		chanrecv1 和 chanrecv2 传入的 block 参数的值是 true，都是阻塞方式
		分析 chanrecv 的实现的时候，不考虑 block=false 的情况
	第一部分
		chan 为 nil 的情况
		和 send 一样，从 nil chan 中接收（读取、获取）数据时，调用者会被永远阻塞
	第二部分
		先忽略，因为不是这次要分析的场景
	第三部分
		chan 已经被 close 的情况
		如果 chan 已经被 close 了，并且队列中没有缓存的元素，那么返回 true、false
	第四部分
		处理 sendq 队列中有等待者的情况
		这个时候，如果 buf 中有数据，优先从 buf 中读取数据，否则直接从等待队列中弹出一个 sender，把它的数据复制给这个 receiver
	第五部分
		处理没有等待的 sender 的情况
		这个是和 chansend 共用一把大锁，所以不会有并发的问题
		如果 buf 有元素，就取出一个元素给 receiver
	第六部分
		处理 buf 中没有元素的情况
		如果没有元素，那么当前的 receiver 就会被阻塞，直到它从 sender 中接收了数据，或者是 chan 被 close，才返回
close
	通过 close 函数，可以把 chan 关闭，编译器会替换成 closechan 方法的调用
		如果 chan 为 nil，close 会 panic；如果 chan 已经 closed，再次 close 也会 panic
		否则的话，如果 chan 不为 nil，chan 也没有 closed
		就把等待队列中的 sender（writer）和 receiver（reader）从队列中全部移除并唤醒

使用 Channel 容易犯的错误
	根据 2019 年第一篇全面分析 Go 并发 Bug 的论文，那些知名的 Go 项目中使用 Channel 所犯的 Bug 反而比传统的并发原语的 Bug 还要多
		主要有两个原因：
		一个是，Channel 的概念还比较新，程序员还不能很好地掌握相应的使用方法和最佳实践
		第二个是，Channel 有时候比传统的并发原语更复杂，使用起来很容易顾此失彼
		link：论文
	使用 Channel 最常见的错误是 panic 和 goroutine 泄漏
	panic
		1. close 为 nil 的 chan
		2. send 已经 close 的 chan
		3. close 已经 close 的 chan
	goroutine 内存泄漏
		示例
			func process(timeout time.Duration) bool {
				ch := make(chan bool)
				go func() {
					time.Sleep((timeout + time.Second)) // 模拟处理耗时的业务
					ch <- true                          // block
					fmt.Println("exit goroutine")
				}()
				select {
				case result := <-ch:
					return result
				case <-time.After(timeout):
					return false
				}
			}
		分析：主 goroutine 接收到任务处理完成的通知，或者超时后就返回了
			如果发生超时，process 函数就返回了，这就会导致 unbuffered 的 chan 从来就没有被读取
			unbuffered chan 必须等 reader 和 writer 都准备好了才能交流，否则就会阻塞
			超时导致未读，结果就是子 goroutine 就阻塞在第 7 行永远结束不了，进而导致 goroutine 泄漏
		解决方案
			将 unbuffered chan 改成容量为 1 的 chan，这样 'ch <- true ' 就不会被阻塞了
	Channel vs 并发原语
		Go 的开发者极力推荐使用 Channel
		不过，这两年，大家意识到，Channel 并不是处理并发问题的“银弹”，有时候使用并发原语更简单，而且不容易出错
	选择总结
		1. 共享资源的并发访问使用传统并发原语
		2. 复杂的任务编排和消息传递使用 Channel
		3. 消息通知机制使用 Channel，除非只想 signal 一个 goroutine，才使用 Cond
		4. 简单等待所有任务的完成用 WaitGroup，也有 Channel 的推崇者用 Channel，都可以
		5. 需要和 Select 语句结合，使用 Channel
		6. 需要和超时配合时，使用 Channel 和 Context

它们踩过的坑
	etcd issue 6857：一个程序 hang 住的问题
		描述
			在异常情况下，没有往 chan 实例中填充所需的元素，导致等待者永远等待
			具体来说，Status 方法的逻辑是生成一个 chan Status，然后把这个 chan 交给其它的 goroutine 去处理和写入数据
			最后，Status 返回获取的状态信息
		分析
			不幸的是，如果正好节点停止了，没有 goroutine 去填充这个 chan，会导致方法 hang 在返回的那一行上
			解决办法就是，在等待 status chan 返回元素的同时，也检查节点是不是已经停止了（done 这个 chan 是不是 close 了）
		图示
			13.channel_01_demo_etcd_01.jpg
		扩展
			感觉这个修改还是有问题的
			问题就在于，如果程序执行了 466 行，成功地把 c 写入到 Status 待处理队列后
			执行到第 467 行时，如果停止了这个节点，那么，这个 Status 方法还是会阻塞在第 467 行
	etcd issue 5505
		虽然没有任何的 Bug 描述，但是从修复内容上看，它是一个往已经 close 的 chan 写数据导致 panic 的问题
	etcd issue 11256：因为 unbuffered chan goroutine 泄漏的问题
		描述
			TestNodeProposeAddLearnerNode 方法中一开始定义了一个 unbuffered 的 chan，也就是 applyConfChan
			然后启动一个子 goroutine，这个子 goroutine 会在循环中执行业务逻辑，并且不断地往这个 chan 中添加一个元素
			TestNodeProposeAddLearnerNode 方法的末尾处会从这个 chan 中读取一个元素
		问题
			这段代码在 for 循环中就往此 chan 中写入了一个元素，结果导致 TestNodeProposeAddLearnerNode 从这个 chan 中读取到元素就返回了
			悲剧的是，子 goroutine 的 for 循环还在执行，阻塞在下图中红色的第 851 行，并且一直 hang 在那里
		图示
			13.channel_01_demo_etcd_02.jpg
		修复
			只要改动一下 applyConfChan 的处理逻辑就可以了：
			只有子 goroutine 的 for 循环中的主要逻辑完成之后，才往 applyConfChan 发送一个元素
			这样，TestNodeProposeAddLearnerNode 收到通知继续执行，子 goroutine 也不会被阻塞住了
	etcd issue 9956
		往一个已 close 的 chan 发送数据，其实它是 grpc 的一个 bug（grpc issue 2695）
			修复办法就是不 close 这个 chan 就好了
		图示
			13.channel_01_demo_etcd_02.jpg

总结
	chan 的值和状态有多种情况
		而不同的操作（send、recv、close）又可能得到不同的结果，这是使用 chan 类型时经常让人困惑的地方
	关注点
		那些 panic 的情况，另外还要掌握那些会 block 的场景，它们是导致死锁或者 goroutine 泄露的罪魁祸首
	注意
		只要一个 chan 还有未读的数据，即使把它 close 掉，你还是可以继续把这些未读的数据消费完，之后才是读取零值数据
	图示
		13.channel_01_error_panic.jpg

思考
	1.有一道经典的使用 Channel 进行任务编排的题，你可以尝试做一下：有四个 goroutine，编号为 1、2、3、4
		每秒钟会有一个 goroutine 打印出它自己的编号，要求你编写一个程序，让输出的编号总是按照 1、2、3、4、1、2、3、4、……的顺序打印出来
	2.chan T 是否可以给 <- chan T 和 chan<- T 类型的变量赋值？反过来呢？

重要
	1.不要在 Unlock 后 Sleep，否则死锁
	2.先 Unlock，再 Signal/Broadcast
*/

// ChannelMapGoQueue 方式六
func ChannelMapGoQueue() {
	const N = 4
	//chs := [N]chan struct{}{}	// 必须初始化
	chs := [N]chan struct{}{
		make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{}),
	}
	var wg sync.WaitGroup
	wg.Add(N * N)
	for i := 1; i <= N; i++ {
		for j := 0; j < N; j++ { // 每个任务 i 都有 N 个 go 抢，谁抢到谁执行
			go func(i, j int) {
				defer wg.Done()
				newChanTask(i, j, chs[i-1], chs[i%N])
			}(i, j)
		}
	}
	chs[0] <- struct{}{}
	//chs[0] <- struct{}{}
	wg.Wait()
}
func newChanTask(i, j int, r <-chan struct{}, w chan<- struct{}) {
	for {
		csp := <-r
		fmt.Println(i, j)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)+300))
		w <- csp
	}
}

// ChannelQueue 方式五：不是自己的数据，就再放回去
// ugly 且有问题
func ChannelQueue() {
	const N = 4
	var wg sync.WaitGroup
	wg.Add(N)
	ch := make(chan int, 1)
	for i := 1; i <= N; i++ {
		go func(i int) {
			defer wg.Done()
			for {
				select {
				case v := <-ch:
					if v == i {
						fmt.Println(i)
						ch <- v%N + 1
						time.Sleep(time.Second * 2)
					} else {
						ch <- v
					}
				default:
					time.Sleep(time.Millisecond * 200)
				}
			}
		}(i)
	}
	ch <- 1
	wg.Wait()
}

// AtomicQueue 方式四
func AtomicQueue() {
	const N = 4
	var cnt uint32 = 1
	next := func(v uint32, f func(i uint32)) { // next 函数实现了一种自旋（spinning：除非发现条件已满足，否则它会不断地进行检查）
		for {
			//if atomic.LoadUint32(&cnt) == v {
			//	f(v)
			//	atomic.CompareAndSwapUint32(&cnt, v, v%N+1)
			//}
			//time.Sleep(time.Millisecond * 1000)
			if atomic.CompareAndSwapUint32(&cnt, v, v%N+1) {
				f(v)
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)+500))
			}
			time.Sleep(time.Millisecond * 200)
		}
	}
	fn := func(i uint32) {
		fmt.Println(i)
	}
	for i := uint32(1); i <= N; i++ {
		//atomic.AddUint32(&cnt, ^uint32(0))
		go func(i uint32) {
			next(i, fn)
		}(i)
	}
	//next(math.MaxInt32, fn)
	select {}
}

// CondQueue 方式一二三
func CondQueue() {
	const N = 4
	c := sync.NewCond(&Mutex{})
	var cnt int32
	f := func(i int32) {
		c.L.Lock()
		c.Wait()
		time.Sleep(time.Millisecond * 500)
		fmt.Println(i)
		c.L.Unlock()
		c.Signal()
		//time.Sleep(time.Millisecond * 500) // 不要在 Unlock 后 Sleep，否则死锁
	}
	for i := int32(1); i <= N; i++ {
		// 方式一
		//go func(i int32, f func(int32)) {
		//	for !atomic.CompareAndSwapInt32(&cnt, int32(i-1), int32(i)) {
		//		time.Sleep(time.Millisecond * 100)
		//	}
		//	for {
		//		f(i)
		//	}
		//}(i, f)

		// 方式二
		go func(f func(int32)) {
			v := atomic.AddInt32(&cnt, 1)
			for {
				f(v)
			}
		}(f)

		// 方式三
		//go func(i int32) {
		//	v := atomic.AddInt32(&cnt, 1) // 可以让 cnt 初始值为 5，然后 -1
		//	//var x uint32
		//	//y := atomic.AddUint32(&x, ^uint32(0))
		//	for {
		//		c.L.Lock()
		//		c.Wait()
		//		if atomic.LoadInt32(&cnt) == i {
		//			atomic.StoreInt32(&cnt, v%N+1)
		//			time.Sleep(time.Millisecond * 500)
		//			fmt.Println(v)
		//		}
		//		c.L.Unlock()
		//		c.Signal()
		//		//c.Broadcast()
		//	}
		//}(i)
	}
	for {
		//fmt.Println(atomic.LoadInt32(&cnt))
		if atomic.LoadInt32(&cnt) == N {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	c.Signal()
	select {}
}
