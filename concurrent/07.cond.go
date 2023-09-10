package concurrent

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

/*
Cond：条件变量的实现机制及避坑指南

Java面试问题：等待 / 通知（wait/notify）机制
	请实现一个限定容量的队列（queue），当队列满或者空的时候，利用等待 / 通知机制实现阻塞或者唤醒
Go
	可以实现一个类似的限定容量的队列，而且实现起来也比较简单，只要用条件变量（Cond）并发原语就可以
	Cond 并发原语相对来说不是那么常用，但是在特定的场景使用会事半功倍，比如你需要在唤醒一个或者所有的等待者做一些检查操作的时候

Go 标准库的 Cond
	Go 标准库提供 Cond 原语的目的
		为等待 / 通知场景下的并发问题提供支持
		Cond 通常应用于等待某个条件的一组 goroutine，等条件变为 true 的时候，其中一个 goroutine 或者所有的 goroutine 都会被唤醒执行
	原理
		Cond 是和某个条件相关，这个条件需要一组 goroutine 协作共同完成
		在条件还没有满足的时候，所有等待这个条件的 goroutine 都会被阻塞住
		只有这一组 goroutine 通过协作达到了这个条件，等待的 goroutine 才可能继续进行下去
	等待的条件
		可以是某个变量达到了某个阈值或者某个时间点，也可以是一组变量分别都达到了某个阈值，还可以是某个对象的状态满足了特定的条件
		总结来讲，等待的条件是一种可以用来计算结果是 true 还是 false 的条件
	开发中
		真正使用 Cond 的场景比较少，因为一旦遇到需要使用 Cond 的场景，我们更多地会使用 Channel 的方式去实现，因为那才是更地道的 Go 语言的写法
		甚至 Go 的开发者有个“把 Cond 从标准库移除”的提议（issue 21165）
		而有的开发者认为，Cond 是唯一难以掌握的 Go 并发原语

Cond 的基本用法
	初始化
		标准库中的 Cond 并发原语初始化的时候，需要关联一个 Locker 接口的实例，一般我们使用 Mutex 或者 RWMutex
	API
		type Cond
			func NeWCond(l Locker) *Cond
			func (c *Cond) Broadcast()
			func (c *Cond) Signal()
			func (c *Cond) Wait()
		Cond 关联的 Locker 实例可以通过 c.L 访问，它内部维护着一个先入先出的等待队列
	Signal 方法
		允许调用者 Caller 唤醒一个等待此 Cond 的 goroutine
			如果此时没有等待的 goroutine，显然无需通知 waiter
			如果 Cond 等待队列中有一个或者多个等待的 goroutine，则需要从等待队列中移除第一个 goroutine 并把它唤醒
		在其他编程语言中，比如 Java 语言中，Signal 方法也被叫做 notify 方法
		调用 Signal 方法时，不强求你一定要持有 c.L 的锁
	Broadcast 方法
		允许调用者 Caller 唤醒所有等待此 Cond 的 goroutine
			如果此时没有等待的 goroutine，显然无需通知 waiter
			如果 Cond 等待队列中有一个或者多个等待的 goroutine，则清空所有等待的 goroutine，并全部唤醒
		在其他编程语言中，比如 Java 语言中，Broadcast 方法也被叫做 notifyAll 方法。
		同样地，调用 Broadcast 方法时，也不强求你一定持有 c.L 的锁
	Wait 方法
		会把调用者 Caller 放入 Cond 的等待队列中并阻塞，直到被 Signal 或者 Broadcast 的方法从等待队列中移除并唤醒
		调用 Wait 方法时必须要持有 c.L 的锁
	通用方法名
		Go 实现的 sync.Cond 的方法名是 Wait、Signal 和 Broadcast，这是计算机科学中条件变量的通用方法名
		比如，C 语言中对应的方法名是 pthread_cond_wait、pthread_cond_signal 和 pthread_cond_broadcast
	示例：百米赛跑
		CondRunner & TestCondRunner
		然而示例代码没有实现单独唤醒“裁判”，最后裁判 Broadcast 所有运动员，最后运动员执行“赛跑”
	Cond 的复杂在于
		一，这段代码有时候需要加锁，有时候可以不加
		二，Wait 唤醒后需要检查条件
		三，条件变量的更改，其实是需要原子操作或者互斥锁保护的
		所以，有的开发者会认为，Cond 是唯一难以掌握的 Go 并发原语

Cond 的实现原理
	Cond 的实现非常简单，或者说复杂的逻辑已经被 Locker 或者 runtime 的等待队列实现了
	runtime/sema.go
		runtime_notifyListXXX 是运行时实现的方法，实现了一个等待 / 通知的队列
		参见源码 runtime/sema.go
	copyChecker
		是一个辅助结构，可以在运行时检查 Cond 是否被复制使用
	copyChecker vs noCopy
		noCopy，是一个辅助的、用来帮助 vet 检查用的类型
		而Cond还有个copyChecker 是一个辅助结构，可以在运行时检查 Cond 是否被复制使用
		nocpoy是静态检查，copyChecker是运行时检查，分别应用不同阶段。而像 Mutex 只有静态检查
	Signal 和 Broadcast
		只涉及到 notifyList 数据结构，不涉及到锁
	Wait
		把调用者加入到等待队列时会释放锁，在被唤醒之后还会请求锁
		在阻塞休眠期间，调用者是不持有锁的，这样能让其他 goroutine 有机会检查或者更新等待变量

使用 Cond 的 2 个常见错误
	调用 Wait 的时候没有加锁：Cond 最常见的使用错误
	示例：CondRunner
		注释 c.L.Lock() 和 c.L.Unlock()
		报错
			fatal error: sync: unlock of unlocked mutex
		原因分析
			cond.Wait 方法的实现是，把当前调用者加入到 notify 队列之中后会释放锁，然后一直等待
			如果不释放锁，其他 Wait 的调用者就没有机会加入到 notify 队列中了
			等调用者被唤醒之后，又会去争抢这把锁
			如果调用 Wait 之前不加锁的话，就有可能 Unlock 一个未加锁的 Locker
		切记
			调用 cond.Wait 方法之前一定要加锁
	没有检查条件是否满足程序就继续执行了
		只调用了一次 Wait，没有检查等待条件是否满足，结果条件没满足，程序就继续执行了
	示例：CondRunner
		注释 for ready != 10 {
		原因
			误以为 Cond 的使用，就像 WaitGroup 那样调用一下 Wait 方法等待那么简单
			每一个运动员准备好之后都会唤醒所有的等待者，也就是这里的裁判员
			比如第一个运动员准备好后就唤醒了裁判员，结果这个裁判员傻傻地没做任何检查，以为所有的运动员都准备好了，就继续执行了
		现象
			可能只有几个运动员准备好之后程序就运行完了，而不是我们期望的所有运动员都准备好才进行下一步
		切记
			waiter goroutine 被唤醒不等于等待条件被满足，只是有 goroutine 把它唤醒了而已
			等待条件有可能已经满足了，也有可能不满足，我们需要进一步检查
			你也可以理解为，等待者被唤醒，只是得到了一次检查的机会而已

知名项目中 Cond 的使用
	Cond 在实际项目中被使用的机会比较少，原因总结起来有两个
		第一，同样的场景我们会使用其他的并发原语来替代
			Go 特有的 Channel 类型，有一个应用很广泛的模式就是通知机制，这个模式使用起来也特别简单
			所以很多情况下，我们会使用 Channel 而不是 Cond 实现 wait/notify 机制
		第二，对于简单的 wait/notify 场景，比如等待一组 goroutine 完成之后继续执行余下的代码，我们会使用 WaitGroup 来实现
			因为 WaitGroup 的使用方法更简单，而且不容易出错
			比如，上面百米赛跑的问题，就可以很方便地使用 WaitGroup 来实现
	sync.Cond 的路越走越窄
		标准库内部有几个地方使用了 Cond，比如 io/pipe.go 等，后来都被其他的并发原语（比如 Channel）替换了
	Cond vs Channel：忠实“粉丝”坚持使用 Cond 的原因在于 Cond 有三点特性是 Channel 无法替代的
		Cond 和一个 Locker 关联，可以利用这个 Locker 对相关的依赖条件更改提供保护
		Cond 可以同时支持 Signal 和 Broadcast 方法，而 Channel 只能同时支持其中一种
		Cond 的 Broadcast 方法可以被重复调用
			等待条件再次变成不满足的状态后，我们又可以调用 Broadcast 再次唤醒等待的 goroutine
			这也是 Channel 不能支持的，Channel 被 close 掉了之后不支持再 open
	Kubernetes：开源项目中鲜有的Cond案例
		Kubernetes 项目中定义了优先级队列 PriorityQueue 这样一个数据结构，用来实现 Pod 的调用
			它内部有三个 Pod 的队列，即 activeQ、podBackoffQ 和 unschedulableQ
			其中 activeQ 就是用来调度的活跃队列（heap）
			Pop 方法调用的时候，如果这个队列为空，并且这个队列没有 Close 的话，会调用 Cond 的 Wait 方法等待
			在调用 Wait 方法的时候，调用者是持有锁的，并且被唤醒的时候检查等待条件（队列是否为空）
		Pop()	// 从队列中取出一个元素
			func (p *PriorityQueue) Pop() (*framework.QueuedPodInfo, error) {
				p.lock.Lock()
				defer p.lock.Unlock()
				for p.activeQ.Len() == 0 { // 如果队列为空
					if p.closed {
						return nil, fmt.Errorf(queueClosed)
					}
					p.cond.Wait() // 等待，直到被唤醒
				}
				......
				return pInfo, err
			}
		当 activeQ 增加新的元素时，会调用条件变量的 Boradcast 方法，通知被 Pop 阻塞的调用者
		Add()	// 增加元素到队列中
			func (p *PriorityQueue) Add(pod *v1.Pod) error {
				p.lock.Lock()
				defer p.lock.Unlock()
				pInfo := p.newQueuedPodInfo(pod)
				if err := p.activeQ.Add(pInfo); err != nil {//增加元素到队列中
					klog.Errorf("Error adding pod %v to the scheduling queue: %v", nsNameFor
					return err
				}
				......
				p.cond.Broadcast() //通知其它等待的goroutine，队列中有元素了
				return nil
			}
		这个优先级队列被关闭的时候，也会调用 Broadcast 方法，避免被 Pop 阻塞的调用者永远 hang 住
		Close()
			func (p *PriorityQueue) Close() {
				p.lock.Lock()
				defer p.lock.Unlock()
				close(p.stop)
				p.closed = true
				p.cond.Broadcast() //关闭时通知等待的goroutine，避免它们永远等待
			}

总结
	工程中
		处理等待 / 通知的场景时，我们常常会使用 Channel 替换 Cond，因为 Channel 类型使用起来更简洁，而且不容易出错
		但是对于需要重复调用 Broadcast 的场景，比如上面 Kubernetes 的例子
		每次往队列中成功增加了元素后就需要调用 Broadcast 通知所有的等待者，使用 Cond 就再合适不过了
	Cond 常见错误
		Wait 调用需要加锁，以及被唤醒后一定要检查条件是否真的已经满足
	Cond vs WaitGroup
		WaitGroup 是主 goroutine 等待确定数量的子 goroutine 完成任务
		而 Cond 是等待某个条件满足，这个条件的修改可以被任意多的 goroutine 更新
		而且 Cond 的 Wait 不关心也不知道其他 goroutine 的数量，只关心等待条件
		而且 Cond 还有单个通知的机制，也就是 Signal 方法

思考
	1.一个 Cond 的 waiter 被唤醒的时候，为什么需要再检查等待条件，而不是唤醒后进行下一步？
	2.你能否利用 Cond 实现一个容量有限的 queue？
		参考 Kubernetes 案例
	3.Kubernetes 案例中，为什么使用 Cond 这个并发原语，能不能换成 Channel 实现呢？
*/

func CondRunner() {
	c := sync.NewCond(&sync.Mutex{})
	var ready int
	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Duration(rand.Int63n(3)) * time.Second)
			c.L.Lock()
			ready++
			c.L.Unlock()
			log.Printf("运动员 #%d 已准备就绪\n", i)
			c.Broadcast() // 广播唤醒所有等待者
		}(i)
	}
	c.L.Lock()        // 调用 Wait 的时候没有加锁：Cond 最常见的使用错误
	for ready != 10 { // Lock 放在 for 外面的原因是可以利用锁保护共享数据的读写。wait总是需要锁
		c.Wait()
		fmt.Println("裁判员被唤醒一次")
	}
	c.L.Unlock() //所有的运动员是否就绪
	log.Println("所有运动员都准备就绪。比赛开始，3,2,1 ...")
}
