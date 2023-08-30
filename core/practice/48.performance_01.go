package practice

/*
程序性能分析基础（上）

Go 语言为程序开发者们提供了丰富的性能分析 API，和非常好用的标准工具
	API主要存在于三个代码包
	1. runtime/pprof
	2. net/http/pprof
	3. runtime/trace
	另外，runtime代码包中还包含了一些更底层的 API
	它们可以被用来收集或输出 Go 程序运行过程中的一些关键指标，并帮助我们生成相应的概要文件以供后续分析时使用
标准工具
	主要有go tool pprof和go tool trace
		可以解析概要文件中的信息，并以人类易读的方式把这些信息展示出来
	go test
		go test命令也可以在程序测试完成后生成概要文件
		可以很方便地使用前面那两个工具读取概要文件，并对被测程序的性能加以分析
在 Go 语言中，用于分析程序性能的概要文件有三种
	CPU 概要文件（CPU Profile）
	内存概要文件（Mem Profile）
	阻塞概要文件（Block Profile）
这些概要文件中包含的都是：在某一段时间内，对 Go 程序的相关指标进行多次采样后得到的概要信息
	CPU 概要文件（CPU Profile）
		每一段独立的概要信息都记录着，在进行某一次采样的那个时刻，CPU 上正在执行的 Go 代码
	内存概要文件（Mem Profile）
		每一段概要信息都记载着，在某个采样时刻，正在执行的 Go 代码以及堆内存的使用情况，这里包含已分配和已释放的字节数量和对象数量
	阻塞概要文件（Block Profile）
		每一段概要信息，都代表着 Go 程序中的一个 goroutine 阻塞事件
go tool pprof
	在默认情况下，这些概要文件中的信息并不是普通的文本，它们都是以二进制的形式展现的
	如果你使用一个常规的文本编辑器查看它们的话，那么肯定会看到一堆“乱码”
	可以通过 go tool pprof 进入一个基于命令行的交互式界面，并对指定的概要文件进行查阅
	示例
		$ go tool pprof cpuprofile.out
		Type: cpu
		Time: Nov 9, 2018 at 4:31pm (CST)
		Duration: 7.96s, Total samples = 6.88s (86.38%)
		Entering interactive mode (type "help" for commands, "o" for options)
		(pprof)
protocol buffers
	概要文件中的信息不是普通的文本，它们是通过 protocol buffers 生成的二进制数据流，或者说字节流
	概括来讲，protocol buffers 是一种数据序列化协议，同时也是一个序列化工具
	protocol buffers 定义和实现了一种“可以让数据在结构形态和扁平形态之间互相转换”的方式
	序列化
		它可以把一个值，比如一个结构体或者一个字典，转换成一段字节流
	反序列化
		把经过它生成的字节流反向转换为程序中的一个值
	优势
		比如，它可以在序列化数据的同时对数据进行压缩，所以它生成的字节流，通常都要比相同数据的其他格式（例如 XML 和 JSON）占用的空间明显小很多
		又比如，它既能让我们自己去定义数据序列化和结构化的格式，也允许我们在保证向后兼容的前提下去更新这种格式
		正因为这些优势，Go 语言从 1.8 版本开始，把所有 profile 相关的信息生成工作都交给 protocol buffers 来做
	用途
		Protocol buffers 的用途非常广泛，并且在诸如数据存储、数据传输等任务中有着很高的使用率

问题：怎样让程序对 CPU 概要信息进行采样？
	需要用到runtime/pprof包中的 API
	在我们想让程序开始对 CPU 概要信息进行采样的时候，需要调用这个代码包中的StartCPUProfile函数
	而在停止采样的时候则需要调用该包中的StopCPUProfile函数
问题解析
	赫兹
		也称 Hz，是从英文单词“Hertz”（一个英文姓氏）音译过来的一个中文词。它是 CPU 主频的基本单位
		CPU 的主频指的是，CPU 内核工作的时钟频率，也常被称为 CPU clock speed
		这个时钟频率的倒数即为时钟周期（clock cycle），也就是一个 CPU 内核执行一条运算指令所需的时间，单位是秒
	示例
		主频为1000Hz 的 CPU，它的单个内核执行一条运算指令所需的时间为0.001秒，即1毫秒
		又例如，我们现在常用的3.2GHz 的多核 CPU，其单个内核在1个纳秒的时间里就可以至少执行三条运算指令
	runtime/pprof.StartCPUProfile函数
		在被调用的时候，先会去设定 CPU 概要信息的采样频率，并会在单独的 goroutine 中进行 CPU 概要信息的收集和输出
		注意，StartCPUProfile函数设定的采样频率总是固定的，即：100赫兹。也就是说，每秒采样100次，或者说每10毫秒采样一次
		StartCPUProfile函数设定的 CPU 概要信息采样频率，相对于现代的 CPU 主频来说是非常低的。两个原因：
			一方面，过高的采样频率会对 Go 程序的运行效率造成很明显的负面影响
				因此，runtime包中SetCPUProfileRate函数在被调用的时候，会保证采样频率不超过1MHz（兆赫），也就是说，它只允许每1微秒最多采样一次
				StartCPUProfile函数正是通过调用这个函数来设定 CPU 概要信息的采样频率的
			另一方面，经过大量的实验，Go 语言团队发现 100Hz 是一个比较合适的设定
				因为这样做既可以得到足够多、足够有用的概要信息，又不至于让程序的运行出现停滞
				另外，操作系统对高频采样的处理能力也是有限的，一般情况下，超过 500Hz 就很可能得不到及时的响应了
		在StartCPUProfile函数执行之后，一个新启用的 goroutine 将会负责执行 CPU 概要信息的收集和输出
		直到runtime/pprof包中的StopCPUProfile函数被成功调用
	StopCPUProfile函数
		也会调用runtime.SetCPUProfileRate函数，并把参数值（采样频率）设为0
			这会让针对 CPU 概要信息的采样工作停止
		同时，它也会给负责收集 CPU 概要信息的代码一个“信号”，以告知收集工作也需要停止了
			在接到这样的“信号”之后，那部分程序将会把这段时间内收集到的所有 CPU 概要信息，全部写入到我们在调用StartCPUProfile函数的时候指定的写入器中
			只有在上述操作全部完成之后，StopCPUProfile函数才会返回

总结
	与程序性能分析有关的 API 主要存在于runtime、runtime/pprof和net/http/pprof这几个代码包中
		它们可以帮助我们收集相应的性能概要信息，并把这些信息输出到我们指定的地方
	Go 语言的运行时系统会根据要求对程序的相关指标进行多次采样，并对采样的结果进行组织和整理，最后形成一份完整的性能分析报告
		这份报告就是我们一直在说的概要信息的汇总
	一般情况下，我们会把概要信息输出到文件
		根据概要信息的不同，概要文件的种类主要有三个，分别是
		CPU 概要文件（CPU Profile）、内存概要文件（Mem Profile）和阻塞概要文件（Block Profile）
*/
