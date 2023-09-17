package concurrent

import (
	"context"
	"fmt"
	"github.com/marusama/cyclicbarrier"
	"golang.org/x/sync/semaphore"
	"log"
	"sort"
	"sync"
	"time"
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
	循环栅栏（CyclicBarrier），它常常应用于重复进行一组 goroutine 同时执行的场景中
		CyclicBarrier 允许一组 goroutine 彼此等待，到达一个共同的执行点
		同时，因为它可以被重复使用，所以叫循环栅栏。具体的机制是，大家都在栅栏前等待，等全部都到齐了，就抬起栅栏放行
	Java 实现
		事实上，这个 CyclicBarrier 是参考 Java CyclicBarrier 和 C# Barrier的功能实现的
		Java 提供了 CountDownLatch（倒计时器）和 CyclicBarrier（循环栅栏）两个类似的用于保证多线程到达同一个执行点的类
		只不过前者是到达 0 的时候放行，后者是到达某个指定的数的时候放行
		C# Barrier 功能也是类似的
		link：Java CyclicBarrier
		link：C# Barrier
	CyclicBarrier vs WaitGroup
		vs
			CyclicBarrier 和 WaitGroup 的功能有点类似，确实是这样
			不过，CyclicBarrier 更适合用在“固定数量的 goroutine 等待同一个执行点”的场景中
			而且在放行 goroutine 之后，CyclicBarrier 可以重复利用
			不像 WaitGroup 重用的时候，必须小心翼翼避免 panic
		图示
			17.cyclicbarrier_waitgroup.jpg
		CyclicBarrier
			处理可重用的多 goroutine 等待同一个执行点的场景的时候，CyclicBarrier 和 WaitGroup 方法调用的对应关系
			使用 WaitGroup 实现的话，调用比较复杂，不像 CyclicBarrier 那么清爽
			更重要的是，如果想重用 WaitGroup，你还要保证，将 WaitGroup 的计数值重置到 n 的时候不会出现并发问题
		WaitGroup
			WaitGroup 更适合用在“一个 goroutine 等待一组 goroutine 到达同一个执行点”的场景中，或者是不需要重用的场景中
	三方库
		github.com/marusama/cyclicbarrier
实现原理
	CyclicBarrier 有两个初始化方法
		1. New 方法，它只需要一个参数，来指定循环栅栏参与者的数量
		2. NewWithAction，它额外提供一个函数，可以在每一次到达执行点的时候执行一次
			具体的时间点是在最后一个参与者到达之后，但是其它的参与者还未被放行之前
			我们可以利用它，做放行之前的一些共享状态的更新等操作
		func New(parties int) CyclicBarrier
		func NewWithAction(parties int, barrierAction func() error) CyclicBarrier
	CyclicBarrier 接口
		type CyclicBarrier interface {
			Await(ctx context.Context) error // 等待所有的参与者到达，如果被ctx.Done()中断，会返回ErrBrokenBarrier
			Reset()                          // 重置循环栅栏到初始化状态。如果当前有等待者，那么它们会返回ErrBrokenBarrier
			GetNumberWaiting() int           // 返回当前等待者的数量
			GetParties() int                 // 参与者的数量
			IsBroken() bool                  // 循环栅栏是否处于中断状态
		}
	循环栅栏的使用
		循环栅栏的参与者只需调用 Await 等待，等所有的参与者都到达后，再执行下一步
		当执行下一步的时候，循环栅栏的状态又恢复到初始的状态了，可以迎接下一轮同样多的参与者
并发趣题：一氧化二氢制造工厂
	非常经典的并发编程的题目，非常适合使用循环栅栏
	题目描述
		有一个名叫大自然的搬运工的工厂，生产一种叫做一氧化二氢的神秘液体。这种液体的分子是由一个氧原子和两个氢原子组成的，也就是水
		这个工厂有多条生产线，每条生产线负责生产氧原子或者是氢原子，每条生产线由一个 goroutine 负责
		这些生产线会通过一个栅栏，只有一个氧原子生产线和两个氢原子生产线都准备好，才能生成出一个水分子，否则所有的生产线都会处于等待状态
		也就是说，一个水分子必须由三个不同的生产线提供原子，而且水分子是一个一个按照顺序产生的，每生产一个水分子，就会打印出 HHO、HOH、OHH 三种形式的其中一种
		HHH、OOH、OHO、HOO、OOO 都是不允许的
		生产线中氢原子的生产线为 2N 条，氧原子的生产线为 N 条
	实现
		定义一个 H2O 辅助数据类型，它包含两个信号量的字段和一个循环栅栏
			1. semaH 信号量：控制氢原子。一个水分子需要两个氢原子，所以，氢原子的空槽数资源数设置为 2
			2. semaO 信号量：控制氧原子。一个水分子需要一个氧原子，所以资源数的空槽数设置为 1
			3. 循环栅栏：等待两个氢原子和一个氧原子填补空槽，直到任务完成
		各条流水线的处理情况
			流水线分为氢原子处理流水线和氧原子处理流水线
		氢原子的流水线：如果有可用的空槽，氢原子的流水线的处理方法是 hydrogen
			hydrogen 方法就会占用一个空槽（h2o.semaH.Acquire），输出一个 H 字符，然后等待栅栏放行
			等其它的 goroutine 填补了氢原子的另一个空槽和氧原子的空槽之后，程序才可以继续进行
		氧原子的流水线：氧原子的流水线处理方法是 oxygen
			oxygen 方法是等待氧原子的空槽，然后输出一个 O，就等待栅栏放行
			放行后，释放氧原子空槽位
		栅栏
			在栅栏放行之前，只有两个氢原子的空槽位和一个氧原子的空槽位
			只有等栅栏放行之后，这些空槽位才会被释放
			栅栏放行，就意味着一个水分子组成成功
		为什么能保证是 HHO
			关键代码
				func NewH2O() *H2O {
					return &H2O{
						semaH: semaphore.NewWeighted(2), // 氢原子需要两个
						semaO: semaphore.NewWeighted(1), // 氧原子需要一个
						b:     cyclicbarrier.New(3),     // 需要三个原子才能合成
					}
				}
			信号量的定义资源数量的定义
				2 和 1 制约了 3 个原子的组成情况
				如果是 3 和 3，那么可以组成任意的原子组合
	代码：WaterFactory

总结
	用 WaitGroup 来实现这个水分子制造工厂的例子
		使用 WaitGroup 非常复杂，而且，重用和 Done 方法的调用有并发的问题，程序可能 panic
		远远没有使用循环栅栏 CyclicBarrier 更加简单直接
		示例：H2OWG
	建议
		了解一些并发原语，甚至是从其它编程语言、操作系统中学习更多的并发原语
		在面对并发场景的时候，你也能更加游刃有余

思考
	1.你觉得，SingleFlight 能不能合并并发的写操作呢？
	2.如果大自然的搬运工工厂生产的液体是双氧水（双氧水分子是两个氢原子和两个氧原子），你又该怎么实现呢？
*/

// H2OWG ==========一氧化二氢制造工厂 WaitGroup==========
type H2OWG struct {
	semaH *semaphore.Weighted
	semaO *semaphore.Weighted
	wg    sync.WaitGroup //将循环栅栏替换成WaitGroup
}

func NewH2OWG() *H2OWG {
	var wg sync.WaitGroup
	wg.Add(3)
	return &H2OWG{
		semaH: semaphore.NewWeighted(2),
		semaO: semaphore.NewWeighted(1),
		wg:    wg,
	}
}
func (h2o *H2OWG) hydrogenWG(releaseHydrogen func()) {
	h2o.semaH.Acquire(context.Background(), 1)
	releaseHydrogen()
	// 标记自己已达到，等待其它goroutine到达
	h2o.wg.Done()
	h2o.wg.Wait()
	h2o.semaH.Release(1)
}
func (h2o *H2OWG) oxygenWG(releaseOxygen func()) {
	h2o.semaO.Acquire(context.Background(), 1)
	releaseOxygen()
	// 标记自己已达到，等待其它goroutine到达
	h2o.wg.Done()
	h2o.wg.Wait()
	//都到达后重置wg
	h2o.wg.Add(3)
	h2o.semaO.Release(1)
}

// WaterFactory ==========一氧化二氢制造工厂==========
func WaterFactory() {
	//用来存放水分子结果的channel
	var ch chan string
	releaseHydrogen := func() {
		ch <- "H"
	}
	releaseOxygen := func() {
		ch <- "O"
	}
	// 300个原子，300个goroutine,每个goroutine并发的产生一个原子
	var N = 10
	ch = make(chan string, N*3)
	h2o := NewH2O()

	var wg sync.WaitGroup // 用来等待所有的goroutine完成
	wg.Add(N * 3)
	for i := 0; i < 2*N; i++ { // 200个氢原子goroutine
		go func() {
			//time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			h2o.hydrogen(releaseHydrogen)
			wg.Done()
		}()
	}

	for i := 0; i < 1*N; i++ { // 100个氧原子goroutine
		go func() {
			//time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			h2o.oxygen(releaseOxygen)
			wg.Done()
		}()
	}
	wg.Wait() //等待所有的goroutine执行完

	if len(ch) != N*3 { // 结果中肯定是300个原子
		log.Fatalf("expect %d atom but got %d", N*3, len(ch))
	}
	// 每三个原子一组，分别进行检查。要求这一组原子中必须包含两个氢原子和一个氧原子，这样才能生产水
	var s = make([]string, 3)
	for i := 0; i < N; i++ {
		s[0] = <-ch
		s[1] = <-ch
		s[2] = <-ch
		sort.Strings(s)
		water := s[0] + s[1] + s[2]
		fmt.Println(water)
		//if water != "HHO" {
		//	log.Fatalf("expect a water molecule but got %s", water)
		//}
	}
}

// H2O 定义水分子合成的辅助数据结构
type H2O struct {
	semaH *semaphore.Weighted         // 氢原子的信号量
	semaO *semaphore.Weighted         // 氧原子的信号量
	b     cyclicbarrier.CyclicBarrier // 循环栅栏，用来控制合成
}

func NewH2O() *H2O {
	return &H2O{
		semaH: semaphore.NewWeighted(3), //氢原子需要两个
		semaO: semaphore.NewWeighted(3), // 氧原子需要一个
		b:     cyclicbarrier.New(3),     // 需要三个原子才能合成
	}
}
func (h2o *H2O) hydrogen(releaseHydrogen func()) {
	h2o.semaH.Acquire(context.Background(), 1)
	releaseHydrogen()                 // 输出H
	h2o.b.Await(context.Background()) //等待栅栏放行
	h2o.semaH.Release(1)              // 释放氢原子空槽
}
func (h2o *H2O) oxygen(releaseOxygen func()) {
	h2o.semaO.Acquire(context.Background(), 1)
	releaseOxygen()                   // 输出O
	h2o.b.Await(context.Background()) //等待栅栏放行
	h2o.semaO.Release(1)              // 释放氢原子空槽
}

// CyclicBarrierSimpleExample ==========marusama/cyclicbarrier example==========
func CyclicBarrierSimpleExample() {
	// create a barrier for 10 parties with an action that increments counter
	// this action will be called each time when all goroutines reach the barrier
	cnt := 0
	b := cyclicbarrier.NewWithAction(10, func() error {
		cnt++
		return nil
	})

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ { // create 10 goroutines (the same count as barrier parties)
		wg.Add(1)
		go func() {
			for j := 0; j < 5; j++ {

				// do some hard work 5 times
				time.Sleep(100 * time.Millisecond)

				//fmt.Println(b.GetParties(), b.GetNumberWaiting(), cnt)
				err := b.Await(context.TODO()) // ..and wait for other parties on the barrier.
				// Last arrived goroutine will do the barrier action
				// and then pass all other goroutines to the next round
				if err != nil {
					panic(err)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println(cnt) // cnt = 5, it means that the barrier was passed 5 times
}
