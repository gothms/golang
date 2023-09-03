package concurrent

/*
分组操作：处理一组子任务，该用什么并发原语？

分组编排
	共享资源保护、任务编排和消息传递是 Go 并发编程中常见的场景
	而分组执行一批相同的或类似的任务则是任务编排中一类情形
	分组编排的一些常用场景和并发原语，包括 ErrGroup、gollback、Hunch 和 schedgroup

ErrGroup

基本用法
	1.WithContext

	2.Go

	3.Wait


ErrGroup 使用例子
简单例子：返回第一个错误




更进一步，返回所有子任务的错误


任务执行流水线 Pipeline

扩展库





其它实用的 Group 并发原语
SizedGroup/ErrSizedGroup



gollback













Hunch

	All 方法

	Take 方法

	Last 方法

	Retry 方法

	Waterfall 方法

schedgroup







总结


思考
	官方扩展库 ErrGroup 没有实现可以取消子任务的功能，请你课下可以自己去实现一个子任务可取消的 ErrGroup
*/
