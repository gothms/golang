package concurrent

/*
项目地址
	https://github.com/smallnest/dive-to-gosync-workshop/tree/master
1-11
	基本并发原语
	互斥锁Mutex、读写锁RWMutex
	并发编排WaitGroup
	条件变量Cond
12
	原子操作
13-15
	Channel
16-18
	扩展并发原语
19-20
	分布式并发原语

Go 并发5大问题
	1. 在面对并发难题时，感觉无从下手，不知道该用什么并发原语来解决问题
	2. 如果多个并发原语都可以解决问题，那么，究竟哪个是最优解呢？比如说是用互斥锁，还是用 Channel
	3. 不知道如何编排并发任务。并发编程不像是传统的串行编程，程序的运行存在着很大的不确定性
		这个时候，就会面临一个问题，怎么才能让相应的任务按照你设想的流程运行呢？
	4. 有时候，按照正常理解的并发方式去实现的程序，结果莫名其妙就 panic 或者死锁了，排查起来非常困难
	5. 已知的并发原语都不能解决并发问题，程序写起来异常复杂，而且代码混乱，容易出错
Go 并发编程知识主线
	基本并发原语：Mutex、RWMutex、Waitgroup、Cond、Pool、Context 等标准库并发原语
		这些都是传统的并发原语，在其它语言中也很常见，是我们在并发编程中常用的类型
	原子操作：Go 标准库中提供的原子操作
		原子操作是其它并发原语的基础，学会了你就可以自己创造新的并发原语
	Channel：Channel 类型是 Go 语言独特的类型，因为比较新，所以难以掌握
		全方位地学习 Channel 类型，你不仅能掌握它的基本用法，而且还能掌握它的处理场景和应用模式，避免踩坑
	扩展并发原语：目前来看，Go 开发组不准备在标准库中扩充并发原语了
		但是还有一些并发原语应用广泛，比如信号量、SingleFlight、循环栅栏、ErrGroup 等
		掌握了它们，就可以在处理一些并发问题时，取得事半功倍的效果
	分布式并发原语：分布式并发原语是应对大规模的应用程序中并发问题的并发类型
		使用 etcd 实现的一些分布式并发原语
		比如 Leader 选举、分布式互斥锁、分布式读写锁、分布式队列等，在处理分布式场景的并发问题时，特别有用
学习主线
	基础用法
	实现原理
	易错场景
	知名项目中的 Bug
Go 并发大方向
	原则
		任务编排用 Channel，共享资源保护用传统并发原语
	打破原则
		但是，如果你想要在 Go 并发编程的道路上向前走，就不能局限于这个原则
		实际上，针对同一种场景，也许存在很多并发原语都适用的情况，但是一定是有最合适的那一个
		所以，你必须非常清楚每种并发原语的实现机制和适用场景，千万不要被网上的一些文章误导，万事皆用 Channel
Go 并发原语源代码设计
	Mutex 为了公平性考量的设计
	sync.Map 为提升性能做的设计
	以及很多并发原语的异常状况的处理方式。这些异常状况，常常是并发编程中程序 panic 的原因
创造出自己需要的并发原语
	对既有的并发原语进行组合，使用两个、三个或者更多的并发原语去解决问题
		比如说，我们可以通过信号量和 WaitGroup 组合成一个新的并发原语，这个并发原语可以使用有限个 goroutine 并发处理子任务
	“无中生有”，根据已经掌握的并发原语的设计经验，创造出合适的新的并发原语，以应对一些特殊的并发问题
		比如说，标准库中并没有信号量，你可以自己创造出这个类型
	对 Go 并发原语的掌握已经出神入化了
3个目标
	建立起一个丰富的并发原语库
	熟知每一种并发原语的实现机制和适用场景
	能够创造出自己需要的并发原语
*/
