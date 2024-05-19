package modes

/*
Go编程模式：Pipeline

Pipeline：一种把各种命令拼接起来完成一个更强功能的技术方法
	1)现在的
		流式处理
		函数式编程
		应用网关对微服务进行简单的API编排
	其实都是受 Pipeline 这种技术方式的影响
	2)Pipeline 很容地把代码按单一职责的原则拆分成多个高内聚低耦合的小模块
		然后轻松地把它们拼装起来，完成比较复杂的功能

HTTP 处理
	07 节 decorator HTTP 相关的一个示例
	E:\gothmslee\golang\modes\07.decoration.go
Channel 管理
	若要写一个 泛型的 Pipeline 框架并不容易，可以使用 Go Generation 实现
	但 Go 语言最具特色的 Go Routine 和 Channel 完全可以用来构造这种编程
		参考：Rob Pike 在 Go Concurrency Patterns: Pipelines and cancellation 这篇博客中介绍了一种编程模式
		https://blog.golang.org/pipelines
	Channel 转发函数





 


Fan In/Fan Out：一对多/多对一
	**看原文的代码**
		首先，我们制造了从 1 到 10000 的数组
		然后，把这堆数组全部 echo到一个 Channel 里—— in
		此时，生成 5 个 Channel，接着都调用 sum(prime(in)) ，于是，每个 Sum 的 Go Routine 都会开始计算和
		最后，再把所有的结果再求和拼起来，得到最终的结果
*/

// EchoFunc Channel 转发函数
type EchoFunc func([]int) <-chan int
type PipeFunc func(<-chan int) <-chan int

func pipeline(nums []int, echo EchoFunc, pipeFns ...PipeFunc) <-chan int {
	ch := echo(nums)
	for i := range pipeFns {
		ch = pipeFns[i](ch)
	}
	return ch
}
