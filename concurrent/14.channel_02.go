package concurrent

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

/*
Channel：透过代码看典型的应用模式

使用反射操作 Channel
	通过反射的方式执行 select 语句，在处理很多的 case clause，尤其是不定长的 case clause 的时候，非常有用
	需求分析
		chan 的数量在编译的时候是不定的，在运行的时候需要处理一个 slice of chan
	reflect.Select 函数
		Go 的 select 是伪随机的，它可以在执行的 case 中随机选择一个 case，并把选择的这个 case 的索引（chosen）返回
		如果没有可用的 case 返回，会返回一个 bool 类型的返回值，这个返回值用来表示是否有 case 成功被选择
		如果是 recv case，还会返回接收的元素
		func Select(cases []SelectCase) (chosen int, recv Value, recvOK bool)
	示例
		ChannelReflectSelect & TestChannelReflectSelect

典型的应用场景
消息交流
	简介
		从 chan 的内部实现看，它是以一个循环队列的方式存放数据，所以，它有时候也会被当成线程安全的队列和 buffer 使用
		一个 goroutine 可以安全地往 Channel 中塞数据，另外一个 goroutine 可以安全地从 Channel 中读取数据，goroutine 就可以安全地实现信息交流了
	案例：worker 池
		描述
			Marcio Castilho 在 '使用 Go 每分钟处理百万请求' 这篇文章中，就介绍了他们应对大并发请求的设计
			他们将用户的请求放在一个 chan Job 中，这个 chan Job 就相当于一个待处理任务队列
			除此之外，还有一个 chan chan Job 队列，用来存放可以处理任务的 worker 的缓存队列
			link：http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/
		原理
			dispatcher 会把待处理任务队列中的任务放到一个可用的缓存队列中，worker 会一直处理它的缓存队列
			通过使用 Channel，实现了一个 worker 池的任务处理中心，并且解耦了前端 HTTP 请求处理和后端任务处理的逻辑
		三方 worker 池：参见 10.pool.go
			它们全部都是通过 Channel 实现的，这是 Channel 的一个常见的应用场景
			worker 池的生产者和消费者的消息交流都是通过 Channel 实现的
	案例：etcd
		etcd 中的 node 节点的实现，包含大量的 chan 字段
		比如 recvc 是消息处理的 chan，待处理的 protobuf 消息都扔到这个 chan 中，node 有一个专门的 run goroutine 负责处理这些消息
		14.channel_02_demo_etcd.jpg
数据传递
	类比
		“击鼓传花”的游戏很多人都玩过，花从一个人手中传给另外一个人，就有点类似流水线的操作
		这个花就是数据，花在游戏者之间流转，这就类似编程中的数据传递
	示例：13.channel_01.go 思考
		ChannelMapDataTrans
		为了实现顺序的数据传递，我们可以定义一个令牌的变量，谁得到令牌，谁就可以打印一次自己的编号，同时将令牌传递给下一个 goroutine
	场景特点
		当前持有数据的 goroutine 都有一个信箱，信箱使用 chan 实现
		goroutine 只需要关注自己的信箱中的数据，处理完毕后，就把结果发送到下一家的信箱中
信号通知
	实现 wait/notify 的设计模式
		chan 类型有这样一个特点：
		chan 如果为空，那么，receiver 接收数据的时候就会阻塞等待，直到 chan 被关闭或者有新的数据到来
		利用这个机制，我们可以实现 wait/notify 的设计模式
	Cond
		并发原语 Cond 也能实现这个功能
		但是，Cond 使用起来比较复杂，容易出错，而使用 chan 实现 wait/notify 模式，就方便多了
	退出前清理：doCleanup
		除了正常的业务处理时的 wait/notify，我们经常碰到的一个场景，就是程序关闭的时候，需要在退出之前做一些清理（doCleanup 方法）的动作
		这个时候，我们经常要使用 chan
	示例：优雅退出
		func main() {
			go func() {
				...... // 执行业务处理
			}()

			// 处理CTRL+C等中断信号
			termChan := make(chan os.Signal)
			signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
			<-termChan

			// 执行退出之前的清理动作
			doCleanup()

			fmt.Println("优雅退出")
		}
	“两阶段”优雅退出
		需求
			有时候，doCleanup 可能是一个很耗时的操作，比如十几分钟才能完成，如果程序退出需要等待这么长时间，用户是不能接受的
			所以，在实践中，我们需要设置一个最长的等待时间。只要超过了这个时间，程序就不再等待，可以直接退出
		退出的时候分为两个阶段
			1. closing，代表程序退出，但是清理工作还没做
			2. closed，代表清理工作已经做完
		示例代码
			ChannelShutdownDoCleanup & main/channel.go
			CTRL+c 需要在 main 中测试
锁
	使用 chan 也可以实现互斥锁
	happens-before 的关系：是指事件发生的先后顺序关系
		在 chan 的内部实现中，就有一把互斥锁保护着它的所有字段
		从外在表现上，chan 的发送和接收之间也存在着 happens-before 的关系，保证元素放进去之后，receiver 才能读取到
	使用 chan 实现互斥锁，至少有两种方式
		一种方式是先初始化一个 capacity 等于 1 的 Channel，然后再放入一个元素
			这个元素就代表锁，谁取得了这个元素，就相当于获取了这把锁
		另一种方式是，先初始化一个 capacity 等于 1 的 Channel
			它的“空槽”代表锁，谁能成功地把元素发送到这个 Channel，谁就获取了这把锁
	示例：第一种方式
		ChannelMutexDemo
	原理
		在初始化这个锁的时候往 Channel 中先塞入一个元素，谁把这个元素取走，谁就获取了这把锁，把元素放回去，就是释放了锁
		元素在放回到 chan 之前，不会有 goroutine 能从 chan 中取出元素的，这就保证了互斥性
		利用 select+chan 的方式，很容易实现 TryLock、Timeout 的功能
		具体来说就是，在 select 语句中，我们可以使用 default 实现 TryLock，使用一个 Timer 来实现 Timeout 的功能
	优势
		channel 实现互斥锁的优势是 TryLock 和 Timeout，如果不需要这些特性，建议使用 Mutex
		现在 sync.Mutex 也有 TryLock 了，而且有了公平性的饥饿标记 mutexStarving
任务编排
	数据传递
		消息交流的场景是一个特殊的任务编排的场景，这个“击鼓传花”的模式也被称为流水线模式
	等待模式
		我们可以利用 WaitGroup 实现等待模式：启动一组 goroutine 执行任务，然后等待这些任务都完成
		其实，我们也可以使用 chan 实现 WaitGroup 的功能
	编排
		编排既指安排 goroutine 按照指定的顺序执行，也指多个 chan 按照指定的方式组合处理的方式
		goroutine 的编排类似“击鼓传花”的例子，我们通过编排数据在 chan 之间的流转，就可以控制 goroutine 的执行
	五种编排方式
		分别是 Or-Done 模式、扇入模式、扇出模式、Stream 和 map-reduce
Or-Done 模式
	Or-Done 模式是信号通知模式中更宽泛的一种模式
	“信号通知模式”
		我们会使用“信号通知”实现某个任务执行完成后的通知机制
		在实现时，我们为这个任务定义一个类型为 chan struct{}类型的 done 变量
		等任务结束后，我们就可以 close 这个变量，然后，其它 receiver 就会收到这个通知
	Or-Done
		这是有一个任务的情况，如果有多个任务，只要有任意一个任务执行完，我们就想获得这个信号，这就是 Or-Done 模式
		比如，你发送同一个请求到多个微服务节点，只要任意一个微服务节点返回结果，就算成功
	示例：递归 / 反射 / goroutine
		OrDone & TestOrDone
		1. 当 chan 的数量大于 2 时，使用递归的方式等待信号
		2. 在 chan 数量比较多的情况下，递归并不是一个很好的解决方式，而使用 reflect.Select
			反射方式避免了深层递归的情况，可以处理有大量 chan 的情况
		3. 最笨的一种方法就是为每一个 Channel 启动一个 goroutine
			不过这会启动非常多的 goroutine，太多的 goroutine 会影响性能，所以不太常用
扇入模式
	简介
		扇入借鉴了数字电路的概念，它定义了单个逻辑门能够接受的数字信号输入最大量的术语
		一个逻辑门可以有多个输入，一个输出
	Channel 扇入模式
		在软件工程中，模块的扇入是指有多少个上级模块调用它
		而对于我们这里的 Channel 扇入模式来说，就是指有多个源 Channel 输入、一个目的 Channel 输出的情况
		扇入比就是源 Channel 数量比 1
	实现思路
		每个源 Channel 的元素都会发送给目标 Channel
		相当于目标 Channel 的 receiver 只需要监听目标 Channel，就可以接收所有发送给源 Channel 的数据
		扇入模式也可以使用反射、递归，或者是用最笨的每个 goroutine 处理一个 Channel 的方式来实现
	示例
		FanInReflect：reflect 方式
		FanInRec：递归方式
扇出模式
	扇出模式是和扇入模式相反的
	扇出模式 & 观察者模式
		扇出模式只有一个输入源 Channel，有多个目标 Channel，扇出比就是 1 比目标 Channel 数的值
		经常用在设计模式中的观察者模式中（观察者设计模式定义了对象间的一种一对多的组合关系。这样一来，一个对象的状态发生变化时，所有依赖于它的对象都会得到通知并自动刷新）
		在观察者模式中，数据变动后，多个观察者都会收到这个变更信号
	示例
		FanOut
		也可以尝试使用反射方式来实现
	补充
		可以使用 waitgroup，等待异步的 go 全部写入完
Stream
	把 Channel 当作流式管道使用
		也就是把 Channel 看作流（Stream），提供跳过几个元素，或者是只取其中的几个元素等方法
	示例：AsStream
		func AsStream(done <-chan struct{}, values ...interface{}) <-chan interface{} { // 把一个数据 slice 转换成流
			s := make(chan interface{}) //创建一个unbuffered的channel
			go func() {                 // 启动一个goroutine，往s中塞数据
				defer close(s)             // 退出时关闭chan
				for _, v := range values { // 遍历数组
					select {
					case <-done:
						return
					case s <- v: // 将数组元素塞入到chan中
					}
				}
			}()
			return s
		}
	实现流的方法
		takeN：只取流中的前 n 个数据
		takeFn：筛选流中的数据，只保留满足条件的数据
		takeWhile：只取前面满足条件的数据，一旦不满足条件，就不再取
		skipN：跳过流中前几个数据
		skipFn：跳过满足条件的数据
		skipWhile：跳过前面满足条件的数据，一旦不满足条件，当前这个元素和以后的元素都会输出给 Channel 的 receiver
	示例
		TakeN
map-reduce
	简介
		map-reduce 是一种处理数据的方式，最早是由 Google 公司研究提出的一种面向大规模数据处理的并行计算模型和方法
		开源的版本是 hadoop，前几年比较火
	map-reduce
		map-reduce 分为两个步骤
		第一步是映射（map），处理队列中的数据
		第二步是规约（reduce），把列表中的每一个元素按照一定的处理方式处理成结果，放入到结果队列中
		就像做汉堡一样，map 就是单独处理每一种食材，reduce 就是从每一份食材中取一部分，做成一个汉堡
	示例：并不是分布式的 map-reduce，而是单机单进程的 map-reduce 方法
		MapChanReduce & TestMapChanReduce

总结
	虽然 Channel 最初是基于 CSP 设计的用于 goroutine 之间的消息传递的一种数据类型
		但是，除了消息传递这个功能之外，大家居然还演化出了各式各样的应用模式
		我不确定 Go 的创始人在设计这个类型的时候，有没有想到这一点，但是，我确实被各位大牛利用 Channel 的各种点子折服了
		比如有人实现了一个基于 TCP 网络的分布式的 Channel
		link：https://github.com/chrislusf/netchan
	在使用 Go 开发程序的时候，你也不妨多考虑考虑是否能够使用 chan 类型，看看你是不是也能创造出别具一格的应用模式

思考
	利用 chan 实现互斥锁，如果 buffer 设置的不是 1，而是一个更大的值，会出现什么状况吗？能解决什么问题吗？

补充
	反射的效率很低
		Francesc Campoy 有一篇文章专门做了测试，analyzing the performance of go functions with benchmarks
	sync.Once 用于 close(chan)
		避免重复 close panic
	在看Marcio Castilho 在 使用 Go 每分钟处理百万请求的例子的时候我很困惑，我认为双层管道是没有意义的，生产者直接发送job给消费者，额定数量的消费者直接进行并发接收并处理就可以达到同样控制并发的效果
	为什么非要消费者把一个管道交给生产者，生产者在把job通过管道传递给消费者。我想请问一下这个步骤的作用是什么？
		作者回复: 我也觉得他的实现方案有问题，有点感觉是从其它语言的实现复制过来的，有点设计过度
		我和你的设计方案也有些不同，你的方案类似scala的actor模式。这种模式有个问题，如果某个actor处理慢了，可能会导致任务堆积在这个actor上
		go 很常规的操作就是用一个 channel 实现 buffer ,各个 goroutine 都从这个 buffer 获取任务就好了
	channel适不适用于传输大文件？还有就是传输大文件的性能情况是怎样的？
		channel是线程内的，你怎么有传递大文件的需求，什么是大文件传输？
		channel可以传递大的byte slice，没有性能问题。要是大数组就不合适了


*/

// MapChanReduce ==========map-reduce 示例==========
func MapChanReduce() {
	in := asStream(nil)
	// map操作: 乘以10
	mapFn := func(v interface{}) interface{} {
		return v.(int) * 10
	}
	// reduce操作: 对map的结果进行累加
	reduceFn := func(r, v interface{}) interface{} {
		return r.(int) + v.(int)
	}
	sum := reduce(mapChan(in, mapFn), reduceFn) //返回累加结果
	fmt.Println(sum)
}
func asStream(done <-chan struct{}) <-chan interface{} { // 生成一个数据流
	s := make(chan interface{})
	values := []int{1, 2, 3, 4, 5}
	go func() {
		defer close(s)
		for _, v := range values { // 从数组生成
			select {
			case <-done:
				return
			case s <- v:
			}
		}
	}()
	return s
}
func mapChan(in <-chan interface{}, fn func(interface{}) interface{}) <-chan interface{} {
	out := make(chan interface{}) //创建一个输出chan
	if in == nil {                // 异常检查
		close(out)
		return out
	}
	go func() { // 启动一个goroutine,实现map的主要逻辑
		defer close(out)
		for v := range in { // 从输入chan读取数据，执行业务操作，也就是map操作
			out <- fn(v)
		}
	}()
	return out
}
func reduce(in <-chan interface{}, fn func(r, v interface{}) interface{}) interface{} {
	if in == nil { // 异常检查
		return nil
	}
	out := <-in         // 先读取第一个元素
	for v := range in { // 实现reduce的主要逻辑
		out = fn(out, v)
	}
	return out
}

// TakeN ==========Stream 示例==========
func TakeN(done <-chan struct{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{}) // 创建输出流
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ { // 只读取前num个元素
			select {
			case <-done:
				return
			case takeStream <- <-valueStream: //从输入流中读取元素
			}
		}
	}()
	return takeStream
}

// FanOut ==========扇出模式示例==========
func FanOut(ch <-chan interface{}, out []chan interface{}, async bool) {
	go func() {
		defer func() { //退出时关闭所有的输出chan
			for i := 0; i < len(out); i++ {
				close(out[i])
			}
		}()
		for v := range ch { // 从输入chan中读取数据
			v := v // 因为异步，所以把值记录下来
			for i := 0; i < len(out); i++ {
				i := i
				if async { //异步
					go func() {
						out[i] <- v // 放入到输出chan中,异步方式
					}()
				} else {
					out[i] <- v // 放入到输出chan中，同步方式
				}
			}
		}
	}()
}

// FanInReflect ==========扇入模式，reflect 示例==========
func FanInReflect(chs ...<-chan any) <-chan any {
	out := make(chan any)
	go func() {
		defer close(out)
		var cases []reflect.SelectCase // 构造SelectCase slice
		for _, c := range chs {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}
		for n := len(cases); n > 0; n = len(cases) { // 循环，从cases中选择一个可用的
			chosen, recv, ok := reflect.Select(cases)
			if !ok { // 此channel已经close
				cases[chosen], cases[n-1] = cases[n-1], cases[chosen]
				cases = cases[:n-1]
				continue
			}
			out <- recv
		}
	}()
	return out
}

// FanInRec ==========扇入模式，recursion 示例==========
func FanInRec(chans ...<-chan interface{}) <-chan interface{} {
	switch len(chans) {
	case 0:
		c := make(chan interface{})
		close(c)
		return c
	case 1:
		return chans[0]
	case 2:
		return mergeTwo(chans[0], chans[1])
	default:
		m := len(chans) / 2
		return mergeTwo(
			FanInRec(chans[:m]...),
			FanInRec(chans[m:]...))
	}
}
func mergeTwo(a, b <-chan interface{}) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		for a != nil || b != nil { //只要还有可读的chan
			select {
			case v, ok := <-a:
				if !ok { // a 已关闭，设置为nil
					a = nil
					continue
				}
				c <- v
			case v, ok := <-b:
				if !ok { // b 已关闭，设置为nil
					b = nil
					continue
				}
				c <- v
			}
		}
	}()
	return c
}

// OrDone ==========OrDone示例==========
func OrDone(channels ...<-chan interface{}) <-chan interface{} {
	// 特殊情况，只有零个或者1个chan
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}
	orDone := make(chan interface{})

	// 递归
	//go func() {
	//	defer close(orDone)
	//
	//	switch len(channels) {
	//	case 2: // 2个也是一种特殊情况
	//		select {
	//		case <-channels[0]:
	//		case <-channels[1]:
	//		}
	//	default: //超过两个，二分法递归处理
	//		m := len(channels) / 2
	//		select {
	//		case <-OrDone(channels[:m]...):
	//		case <-OrDone(channels[m:]...):
	//		}
	//	}
	//}()

	// reflect.Select
	go func() {
		defer close(orDone) // 利用反射构建SelectCase
		var cases []reflect.SelectCase
		for _, c := range channels {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}
		reflect.Select(cases) // 随机选择一个可用的case
	}()
	return orDone
}

// ChannelMutex ==========使用chan实现互斥锁==========
type ChannelMutex struct {
	ch chan struct{}
}

func NewMutex() *ChannelMutex { // 使用锁需要初始化
	mu := &ChannelMutex{make(chan struct{}, 1)}
	mu.ch <- struct{}{}
	return mu
}

func (m *ChannelMutex) Lock() { // 请求锁，直到获取到
	<-m.ch
}

func (m *ChannelMutex) Unlock() { // 解锁
	select {
	case m.ch <- struct{}{}:
	default:
		panic("unlock of unlocked mutex")
	}
}

func (m *ChannelMutex) TryLock() bool { // 尝试获取锁
	select {
	case <-m.ch:
		return true
	default:
	}
	return false
}

func (m *ChannelMutex) LockTimeout(timeout time.Duration) bool { // 加入一个超时的设置
	timer := time.NewTimer(timeout)
	select {
	case <-m.ch:
		timer.Stop()
		return true
	case <-timer.C:
	}
	return false
}

func (m *ChannelMutex) IsLocked() bool { // 锁是否已被持有
	return len(m.ch) == 0
}

func ChannelMutexDemo() {
	m := NewMutex()
	ok := m.TryLock()
	fmt.Printf("locked %v\n", ok)
	ok = m.TryLock()
	fmt.Printf("locked %v\n", ok)
}

// ChannelShutdownDoCleanup ==========“两阶段”优雅退出==========
func ChannelShutdownDoCleanup() {
	var closing = make(chan struct{})
	var closed = make(chan struct{})

	go func() {
		// 模拟业务处理
		for {
			select {
			case <-closing:
				return
			default:
				// ....... 业务计算
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// 处理CTRL+C等中断信号
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	close(closing)
	// 执行退出之前的清理动作
	go doCleanup(closed)

	// 原示例代码
	//select {
	//case <-closed:
	//case <-time.After(time.Second):
	//	fmt.Println("清理超时，不等了")
	//}
	//fmt.Println("优雅退出")

	// 修改后示例代码：调用Server.Shutdown graceful结束
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	select {
	case <-closed:
		//case <-time.After(time.Second):
		//	fmt.Println("清理超时，不等了")
	case <-timeoutCtx.Done():
		// http.Server Shutdown
		server := &http.Server{
			Handler: nil, // TODO
			Addr:    ":8080",
		}
		if err := server.Shutdown(timeoutCtx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		close(closed)
	}
	fmt.Println("优雅退出")
}

func doCleanup(closed chan struct{}) {
	time.Sleep((time.Minute))
	close(closed)
}

// Token ==========数据传递示例==========
type Token struct{} // 令牌类型（Token）
// 定义一个创建 worker 的方法，这个方法会从它自己的 chan 中读取令牌
func newWorker(id int, ch chan Token, nextCh chan Token) {
	for { // 哪个 goroutine 取得了令牌，就可以打印出自己编号，因为需要每秒打印一次数据
		token := <-ch         // 取得令牌
		fmt.Println((id + 1)) // id从1开始
		time.Sleep(time.Second)
		nextCh <- token
	}
}
func ChannelMapDataTrans() {
	chs := []chan Token{make(chan Token), make(chan Token), make(chan Token), make(chan Token)}
	for i := 0; i < 4; i++ { // 创建4个worker
		go newWorker(i, chs[i], chs[(i+1)%4]) // 启动每个 worker 的 goroutine
	}
	chs[0] <- struct{}{} //首先把令牌交给第一个worker
	select {}
}

// ChannelReflectSelect ==========reflect.Select 示例==========
func ChannelReflectSelect() {
	var (
		ch1   = make(chan int, 10)
		ch2   = make(chan int, 10)
		cases = createCases(ch1, ch2)
	)
	for i := 0; i < 10; i++ {
		chosen, recv, ok := reflect.Select(cases)
		if recv.IsValid() {
			fmt.Println("recv:", i, cases[chosen].Dir, recv, ok)
		} else {
			fmt.Println("send:", i, cases[chosen].Dir, ok)
		}
	}

	// 读已关闭的 channel
	ch := make(chan int, 1)
	ch <- 99
	rs := []reflect.SelectCase{
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		},
	}
	<-ch
	close(ch)
	chosen, recv, ok := reflect.Select(rs)
	fmt.Println(chosen, recv, ok) // 0 0 false
}

// createCases 函数分别为每个 chan 生成了 recv case 和 send case，并返回一个 reflect.SelectCase 数组
func createCases(chs ...chan int) []reflect.SelectCase {
	var cases []reflect.SelectCase
	for _, ch := range chs { // 创建 recv case
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}
	for i, ch := range chs { // 创建 send case
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(ch),
			Send: reflect.ValueOf(i + 100),
		})
	}
	return cases
}
