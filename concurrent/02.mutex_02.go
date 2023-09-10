package concurrent

/*
Mutex：庖丁解牛看实现

Mutex 演进历史
	从一个简单易于理解的互斥锁的实现，到一个非常复杂的数据结构，这是一个逐步完善的过程
	Go 开发者们做了种种努力，精心设计，展现了匠心和精益求精的精神
设计并发数据接口
	逐步提升性能和公平性
	逐步设计一个完善的同步原语，并对复杂度、性能、结构设计的权衡考量有新的认识
Mutex 架构演进四个阶段
	02.mutex_02_steps.jpg
	“初版”
		Mutex 使用一个 flag 来表示锁是否被持有
	“给新人机会”
		照顾到新来的 goroutine，所以会让新的 goroutine 也尽可能地先获取到锁
	“多给些机会”
		照顾新来的和被唤醒的 goroutine，但是会带来饥饿问题
	“解决饥饿”
		加入了饥饿的解决方案

初版的互斥锁
	Russ Cox 在 2008 年提交的第一版 Mutex
		可以通过一个 flag 变量，标记当前的锁是否被某个 goroutine 持有
		如果这个 flag 的值是 1，就代表锁已经被持有，那么，其它竞争的 goroutine 只能等待
		如果这个 flag 的值是 0，就可以通过 CAS（compare-and-swap，或者 compare-and-set）将这个 flag 设置为 1，标识锁被当前的这个 goroutine 持有了
	CAS 指令简介
		CAS 指令将给定的值和一个内存地址中的值进行比较
		如果它们是同一个值，就使用新值替换内存地址中的值，这个操作是原子性的
		CAS 是实现互斥锁和同步原语的基础
	原子性
		原子性保证这个指令总是基于最新的值进行计算，如果同时有其它线程已经修改了这个值，那么，CAS 会返回失败
	示例代码：// ==========初版的互斥锁==========
		Mutex 结构体包含两个字段
			字段 key：是一个 flag，用来标识这个排外锁是否被某个 goroutine 所持有，如果 key 大于等于 1，说明这个排外锁已经被持有
			字段 sema：是个信号量变量，用来控制等待 goroutine 的阻塞休眠和唤醒
		图示
			02.mutex_02_step_1.jpg
		Lock
			调用 Lock 请求锁的时候，通过 xadd 方法进行 CAS 操作，xadd 方法通过循环执行 CAS 操作直到成功，保证对 key 加 1 的操作成功完成
			如果比较幸运，锁没有被别的 goroutine 持有，那么，Lock 方法成功地将 key 设置为 1，这个 goroutine 就持有了这个锁
			如果锁已经被别的 goroutine 持有了，那么，当前的 goroutine 会把 key 加 1
			而且还会调用 semacquire 方法，使用信号量将自己休眠，等锁释放的时候，信号量会将它唤醒
		Unlock
			持有锁的 goroutine 调用 Unlock 释放锁时，它会将 key 减 1
			如果当前没有其它等待这个锁的 goroutine，这个方法就返回了
			但是，如果还有等待此锁的其它 goroutine，那么，它会调用 semrelease 方法，利用信号量唤醒等待锁的其它 goroutine 中的一个
		实现原理
			初版的 Mutex 利用 CAS 原子操作，对 key 这个标志量进行设置
			key 不仅仅标识了锁是否被 goroutine 所持有，还记录了当前持有和等待获取锁的 goroutine 的数量
	危险功能
		释放锁
			Unlock 方法可以被任意的 goroutine 调用释放锁，即使是没持有这个互斥锁的 goroutine，也可以进行这个操作
			这是因为，Mutex 本身并没有包含持有这把锁的 goroutine 的信息，所以，Unlock 也不会对此进行检查
			Mutex 的这个设计一直保持至今
		危险原因
			其它 goroutine 可以强制释放锁，这是一个非常危险的操作
			因为在临界区的 goroutine 可能不知道锁已经被释放了，还会继续执行临界区的业务操作，这可能会带来意想不到的结果
			因为这个 goroutine 还以为自己持有锁呢，有可能导致 data race 问题
		“谁申请，谁释放”原则
			在使用 Mutex 的时候，必须要保证 goroutine 尽可能不去释放自己未持有的锁，一定要遵循“谁申请，谁释放”的原则
			在真实的实践中，我们使用互斥锁的时候，很少在一个方法中单独申请锁，而在另外一个方法中单独释放锁，一般都会在同一个方法中获取锁和释放锁
		其他语言
			其它语言（比如 Java 语言）的互斥锁的实现，Mutex 这一点和其它语言的互斥锁不同，所以，如果是从其它语言转到 Go 语言开发的同学，一定要注意
		错误示例
			经常会基于性能的考虑，及时释放掉锁，所以在一些 if-else 分支中加上释放锁的代码，代码看起来很臃肿
			而且，在重构的时候，也很容易因为误删或者是漏掉而出现死锁的现象
				type Foo struct {
					mu sync.Mutex
					count int
				}
				func (f *Foo) Bar() {
					f.mu.Lock()
					if f.count < 1000 {
						f.count += 3
						f.mu.Unlock() // 此处释放锁
						return
					}
					f.count++
					f.mu.Unlock() // 此处释放锁
					return
				}
		defer
			从 1.14 版本起，Go 对 defer 做了优化，采用更有效的内联方式，取代之前的生成 defer 对象到 defer chain 中
			defer 对耗时的影响微乎其微了，所以基本上修改成下面简洁的写法也没问题：
				func (f *Foo) Bar() {
					f.mu.Lock()
					defer f.mu.Unlock()
					if f.count < 1000 {
						f.count += 3
						return
					}
					f.count++
					return
				}
			好处就是 Lock/Unlock 总是成对紧凑出现，不会遗漏或者多调用，代码更少
		尽早释放锁
			但是，如果临界区只是方法中的一部分，为了尽快释放锁，还是应该第一时间调用 Unlock
			而不是一直等到方法返回时才释放
	新的问题
		初版的 Mutex 实现之后，Go 开发组又对 Mutex 做了一些微调
			比如把字段类型变成了 uint32 类型
			调用 Unlock 方法会做检查
			使用 atomic 包的同步原语执行原子操作等
			这些小的改动，都不是核心功能
		但是，初版的 Mutex 实现有一个问题：请求锁的 goroutine 会排队等待获取互斥锁
			虽然这貌似很公平，但是从性能上来看，却不是最优的
			因为如果我们能够把锁交给正在占用 CPU 时间片的 goroutine 的话，那就不需要做上下文的切换，在高并发的情况下，可能会有更好的性能
给新人机会
	Go 开发者在 2011 年 6 月 30 日的 commit 中对 Mutex 做了一次大的调整，调整后的 Mutex 实现如下
		虽然 Mutex 结构体还是包含两个字段，但是第一个字段已经改成了 state，它的含义也不一样了
			type Mutex struct {
				state int32
				sema uint32
			}
			const (
				mutexLocked = 1 << iota // mutex is locked
				mutexWoken
				mutexWaiterShift = iota
			)
		图示
			02.mutex_02_step_2.jpg
		state 是一个复合型的字段，一个字段包含多个意义，这样可以通过尽可能少的内存来实现互斥锁
			这个字段的第一位（最小的一位）来表示这个锁是否被持有，第二位代表是否有唤醒的 goroutine，剩余的位数代表的是等待此锁的 goroutine 数
			所以，state 这一个字段被分成了三部分，代表三个数据
	Lock
		对字段 state 的操作，和代码逻辑都变复杂
			func (m *Mutex) Lock() {
				// Fast path: 幸运case，能够直接获取到锁
				if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
					return
				}
				awoke := false	// 新来的 goroutine
				for {
					old := m.state	// 先保存当前锁的状态，且被唤醒的 goroutine 也会重新获取 state
					new := old | mutexLocked // 新状态加锁：即加锁状态
					if old&mutexLocked != 0 {	// 锁还没被释放
						new = old + 1<<mutexWaiterShift //等待者数量加一
					}
					if awoke {	// 被唤醒的 goroutine
						new &^= mutexWoken	// 新状态清除唤醒标志
					}
					if atomic.CompareAndSwapInt32(&m.state, old, new) {//设置新状态
						if old&mutexLocked == 0 { // 锁原状态未加锁
							break	// 只会有一个 go cas 成功，获取到锁
						}
						runtime.Semacquire(&m.sema) // 请求信号量
						awoke = true	// 被唤醒
					}
				}
			}
		首先是通过 CAS 检测 state 字段中的标志，如果没有 goroutine 持有锁，也没有等待持有锁的 gorutine
			那么，当前的 goroutine 就很幸运，可以直接获得锁
			这也是注释中的 Fast path 的意思
		awoke := false
			如果不够幸运，state 不是零值，那么就通过一个循环进行检查
			最新版的 Mutex 也是类似的实现逻辑
		给新人机会
			如果想要获取锁的 goroutine 没有机会获取到锁，就会进行休眠
			但是在锁释放唤醒之后，它并不能像先前一样直接获取到锁，还是要和正在请求锁的 goroutine 进行竞争
			这会给后来请求锁的 goroutine 一个机会，也让 CPU 中正在执行的 goroutine 有更多的机会获取到锁，在一定程度上提高了程序的性能
		awoke = true
			for 循环是不断尝试获取锁，如果获取不到，就通过 runtime.Semacquire(&m.sema) 休眠
			休眠醒来之后 awoke 置为 true，尝试争抢
		new := old | mutexLocked
			将当前的 flag 设置为加锁状态，如果能成功地通过 CAS 把这个新值赋予 state，就代表抢夺锁的操作成功了
			需要注意的是，如果成功地设置了 state 的值，但是之前的 state 是有锁的状态
			那么，state 只是清除 mutexWoken 标志或者增加一个 waiter 而已
		goroutine 不同来源不同状态下的处理逻辑
			请求锁的 goroutine 有两类，一类是新来请求锁的 goroutine，另一类是被唤醒的等待请求锁的 goroutine
			锁的状态也有两种：加锁和未加锁
			图示：02.mutex_02_step_2_state.jpg
	Unlock
		释放锁的逻辑
			func (m *Mutex) Unlock() {
				// Fast path: drop lock bit.
				new := atomic.AddInt32(&m.state, -mutexLocked) //去掉锁标志
				if (new+mutexLocked)&mutexLocked == 0 { //本来就没有加锁
					panic("sync: unlock of unlocked mutex")	// 解锁 panic
				}
				old := new
				for {
					if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken) != 0
						return	// 没有 waiter || 有唤醒的 go 或已被加锁
					}
					new = (old - 1<<mutexWaiterShift) | mutexWoken // 新状态，准备唤醒goroutine
					if atomic.CompareAndSwapInt32(&m.state, old, new) {	// Unlock 写入 state，都会写入 mutexWoken
						runtime.Semrelease(&m.sema)	// 唤醒休眠中的 go
						return
					}
					old = m.state
				}
			}
		先是尝试将持有锁的标识设置为未加锁的状态
			通过减 1 而不是将标志位置零的方式实现
			还会检测原来锁的状态是否已经未加锁的状态，如果是 Unlock 一个未加锁的 Mutex 会直接 panic
		此时还可能有一些等待这个锁的 goroutine（也称为 waiter）需要通过信号量的方式唤醒它们中的一个
			第一种情况
				如果没有其它的 waiter，说明对这个锁的竞争的 goroutine 只有一个，那就可以直接返回了
				如果这个时候有唤醒的 goroutine，或者是又被别人加了锁，那当前的这个 goroutine 就可以放心返回了
			第二种情况，如果有等待者，并且没有唤醒的 waiter，那就需要唤醒一个等待的 waiter
				在唤醒之前，需要将 waiter 数量减 1，并且将 mutexWoken 标志设置上，这样，Unlock 就可以返回了
	主要改动
		新来的 goroutine 也有机会先获取到锁，甚至一个 goroutine 可能连续获取到锁，打破了先来先得的逻辑
		这一版的 Mutex 已经给新来请求锁的 goroutine 一些机会，让它参与竞争，没有空闲的锁或者竞争失败才加入到等待队列中
多给些机会
	2015 年 2 月的改动中，如果新来的 goroutine 或者是被唤醒的 goroutine 首次获取不到锁
		它们就会通过自旋（spin，通过循环不断尝试，spin 的逻辑是在 runtime 实现的）的方式，尝试检查锁是否被释放
	在尝试一定的自旋次数后，再执行原来的逻辑
		func (m *Mutex) Lock() {
			// Fast path: 幸运之路，正好获取到锁
			if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
				return
			}
			awoke := false
			iter := 0
			for { // 不管是新来的请求锁的goroutine, 还是被唤醒的goroutine，都不断尝试请求锁
				old := m.state            // 先保存当前锁的状态
				new := old | mutexLocked  // 新状态设置加锁标志
				if old&mutexLocked != 0 { // 锁还没被释放
					if runtime_canSpin(iter) { // 还可以自旋
						if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
							atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
							awoke = true
						}
						runtime_doSpin()
						iter++
						continue // 自旋，再次尝试请求锁
					}
					new = old + 1<<mutexWaiterShift
				}
				if awoke { // 唤醒状态
					if new&mutexWoken == 0 {
						panic("sync: inconsistent mutex state")
					}
					new &^= mutexWoken // 新状态清除唤醒标记
				}
				if atomic.CompareAndSwapInt32(&m.state, old, new) {
					if old&mutexLocked == 0 { // 旧状态锁已释放，新状态成功持有了锁，直接
						break
					}
					runtime_Semacquire(&m.sema) // 阻塞等待
					awoke = true                // 被唤醒
					iter = 0
				}
			}
		}
	spin
		for 循环会重新检查锁是否释放
		对于临界区代码执行非常短的场景来说，这是一个非常好的优化
		因为临界区的代码耗时很短，锁很快就能释放，而抢夺锁的 goroutine 不用通过休眠唤醒方式等待调度，直接 spin 几次，可能就获得了锁
解决饥饿
	饥饿问题
		因为新来的 goroutine 也参与竞争，有可能每次都会被新来的 goroutine 抢到获取锁的机会
		在极端情况下，等待中的 goroutine 可能会一直获取不到锁
	优化历史
		2016 年 Go 1.9 中 Mutex 增加了饥饿模式，让锁变得更公平，不公平的等待时间限制在 1 毫秒
			并且修复了一个大 Bug：总是把唤醒的 goroutine 放在等待队列的尾部，会导致更加不公平的等待时间
		2018 年，Go 开发者将 fast path 和 slow path 拆成独立的方法，以便内联，提高性能
		2019 年也有一个 Mutex 的优化，虽然没有对 Mutex 做修改
			但是，对于 Mutex 唤醒后持有锁的那个 waiter，调度器可以有更高的优先级去执行，这已经是很细致的性能优化了
	Mutex 宗旨：现在的 Mutex 代码已经复杂得接近不可读的状态了。为了一个貌似很小的 feature 不得不将代码变得非常复杂
		Mutex 绝不容忍一个 goroutine 被落下，永远没有机会获取锁
		不抛弃不放弃是它的宗旨，而且它也尽可能地让等待较长的 goroutine 更有机会获取到锁
	代码
		02_mutex.go
	增加饥饿模式
		将饥饿模式的最大等待时间阈值设置成了 1 毫秒
		一旦等待者等待的时间超过了这个阈值，Mutex 的处理就有可能进入饥饿模式，优先让等待者先获取到锁
		新来的 go 主动谦让一下，给饥饿 go 一些机会
	公平性
		通过加入饥饿模式，可以避免把机会全都留给新来的 goroutine，保证了请求锁的 goroutine 获取锁的公平性
		对于我们使用锁的业务代码来说，不会有业务一直等待锁不被处理

饥饿模式和正常模式
	Mutex 可能处于两种操作模式下
		正常模式和饥饿模式
	fast path
		请求锁时调用的 Lock 方法中一开始是 fast path，这是一个幸运的场景，当前的 goroutine 幸运地获得了锁，没有竞争，直接返回
		否则就进入了 lockSlow 方法
		这样的设计，方便编译器对 Lock 方法进行内联，你也可以在程序开发中应用这个技巧
	正常模式
		正常模式下，waiter 都是进入先入先出队列，被唤醒的 waiter 并不会直接持有锁，而是要和新来的 goroutine 进行竞争
		新来的 goroutine 有先天的优势，它们正在 CPU 中运行，可能它们的数量还不少
		所以，在高并发情况下，被唤醒的 waiter 可能比较悲剧地获取不到锁，这时，它会被插入到队列的前面
		如果 waiter 获取不到锁的时间超过阈值 1 毫秒，那么，这个 Mutex 就进入到了饥饿模式
	饥饿模式
		在饥饿模式下，Mutex 的拥有者将直接把锁交给队列最前面的 waiter
		新来的 goroutine 不会尝试获取锁，即使看起来锁没有被持有，它也不会去抢，也不会 spin，它会乖乖地加入到等待队列的尾部
	模式切换
		如果拥有 Mutex 的 waiter 发现下面两种情况的其中之一，它就会把这个 Mutex 转换成正常模式:
		此 waiter 已经是队列中的最后一个 waiter 了，没有其它的等待锁的 goroutine 了
		此 waiter 的等待时间小于 1 毫秒
	vs
		正常模式拥有更好的性能，因为即使有等待抢锁的 waiter，goroutine 也可以连续多次获取到锁
		饥饿模式是对公平性和性能的一种平衡，它避免了某些 goroutine 长时间的等待锁
		在饥饿模式下，优先对待的是那些一直在等待的 waiter
	代码分析：逐步分析下 Mutex 代码的关键行，彻底搞清楚饥饿模式的细节
		02.mutex.go

总结
	初版的 Mutex 设计非常简洁，充分展示了 Go 创始者的简单、简洁的设计哲学
		但是，随着大家的使用，逐渐暴露出一些缺陷，为了弥补这些缺陷，Mutex 不得不越来越复杂
	Go 创始者的哲学
		强调 GO 语言和标准库的稳定性，新版本要向下兼容，用新的版本总能编译老的代码
		Go 语言从出生到现在已经 10 多年了，这个 Mutex 对外的接口却没有变化，依然向下兼容
		即使现在 Go 出了两个版本，每个版本也会向下兼容，保持 Go 语言的稳定性，你也能领悟他们软件开发和设计的思想
	公平性
		为了一个程序 20% 的特性，你可能需要添加 80% 的代码，这也是程序越来越复杂的原因
		所以，最开始的时候，如果能够有一个清晰而且易于扩展的设计，未来增加新特性时，也会更加方便

思考
	1.目前 Mutex 的 state 字段有几个意义，这几个意义分别是由哪些字段表示的？
	2.等待一个 Mutex 的 goroutine 数最大是多少？是否能满足现实的需求？
		2^28 - 1，int32 的最高位表示负数

补充
给新人机会
	流程
		1.取号排队
			获取state，从 mutexWaiterShift 取号排队
		2.排队成功：获取锁 / 休眠
			Lock 写入时，新 go 可能写入 mutexWoken（上一个持有者“告知”的），而 awoke 的 go 不会写入 mutexWoken
		3.醒来后重新去取号
	重点：mutexWoken
		old>>mutexWaiterShift == 0 说明休眠的 go 为 0，其他的 go 都是没有休眠的
		每个waiter在休眠前，要么写入 mutexWoken（新来的 go 有可能写入），要么不写入 mutexWoken（awoke 的 go 肯定不写入）
		当新来的 go 和唤醒的 go，“一换一”（waiter和休眠位置互换）时，mutexWoken 为 2？
		...
*/

// ==========初版的互斥锁==========
//// CAS操作，当时还没有抽象出atomic包
//func cas(val *int32, old, new int32) bool
//func semacquire(*int32)
//func semrelease(*int32)
//
//// 互斥锁的结构，包含两个字段
//type Mutex struct {
//	key  int32 // 锁是否被持有的标识
//	sema int32 // 信号量专用，用以阻塞/唤醒goroutine
//}
//
//// 保证成功在val上增加delta的值
//func xadd(val *int32, delta int32) (new int32) {
//	for {
//		v := *val
//		if cas(val, v, v+delta) {
//			return v + delta
//		}
//	}
//	panic("unreached")
//}
//
//// 请求锁
//func (m *Mutex) Lock() {
//	if xadd(&m.key, 1) == 1 { //标识加1，如果等于1，成功获取到锁
//		return
//	}
//	semacquire(&m.sema) // 否则阻塞等待
//}
//func (m *Mutex) Unlock() {
//	if xadd(&m.key, -1) == 0 { // 将标识减去1，如果等于0，则没有其它等待者
//		return
//	}
//	semrelease(&m.sema) // 唤醒其它阻塞的goroutine
//}
