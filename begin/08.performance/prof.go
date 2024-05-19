package performance

/*
通过文件方式输出 Profile
	灵活性高，适用于特定代码段的分析
	通过手动调用 runtime/pprof 的 API
	API相关文档 https://studygolang.com/static/pkgdoc/pkg/runtime_pprof.htm
	$ go tool pprof [binary][binary.prof]
		[binary]：二进制
		[binary.prof]：要查看的 prof

测试：Go 内置了 runtime/pprof 包
	CPU
	程序堆
	Goroutine：pprof.Lookup("goroutine")
		Go 支持的多种 Profile
			go help testflag
			Lookup 不同的 Tag：https://go.dev/src/runtime/pprof/pprof.go
	创建输出文件：选择目录
		f, err := os.Create("cpu.prof")
		f, err := os.Create("prof/cpu.prof")

查看：$ go tool pprof prof cpu.prof：prof 不是指路径
	top：查看 top 的情况
		 flat  flat%   sum%        cum   cum%
		 440ms 44.00% 44.00%      540ms 54.00%  math/rand.(*Rand).Int31n
		 310ms 31.00% 75.00%      920ms 92.00%  golang/performance.fillMatrix
	top -cum：按 cum 排序
	list fillMatrix：详细分析 fillMatrix 函数，可以模糊匹配
	svg：生成 svg 图，需要安装 Graphviz
		使用浏览器打开 .svg
			CPU 时间：CPU耗时
			Wall time：挂钟时间
	go-torch：
		flamegraph.pl -h：测试
		问题：go-torch.exe 未生成
		解决：不使用 glide 时，可在 $GOPATH$\bin 目录下生成 go-torch.exe
	go-torch cpu.prof：
		报错：
		FATAL[02:53:12] Failed: could not generate flame graph: fork/exec E:\gothmslee\bin\flamegraph.pl: %1 is not a valid Win32 application.
		解决：
		安装 perl 并重启 goland，$ perl -v 查看版本信息
	exit：退出 prof
测试时生成 prof：
	go test -bench=. --cpuprofile=cpu.prof
	go test -bench=. --memprofile=mem.prof
	go test -bench=. --blockprofile=block.prof
	go tool pprof cpu.prof

	go help testflag

	示例：
	go test -bench=. --cpuprofile=cpu.prof
	go tool pprof cpu.prof
	go-torch cpu.prof
*/
