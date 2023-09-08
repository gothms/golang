package concurrent

/*
源码阅读工具
	https://mp.weixin.qq.com/s/E2TL_kcbVcRJ0CnxwbXWLw

拷贝
	sync 包的同步原语在使用后不能复制
	“不能”作为参数传递
开箱即用
	Mutex
	RWMutex
	WaitGroup：可重用

必须初始化
	Channel

工具
	-race
	vet
	pprof

三方库
	petermattis/goid
		获取 goroutine id，可以支持多个 Go 版本的 goroutine
		https://github.com/petermattis/goid
	go-deadlock
		死锁检测工具
		https://github.com/sasha-s/go-deadlock
	go-tools
		死锁检测工具
		https://github.com/dominikh/go-tools
	Hugo
		著名的静态网站生成工具
	vitess
		YouTube 开源的知名项目 vitess 中提供了 bucketpool 的实现，它提供了更加通用的多层 buffer 池
	bytebufferpool
		fasthttp 作者 valyala 提供的一个 buffer 池，基本功能和 sync.Pool 相同
		底层也是使用 sync.Pool 实现的，包括会检测最大的 buffer，超过最大尺寸的 buffer，就会被丢弃。提供了校准机制
	oxtoacart/bpool
		保持池子中元素的数量，一旦 Put 的数量多于它的阈值，就会自动丢弃，而 sync.Pool 是一个没有限制的池子，只要 Put 就会收进去
		基于 Channel 实现
	fatih/pool
		最常用的一个 TCP 连接池
	gomemcache
		Brad Fitzpatrick 是知名缓存库 Memcached 的原作者，前 Go 团队成员
		gomemcache 是他使用 Go 开发的 Memchaced 的客户端，其中也用了连接池的方式池化 Memcached 的连接
	fasthttp
		fasthttp 中的 Worker Pool，TCP 连接池实现
		ps：Worker Pool 推荐
			gammazero/workerpool
			ivpusic/grpool
			dpaks/goworkers
	uber-go/atomic
		定义和封装了几种与常见类型相对应的原子操作类型，这些类型提供了原子操作的方法
	marusama/semaphore
		实现了一个可以动态更改资源容量的信号量，也是一个非常有特色的实现
	marusama/cyclicbarrier
		github.com/marusama/cyclicbarrier
	18 及后面的三方库尚未记录



*/
