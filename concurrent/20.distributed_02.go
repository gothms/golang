package concurrent

/*
在分布式环境中，队列、栅栏和STM该如何实现？

分布式队列和优先级队列
	站在 etcd 的肩膀上，利用 etcd 提供的功能实现分布式队列
		etcd 集群的可用性由 etcd 集群的维护者来保证，我们不用担心网络分区、节点宕机等问题
		我们可以把这些通通交给 etcd 的运维人员，把我们自己的关注点放在使用上
		etcd 通过 github.com/coreos/etcd/contrib/recipes 包提供了分布式队列这种数据结构
	创建分布式队列
		只需要传入 etcd 的 client 和这个队列的名字
		func NewQueue(client *v3.Client, keyPrefix string) *Queue
	API：这个队列只有两个方法，分别是出队和入队，队列中的元素是字符串类型
		func (q *Queue) Enqueue(val string) error	// 入队
		func (q *Queue) Dequeue() (string, error)	// 出队
		如果这个分布式队列当前为空，调用 Dequeue 方法的话，会被阻塞，直到有元素可以出队才返回
		etcd 的分布式队列是一种多读多写的队列，所以，你也可以启动多个写节点和多个读节点
	示例
		命令操作
			首先，我们启动一个程序，它会从命令行读取你的命令，然后执行
			你可以输入push <value>，将一个元素入队，输入pop，将一个元素弹出
			另外，你还可以使用这个程序启动多个实例，用来模拟分布式的环境
		测试代码
			==========Etcd 测试 Queue==========
			打开两个终端，分别执行这个程序
			在第一个终端中执行入队操作，在第二个终端中执行出队操作，并且观察一下出队、入队是否正常
	优先级队列（PriorityQueue）
		它的用法和队列类似，也提供了出队和入队的操作
		只不过，在入队的时候，除了需要把一个值加入到队列，我们还需要提供 uint16 类型的一个整数，作为此值的优先级，优先级高的元素会优先出队
		源码：contrib/recipes/priority_queue.go
	示例
		可以在一个节点输入一些不同优先级的元素，在另外一个节点读取出来
		==========Etcd 测试 PriorityQueue==========

分布式栅栏
	循环栅栏 CyclicBarrier 和标准库中的 WaitGroup，本质上是同一类并发原语
		都是等待同一组 goroutine 同时执行，或者是等待同一组 goroutine 都完成
	分布式场景
		一组节点协同工作，共同等待一个信号，在信号未出现前，这些节点会被阻塞住
		而一旦信号出现，这些阻塞的节点就会同时开始继续执行下一步的任务
	etcd 提供了相应的分布式并发原语
		Barrier：分布式栅栏
			如果持有 Barrier 的节点释放了它，所有等待这个 Barrier 的节点就不会被阻塞，而是会继续执行
		DoubleBarrier：计数型栅栏
			在初始化计数型栅栏的时候，我们就必须提供参与节点的数量，当这些数量的节点都 Enter 或者 Leave 的时候，这个栅栏就会放开
			所以，我们把它称为计数型栅栏
Barrier：分布式栅栏
	contrib/recipes/barrier.go
	创建
		func NewBarrier(client *v3.Client, key string) *Barrier
	API
		func (b *Barrier) Hold() error
		func (b *Barrier) Release() error
		func (b *Barrier) Wait() error
		Hold 方法是创建一个 Barrier
			如果 Barrier 已经创建好了，有节点调用它的 Wait 方法，就会被阻塞
		Release 方法是释放这个 Barrier，也就是打开栅栏
			如果使用了这个方法，所有被阻塞的节点都会被放行，继续执行
		Wait 方法会阻塞当前的调用者，直到这个 Barrier 被 release
			如果这个栅栏不存在，调用者不会被阻塞，而是会继续执行
	示例
		模拟使用栅栏
			以在一个终端中运行这个程序，执行"hold""release"命令，模拟栅栏的持有和释放
			在另外一个终端中运行这个程序，不断调用"wait"方法，看是否能正常地跳出阻塞继续执行
		测试代码
			==========Etcd 测试 Barrier==========
DoubleBarrier：计数型栅栏
	contrib/recipes/double_barrier.go
	初始化的时候需要提供一个计数 count
		func NewDoubleBarrier(client *clientv3.Client, key string, count int) *DoubleBarrier
	API
		func (b *DoubleBarrier) Enter() error
		func (b *DoubleBarrier) Leave() error
	Enter
		当调用者调用 Enter 时，会被阻塞住，直到一共有 count（初始化这个栅栏的时候设定的值）个节点调用了 Enter，这 count 个被阻塞的节点才能继续执行
		所以，你可以利用它编排一组节点，让这些节点在同一个时刻开始执行任务
	Leave
		如果你想让一组节点在同一个时刻完成任务，就可以调用 Leave 方法
		节点调用 Leave 方法的时候，会被阻塞，直到有 count 个节点，都调用了 Leave 方法，这些节点才能继续执行
	示例
		模拟使用
			起两个节点，同时执行 Enter 方法，看看这两个节点是不是先阻塞，之后才继续执行
			然后，你再执行 Leave 方法，也观察一下，是不是先阻塞又继续执行的
		测试代码
	小结
		分布式栅栏和计数型栅栏控制的是不同节点、不同进程的执行
		当你需要协调一组分布式节点在某个时间点同时运行的时候，可以考虑 etcd 提供的这组并发原语

STM（Software Transactional Memory，软件事务内存）
	事务
		在开发基于数据库的应用程序的时候，我们经常用到事务
		事务就是要保证一组操作要么全部成功，要么全部失败
	etcd 的事务
		etcd 提供了在一个事务中对多个 key 的更新功能，这一组 key 的操作要么全部成功，要么全部失败
		etcd 的事务实现方式是基于 CAS 方式实现的，融合了 Get、Put 和 Delete 操作
	etcd 的事务操作
		分为条件块、成功块和失败块，条件块用来检测事务是否成功，如果成功，就执行 Then(...)，如果失败，就执行 Else(...)
		Txn().If(cond1, cond2, ...).Then(op1, op2, ...,).Else(op1’, op2’, …)
	示例：利用 etcd 的事务实现转账，从账户 from 向账户 to 转账 amount
		==========Etcd 测试 Txn==========
		虽然可以利用 etcd 实现事务操作，但是逻辑还是比较复杂的
		因为事务使用起来非常麻烦，所以 etcd 又在这些基础 API 上进行了封装，新增了一种叫做 STM 的操作，提供了更加便利的方法
	要使用 STM，你需要先编写一个 apply 函数，这个函数的执行是在一个事务之中的
		apply func(STM) error
		这个方法包含一个 STM 类型的参数，它提供了对 key 值的读写操作
	STM API：clientv3/concurrency/stm.go
		type STM interface {
			Get(key ...string) string
			Put(key, val string, opts ...v3.OpOption)
			Rev(key string) int64
			Del(key string)
		}
	使用
		使用 etcd STM 的时候，我们只需要定义一个 apply 方法，比如说转账方法 exchange
		然后通过 concurrency.NewSTM(cli, exchange)，就可以完成转账事务的执行了
	示例
		问题描述
			创建了 5 个银行账号，然后随机选择一些账号两两转账
			在转账的时候，要把源账号一半的钱要转给目标账号
			例子启动了 10 个 goroutine 去执行这些事务，每个 goroutine 要完成 100 个事务
			为了确认事务是否出错了，我们最后要校验每个账号的钱数和总钱数。总钱数不变，就代表执行成功了
		测试代码
			==========Etcd 测试 STM==========
	小结
		利用 etcd 做存储时，是可以利用 STM 实现事务操作的
		一个事务可以包含多个账号的数据更改操作，事务能够保证这些更改要么全成功，要么全失败

总结
	Etcd、Redis、MySQL
		也可以使用 Redis 实现分布式锁，或者是基于 MySQL 实现分布式锁，这也是常用的选择
		对于大厂来说，选择起来是非常简单的，只需要看看厂内提供了哪个基础服务，哪个更稳定些
		对于没有 etcd、Redis 这些基础服务的公司来说，很重要的一点，就是自己搭建一套这样的基础服务，并且运维好
		这就需要考察你们对 etcd、Redis、MySQL 的技术把控能力了，哪个用得更顺手，就用哪个
	建议
		一般来说，我不建议你自己去实现分布式原语，最好是直接使用 etcd、Redis 这些成熟的软件提供的功能
		这也意味着，我们将程序的风险转嫁到了这些基础服务上，这些基础服务必须要能够提供足够的服务保障

思考
	1. 部署一个 3 节点的 etcd 集群，测试一下分布式队列的性能
	2. etcd 提供的 STM 是分布式事务吗？
*/
