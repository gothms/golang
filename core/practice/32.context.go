package practice

/*
context.Context类型

需求分析
	在使用WaitGroup值的时候，不能在一开始就确定执行子任务的 goroutine 的数量
	那么使用WaitGroup值来协调它们和分发子任务的 goroutine，就是有一定风险的，如并发地调用该值的Add方法，那么就很可能会引发 panic
	一个解决方案是：分批地启用执行子任务的 goroutine
示例：demo67.go
	在严格遵循“保证其计数周期的完整性”规则的前提下，分批地启用执行子任务的 goroutine
	最简单的方式就是使用for循环来作为辅助

问题：怎样使用context包中的程序实体，实现一对多的 goroutine 协作流程？
	context.Context
		在 Go 1.7 发布时才被加入到标准库的.而后，标准库中的很多其他代码包都为了支持它而进行了扩展
		包括：os/exec包、net包、database/sql包，以及runtime/pprof包和runtime/trace包，等
	同步工具
		Context类型之所以受到了标准库中众多代码包的积极支持，主要是因为它是一种非常通用的同步工具
		它的值不但可以被任意地扩散，而且还可以被用来传递额外的信息和信号
		Context类型可以提供一类代表上下文的值。此类值是并发安全的，也就是说它可以被传播给多个 goroutine
	接口类型
		由于Context类型实际上是一个接口类型，而context包中实现该接口的所有私有类型，都是基于某个数据类型的指针类型
		所以，如此传播并不会影响该类型值的功能和安全
	上下文“树”
		Context类型的值是可以繁衍的，这意味着我们可以通过一个Context值产生出任意个子值
		这些子值可以携带其父值的属性和数据，也可以响应我们通过其父值传达的信号
		正因为如此，所有的Context值共同构成了一颗代表了上下文全貌的树形结构
		这棵树的树根（即上下文根节点）是一个已经在context包中预定义好的Context值，它是全局唯一的
		通过调用context.Background函数，我们就可以获取到它
	context.Background
		上下文根节点仅仅是一个最基本的支点，它不提供任何额外的功能
		也就是说，它既不可以被撤销（cancel），也不能携带任何数据
	四个“繁衍”Context值的函数
		WithCancel、WithDeadline、WithTimeout和WithValue
		这些函数的第一个参数的类型都是context.Context，而名称都为parent
		即这个参数对应的都是它们将会产生的Context值的父值
		WithCancel：用于产生一个可撤销的parent的子值，退出 Context
		WithDeadline：产生一个会定时撤销的parent的子值，有截止时间的 Context
		WithTimeout：产生一个会定时撤销的parent的子值，有超时时间的 Context
		WithValue：产生一个会携带额外数据的parent的子值

知识扩展
问题 1：“可撤销的”在context包中代表着什么？“撤销”一个Context值又意味着什么？
	Done 方法：感知撤销信号
		func (c *cancelCtx) Done() <-chan struct{}
		返回一个元素类型为struct{}的接收通道，其用途并不是传递元素值，而是让调用方去感知“撤销”当前Context值的那个信号
		一旦当前的Context值被撤销，这里的接收通道就会被立即关闭
		对于一个未包含任何元素值的通道来说，它的关闭会使任何针对它的接收操作立即结束
		正因为如此，基于调用表达式cxt.Done()的接收操作，才能够起到感知撤销信号的作用
	Err() error：撤销原因
		值只可能等于context.Canceled变量的值，或者context.DeadlineExceeded变量的值
		Canceled 用于表示手动撤销，DeadlineExceeded 表示给定的过期时间已到，而导致的撤销
	context.WithCancel函数
		产生一个可撤销的Context值时，还会获得一个用于触发撤销信号的函数
		通过调用这个函数，我们就可以触发针对这个Context值的撤销信号
		一旦触发，撤销信号就会立即被传达给这个Context值，并由它的Done方法的结果值（一个接收通道）表达出来
		撤销函数只负责触发信号，而对应的可撤销的Context值也只负责传达信号，它们都不会去管后边具体的“撤销”操作
		实际上，我们的代码可以在感知到撤销信号之后，进行任意的操作，Context值对此并没有任何的约束
	“撤销”最原始的含义
		终止程序针对某种请求（比如 HTTP 请求）的响应，或者取消对某种指令（比如 SQL 指令）的处理
		这也是 Go 语言团队在创建context代码包，和Context类型时的初衷
		net包和database/sql包的 API 和源码，了解它们在这方面的典型应用
问题 2：撤销信号是如何在上下文树中传播的？
	WithCancel、WithDeadline和WithTimeout都是被用来基于给定的Context值产生可撤销的子值的
		WithCancel
			context包的WithCancel函数在被调用后会产生两个结果值
			第一个结果值就是那个可撤销的Context值，而第二个结果值则是用于触发撤销信号的函数
		撤销流程
			在撤销函数被调用之后，对应的Context值会先关闭它内部的接收通道，也就是它的Done方法会返回的那个通道
			然后，它会向它的所有子值（或者说子节点）传达撤销信号。这些子值会如法炮制，把撤销信号继续传播下去
			最后，这个Context值会断开它与其父值之间的关联
		WithDeadline函数或者WithTimeout函数
			两者生成的Context值也是可撤销的
			它们不但可以被手动撤销，还会依据在生成时被给定的过期时间，自动地进行定时撤销。这里定时撤销的功能是借助它们内部的计时器来实现的
			当过期时间到达时，这两种Context值的行为与Context值被手动撤销时的行为是几乎一致的，只不过前者会在最后停止并释放掉其内部的计时器
	context.WithValue
		通过调用context.WithValue函数得到的Context值是不可撤销的
		撤销信号在被传播时，若遇到它们则会直接跨过，并试图将信号直接传给它们的子值
问题 3：怎样通过Context值携带数据？怎样从中获取数据？
	存储
		WithValue函数在产生新的Context值的时候需要三个参数，即：父值、键和值
		与“字典对于键的约束”类似，这里键的类型必须是可判等的
		Context值并不是用字典来存储键和值
	获取数据
		Value方法就是被用来获取数据的
		调用含数据的Context值的Value方法时，它会先判断给定的键，是否与当前值中存储的键相等，如果相等就把该值中存储的值直接返回，否则就到其父值中继续查找
		如果其父值中仍然未存储相等的键，那么该方法就会沿着上下文根节点的方向一路查找下去
		除了含数据的Context值以外，其他几种Context值都是无法携带数据的。因此，Context值的Value方法在沿路查找的时候，会直接跨过那几种值
			如果我们调用的Value方法的所属值本身就是不含数据的，那么实际调用的就将会是其父辈或祖辈的Value方法
			Context值的实际类型，都属于结构体类型，并且它们都是通过“将其父值嵌入到自身”，来表达父子关系的
	Context接口并没有提供改变数据的方法
		在通常情况下，我们只能通过在上下文树中添加含数据的Context值来存储新的数据，或者通过撤销此种值的父值丢弃掉相应的数据
		如果你存储在这里的数据可以从外部改变，那么必须自行保证安全

总结
	Context类型
		是一个可以帮助我们实现多 goroutine 协作流程的同步工具
		还可以通过此类型的值传达撤销信号或传递数据
	Context类型的实际值大体上分为三种
		根Context值、可撤销的Context值和含数据的Context值
		根
			所有的Context值共同构成了一颗上下文树
			这棵树的作用域是全局的，而根Context值就是这棵树的根。它是全局唯一的，并且不提供任何额外的功能
		可撤销的Context值又分为
			只可手动撤销的Context值，和可以定时撤销的Context值
			可以通过生成它们时得到的撤销函数来对其进行手动的撤销
			定时撤销的时间必须在生成时就完全确定，并且不能更改。不过，可以在过期时间达到之前，对其进行手动的撤销
		撤销
			撤销”这个操作是Context值能够协调多个 goroutine 的关键所在。撤销信号总是会沿着上下文树叶子节点的方向传播开来
	valueCtx
		每个valueCtx值都可以存储一对键和值
		在我们调用它的Value方法的时候，它会沿着上下文树的根节点的方向逐个值的进行查找
		如果发现相等的键，它就会立即返回对应的值，否则将在最后返回nil
	含数据的Context值不能被撤销，而可撤销的Context值又无法携带数据
		由于它们共同组成了一个有机的整体（即上下文树），所以在功能上要比sync.WaitGroup强大得多

思考
	Context值在传达撤销信号的时候是广度优先的，还是深度优先的？其优势和劣势都是什么？
A
	深度优先
	优势和劣势都是
		直接分支的产生时间越早，其中的所有子节点就会越先接收到信号
		至于什么时候是优势、什么时候是劣势还要看具体的应用场景
	例如
		如果子节点的存续时间与资源的消耗是正相关的，那么这可能就是一个优势
		但是，如果每个分支中的子节点都很多，而且各个分支中的子节点的产生顺序并不依从于分支的产生顺序，那么这种优势就很可能会变成劣势
		最终的定论还是要看测试的结果
*/
