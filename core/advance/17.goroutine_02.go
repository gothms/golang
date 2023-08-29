package advance

import (
	"fmt"
	"sync/atomic"
	"time"
)

/*
go语句及其执行规则（下）

知识扩展
问题 1：怎样才能让主 goroutine 等待其他 goroutine？
 1. time.Sleep
    time.Sleep函数会在被调用时用当前的绝对时间，再加上相对时间计算出在未来的恢复运行时间
    一旦到达恢复运行时间，当前的 goroutine 就会从“睡眠”中醒来，并开始继续执行后边的代码
 2. channel
    在每个手动启用的 goroutine 即将运行完毕的时候，我们都要向该通道发送一个值
    通道类型为 struct{}，struct{}类型值的表示法只有一个，即：struct{}{}
    且它占用的内存空间是0字节，这个值在整个 Go 程序中永远都只会存在一份
    3.sync.WaitGroup
问题 2：怎样让我们启用的多个 goroutine 按照既定的顺序运行？
	示例：GoroutineTest2()
		参考方案
		异步发起的go函数得到了同步地（或者说按照既定顺序地）执行
	sync/atomic：原子操作
		操作变量count的时候使用的都是原子操作
		这是由于trigger函数会被多个 goroutine 并发地调用，所以它用到的非本地变量count，就被多个用户级线程共用了
		因此，对它的操作就产生了竞态条件（race condition），破坏了程序的并发安全性

总结
	go函数的实际执行顺序往往与其所属的go语句的执行顺序（或者说 goroutine 的启用顺序）不同
	而且默认情况下的执行顺序是不可预知的

思考
	runtime包中提供了哪些与模型三要素 G、P 和 M 相关的函数？
A
	查阅文档
	https://golang.google.cn/pkg/runtime/
	https://studygolang.com/pkgdoc
*/

// GoroutineTest1 16_goroutine_01
func GoroutineTest1() {
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println(i)
		}()
	}
}

// GoroutineTest2 vs GoroutineTest1
// 传参保证每个 goroutine 都可以拿到一个唯一的整数
// 在go语句被执行时，我们传给go函数的参数i会先被求值，如此就得到了当次迭代的序号
// 之后，无论go函数会在什么时候执行，这个参数值都不会变
func GoroutineTest2() {
	for i := 0; i < 10; i++ {
		go func(j int) {
			fmt.Println(j) // 打印 0-9，但基本打印不完
		}(i) // 打印的一定会是那个当次迭代的序号
	}
}

// GoroutineSync 启用的多个 goroutine 按照既定的顺序运行
func GoroutineSync() {
	count := uint32(0)                     // 让count变量成为一个信号，它的值总是下一个可以调用打印函数的go函数的序号
	trigger := func(i uint32, fn func()) { // trigger函数实现了一种自旋（spinning）。除非发现条件已满足，否则它会不断地进行检查
		for {
			// 选用的原子操作函数对被操作的数值的类型有约束，所以才都使用 uint32
			if n := atomic.LoadUint32(&count); n == i { // 操作变量count的时候使用的都是原子操作
				//if count == i {	// 替代上一行，效果一样，why？
				fn()
				atomic.AddUint32(&count, 1) // +1
				break                       // 显式地退出当前的循环
			}
			time.Sleep(time.Nanosecond) // 可注释
		}
	}
	for i := uint32(0); i < 10; i++ {
		go func(i uint32) {
			fn := func() { // 无参数声明也无结果声明的函数类型
				fmt.Println(i)
			}
			trigger(i, fn)
		}(i) // 保证 i 一定会是那个当次迭代的序号
	}
	trigger(10, func() {}) // 阻塞
}
