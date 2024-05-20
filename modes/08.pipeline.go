package modes

import (
	"fmt"
	"math"
	"sync"
)

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
		echo()函数，它会把一个整型数组放到一个 Channel 中，并返回这个 Channel
			func echo(nums []int) <-chan int {
			  out := make(chan int)
			  go func() {
				for _, n := range nums {
				  out <- n
				}
				close(out)
			  }()
			  return out
			}
		依照 echo 函数，实现（可以通过 Map/Reduce 编程模式或者是 Go Generation 的方式合并一下）
			平方函数
			过滤奇数函数
			求和函数
		也可以使用一个代理函数来完成
			type EchoFunc func([]int) <-chan int
			type PipeFunc func(<-chan int) <-chan int

			func pipeline(nums []int, echo EchoFunc, pipeFns ...PipeFunc) <-chan int {
				ch := echo(nums)
				for i := range pipeFns {
					ch = pipeFns[i](ch)
				}
				return ch
			}

Fan In/Fan Out
	动用 Go 语言的 Go Routine 和 Channel 还有一个好处，就是可以写出 1 对多，或多对 1 的 Pipeline，也就是 Fan In/ Fan Out
	通过并发的方式对一个很长的数组中的质数进行求和运算，我们想先把数组分段求和，然后再把它们集中起来
		首先，我们制造了从 1 到 10000 的数组
		然后，把这堆数组全部 echo到一个 Channel 里—— in
		此时，生成 5 个 Channel，接着都调用 sum(prime(in)) ，于是，每个 Sum 的 Go Routine 都会开始计算和
		最后，再把所有的结果再求和拼起来，得到最终的结果
	代码：func FanInOutTest()
	图示：08.pipeline.jpg

参考
	Go Concurrency Patterns – Rob Pike – 2012 Google I/O presents the basics of Go‘s concurrency primitives and several ways to apply them.
		https://www.youtube.com/watch?v=f6kdp27TYZs
	Advanced Go Concurrency Patterns – Rob Pike – 2013 Google I/O
		https://blog.golang.org/advanced-go-concurrency-patterns
		covers more complex uses of Go’s primitives, especially select.
			https://blog.golang.org/advanced-go-concurrency-patterns
	Squinting at Power Series – Douglas McIlroy's paper
		https://swtch.com/~rsc/thread/squint.pdf
		shows how Go-like concurrency provides elegant support for complex calculations.
		https://swtch.com/~rsc/thread/squint.pdf
*/

// FanInOutTest 通过并发的方式对一个很长的数组中的质数进行求和运算，我们想先把数组分段求和，然后再把它们集中起来
func FanInOutTest() {
	nums := makeRange(1, 10000)
	in := echo(nums)

	const nProcess = 5
	var chans [nProcess]<-chan int
	for i := range chans {
		chans[i] = sum(prime(in))
	}

	for n := range sum(merge(chans[:])) {
		fmt.Println(n)
	}
}

func echo(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}
func sum(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		var sum = 0
		for n := range in {
			sum += n
		}
		out <- sum
		close(out)
	}()
	return out
}
func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
func is_prime(value int) bool {
	for i := 2; i <= int(math.Floor(float64(value)/2)); i++ {
		if value%i == 0 {
			return false
		}
	}
	return value > 1
}

func prime(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if is_prime(n) {
				out <- n
			}
		}
		close(out)
	}()
	return out
}
func merge(cs []<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan int) {
			for n := range c {
				out <- n
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
