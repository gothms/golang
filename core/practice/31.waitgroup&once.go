package practice

/*
sync.WaitGroup和sync.Once

使用 channel 方式协作多 goroutine
	func coordinateWithChan() {
		sign := make(chan struct{}, 2)
		num := int32(0)
		fmt.Printf("The number: %d [with chan struct{}]\n", num)
		max := int32(10)
		go addNum(&num, 1, max, func() {
			sign <- struct{}{}
		})
		go addNum(&num, 2, max, func() {
			sign <- struct{}{}
		})
		<-sign
		<-sign
	}
sync.WaitGroup类型
	更加适合实现一对多的 goroutine 协作流程
	是开箱即用的，也是并发安全的
	与其他的几个同步工具一样，它一旦被真正使用就不能被复制了
	WaitGroup类型拥有三个指针方法：Add、Done和Wait
“原理”
	想象该类型中有一个计数器，它的默认值是0
	Add
		可以通过调用该类型值的Add方法来增加，或者减少这个计数器的值
		一般情况下，用这个方法来记录需要等待的 goroutine 的数量
	Done
		相对应的，这个类型的Done方法，用于对其所属值中计数器的值进行减一操作
		我们可以在需要等待的 goroutine 中，通过defer语句调用它
	Wait
		功能是，阻塞当前的 goroutine，直到其所属值中的计数器归零
		如果在该方法被调用的时候，那个计数器的值就是0，那么它将不会做任何事情

问题：sync.WaitGroup类型值中计数器的值可以小于0吗？
	不可以
问题解析
	在调用Add方法的时候是可以传入一个负数的
		WaitGroup值中计数器的值不能小于0，是因为这样会引发一个 panic
		不适当地调用这类值的Done方法和Add方法都会如此
	其他 panic
		尽早地增加其计数器的值
			如果我们对它的Add方法的首次调用，与对它的Wait方法的调用是同时发起的
			比如，在同时启用的两个 goroutine 中，分别调用这两个方法，那么就有可能会让这里的Add方法抛出一个 panic
			这种情况不太容易复现，但是尽早地增加其计数器的值，还是非常有必要的
		WaitGroup值是可以被复用的，但需要保证其计数周期的完整性
			只要计数器的值始于0又归为0，就可以被视为一个计数周期
			在一个此类值的生命周期中，它可以经历任意多个计数周期。但是，只有在它走完当前的计数周期之后，才能够开始下一个计数周期
				如果一个此类值的Wait方法在它的某个计数周期中被调用，那么就会立即阻塞当前的 goroutine，直至这个计数周期完成
			如果在一个此类值的Wait方法被执行期间，跨越了两个计数周期，那么就会引发一个 panic
			示例：
				在当前的 goroutine 因调用此类值的Wait方法，而被阻塞的时候，另一个 goroutine 调用了该值的Done方法，并使其计数器的值变为了0
				这会唤醒当前的 goroutine，并使它试图继续执行Wait方法中其余的代码
				但在这时，又有一个 goroutine 调用了它的Add方法，并让其计数器的值又从0变为了某个正整数
				此时，这里的Wait方法就会立即抛出一个 panic
		小结：WaitGroup值的使用禁忌
			不要把增加其计数器值的操作和调用其Wait方法的代码，放在不同的 goroutine 中执行
			要杜绝对同一个WaitGroup值的两种操作的并发执行
		参见：sync -> waitgroup_test.go -> TestWaitGroupMisuse

知识扩展
问题：sync.Once类型值的Do方法是怎么保证只执行参数函数一次的？
	sync.Once类型
		属于结构体类型，同样也是开箱即用和并发安全的
		该类型中包含了一个sync.Mutex类型的字段，所以，复制该类型的值也会导致功能的失效
	Do
		Once类型的Do方法只接受一个参数，这个参数的类型必须是func()
		该方法的功能并不是对每一种参数函数都只执行一次，而是只执行“首次被调用时传入的”那个函数，并且之后不会再执行任何参数函数
		所以，有多个只需要执行一次的函数，那么就应该为它们中的每一个都分配一个sync.Once类型的值
	done
		uint32类型的字段。它的作用是记录其所属值的Do方法被调用的次数
		不过，该字段的值只可能是0或者1。一旦Do方法的首次调用完成，它的值就会从0变为1
	单例模式
		通过 atomic.LoadUint32(*Once.done) 和 mutex.Lock() 实现
		两次判断 *Once.done == 0，被统称为（跨临界区的）“双重检查”
		第一次判断：若条件不满足则立即返回，这通常被称为“快路径”，或者叫做“快速失败路径
		第二次判断：互斥锁保证串行，被称为“慢路径”或者“常规路径”
	Do 方法在功能方面的两个特点
		1. 阻塞
			由于Do方法只会在参数函数执行结束之后把done字段的值变为1
			因此，如果参数函数的执行需要很长时间或者根本就不会结束（比如执行一些守护任务），那么就有可能会导致相关 goroutine 的同时阻塞
			除了那个抢先执行了参数函数的 goroutine 之外，其他的 goroutine 都会被阻塞在锁定该Once值的互斥锁m的那行代码上
		2. done 字段值变为 1，导致功能缺失
			Do方法在参数函数执行结束后，对done字段的赋值用的是原子操作，并且，这一操作是被挂在defer语句中的
			因此，不论参数函数的执行会以怎样的方式结束，done字段的值都会变为1
			即使这个参数函数没有执行成功（比如引发了一个 panic），我们也无法使用同一个Once值重新执行它了
			所以，如果你需要为参数函数的执行设定重试机制，那么就要考虑Once值的适时替换问题

总结
	sync代码包的WaitGroup类型和Once类型都是非常易用的同步工具。它们都是开箱即用和并发安全的
	WaitGroup
		可以很方便地实现一对多的 goroutine 协作流程
		即：一个分发子任务的 goroutine，和多个执行子任务的 goroutine，共同来完成一个较大的任务
		panic：
		千万不要让其中的计数器的值小于0，否则就会引发 panic
		用“先统一Add，再并发Done，最后Wait”这种标准方式，来使用WaitGroup值
		不要在调用Wait方法的同时，并发地通过调用Add方法去增加其计数器的值，因为这也有可能引发 panic
	Once
		只要传入某个Do方法的参数函数没有结束执行，任何之后调用该方法的 goroutine 就都会被阻塞
	Once类型使用互斥锁和原子操作实现了功能，而WaitGroup类型中只用到了原子操作
		所以可以说，它们都是更高层次的同步工具
		sync包中的其他高级同步工具，都是基于基本的通用工具，实现了某一种特定的功能

思考
	在使用WaitGroup值实现一对多的 goroutine 协作流程时，怎样才能让分发子任务的 goroutine 获得各个子任务的具体执行结果？
A
	可以考虑使用锁 + 容器（数组、切片或字典等），也可以考虑使用通道
	或许也可以用上golang.org/x/sync/errgroup代码包中的程序实体
*/
