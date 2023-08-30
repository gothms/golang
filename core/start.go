package core

/*
代码
	https://github.com/hyper0x/Golang_Puzzlers

07-15
	Go 语言（全部）数据类型
07-22
	Go 语言的所有内建数据类型，以及非常有特色的那些流程和语句
	sync.Map：高级数据结构
23-25
	测试
26-49
	Go 语言标准库中一些核心代码包
26-35：同步工具
	Channel
	互斥锁：sync.Locker Mutex & RWMutex
	条件变量：sync.Cond
	原子操作：sync/atomic
	sync.WaitGroup & sync.Once
	context.Context
	sync.Pool

	sync包中的高级同步工具，都是基于基本的通用工具，实现了某一种特定的功能
36-49
	Go 语言标准库中常用代码包

数据类型
	切片：底层基于数组
	通道：传递数据
	函数：一等类型
	结构体：可面向对象
	接口：无侵入实现
	...
语法
	go：异步编程
	defer：函数最后关卡
	switch：可类型判断
	select：多通道操作
	panic & recover：特色异常处理函数
测试
	独立的测试源码文件
	三种功用不同的测试函数
	专用的testing代码包
	功能强大的go test命令
并发编程工具
	互斥锁 & 读写锁：sync.Locker Mutex & RWMutex
	条件变量：sync.Cond
	原子操作：sync/atomic
	Go 特有数据类型
		sync.Once：单次执行
		sync.Pool：临时对象池
		sync.WaitGroup：多 goroutine 协作流程
		context.Context：多 goroutine 协作流程
	sync.Map：并发安全字典
*/
