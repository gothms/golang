package advance

import (
	"errors"
	"fmt"
)

/*
panic函数、recover函数以及defer语句（下）

知识扩展
问题 1：怎样让 panic 包含一个值，以及应该让它包含什么样的值？
	panic 的唯一参数
		在调用panic函数时，把某个值作为参数传给该函数就可以了
		由于panic函数的唯一一个参数是空接口（即interface{}）类型的，所以从语法上讲，它可以接受任何类型的值
		最好传入error类型的错误值，或者其他的可以被有效序列化的值
		这里的“有效序列化”指的是，可以更易读地去表示形式转换
	示例：有效序列化
		对于fmt包下的各种打印函数来说，error类型值的Error方法与其他类型值的String方法是等价的，它们的唯一结果都是string类型的
		在通过占位符%s打印这些值的时候，它们的字符串表示形式分别都是这两种方法产出的
		程序日志：
			一旦程序异常了，我们就一定要把异常的相关信息记录下来，这通常都是记到程序日志里
			我们在为程序排查错误的时候，首先要做的就是查看和解读程序日志
			而最常用也是最方便的日志记录方式，就是记下相关值的字符串表示形式
			如果你觉得某个值有可能会被记到日志里，那么就应该为它关联String方法
			如果这个值是error类型的，那么让它的Error方法返回你为它定制的字符串表示形式就可以了
		fmt.Sprintf，以及fmt.Fprintf这类可以格式化并输出参数的函数
		不过，它们在功能上，肯定远不如我们自己定义的Error方法或者String方法。因此，为不同的数据类型分别编写这两种方法总是首选
	panic 与日志
		相同的道理，在程序崩溃的时候，panic 包含的那个值字符串表示形式会被打印出来
		我们还可以施加某种保护措施，避免程序的崩溃
		这个时候，panic 包含的值会被取出，而在取出之后，它一般都会被打印出来或者记录到日志里
问题 2：怎样施加应对 panic 的保护措施，从而避免程序崩溃？
	recover
		Go 语言的内建函数recover专用于恢复 panic，或者说平息运行时恐慌
		recover函数无需任何参数，并且会返回一个空接口类型的值
		如果用法正确，这个值实际上就是即将恢复的 panic 包含的值
		并且，如果这个 panic 是因我们调用panic函数而引发的，那么该值同时也会是我们此次调用panic函数时，传入的参数值副本
	错误调用 recover
		示例：RecoverWrongCall()
		panic 一旦发生，控制权就会讯速地沿着调用栈的反方向传播
		所以，在panic函数调用之后的代码，根本就没有执行的机会
		先调用recover函数，再调用panic函数会怎么样呢？
		显然也是不行的，因为，如果在我们调用recover函数时未发生 panic，那么该函数就不会做任何事情，并且只会返回一个nil
	defer 调用 & 表达式限制
		限制
			针对 Go 语言内建函数的调用表达式，以及针对unsafe包中的函数的调用表达式
			对于go语句中的调用表达式，限制也是一样的
		调用
			defer语句就是被用来延迟执行代码的。延迟到该语句所在的函数即将执行结束的那一刻，无论结束执行的原因是什么
				即使导致它执行结束的原因是一个 panic 也会是这样
			被调用的函数可以是有名称的，也可以是匿名的。可以把这里的函数叫做defer函数或者延迟函数
			注意，被延迟执行的是defer函数，而不是defer语句
问题 3：如果一个函数中有多条defer语句，那么那几个defer函数调用的执行顺序是怎样的？
	在同一个函数中，defer函数调用的执行顺序与它们分别所属的defer语句的出现顺序（更严谨地说，是执行顺序）完全相反
	for & defer
		同一条defer语句每被执行一次，其中的defer函数调用就会产生一次，而且，这些函数调用同样不会被立即执行
		在defer语句每次执行的时候，Go 语言会把它携带的defer函数及其参数值另行存储到一个队列中
		这个队列与该defer语句所属的函数是对应的，并且，它是先进后出（FILO）的，相当于一个栈
		在需要执行某个函数中的defer函数调用的时候，Go 语言会先拿到对应的队列，然后从该队列中一个一个地取出defer函数及其参数值，并逐个执行调用

总结
	recover
		recover函数专用于恢复 panic，并且调用即恢复
		它在被调用时会返回一个空接口类型的结果值。如果在调用它时并没有 panic 发生，那么这个结果值就会是nil
		而如果被恢复的 panic 是我们通过调用panic函数引发的，那么它返回的结果值就会是我们传给panic函数参数值的副本
	defer
		对recover函数的调用只有在defer语句中才能真正起作用。defer语句是被用来延迟执行代码的
		它会让其携带的defer函数的调用延迟执行，并且会延迟到该defer语句所属的函数即将结束执行的那一刻
		同一条defer语句每被执行一次，就会产生一个延迟执行的defer函数调用

思考
	我们可以在defer函数中恢复 panic，那么可以在其中引发 panic 吗？
*/

func RecoverWrongCall() {
	fmt.Println("Enter function main.")
	// 引发 panic。
	panic(errors.New("something wrong"))
	p := recover()
	fmt.Printf("panic: %s\n", p) // panic: something wrong
	fmt.Println("Exit function main.")
}

func DeferStack() {
	// last defer
	// defer in for [2]
	// defer in for [1]
	// defer in for [0]
	// first defer
	defer fmt.Println("first defer")
	for i := 0; i < 3; i++ {
		defer fmt.Printf("defer in for [%d]\n", i)
	}
	defer fmt.Println("last defer")
}
