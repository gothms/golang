package concurrent

/*
在分布式环境中，Leader选举、互斥锁和读写锁该如何实现？

在进程内使用
	1-18 中的并发原语都是在进程内使用的
	也就是我们常见的一个运行程序为了控制共享资源、实现任务编排和进行消息传递而提供的控制类型
分布式的并发原语
	它们控制的资源或编排的任务分布在不同进程、不同机器上
	分布式的并发原语实现更加复杂，因为在分布式环境中，网络状况、服务状态都是不可控的
	不过还好有相应的软件系统去做这些事情
	这些软件系统会专门去处理这些节点之间的协调和异常情况，并且保证数据的一致性
	我们要做的就是在它们的基础上实现我们的业务
常用来做协调工作的软件系统
	Zookeeper、etcd、Consul 之类的软件
	Zookeeper 为 Java 生态群提供了丰富的分布式并发原语（通过 Curator 库），但是缺少 Go 相关的并发原语库
	Consul 在提供分布式并发原语这件事儿上不是很积极，而 etcd 就提供了非常好的分布式并发原语
	比如分布式互斥锁、分布式读写锁、Leader 选举，等
前提
	既然我们依赖 etcd，那么，在生产环境中要有一个 etcd 集群
	而且应该保证这个 etcd 集群是 7*24 工作的

Leader 选举
	主从架构
		主从架构中的服务节点分为主（Leader、Master）和从（Follower、Slave）两种角色，实际节点包括 1 主 n 从，一共是 n+1 个节点
		主节点常常执行写操作，从节点常常执行读操作
		如果读写都在主节点，从节点只是提供一个备份功能的话，那么，主从架构就会退化成主备模式架构
	Leader 选举常常用在主从架构的系统中
		主从架构中最重要的是如何确定节点的角色，也就是，到底哪个节点是主，哪个节点是从
		在同一时刻，系统中不能有两个主节点，否则，如果两个节点都是主，都执行写操作的话，就有可能出现数据不一致的情况
		所以，我们需要一个选主机制，选择一个节点作为主节点，这个过程就是 Leader 选举
	Leader 选举场景
		当主节点宕机或者是不可用时，就需要新一轮的选举，从其它的从节点中选择出一个节点，让它作为新主节点
		宕机的原主节点恢复后，可以变为从节点，或者被摘掉
	etcd 选举 Leader
		可以通过 etcd 基础服务来实现 leader 选举
		具体点说，我们可以将 Leader 选举的逻辑交给 etcd 基础服务，这样，我们只需要把重心放在业务开发上
		etcd 基础服务可以通过多节点的方式保证 7*24 服务，所以，我们也不用担心 Leader 选举不可用的问题
		图示 19.distributed_etcd.jpg
	准备工作：go get github.com/coreos/etcd
		先部署一个 etcd 的集群，或者部署一个 etcd 节点做测试
	实现一个测试分布式程序的框架
		它会先从命令行中读取命令，然后再执行相应的命令。你可以打开两个窗口，模拟不同的节点，分别执行不同的命令
		==========Etcd 测试==========

选举
	三个和选主相关的方法
		如果你的业务集群还没有主节点，或者主节点宕机了，你就需要发起新一轮的选主操作，主要会用到 Campaign 和 Proclaim
		如果你需要主节点放弃主的角色，让其它从节点有机会成为主节点，就可以调用 Resign 方法
	第一个方法是 Campaign：它的作用是，把一个节点选举为主节点，并且会设置一个值
		func (e *Election) Campaign(ctx context.Context, val string) error
		这是一个阻塞方法，在调用它的时候会被阻塞，直到满足下面的三个条件之一，才会取消阻塞
			1. 成功当选为主
			2. 此方法返回错误
			3. ctx 被取消
	第二个方法是 Proclaim：它的作用是，重新设置 Leader 的值，但是不会重新选主，这个方法会返回新值设置成功或者失败的信息
		func (e *Election) Proclaim(ctx context.Context, val string) error
	第三个方法是 Resign：开始新一次选举。这个方法会返回新的选举成功或者失败的信息
		func (e *Election) Resign(ctx context.Context) (err error)
	示例：测试这三个方法
		启动两个节点，执行和这三个方法相关的命令
		==========Etcd 测试 Campaign、Proclaim、Resign==========
查询
	除了选举 Leader，程序在启动的过程中，或者在运行的时候，还有可能需要查询当前的主节点是哪一个节点？主节点的值是什么？版本是多少？
		不光是主从节点需要查询和知道哪一个节点，在分布式系统中，还有其它一些节点也需要知道集群中的哪一个节点是主节点，哪一个节点是从节点
		这样它们才能把读写请求分别发往相应的主从节点上
	Leader 方法：查询当前 Leader 的方法 Leader，如果当前还没有 Leader，就返回一个错误，你可以使用这个方法来查询主节点信息
		func (e *Election) Leader(ctx context.Context) (*v3.GetResponse, error)
	每次主节点的变动都会生成一个新的版本号，你还可以查询版本号信息（Rev 方法），了解主节点变动情况
		func (e *Election) Rev() int64
	可以在测试完选主命令后，测试查询命令（query、rev）
监控
	如果主节点变化了，我们需要得到最新的主节点信息
	Observe 方法：监控主的变化
		func (e *Election) Observe(ctx context.Context) <-chan v3.GetResponse
		它会返回一个 chan，显示主节点的变动信息。需要注意的是，它不会返回主节点的全部历史变动信息，而是只返回最近的一条变动信息以及之后的变动信息
	测试
		==========Etcd 测试监控命令（Observe）==========
	小结
		etcd 提供了选主的逻辑，而你要做的就是利用这些方法，让它们为你的业务服务
		在使用的过程中，你还需要做一些额外的设置，比如查询当前的主节点、启动一个 goroutine 阻塞调用 Campaign 方法，等
		虽然你需要做一些额外的工作，但是跟自己实现一个分布式的选主逻辑相比，大大地减少了工作量

互斥锁
	Mutex、RWMutex 等互斥锁，都是用来保护同一进程内的共享资源的
		而分布在不同机器中的不同进程内的 goroutine，如何利用分布式互斥锁来保护共享资源呢
	互斥锁的应用场景和主从架构的应用场景不太一样
		使用互斥锁的不同节点是没有主从这样的角色的，所有的节点都是一样的，只不过在同一时刻，只允许其中的一个节点持有锁
	互斥锁相关的两个原语
		即 Locker 和 Mutex
Locker
	etcd 提供了一个简单的 Locker 原语，它类似于 Go 标准库中的 sync.Locker 接口，也提供了 Lock/UnLock 的机制
		func NewLocker(s *Session, pfx string) sync.Locker
		返回值是一个 sync.Locker，它只有 Lock/Unlock 两个方法
	示例
		==========Etcd 测试 Locker==========
		可以同时在两个终端中运行这个测试程序
		它们获得锁是有先后顺序的，一个节点释放了锁之后，另外一个节点才能获取到这个分布式锁
Mutex
	Locker 是基于 Mutex 实现的，只不过，Mutex 提供了查询 Mutex 的 key 的信息的功能
	示例
		==========Etcd 测试 Mutex==========
		Mutex 并没有实现 sync.Locker 接口，它的 Lock/Unlock 方法需要提供一个 context.Context 实例做参数
		这也就意味着，在请求锁的时候，你可以设置超时时间，或者主动取消请求
读写锁
	RWMutex
		互斥锁 Mutex 是在 github.com/coreos/etcd/clientv3/concurrency 包中提供的
		读写锁 RWMutex 却是在 github.com/coreos/etcd/contrib/recipes 包中提供的
	etcd 提供的分布式读写锁的功能和标准库的读写锁的功能是一样的
		只不过，etcd 提供的读写锁，可以在分布式环境中的不同的节点使用
		它提供的方法也和标准库中的读写锁的方法一致，分别提供了 RLock/RUnlock、Lock/Unlock 方法
	示例
		==========Etcd 测试 RWMutex==========

总结
	自己实现分布式环境的并发原语，是相当困难的一件事
		因为你需要考虑网络的延迟和异常、节点的可用性、数据的一致性等多种情况
	考虑异常的情况
		比如网络断掉等
		同时，分布式并发原语需要网络之间的通讯，所以会比使用标准库中的并发原语耗时更长

思考
	1. 如果持有互斥锁或者读写锁的节点意外宕机了，它持有的锁会不会被释放？
	2. etcd 提供的读写锁中的读和写有没有优先级？

补充
	源码地址：https://github.com/smallnest/distributed
*/
