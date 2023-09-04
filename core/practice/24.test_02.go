package practice

/*
测试的基本规则和流程（下）

知识扩展
问题 1：怎样解释功能测试的测试结果？
	$ go test golang/core/advance/test
		导入路径为 golang/core/advance/test 的代码包进行测试
		ok      golang/core/basic/test  0.034s
		释义：
			ok表示此次测试成功，也就是说没有发现测试结果不如预期的情况
			golang/core/basic/test，被测代码包的导入路径
			0.034s，此次对该代码包的测试所耗费的时间，即 34 毫秒
		cached：
			由于测试代码与被测代码都没有任何变动，所以go test命令直接把之前缓存测试成功的结果打印出来了
	go env GOCACHE
		C:\Users\sc\AppData\Local\go-build
		查看缓存目录的路径
		缓存的数据总是能够正确地反映出当时的各种源码文件、构建环境、编译器选项等等的真实情况
			一旦有任何变动，缓存数据就会失效，go 命令就会再次真正地执行操作
		go 命令会定期地删除最近未使用的缓存数据
			go clean -cache：手动删除所有的缓存数据
		对于测试成功的结果，go 命令也是会缓存的
			运行go clean -testcache将会删除所有的测试结果缓存
			不过，这样做肯定不会删除任何构建结果缓存
		此外，设置环境变量GODEBUG的值也可以稍稍地改变 go 命令的缓存行为
			比如，设置值为gocacheverify=1将会导致 go 命令绕过任何的缓存数据，而真正地执行操作并重新生成所有结果，然后再去检查新的结果与现有的缓存数据是否一致
		并不用在意缓存数据的存在，因为它们肯定不会妨碍go test命令打印正确的测试结果
	t.Log & t.Logf ...
		打印常规的测试日志，只不过当测试成功的时候，go test命令就不会打印这类日志了
	-v
		在测试结果中看到所有的常规测试日志
	测试失败
		t.Fail：Error
			调用t.Fail方法时，虽然当前的测试函数会继续执行下去，但是结果会显示该测试失败
			对于失败测试的结果，go test命令并不会进行缓存，所以，这种情况下的每次测试都会产生全新的结果
			如果测试失败了，那么go test命令将会导致：失败的测试函数中的常规测试日志一并被打印出来
		t.FailNow：Fatal
			让某个测试函数在执行的过程中立即失败
			在t.FailNow()执行之后，当前函数会立即终止执行，该行代码/异常之后的所有代码都会失去执行机会
		t.Error方法或者t.Errorf方法
			在测试失败的同时打印失败测试日志
			它们相当于t.Log/t.Logf方法和t.Fail方法的连续调用
		t.Fatal方法和t.Fatalf方法
			在打印失败错误日志之后立即终止当前测试函数的执行并宣告测试失败
			相当于它们在最后都调用了t.FailNow方法
问题 2：怎样解释性能测试的测试结果？
	$ go test -bench="." -run=^$ golang/core/practice/test
		-bench=.
			.表明需要执行任意名称的性能测试函数，函数名称还是要符合 Go 程序测试的基本规则的
		-run=^$
			表明需要执行哪些功能测试函数，这同样也是以函数名称为依据的
			标记的值^$意味着：只执行名称为空的功能测试函数，换句话说，不执行任何功能测试函数
		正则表达式
			. 和 ^$，这两个标记的值都是正则表达式。实际上，它们只能以正则表达式为值
		-run
			运行go test命令的时候不加-run标记，那么就会使它执行被测代码包中的所有功能测试函数
	BenchmarkStringBuilder-8         1927489               635.7 ns/op
		BenchmarkStringBuilder-8
			单个性能测试的名称，它表示命令执行了性能测试函数 BenchmarkStringBuilder
			并且当时所用的最大 P 数量为8，相当于可以同时运行 goroutine 的逻辑 CPU 的最大个数
		逻辑 CPU
			也可以被称为 CPU 核心，但它并不等同于计算机中真正的 CPU 核心
			只是 Go 语言运行时系统内部的一个概念，代表着它同时运行 goroutine 的能力
			一台计算机的 CPU 核心的个数，意味着它能在同一时刻执行多少条程序指令，代表着它并行处理程序指令的能力
		runtime.GOMAXPROCS & -cpu
			可以通过调用 runtime.GOMAXPROCS函数改变最大 P 数量
			也可以在运行go test命令时，加入标记-cpu来设置一个最大 P 数量的列表，以供命令在多次测试时使用
		执行次数 1927489：b.N
			go test命令最后一次执行性能测试函数的时候，被测函数被执行的实际次数
			go test命令在执行性能测试函数的时候会给它一个正整数，若该测试函数的唯一参数的名称为b，则该正整数就由b.N代表
			go test命令会先尝试把b.N设置为1，然后执行测试函数
			如果测试函数的执行时间没有超过上限，此上限默认为 1 秒，那么命令就会改大b.N的值，然后再次执行测试函数，如此往复，直到这个时间大于或等于上限为止
			当某次执行的时间大于或等于上限时，b.N的值就会被包含在测试结果中，也就是 1927489
			执行次数，指的是被测函数的执行次数，而不是性能测试函数的执行次数
		635.7 ns/op
			表明单次执行被测函数的平均耗时为 635.7 纳秒
			是通过将最后一次执行测试函数时的执行时间，除以（被测函数的）执行次数而得出的

思考
	在编写示例测试函数的时候，我们怎样指定预期的打印内容？

command
	参考 02.flag.go
	go test golang/core/advance/test
		代码包测试
	go test -v -run 22_test.go golang/core/advance/test
		.go文件测试
		等同于：
		cd 到 golang/core/advance/test 目录
		go test -v -run 22_test.go
	go test -v -run TestDeferStack golang/core/advance/test
		指定测试方法测试
		等同于：
		cd 到 golang/core/advance/test 目录
		go test -v -run TestDeferStack 22_test.go
	benchmark
		go test -bench="." -run=^$ golang/core/practice/test -cpu 1,2,4 -count 3
		go test -bench="." -run=^$ golang/core/practice/test -cpu 1,2,4 -count=3
		go test -bench="." golang/core/practice/test
		go test -bench="." 24_test.go -benchmem
		go test -bench=BenchmarkStringBuilder golang/core/practice/test
参考
	https://juejin.cn/post/7131966825474031646
*/
