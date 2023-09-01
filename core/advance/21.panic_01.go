package advance

/*
panic函数、recover函数以及defer语句 （上）

panic
	运行时恐慌，一种在我们意料之外的程序异常
	这种异常只会在程序运行的时候被抛出来
	如果我们没有在程序里添加任何保护措施的话，程序（或者说代表它的那个进程）就会在打印出 panic 的详细情况之后，终止运行
示例：数据越界
	panic: runtime error: index out of range [8] with length 8

	goroutine 1 [running]:
	main.main()
			E:/gothmslee/golang/main/slice.go:27 +0x4b9
解析
	runtime error”的含义是，这是一个runtime代码包中抛出的 panic
		在这个 panic 中，包含了一个runtime.Error接口类型的值
		runtime.Error接口内嵌了error接口并做了一点点扩展，runtime包中有不少它的实现类型
		“panic：”右边的内容，正是这个 panic 包含的runtime.Error类型值的字符串表示形式
	panic 详情中一般还会包含与它的引发原因有关的 goroutine 的代码执行信息
		goroutine 1 [running]，表示有一个 ID 为1的 goroutine 在此 panic 被引发的时候正在运行
		这里的 ID 其实并不重要，因为它只是 Go 语言运行时系统内部给予的一个 goroutine 编号，我们在程序中是无法获取和更改的
	main.main()表明 了这个 goroutine 包装的go函数就是命令源码文件中的那个main函数
		也就是说这里的 goroutine 正是主 goroutine
		并指出的就是这个 goroutine 中的哪一行代码在此 panic 被引发时正在执行
		这包含了此行代码在其所属的源码文件中的行数，以及这个源码文件的绝对路径
		+0x4b9 代表的是：此行代码相对于其所属函数的入口程序计数偏移量，一般情况下它的用处并不大
	“exit status 2”表明我的这个程序是以退出状态码2结束运行的
		在大多数操作系统中，只要退出状态码不是0，都意味着程序运行的非正常结束
		在 Go 语言中，因 panic 导致程序结束运行的退出状态码一般都会是2

问题：从 panic 被引发到程序终止运行的大致过程是什么？
典型回答
	大致过程：某个函数中的某行代码有意或无意地引发了一个 panic
	这时，初始的 panic 详情会被建立起来，并且该程序的控制权会立即从此行代码转移至调用其所属函数的那行代码上，也就是调用栈中的上一级
		这也意味着，此行代码所属函数的执行随即终止。紧接着，控制权并不会在此有片刻停留，它又会立即转移至再上一级的调用代码处
		控制权如此一级一级地沿着调用栈的反方向传播至顶端，也就是我们编写的最外层函数那里
		这里的最外层函数指的是go函数，对于主 goroutine 来说就是main函数
		但是控制权也不会停留在那里，而是被 Go 语言运行时系统收回
	随后，程序崩溃并终止运行，承载程序这次运行的进程也会随之死亡并消失
	与此同时，在这个控制权传播的过程中，panic 详情会被逐渐地积累和完善，并会在程序终止之前被打印出来
问题解析
	有意引发 panic
		Go 语言的内建函数panic是专门用于引发 panic 的。panic函数使程序开发者可以在程序运行期间报告异常
	有意引发panic vs error
		这与从函数返回错误值的意义是完全不同的
		error
			当我们的函数返回一个非nil的错误值时，函数的调用方有权选择不处理，并且不处理的后果往往是不致命的
			“不致命”的意思是，不至于使程序无法提供任何功能（也可以说僵死）或者直接崩溃并终止运行（也就是真死）
		panic
			当一个 panic 发生时，如果我们不施加任何保护措施，那么导致的直接后果就是程序崩溃，就像前面描述的那样，这显然是致命的
	panic 详情会在控制权传播的过程中，被逐渐地积累和完善，并且，控制权会一级一级地沿着调用栈的反方向传播至顶端
		因此，在针对某个 goroutine 的代码执行信息中，调用栈底端的信息会先出现，然后是上一级调用的信息，以此类推，最后才是此调用栈顶端的信息
		深入地了解此过程，以及正确地解读 panic 详情应该是我们的必备技能，这在调试 Go 程序或者为 Go 程序排查错误的时候非常重要

思考
	一个函数怎样才能把 panic 转化为error类型值，并将其作为函数的结果值返回给调用方？
*/