package modes

/*
Go 编程模式：Kubernetes Visitor模式

Kubernetes 的 kubectl 命令中的使用到的一个编程模式：
	kubectl 主要使用到了两个，一个是 Builder，另一个是 Visitor
		Builder 参考 03.functional
	Visitor：面向对象设计模式中一个很重要的设计模式（https://en.wikipedia.org/wiki/Visitor_pattern）
		将算法与操作对象的结构分离的一种方法
		这种分离的实际结果是能够在不修改结构的情况下向现有对象结构添加新操作
		是遵循开放 / 封闭原则的一种方法

一个简单示例
	E:\gothmslee\golang\modes\visitor_01.go
		我们的代码中有一个Visitor的函数定义，还有一个Shape接口，这需要使用 Visitor函数作为参数
		我们的实例的对象 Circle和 Rectangle实现了 Shape 接口的 accept() 方法，这个方法就是等外面给我们传递一个 Visitor
		然后，我们实现两个 Visitor：一个是用来做 JSON 序列化的；另一个是用来做 XML 序列化的
	这段代码的目的就是想解耦数据结构和算法
		虽然使用 Strategy 模式也是可以完成的，而且会比较干净
		但是在有些情况下，多个 Visitor 是来访问一个数据结构的不同部分
		在这种情况下，数据结构有点像一个数据库，而各个 Visitor 会成为一个个的小应用
		kubectl 就是这种情况

Kubernetes 相关背景
	1)Kubernetes 抽象了很多种的 Resource，比如 Pod、ReplicaSet、ConfigMap、Volumes、Namespace、Roles……种类非常繁多
		这些东西构成了 Kubernetes 的数据模型（你可以看看 Kubernetes Resources 地图，了解下有多复杂）
		 Kubernetes Resources 地图：https://github.com/kubernauts/practical-kubernetes-problems/blob/master/images/k8s-resources-map.png
	2)kubectl 是 Kubernetes 中的一个客户端命令，操作人员用这个命令来操作 Kubernetes
		kubectl 会联系到 Kubernetes 的 API Server，API Server 会联系每个节点上的 kubelet ，从而控制每个节点
	3)kubectl 的主要工作是处理用户提交的东西（包括命令行参数、YAML 文件等）
		接着会把用户提交的这些东西组织成一个数据结构体，发送给 API Server
	4)相关的源代码在 src/k8s.io/cli-runtime/pkg/resource/visitor.go 中
		https://github.com/kubernetes/kubernetes/blob/cea1d4e20b4a7886d8ff65f34c6d4f95efcb4742/staging/src/k8s.io/cli-runtime/pkg/resource/visitor.go
	kubectl 的代码比较复杂，不过，简单来说，基本原理就是
		它从命令行和 YAML 文件中获取信息，通过 Builder 模式并把其转成一系列的资源，最后用 Visitor 模式来迭代处理这些 Reources

kubectl 的实现方法
	E:\gothmslee\golang\modes\visitor_02.go
	注意：
		先调用 OtherThingsVisitor 的 Visit 方法
			return v.visitor.Visit(func(info *Info, err error) error {
		由于创建方式：
			info := Info{}
			var v Visitor = &info
			v = LogVisitor{v}
			v = NameVisitor{v}
			v = OtherThingsVisitor{v}
		所以，otherThingsVisitor.visitor = NameVisitor
		依次往上，直到调用 Info 的 Visit 方法
	调用流程：
		otv.Visit
		nv.Visit
		lv.Visit
		info.Visit
			lv.VisitorFunc() -> LogVisitor() before call function
			nv.VisitorFunc() -> NameVisitor() before call function
			otv.VisitorFunc() -> OtherThingsVisitor() before call function
			loadFile->VisitorFunc() ->
				loadFile->VisitorFunc() 返回 -> 设置 info 参数
				otv.VisitorFunc() 返回 -> ==> OtherThings=We are running as remote team.
										OtherThingsVisitor() after call function
				nv.VisitorFunc() 返回 -> ==> Name=Hao Chen, NameSpace=MegaEase
										NameVisitor() after call function
				lv.VisitorFunc() 返回 -> LogVisitor() after call function
	几种功效：
		解耦了数据和程序
		使用了修饰器模式
		还做出了Pipeline的模式
	Visitor 修饰器：用 修饰器模式 重构
		用一个 DecoratedVisitor 的结构来存放所有的VistorFunc函数
		NewDecoratedVisitor 可以把所有的 VisitorFunc转给它，构造 DecoratedVisitor 对象
		DecoratedVisitor实现了 Visit() 方法，里面就是来做一个 for-loop，顺着调用所有的 VisitorFunc
*/
