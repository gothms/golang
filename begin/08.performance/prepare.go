package performance

/*
安装 graphviz
	https://www.graphviz.org/download/
	brew install graphviz
	将 $GOPATH/bin 加入到 $PATH
		Mac OS: 在 .bash_profile 中修改路径
	安装：
		配置环境变量：Path 新增 E:\Program Files\Graphviz\bin
		PS 验证环境变量：dot -version
安装 go-torch(火炬图)
	$ go get -u github.com/uber/go-torch
	下载并复制 flamegraph.pl (https://github.com/brendangregg/FlameGraph) 至 $GOPATH/bin 路径下
	将 $GOPATH/bin 加入 $PATH

	go 1.1 后已内置
*/
