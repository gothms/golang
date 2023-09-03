package concurrent

import (
	"golang.org/x/sync/singleflight"
	"net/http"
)

/*
SingleFlight 和 CyclicBarrier：请求合并和循环栅栏该怎么用？

扩展并发原语：SingleFlight 和 CyclicBarrier
	SingleFlight 的作用是将并发请求合并成一个请求，以减少对下层服务的压力
	CyclicBarrier 是一个可重用的栅栏并发原语，用来控制一组请求同时执行的数据结构
	它们两个并没有直接的关系

请求合并 SingleFlight
	简介
		SingleFlight 是 Go 开发组提供的一个扩展并发原语
		它的作用是，在处理多个 goroutine 同时调用同一个函数的时候，只让一个 goroutine 去调用这个函数
		等到这个 goroutine 返回结果的时候，再把结果返回给这几个同时调用的 goroutine，这样可以减少并发调用的数量
	SingleFlight vs sync.Once
		标准库中的 sync.Once 也可以保证并发的 goroutine 只会执行一次函数 f
		其实，sync.Once 不是只在并发的时候保证只有一个 goroutine 执行函数 f，而是会保证永远只执行一次
		而 SingleFlight 是每次调用都重新执行，并且在多个请求同时调用的时候只有一个执行
		它们两个面对的场景是不同的
		sync.Once 主要是用在单次初始化场景中，而 SingleFlight 主要用在合并并发请求的场景中，尤其是缓存场景
	场景举例
		在面对秒杀等大并发请求的场景，而且这些请求都是读请求时，你就可以把这些请求合并为一个请求
		这样，你就可以将后端服务的压力从 n 降到 1
		尤其是在面对后端是数据库这样的服务的时候，采用 SingleFlight 可以极大地提高性能
	官方库
		标准库：internal/singleflight/singleflight.go
		扩展库：singleflight/singleflight.go
实现原理
	SingleFlight 使用互斥锁 Mutex 和 Map 来实现
		Mutex 提供并发时的读写保护，Map 用来保存同一个 key 的正在处理（in flight）的请求
	SingleFlight 数据结构
		type Group
			func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool)
			func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result
			func (g *Group) Forget(key string)
		Do 方法
			执行一个函数，并返回函数执行的结果
			你需要提供一个 key，对于同一个 key，在同一时间只有一个在执行，同一个 key 并发的请求会等待
			第一个执行的请求返回的结果，就是它的返回结果
			函数 fn 是一个无参的函数，返回一个结果或者 error，而 Do 方法会返回函数执行的结果或者是 error，shared 会指示 v 是否返回给多个请求
		DoChan 方法
			类似 Do 方法，只不过是返回一个 chan
			等 fn 函数执行完，产生了结果以后，就能从这个 chan 中接收这个结果
		Forget 方法
			告诉 Group 忘记这个 key
			这样一来，之后这个 key 请求会执行 f，而不是等待前一个未完成的 fn 函数的结果
	实现
		首先，SingleFlight 定义一个辅助对象 call，这个 call 就代表正在执行 fn 函数的请求或者是已经执行完的请求
		Group 代表 SingleFlight
应用场景
	Go 代码库中有两个地方用到了 SingleFlight
		net/lookup.go
			如果同时有查询同一个 host 的请求，lookupGroup 会把这些请求 merge 到一起，只需要一个请求就可以了
				// lookupGroup merges LookupIPAddr calls together for lookups for the same
				// host. The lookupGroup key is the LookupIPAddr.host argument.
				// The return values are ([]IPAddr, error).
				lookupGroup singleflight.Group
		cmd/go/internal/vcs/vcs.go
			Go 在查询仓库版本信息时，将并发的请求合并成 1 个请求
				func metaImportsForPrefix(importPrefix string, mod ModuleMode, security web.SecurityMode) (*urlpkg.URL, []metaImport, error) {
					// 使用缓存保存请求结果
					setCache := func(res fetchResult) (fetchResult, error) {
						fetchCacheMu.Lock()
						defer fetchCacheMu.Unlock()
						fetchCache[importPrefix] = res
						return res, nil
					}
					// 使用 SingleFlight请求
					resi, _, _ := fetchGroup.Do(importPrefix, func() (resi any, err error) {
						fetchCacheMu.Lock()
						// 如果缓存中有数据，那么直接从缓存中取
						if res, ok := fetchCache[importPrefix]; ok {
							fetchCacheMu.Unlock()
							return res, nil
						}
						fetchCacheMu.Unlock()
						...
				}
			缓存问题：代码中，会把结果放在缓存中，这也是常用的一种解决缓存击穿的例子
				设计缓存问题时，我们常常需要解决缓存穿透、缓存雪崩和缓存击穿问题
				缓存击穿问题是指，在平常高并发的系统中，大量的请求同时查询一个 key 时，如果这个 key 正好过期失效了，就会导致大量的请求都打到数据库上
				这就是缓存击穿
	用 SingleFlight 来解决缓存击穿问题再合适不过了
		因为，这个时候，只要这些对同一个 key 的并发请求的其中一个到数据库中查询，就可以了，这些并发的请求可以共享同一个结果
		因为是缓存查询，不用考虑幂等性问题（分布式环境下的一个常见问题，一般是指我们在进行多次操作时，所得到的结果是一样的，即多次运算结果是一致的）
		事实上，在 Go 生态圈知名的缓存框架 groupcache 中，就使用了较早的 Go 标准库的 SingleFlight 实现
	groupcache
		groupcache 中的 SingleFlight 只有一个方法
			func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error)
		SingleFlight 的作用是，在加载一个缓存项的时候，合并对同一个 key 的 load 的并发请求
			type Group struct {
				//......
				// loadGroup ensures that each key is only fetched once
				// (either locally or remotely), regardless of the number of
				// concurrent callers.
				loadGroup flightGroup
				//......
			}
			func (g *Group) load(ctx context.Context, key string, dest Sink) (value ByteView, err error) {
				viewi, err := g.loadGroup.Do(key, func() (interface{}, error) {
					// 从cache, peer, local尝试查询cache
					return value, nil
				})
				if err == nil {
					value = viewi.(ByteView)
				}
				return
			}
	其他知名项目
		如 Cockroachdb（小强数据库）、CoreDNS（DNS 服务器）等都有 SingleFlight 应用
	小结
		使用 SingleFlight 时，可以通过合并请求的方式降低对下游服务的并发压力，从而提高系统的性能
		常常用于缓存系统中

循环栅栏 CyclicBarrier




实现原理





并发趣题：一氧化二氢制造工厂












总结
	用 WaitGroup 来实现这个水分子制造工厂的例子
		使用 WaitGroup 非常复杂，而且，重用和 Done 方法的调用有并发的问题，程序可能 panic
		远远没有使用循环栅栏 CyclicBarrier 更加简单直接
	建议
		了解一些并发原语，甚至是从其它编程语言、操作系统中学习更多的并发原语
		在面对并发场景的时候，你也能更加游刃有余

思考
	1.你觉得，SingleFlight 能不能合并并发的写操作呢？
	2.如果大自然的搬运工工厂生产的液体是双氧水（双氧水分子是两个氢原子和两个氧原子），你又该怎么实现呢？
*/

func SingleFlight() {
	http.Client{}
	singleflight.Group{}
}
