package performance

/*
安装 graphviz
	https://www.graphviz.org/download/
	brew install graphviz
	将 $GOPATH/bin 加入到 $PATH
		Mac OS: 在 .bash_profile 中修改路径
安装 go-torch
	$ go get -u github.com/uber/go-torch
	下载并复制 flamegraph.pl (https://github.com/brendangregg/FlameGraph) 至 $GOPATH/bin 路径下
	将 $GOPATH/bin 加入 $PATH

	go 1.1 后已内置

性能调优过程
	S->设定优化目标->分析系统瓶颈点->优化瓶颈点->E
		   ↑_______________________丨
常见分析指标
	Wall Time：运行的绝对时间(包括阻塞、等待外部响应...)
	CPU Time：CPU消耗时间
	Block Time：
	Memory allocation：内存分配
	GC times/time spent：GC次数和耗时
*/
