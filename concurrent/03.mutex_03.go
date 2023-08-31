package concurrent

import (
	"fmt"
	"github.com/petermattis/goid"
	"sync"
	"sync/atomic"
)

/*
Mutex：4种易错场景大盘点

Mutex 使用错误场景
	在一些复杂的场景中，比如跨函数调用 Mutex 或者是在重构或者修补 Bug 时误操作
	但是使用 Mutex 时，确实会出现一些 Bug，比如说忘记释放锁、重入锁、复制已使用了的 Mutex 等情况

常见的 4 种错误场景
	Lock/Unlock 不是成对出现
	Copy 已使用的 Mutex
	重入
	死锁
Lock/Unlock 不是成对出现
	Lock/Unlock 没有成对出现，就意味着会出现死锁的情况，或者是因为 Unlock 一个未加锁的 Mutex 而导致 panic
	缺少 Unlock 的常见三种情况
		1. 代码中有太多的 if-else 分支，可能在某个分支中漏写了 Unlock
		2. 在重构的时候把 Unlock 给删除了
		3. Unlock 误写成了 Lock
		锁被获取之后，就不会被释放了，这也就意味着，其它的 goroutine 永远都没机会获取到锁
	缺少 Lock 的场景：一般来说就是误操作删除了 Lock
		比如先前使用 Mutex 都是正常的，结果后来其他人重构代码的时候，由于对代码不熟悉，或者由于开发者的马虎，把 Lock 调用给删除了，或者注释掉了
		示例：
			func foo() {
				var mu sync.Mutex
				defer mu.Unlock()
				fmt.Println("hello world!")
			}
		fatal("sync: unlock of unlocked mutex")
Copy 已使用的 Mutex
	package sync 的同步原语在使用后不能复制
		Mutex 是一个有状态的对象，它的 state 字段记录这个锁的状态
		如果你要复制一个已经加锁的 Mutex 给一个新的变量，那么新的刚初始化的变量居然被加锁了
		这显然不符合你的期望，因为你期望的是一个零值的 Mutex
		关键是在并发环境下，你根本不知道要复制的 Mutex 状态是什么，因为要复制的 Mutex 是由其它 goroutine 并发访问的，状态可能总是在变化
	示例：TestCheckdead
		foo(mutex) // 复制锁
		行在调用 foo 函数的时候，调用者会复制 Mutex 变量 mutex 作为 foo 函数的参数
		不幸的是，复制之前已经使用了这个锁，这就导致，复制的 Mutex 是一个带状态 Mutex
	死锁检查机制：checkdead()方法
		Go 在运行时，有死锁的检查机制（checkdead()方法），它能够发现死锁的 goroutine
		程序运行的时候，死锁检查机制能够发现这种死锁情况并输出错误信息
		错误信息以及错误堆栈
			fatal error: all goroutines are asleep - deadlock!
			E:/gothmslee/golang/concurrent/test/03_test.go:19 +0x68
	vet 工具
		把检查写在 Makefile 文件中，在持续集成的时候跑一跑，这样可以及时发现问题，及时修复
		vet 会对实现 Locker 接口的数据类型做静态检查，一旦代码中有复制使用这种数据类型的情况，就会发出警告
		命令
			go vet 03_test.go
			go vet golang/concurrent/test
		检查结果
			concurrent\test\03_test.go:14:6: call of foo copies lock value: sync.Mutex
			concurrent\test\03_test.go:18:13: foo passes lock by value: sync.Mutex
		在调用 foo 函数的时候发生了 lock value 复制的情况，还提示出问题的代码行数以及 copy lock 导致的错误
	vet 原理
		检查是通过 copylock 分析器静态分析实现的
		这个分析器会分析函数调用、range 遍历、复制、声明、函数返回值等位置，有没有锁的值 copy 的情景，以此来判断有没有问题
		可以说，只要是实现了 Locker 接口，就会被分析
		其实，有些没有实现 Locker 接口的同步原语（比如 WaitGroup），也能被分析
	copylock 源码
		cmd/vendor/golang.org/x/tools/go/analysis/passes/copylock/copylock.go
		var lockerType *types.Interface 的 init 函数，确定什么类型会被分析
重入
	“可重入锁“
		Java 并发包中常用的一个同步原语 ReentrantLock 就是可重入锁
		当一个线程获取锁时，如果没有其它线程拥有这个锁，那么，这个线程就成功获取到这个锁
		之后，如果其它线程再请求这个锁，就会处于阻塞等待的状态
		但是，如果拥有这把锁的线程再请求这把锁的话，不会阻塞，而是成功返回，所以叫可重入锁（有时也叫递归锁）
		只要你拥有这把锁，你可以可着劲儿地调用，比如通过递归实现一些算法，调用者不会阻塞或者死锁
	Mutex 不是可重入的锁
		因为 Mutex 的实现中没有记录哪个 goroutine 拥有这把锁
		理论上，任何 goroutine 都可以随意地 Unlock 这把锁，所以没办法计算重入条件
		所以，一旦误用 Mutex 的重入，就会导致报错
		示例
			TestReentrantLock
		报错
			fatal error: all goroutines are asleep - deadlock!
			程序一直在请求锁，但是一直没有办法获取到锁，结果就是 Go 运行时发现死锁了，没有其它地方能够释放锁让程序运行下去
	实现可重入锁
		关键就是，实现的锁要能记住当前是哪个goroutine 持有这个锁
		方案一：通过 hacker 的方式获取到 goroutine id，记录下获取锁的 goroutine id，它可以实现 Locker 接口
		方案二：调用 Lock/Unlock 方法时，由 goroutine 提供一个 token，用来标识它自己，而不是我们通过 hacker 的方式获取到 goroutine id
			但是，这样一来，就不满足 Locker 接口了
		可重入锁（递归锁）解决了代码重入或者递归调用带来的死锁问题
			同时它也带来了另一个好处，就是我们可以要求，只有持有锁的 goroutine 才能 unlock 这个锁
			这也很容易实现，因为在上面这两个方案中，都已经记录了是哪一个 goroutine 持有这个锁
	方案一：goroutine id
		获取 goroutine id：方式有两种，分别是简单方式和 hacker 方式
		简单方式：runtime.Stack 方法获取栈帧信息，栈帧信息里包含 goroutine id
			func Stack(buf []byte, all bool) int
			第二个参数为 true 会输出所有的 goroutine 信息
			示例：TestGetGoId
		hacker 方式：获取运行时的 g 指针，反解出对应的 g 的结构
			每个运行的 goroutine 结构的 g 指针保存在当前 goroutine 的一个叫做 TLS 对象中
				第一步：我们先获取到 TLS 对象
				第二步：再从 TLS 中获取 goroutine 结构的 g 指针
				第三步：再从 g 指针中取出 goroutine id
			不同版本的 goroutine
				注意，不同 Go 版本的 goroutine 的结构可能不同，所以需要根据 Go 的不同版本进行调整
				当然了，如果想要搞清楚各个版本的 goroutine 结构差异，所涉及的内容又过于底层而且复杂，学习成本太高
			三方库：petermattis/goid
				获取 goroutine id，可以支持多个 Go 版本的 goroutine
				https://github.com/petermattis/goid
			实现和测试：type RecursiveMutex struct & TestRecursiveMutex
				这段代码可以拿来即用，实现非常巧妙。它相当于给 Mutex 打一个补丁，解决了记录锁的持有者的问题
				可以看到，我们用 owner 字段，记录当前锁的拥有者 goroutine 的 id
				recursion 是辅助字段，用于记录重入的次数
				尽管拥有者可以多次调用 Lock，但是也必须调用相同次数的 Unlock，这样才能把锁释放掉
				这是一个合理的设计，可以保证 Lock 和 Unlock 一一对应
		小结
			此方案用 goroutine id 做 goroutine 的标识，我们也可以让 goroutine 自己来提供标识
			不管怎么说，Go 开发者不期望你利用 goroutine id 做一些不确定的东西，所以，他们没有暴露获取 goroutine id 的方法
	方案二：token
		原理
			调用者自己提供一个 token，获取锁的时候把这个 token 传入，释放锁的时候也需要把这个 token 传入
			通过用户传入的 token 替换方案一中 goroutine id，其它逻辑和方案一一致
		实现和测试：TokenRecursiveMutex & TestTokenRecursiveMutex
死锁
	概念
		两个或两个以上的进程（或线程，goroutine）在执行过程中，因争夺共享资源而处于一种互相等待的状态
		如果没有外部干涉，它们都将无法推进下去，此时，我们称系统处于死锁状态或系统产生了死锁
	死锁产生的必要条件。破坏这四个条件中的一个或者几个，可避免死锁
		1. 互斥：至少一个资源是被排他性独享的，其他线程必须处于等待状态，直到资源被释放
		2. 持有和等待：goroutine 持有一个资源，并且还在请求其它 goroutine 持有的资源，也就是常说的“吃着碗里，看着锅里”的意思
		3. 不可剥夺：资源只能由持有它的 goroutine 来释放
		4. 环路等待：一般来说，存在一组等待进程，P={P1，P2，…，PN}，P1 等待 P2 持有的资源，P2 等待 P3 持有的资源
			依此类推，最后是 PN 等待 P1 持有的资源，这就形成了一个环路等待的死结
			哲学家就餐问题：https://zh.wikipedia.org/wiki/%E5%93%B2%E5%AD%A6%E5%AE%B6%E5%B0%B1%E9%A4%90%E9%97%AE%E9%A2%98
	示例
		去派出所开证明，派出所要求物业先证明我是本物业的业主，但是，物业要我提供派出所的证明，才能给我开物业证明，结果就陷入了死锁状态
		可以把派出所和物业看成两个 goroutine，派出所证明和物业证明是两个资源，双方都持有自己的资源而要求对方的资源，而且自己的资源自己持有，不可剥夺
		代码：TestDeadlock
		解决：引入一个第三方的锁，大家都依赖这个锁进行业务处理
			比如现在政府推行的一站式政务服务中心
			或者是解决持有等待问题，物业不需要看到派出所的证明才给开物业证明，等

流行的 Go 开发项目踩坑记
	使用原则
		保证 Lock/Unlock 成对出现，尽可能采用 defer mutex.Unlock 的方式，把它们成对、紧凑地写在一起
	Docker
		简介
			Docker 容器是一个开源的应用容器引擎，开发者可以以统一的方式，把他们的应用和依赖包打包到一个可移植的容器中，然后发布到任何安装了 docker 引擎的服务器上
			Docker 是使用 Go 开发的，也算是 Go 的一个杀手级产品了，它的 Mutex 相关的 Bug 也不少
		issue 36114：https://github.com/moby/moby/pull/36114/files
			一个死锁问题
			原因
				hotAddVHDsAtStart 方法执行的时候，执行了加锁 svm 操作
				但是，在其中调用 hotRemoveVHDsAtStart 方法时，这个 hotRemoveVHDsAtStart 方法也是要加锁 svm 的
				很不幸，Go 标准库中的 Mutex 是不可重入的，所以，代码执行到这里，就出现了死锁的现象
			图示
				03.mutex_03_demo_docker_01.jpg
			解决
				再提供一个不需要锁的 hotRemoveVHDsNoLock 方法，避免 Mutex 的重入
		issue 34881：https://github.com/moby/moby/pull/34881/files
			issue 34881本来是修复 Docker 的一个简单问题，如果节点在初始化的时候，发现自己不是一个 swarm mananger，就快速返回
			原因
				节点发现不满足条件就返回了，但是，c.mu 这个锁没有释放
				这是在重构或者添加新功能的时候经常犯的一个错误，因为不太了解上下文，或者是没有仔细看函数的逻辑，从而导致锁没有被释放
			图示
				03.mutex_03_demo_docker_02.jpg
		其他关于 Mutex 的 issue 或者 pull request
			分别是 36840、37583、35517、35482、33305、32826、30696、29554、29191、28912、26507 等
	Kubernetes
		issue 7236
			1issue 72361 增加 Mutex 为了保护资源
			这是为了解决 data race 问题而做的一个修复，修复方法也很简单，使用互斥锁即可，这也是我们解决 data race 时常用的方法
			图示
				03.mutex_03_demo_kubernetes_01.jpg
		issue 45192：https://github.com/kubernetes/kubernetes/pull/45192/files
			也是一个返回时忘记 Unlock 的典型例子，和 docker issue 34881 犯的错误都是一样的
			图示
				03.mutex_03_demo_kubernetes_02.jpg
		其它的 Mutex 相关的 issue
			比如 71617、70605 等
	gRPC
		简介
			gRPC 是 Google 发起的一个开源远程过程调用 （Remote procedure call）系统
			该系统基于 HTTP/2 协议传输，使用 Protocol Buffers 作为接口描述语言
			它提供 Go 语言的实现
			即使是 Google 官方出品的系统，也有一些 Mutex 的 issue
		issue 795：https://github.com/grpc/grpc-go/pull/795
			一个想不到的 bug，那就是将 Unlock 误写成了 Lock
			图示
				03.mutex_03_demo_grpc.jpg
		其他的为了保护共享资源而添加 Mutex 的 issue
			比如 1318、2074、2542 等
	etcd
		简介
			tcd 是一个非常知名的分布式一致性的 key-value 存储技术，被用来做配置共享和服务发现
		issue 10419：https://github.com/etcd-io/etcd/pull/10419/files
			一个锁重入导致的问题
			描述
				Store 方法内对请求了锁，而调用的 Compact 的方法内又请求了锁，这个时候，会导致死锁，一直等待
				解决办法就是提供不需要加锁的 Compact 方法
			图示
				03.mutex_03_demo_etcd.jpg

总结
	手误和重入导致的死锁，是最常见的使用 Mutex 的 Bug
	死锁检测
		Go 死锁探测工具只能探测整个程序是否因为死锁而冻结了，不能检测出一组 goroutine 死锁导致的某一块业务冻结的情况
		你还可以通过 Go 运行时自带的死锁检测工具，或者是第三方的工具（比如go-deadlock、go-tools）进行检查，这样可以尽早发现一些死锁的问题
		不过，有些时候，死锁在某些特定情况下才会被触发，所以，如果你的测试或者短时间的运行没问题，不代表程序一定不会有死锁问题
	go-deadlock
		https://github.com/sasha-s/go-deadlock
	go-tools
		https://github.com/dominikh/go-tools
	Bug 复现和 pprof
		并发程序最难跟踪调试的就是很难重现，因为并发问题不是按照我们指定的顺序执行的
		由于计算机调度的问题和事件触发的时机不同，死锁的 Bug 可能会在极端的情况下出现
		通过搜索日志、查看日志，我们能够知道程序有异常了，比如某个流程一直没有结束
		这个时候，可以通过 Go pprof 工具分析，它提供了一个 block profiler 监控阻塞的 goroutine
		除此之外，我们还可以查看全部的 goroutine 的堆栈信息，通过它，你可以查看阻塞的 groutine 究竟阻塞在哪一行哪一个对象上了

思考
	查找知名的数据库系统 TiDB 的 issue，看看有没有 Mutex 相关的 issue，看看它们都是哪些相关的 Bug

补充
	Docker issue 34881
		https://github.com/moby/moby/pull/34881/files
		图示 03.mutex_03_github_docker.jpg
*/

// TokenRecursiveMutex Token方式的递归锁
type TokenRecursiveMutex struct {
	sync.Mutex
	token     int64
	recursion int32
}

// Lock 请求锁，需要传入token
func (m *TokenRecursiveMutex) Lock(token int64) {
	if atomic.LoadInt64(&m.token) == token { //如果传入的token和持有锁的token一致
		m.recursion++
		return
	}
	m.Mutex.Lock() // 传入的token不一致，说明不是递归调用
	// 抢到锁之后记录这个token
	atomic.StoreInt64(&m.token, token)
	m.recursion = 1
}

// Unlock 释放锁
func (m *TokenRecursiveMutex) Unlock(token int64) {
	if atomic.LoadInt64(&m.token) != token { // 释放其它token持有的锁
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.token, token))
	}
	m.recursion--         // 当前持有这个锁的token释放锁
	if m.recursion != 0 { // 还没有回退到最初的递归调用
		return
	}
	atomic.StoreInt64(&m.token, 0) // 没有递归调用了，释放锁
	m.Mutex.Unlock()
}

// RecursiveMutex 包装一个Mutex,实现可重入
// hacker 方式
type RecursiveMutex struct {
	sync.Mutex
	owner     int64 // 当前持有锁的goroutine id
	recursion int32 // 这个goroutine 重入的次数
}

func (m *RecursiveMutex) Lock() {
	gid := goid.Get()
	// 如果当前持有锁的goroutine就是这次调用的goroutine,说明是重入
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	// 获得锁的goroutine第一次调用，记录下它的goroutine id,调用次数加1
	atomic.StoreInt64(&m.owner, gid)
	m.recursion = 1
}
func (m *RecursiveMutex) Unlock() {
	gid := goid.Get()
	// 非持有锁的goroutine尝试释放锁，错误的使用
	if atomic.LoadInt64(&m.owner) != gid {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.owner, gid))
	}
	// 调用次数减1
	m.recursion--
	if m.recursion != 0 { // 如果这个goroutine还没有完全释放，则直接返回
		return
	}
	// 此goroutine最后一次调用，需要释放锁
	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()
}
