package practice

/*
基于HTTP协议的网络服务

HTTP 超文本协议
	用net.Dial或net.DialTimeout函数来访问基于 HTTP 协议的网络服务是完全没有问题的
	HTTP 协议是基于 TCP/IP 协议栈的，并且它也是一个面向普通文本的协议
	原则上，我们使用任何一个文本编辑器，都可以轻易地写出一个完整的 HTTP 请求报文
	只要你搞清楚了请求报文的头部（header）和主体（body）应该包含的内容，这样做就会很容易
net/http
	只是访问基于 HTTP 协议的网络服务的话，那么使用net/http代码包中的程序实体来做，显然会更加便捷
http.Get
	http.Get函数会返回两个结果值
		第一个结果值的类型是*http.Response，它是网络服务给我们传回来的响应内容的结构化表示
		第二个结果值是error类型的，它代表了在创建和发送 HTTP 请求，以及接收和解析 HTTP 响应的过程中可能发生的错误
	http.Get函数会在内部使用缺省的 HTTP 客户端，并且调用它的Get方法以完成功能
		这个缺省的 HTTP 客户端是由net/http包中的公开变量DefaultClient代表的，其类型是*http.Client
		它的基本类型也是可以被拿来使用的，甚至它还是开箱即用的
			var httpClient http.Client
			resp, err := httpClient.Get(url)
			等价于：resp, err := http.Get(url)
	http.Client
		一个结构体类型，并且它包含的字段都是公开的
		之所以该类型的零值仍然可用，是因为它的这些字段要么存在着相应的缺省值，要么其零值直接就可以使用，且代表着特定的含义

问题：http.Client类型中的Transport字段代表着什么？
典型回答
	Transport字段代表：向网络服务发送 HTTP 请求，并从网络服务接收 HTTP 响应的操作过程
		该字段的方法RoundTrip应该实现单次 HTTP 事务（基于 HTTP 协议的单次交互）需要的所有步骤
		这个字段是http.RoundTripper接口类型的，它有一个由http.DefaultTransport变量代表的缺省值
		在初始化一个http.Client类型的值的时候，如果没有显式地为该字段赋值，那么这个Client值就会直接使用DefaultTransport
	Timeout字段
		单次 HTTP 事务的超时时间，它是time.Duration类型的。它的零值是可用的，用于表示没有设置超时时间
问题解析
	DefaultTransport
		DefaultTransport的实际类型是*http.Transport(即为http.RoundTripper接口的默认实现)
		这个类型是可以被复用的，也推荐被复用，同时，它也是并发安全的。正因为如此，http.Client类型也拥有着同样的特质
		http.Transport类型，会在内部使用一个net.Dialer类型的值，并且它会把该值的Timeout字段的值，设定为30秒(操作超时)
		在DefaultTransport的值被初始化的时候，这样的Dialer值的DialContext方法会被赋给DefaultTransport的DialContext字段
	http.Transport类型其他的字段，一些是关于操作超时的
		IdleConnTimeout：含义是空闲的连接在多久之后就应该被关闭
		DefaultTransport：会把该字段的值设定为90秒
			如果该值为0，那么就表示不关闭空闲的连接。注意，这样很可能会造成资源的泄露
		ResponseHeaderTimeout：含义是，从客户端把请求完全递交给操作系统到从操作系统那里接收到响应报文头的最大时长
			DefaultTransport并没有设定该字段的值
		ExpectContinueTimeout：含义是，在客户端递交了请求报文头之后，等待接收第一个响应报文头的最长时间
			在客户端想要使用 HTTP 的“POST”方法把一个很大的报文体发送给服务端的时候，它可以先通过发送一个包含了“Expect: 100-continue”的请求报文头，来询问服务端是否愿意接收这个大报文体
			这个字段就是用于设定在这种情况下的超时时间的
			注意，如果该字段的值不大于0，那么无论多大的请求报文体都将会被立即发送出去。这样可能会造成网络资源的浪费
			DefaultTransport把该字段的值设定为了1秒
		TLSHandshakeTimeout：TLS 是 Transport Layer Security 的缩写，可以被翻译为传输层安全
			这个字段代表了基于 TLS 协议的连接在被建立时的握手阶段的超时时间。若该值为0，则表示对这个时间不设限
			DefaultTransport把该字段的值设定为了10秒
	一些与IdleConnTimeout相关的字段
		MaxIdleConns、MaxIdleConnsPerHost以及MaxConnsPerHost
			无论当前的http.Transport类型的值访问了多少个网络服务，MaxIdleConns字段都只会对空闲连接的总数做出限定
			而MaxIdleConnsPerHost字段限定的则是，该Transport值访问的每一个网络服务的最大空闲连接数
		每一个网络服务都会有自己的网络地址，可能会使用不同的网络协议，对于一些 HTTP 请求也可能会用到代理
		Transport值正是通过这三个方面的具体情况，来鉴别不同的网络服务的
			MaxIdleConnsPerHost字段的缺省值，由http.DefaultMaxIdleConnsPerHost变量代表，值为2
			也就是说，在默认情况下，对于某一个Transport值访问的每一个网络服务，它的空闲连接数都最多只能有两个
		与MaxIdleConnsPerHost字段的含义相似的，是MaxConnsPerHost字段
			不过，MaxConnsPerHost限制的是，针对某一个Transport值访问的每一个网络服务的最大连接数，不论这些连接是否是空闲的
			并且，该字段没有相应的缺省值，它的零值表示不对此设限
		DefaultTransport并没有显式地为MaxIdleConnsPerHost和MaxConnsPerHost这两个字段赋值
			但是它却把MaxIdleConns字段的值设定为了100
			在默认情况下，空闲连接的总数最大为100，而针对每个网络服务的最大空闲连接数为2
			注意，上述两个与空闲连接数有关的字段的值应该是联动的，所以，你有时候需要根据实际情况来定制它们
	为什么会出现空闲的连接
		HTTP 协议有一个请求报文头叫做“Connection”。在 HTTP 协议的 1.1 版本中，这个报文头的值默认是“keep-alive”
		在这种情况下的网络连接都是持久连接，它们会在当前的 HTTP 事务完成后仍然保持着连通性，因此是可以被复用的
		既然连接可以被复用，那么就会有两种可能
			一种可能是，针对于同一个网络服务，有新的 HTTP 请求被递交，该连接被再次使用
			另一种可能是，不再有对该网络服务的 HTTP 请求，该连接被闲置。就产生了空闲的连接
		另外，如果分配给某一个网络服务的连接过多的话，也可能会导致空闲连接的产生，因为每一个新递交的 HTTP 请求，都只会征用一个空闲的连接
		所以，为空闲连接设定限制，在大多数情况下都是很有必要的，也是需要斟酌的
	为什么不要彻底地杜绝空闲连接的产生
		可以在初始化Transport值的时候把它的DisableKeepAlives字段的值设定为true
		这时，HTTP 请求的“Connection”报文头的值就会被设置为“close”
		这会告诉网络服务，这个网络连接不必保持，当前的 HTTP 事务完成后就可以断开它了
		如此一来，每当一个 HTTP 请求被递交时，就都会产生一个新的网络连接
		这样做会明显地加重网络服务以及客户端的负载，并会让每个 HTTP 事务都耗费更多的时间
		所以，在一般情况下，我们都不要去设置这个DisableKeepAlives字段
	net.Dialer类型中的字段KeepAlive
		与 HTTP 持久连接并不是一个概念，KeepAlive是直接作用在底层的 socket 上的
		它的背后是一种针对网络连接（TCP 连接）的存活探测机制
		它的值用于表示每间隔多长时间发送一次探测包。当该值不大于0时，则表示不开启这种机制
		DefaultTransport会把这个字段的值设定为30秒

知识扩展
问题：http.Server类型的ListenAndServe方法都做了哪些事情？
	http.Server类型与http.Client是相对应的。http.Server代表的是基于 HTTP 协议的服务端，或者说网络服务
	http.Server类型的ListenAndServe方法
		功能是：监听一个基于 TCP 协议的网络地址，并对接收到的 HTTP 请求进行处理
		这个方法会默认开启针对网络连接的存活探测机制，以保证连接是持久的
		同时，该方法会一直执行，直到有严重的错误发生或者被外界关掉
		当被外界关掉时，它会返回一个由http.ErrServerClosed变量代表的错误值
	ListenAndServe方法主要会做几件事情
		1. 检查当前的http.Server类型的值（以下简称当前值）的Addr字段
			该字段的值代表了当前的网络服务需要使用的网络地址，即：IP 地址和端口号
			如果这个字段的值为空字符串，那么就用":http"代替
			也就是说，使用任何可以代表本机的域名和 IP 地址，并且端口号为80
		2. 通过调用net.Listen函数在已确定的网络地址上启动基于 TCP 协议的监听
		3. 检查net.Listen函数返回的错误值。如果该错误值不为nil，那么就直接返回该值
			否则，通过调用当前值的Serve方法准备接受和处理将要到来的 HTTP 请求
	衍生问题：
		一个是“net.Listen函数都做了哪些事情”
		另一个是“http.Server类型的Serve方法是怎样接受和处理 HTTP 请求的”
	net.Listen函数都做了哪些事情
		1. 解析参数值中包含的网络地址隐含的 IP 地址和端口号
		2. 根据给定的网络协议，确定监听的方法，并开始进行监听
	http.Server类型的Serve方法是怎样接受和处理 HTTP 请求的
		在一个for循环中，网络监听器的Accept方法会被不断地调用，该方法会返回两个结果值
			第一个结果值是net.Conn类型的，它会代表包含了新到来的 HTTP 请求的网络连接
			第二个结果值是代表了可能发生的错误的error类型值
				如果这个错误值不为nil，除非它代表了一个暂时性的错误，否则循环都会被终止
				如果是暂时性的错误，那么循环的下一次迭代将会在一段时间之后开始执行
		如果这里的Accept方法没有返回非nil的错误值
			那么这里的程序将会先把它的第一个结果值包装成一个*http.conn类型的值
			然后通过在新的 goroutine 中调用这个conn值的serve方法，来对当前的 HTTP 请求进行处理
		衍生的细节问题
			比如，这个conn值的状态有几种，分别代表着处理的哪个阶段？
			又比如，处理过程中会用到哪些读取器和写入器，它们的作用分别是什么？
			再比如，这里的程序是怎样调用我们自定义的处理函数的，等

总结
	Transport字段
		代表着单次 HTTP 事务的操作过程。它是http.RoundTripper接口类型的
		它的缺省值由http.DefaultTransport变量代表，其实际类型是*http.Transport
	http.Transport包含的字段
		DefaultTransport中的DialContext字段
		关于操作超时的字段，比如IdleConnTimeout和ExpectContinueTimeout，以及相关的MaxIdleConns和MaxIdleConnsPerHost等

思考
	怎样优雅地停止基于 HTTP 协议的网络服务程序？
*/
