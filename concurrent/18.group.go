package concurrent

import (
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"
)

/*
分组操作：处理一组子任务，该用什么并发原语？

分组编排
	共享资源保护、任务编排和消息传递是 Go 并发编程中常见的场景
	而分组执行一批相同的或类似的任务则是任务编排中一类情形
	分组编排的一些常用场景和并发原语，包括 ErrGroup、gollback、Hunch 和 schedgroup

ErrGroup
	golang.org/x/sync/errgroup 包
		ErrGroup 是 Go 官方提供的一个同步扩展库，非常常用的并发原语
	简介
		我们经常会碰到需要将一个通用的父任务拆成几个小任务并发执行的场景
		其实，将一个大的任务拆成几个小任务并发执行，可以有效地提高程序的并发度
		就像你在厨房做饭一样，你可以在蒸米饭的同时炒几个小菜，米饭蒸好了，菜同时也做好了，很快就能吃到可口的饭菜
	ErrGroup 就是用来应对这种场景的。它和 WaitGroup 有些类似，但是它提供功能更加丰富
		和 Context 集成
		error 向上传播，可以把子任务的错误传递给 Wait 的调用者
基本用法
	ErrGroup 并发原语，底层也是基于 WaitGroup 实现的
		三个方法，分别是 WithContext、Go 和 Wait
	1.WithContext
		创建一个 Group 对象
			func WithContext(ctx context.Context) (*Group, context.Context)
		返回一个 Group 实例，同时还会返回一个使用 context.WithCancel(ctx) 生成的新 Context
			一旦有一个子任务返回错误，或者是 Wait 调用返回，这个新 Context 就会被 cancel
			Group 的零值也是合法的，只不过，你就没有一个可以监控是否 cancel 的 Context 了
		注意，如果传递给 WithContext 的 ctx 参数，是一个可以 cancel 的 Context 的话
			那么，它被 cancel 的时候，并不会终止正在执行的子任务
	2.Go
		执行子任务
			func (g *Group) Go(f func() error)
		传入的子任务函数 f 是类型为 func() error 的函数
			如果任务执行成功，就返回 nil，否则就返回 error，并且会 cancel 那个新的 Context
			一个任务可以分成好多个子任务，而且，可能有多个子任务执行失败返回 error，不过，Wait 方法只会返回第一个错误
	3.Wait
		等所有的子任务都完成后，它才会返回，否则只会阻塞等待
			func (g *Group) Wait() error
		如果有多个子任务返回错误，它只会返回第一个出现的错误，如果所有的子任务都执行成功，就返回 nil
ErrGroup 使用例子
简单例子：返回第一个错误
	示例：ErrGroupDemo01
		子任务 2 会返回执行失败，其它两个执行成功
		在三个子任务都执行后，group.Wait 才会返回第 2 个子任务的错误
	结果
		failed: failed to exec #2
更进一步，返回所有子任务的错误
	示例：ErrGroupDemo02
		使用一个 result slice 保存子任务的执行结果，这样，通过查询 result，就可以知道每一个子任务的结果了
		不仅可以使用 result 记录 error 信息，还可以用它记录计算结果
	结果
		failed: [<nil> failed to exec #2 <nil>]
任务执行流水线 Pipeline
	Go 官方文档中提供 pipeline 例子
		由一个子任务遍历文件夹下的文件，然后把遍历出的文件交给 20 个 goroutine，让这些 goroutine 并行计算文件的 md5
	示例
		源码
			errgroup/errgroup_example_md5all_test.go
			concurrent/18.errgroup_pipeline.go
		TestExampleGroup_pipeline
			多阶段 pipeline 的实现（例子是遍历文件夹和计算 md5 两个阶段）
			控制执行子任务的 goroutine 数量
	应用
		很多公司都在使用 ErrGroup 处理并发子任务，比如 Facebook、bilibili 等公司的一些项目
		但是，这些公司在使用的时候，发现了一些不方便的地方，或者说，官方的 ErrGroup 的功能还不够丰富
		所以，他们都对 ErrGroup 进行了扩展
扩展库
	需求分析
		如果我们无限制地直接调用 ErrGroup 的 Go 方法，就可能会创建出非常多的 goroutine
		太多的 goroutine 会带来调度和 GC 的压力，而且也会占用更多的内存资源
		就像 go#34457 指出的那样，当前 Go 运行时创建的 g 对象只会增长和重用，不会回收
		所以在高并发的情况下，也要尽可能减少 goroutine 的使用
		link：go#34457
	解决方案
		常用的一个手段就是使用 worker pool(goroutine pool)
		或者是类似 containerd/stargz-snapshotter 的方案，使用信号量，信号量的资源的数量就是可以并行的 goroutine 的数量
		但是在这介绍一些其它的手段，比如下面介绍的 bilibili 实现的 errgroup
		link：containerd/stargz-snapshotter
	bilibili/errgroup
		简介
			bilibili 实现了一个扩展的 ErrGroup，可以使用一个固定数量的 goroutine 处理子任务
			如果不设置 goroutine 的数量，那么每个子任务都会比较“放肆地”创建一个 goroutine 并发执行
			link：bilibili/errgroup
			bilibili 官方文档已经很详细地介绍了它的几个扩展功能
		除了可以控制并发 goroutine 的数量，它还提供了 2 个功能
			1. cancel，失败的子任务可以 cancel 所有正在执行任务
			2. recover，而且会把 panic 的堆栈信息放到 error 中，避免子任务 panic 导致的程序崩溃
		缺点
			可优化
				但是，有一点不太好的地方就是，一旦你设置了并发数，超过并发数的子任务需要等到调用者调用 Wait 之后才会执行
				而不是只要 goroutine 空闲下来，就去执行
				如果不注意这一点的话，可能会出现子任务不能及时处理的情况，这是这个库可以优化的一点
			并发问题
				在高并发的情况下，如果任务数大于设定的 goroutine 的数量，并且这些任务被集中加入到 Group 中
				这个库的处理方式是把子任务加入到一个数组中，但是，这个数组不是线程安全的，有并发问题
			图示：g.chs = append(g.chs, f)
				18.group_bili_con.jpg
		示例
			死锁
				运行这个程序的话，你就会发现死锁问题
				因为我们的测试程序是一个简单的命令行工具，程序退出的时候，Go runtime 能检测到死锁问题
				如果是一直运行的服务器程序，死锁问题有可能是检测不出来的，程序一直会 hang 在 Wait 的调用上
			代码
				func BiliDeadlock() {
					var g errgroup.Group
					g.GOMAXPROCS(1) // 只使用一个goroutine处理子任务
					var count int64
					g.Go(func(ctx context.Context) error {
						time.Sleep(time.Second) //睡眠5秒，把这个goroutine占住
						return nil
					})
					total := 10000
					for i := 0; i < total; i++ { // 并发一万个goroutine执行子任务，理论上这些子任务都...
						go func() {
							g.Go(func(ctx context.Context) error {
								atomic.AddInt64(&count, 1)
								return nil
							})
						}()
					}
					// 等待所有的子任务完成。理论上10001个子任务都会被完成
					if err := g.Wait(); err != nil {
						panic(err)
					}
					got := atomic.LoadInt64(&count)
					if got != int64(total) {
						panic(fmt.Sprintf("expect %d but got %d", total, got))
					}
				}
	neilotoole/errgroup
		简介
			它可以直接替换官方的 ErrGroup，方法都一样，原有功能也一样，只不过增加了可以控制并发 goroutine 的功能
			link：neilotoole/errgroup
		API
			type Group
				func WithContext(ctx context.Context) (*Group, context.Context)
				func WithContextN(ctx context.Context, numG, qSize int) (*Group, context.Context)
				func (g *Group) Go(f func() error)
				func (g *Group) Wait() error
		WithContextN
			可以设置并发的 goroutine 数，以及等待处理的子任务队列的大小
			当队列满的时候，如果调用 Go 方法，就会被阻塞，直到子任务可以放入到队列中才返回
			如果你传给这两个参数的值不是正整数，它就会使用 runtime.NumCPU 代替你传入的参数
			当然，你也可以把 bilibili 的 recover 功能扩展到这个库中，以避免子任务的 panic 导致程序崩溃
	facebookgo/errgroup
		基于 WaitGroup
			Facebook 提供的这个 ErrGroup，其实并不是对 Go 扩展库 ErrGroup 的扩展，而是对标准库 WaitGroup 的扩展
			不过，因为它们的名字一样，处理的场景也类似
		API
			type Group
				func (g *Group) Add(delta int)
				func (g *Group) Done()
				func (g *Group) Error(e error)
				func (g *Group) Wait() error
		Wait 方法
			标准库的 WaitGroup 只提供了 Add、Done、Wait 方法，而且 Wait 方法也没有返回子 goroutine 的 error
			而 Facebook 提供的 ErrGroup 提供的 Wait 方法可以返回 error，而且可以包含多个 error
			子任务在调用 Done 之前，可以把自己的 error 信息设置给 ErrGroup
			接着，Wait 在返回的时候，就会把这些 error 信息返回给调用者
		示例
			Error 方法
				func main() {
					var g errgroup.Group
					g.Add(3)
					// 启动第一个子任务,它执行成功
					go func() {
						time.Sleep(5 * time.Second)
						fmt.Println("exec #1")
						g.Done()
					}()
					// 启动第二个子任务，它执行失败
					go func() {
						time.Sleep(10 * time.Second)
						fmt.Println("exec #2")
						g.Error(errors.New("failed to exec #2"))
						g.Done()
					}()
					// 启动第三个子任务，它执行成功
					go func() {
						time.Sleep(15 * time.Second)
						fmt.Println("exec #3")
						g.Done()
					}()
					// 等待所有的goroutine完成，并检查error
					if err := g.Wait(); err == nil {
						fmt.Println("Successfully exec all")
					} else {
						fmt.Println("failed:", err)
					}
				}
			g.Error(errors.New("failed to exec #2"))
				设置 error 信息
				会把这个 error 信息输出出来

其它实用的 Group 并发原语
	几种有趣而实用的 Group 并发原语
		这些并发原语都是控制一组子 goroutine 执行的面向特定场景的并发原语，当你遇见这些特定场景时，就可以参考这些库
	SizedGroup/ErrSizedGroup
	gollback
	Hunch
	schedgroup
SizedGroup/ErrSizedGroup
	go-pkgz/syncs
		提供了两个 Group 并发原语，分别是 SizedGroup 和 ErrSizedGroup
		link：go-pkgz/syncs
	SizedGroup 原理
		SizedGroup 内部是使用信号量和 WaitGroup 实现的
		它通过信号量控制并发的 goroutine 数量，或者是不控制 goroutine 数量，只控制子任务并发执行时候的数量（通过）
	简介
		默认情况下，SizedGroup 控制的是子任务的并发数量，而不是 goroutine 的数量
		在这种方式下，每次调用 Go 方法都不会被阻塞，而是新建一个 goroutine 去执行
		如果想控制 goroutine 的数量，你可以使用 syncs.Preemptive 设置这个并发原语的可选项
		如果设置了这个可选项，但在调用 Go 方法的时候没有可用的 goroutine，那么调用者就会等待，直到有 goroutine 可以处理这个子任务才返回
		这个控制在内部是使用信号量实现的
	示例
		func SizedGroup() {
			// 设置goroutine数是10
			swg := syncs.NewSizedGroup(10)
			// swg := syncs.NewSizedGroup(10, syncs.Preemptive)
			var c uint32
			// 执行1000个子任务，只会有10个goroutine去执行
			for i := 0; i < 1000; i++ {
				swg.Go(func(ctx context.Context) {
					time.Sleep(5 * time.Millisecond)
					atomic.AddUint32(&c, 1)
				})
			}
			// 等待任务完成
			swg.Wait()
			// 输出结果
			fmt.Println(c)
		}
	ErrSizedGroup
		ErrSizedGroup 为 SizedGroup 提供了 error 处理的功能
		它的功能和 Go 官方扩展库的功能一样，就是等待子任务完成并返回第一个出现的 error
	不过，它还提供了额外的功能
		第一个额外的功能，就是可以控制并发的 goroutine 数量，这和 SizedGroup 的功能一样
		第二个功能是，如果设置了 termOnError，子任务出现第一个错误的时候会 cancel Context
			而且后续的 Go 调用会直接返回，Wait 调用者会得到这个错误，这相当于是遇到错误快速返回
			如果没有设置 termOnError，Wait 会返回所有的子任务的错误
		不过，ErrSizedGroup 和 SizedGroup 设计得不太一致的地方是
			SizedGroup 可以把 Context 传递给子任务，这样可以通过 cancel 让子任务中断执行，但是 ErrSizedGroup 却没有实现
	总体来说
		syncs 包提供的并发原语的质量和功能还是非常赞的。不过，目前的 star 只有十几个，这和它的功能严重不匹配
gollback
	简介
		gollback也是用来处理一组子任务的执行的，不过它解决了 ErrGroup 收集子任务返回结果的痛点
		使用 ErrGroup 时，如果你要收到子任务的结果和错误，你需要定义额外的变量收集执行结果和错误，但是这个库可以提供更便利的方式
		link：gollback
		示例 ErrGroupDemo02 中，如果想得到每一个子任务的结果或者 error，我们需要额外提供一个 result slice 进行收集
		使用 gollback 的话，就不需要这些额外的处理了，因为它的方法会把结果和 error 信息都返回
	All 方法
		func All(ctx context.Context, fns ...AsyncFunc) ([]interface{}, []error)
			它会等待所有的异步函数（AsyncFunc）都执行完才返回，而且返回结果的顺序和传入的函数的顺序保持一致
			第一个返回参数是子任务的执行结果，第二个参数是子任务执行时的错误信息
		异步函数的定义
			type AsyncFunc func(ctx context.Context) (interface{}, error)
			ctx 会被传递给子任务。如果你 cancel 这个 ctx，可以取消子任务
		示例
			func gollbackAll() {
				rs, errs := gollback.All( // 调用All方法
					context.Background(),
					func(ctx context.Context) (interface{}, error) {
						time.Sleep(3 * time.Second)
						return 1, nil // 第一个任务没有错误，返回1
					},
					func(ctx context.Context) (interface{}, error) {
						return nil, errors.New("failed") // 第二个任务返回一个错误
					},
					func(ctx context.Context) (interface{}, error) {
						return 3, nil // 第三个任务没有错误，返回3
					},
				)
				fmt.Println(rs) // 输出子任务的结果
				fmt.Println(errs) // 输出子任务的错误信息
			}
	Race 方法
		跟 All 方法类似
			只不过，在使用 Race 方法的时候，只要一个异步函数执行没有错误，就立马返回，而不会返回所有的子任务信息
			如果所有的子任务都没有成功，就会返回最后一个 error 信息
		func Race(ctx context.Context, fns ...AsyncFunc) (interface{}, error)
			如果有一个正常的子任务的结果返回，Race 会把传入到其它子任务的 Context cancel 掉，这样子任务就可以中断自己的执行
	Retry 方法
		Retry 不是执行一组子任务，而是执行一个子任务
			如果子任务执行失败，它会尝试一定的次数，如果一直不成功，就会返回失败错误，如果执行成功，它会立即返回
			如果 retires 等于 0，它会永远尝试，直到成功
		示例
			func gollbackRetry() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				// 尝试5次，或者超时返回
				res, err := gollback.Retry(ctx, 5, func(ctx context.Context) (interface{}, error) {
					return nil, errors.New("failed")
				})
				fmt.Println(res) // 输出结果
				fmt.Println(err) // 输出错误信息
			}
Hunch
	Hunch提供的功能和 gollback 类似，不过它提供的方法更多，而且它提供的和 gollback 相应的方法，也有一些不同
		它定义了执行子任务的函数，这和 gollback 的 AyncFunc 是一样的
		type Executable func(context.Context) (interface{}, error)
	All 方法
		func All(parentCtx context.Context, execs ...Executable) ([]interface{}, error)
		它会传入一组可执行的函数（子任务），返回子任务的执行结果
		和 gollback 的 All 方法不一样的是，一旦一个子任务出现错误，它就会返回错误信息，执行结果（第一个返回参数）为 nil
	Take 方法
		func Take(parentCtx context.Context, num int, execs ...Executable) ([]interface{}, error)
		可以指定 num 参数，只要有 num 个子任务正常执行完没有错误，这个方法就会返回这几个子任务的结果
		一旦一个子任务出现错误，它就会返回错误信息，执行结果（第一个返回参数）为 nil
	Last 方法
		func Last(parentCtx context.Context, num int, execs ...Executable) ([]interface{}, error)
		它只返回最后 num 个正常执行的、没有错误的子任务的结果。一旦一个子任务出现错误，它就会返回错误信息，执行结果（第一个返回参数）为 nil
		比如 num 等于 1，那么，它只会返回最后一个无错的子任务的结果
	Retry 方法
		func Retry(parentCtx context.Context, retries int, fn Executable) (interface{}, error)
		它的功能和 gollback 的 Retry 方法的功能一样，如果子任务执行出错，就会不断尝试，直到成功或者是达到重试上限
		如果达到重试上限，就会返回错误。如果 retries 等于 0，它会不断尝试
	Waterfall 方法
		func Waterfall(parentCtx context.Context, execs ...ExecutableInSequence) (interface{}, error)
		它其实是一个 pipeline 的处理方式，所有的子任务都是串行执行的，前一个子任务的执行结果会被当作参数传给下一个子任务，直到所有的任务都完成，返回最后的执行结果
		一旦一个子任务出现错误，它就会返回错误信息，执行结果（第一个返回参数）为 nil
	小结
		gollback 和 Hunch 是属于同一类的并发原语，对一组子任务的执行结果，可以选择一个结果或者多个结果
		这也是现在热门的微服务常用的服务治理的方法
schedgroup
	和时间相关的处理一组 goroutine 的并发原语 schedgroup
	简介
		schedgroup 是 Matt Layher 开发的 worker pool，可以指定任务在某个时间或者某个时间之后执行
		Matt Layher 也是一个知名的 Gopher，经常在一些会议上分享一些他的 Go 开发经验
		他在 GopherCon Europe 2020 大会上专门介绍了这个并发原语：schedgroup: a timer-based goroutine concurrency primitive
		link：schedgroup
		link：schedgroup: a timer-based goroutine concurrency primitive
	API
		type Group
			func New(ctx context.Context) *Group
			func (g *Group) Delay(delay time.Duration, fn func())
			func (g *Group) Schedule(when time.Time, fn func())
			func (g *Group) Wait() error
	Delay 方法和 Schedule 方法
		它们的功能其实是一样的，都是用来指定在某个时间或者之后执行一个函数
		只不过，Delay 传入的是一个 time.Duration 参数，它会在 time.Now()+delay 之后执行函数，而 Schedule 可以指定明确的某个时间执行
	Wait 方法
		这个方法调用会阻塞调用者，直到之前安排的所有子任务都执行完才返回
		如果 Context 被取消，那么，Wait 方法会返回这个 cancel error
	使用 Wait 方法的注意点
		第一点是，如果调用了 Wait 方法，你就不能再调用它的 Delay 和 Schedule 方法，否则会 panic
		第二点是，Wait 方法只能调用一次，如果多次调用的话，就会 panic
	timer / container/heap
		如果只有几个子任务，使用 timer 不是问题
		但一旦有大量的子任务，而且还要能够 cancel，那么，使用 timer 的话，CPU 资源消耗就比较大了
		所以，schedgroup 在实现的时候，就使用 container/heap，按照子任务的执行时间进行排序
		这样可以避免使用大量的 timer，从而提高性能
	示例
		sg := schedgroup.New(context.Background())
		// 设置子任务分别在100、200、300之后执行
		for i := 0; i < 3; i++ {
			n := i + 1
			sg.Delay(time.Duration(n)*100*time.Millisecond, func() {
				log.Println(n) //输出任务编号
			})
		}
		// 等待所有的子任务都完成
		if err := sg.Wait(); err != nil {
			log.Fatalf("failed to wait: %v", err)
		}

总结
	几种常见的处理一组子任务的并发原语
		包括 ErrGroup、gollback、Hunch、schedgroup，等
		遇到相同的业务场景时，可以考虑使用这些并发原语
	新的并发原语不断出现
		如 go-waitgroup，link：
		了解这些并发原语，构造新的并发原语来处理应对你的特有场景，实现代码重用和业务逻辑简化

思考
	官方扩展库 ErrGroup 没有实现可以取消子任务的功能，请你课下可以自己去实现一个子任务可取消的 ErrGroup
*/

// ErrGroupDemo02 ==========ErrGroup Demo 02==========
func ErrGroupDemo02() {
	var g errgroup.Group
	var result = make([]error, 3)
	// 启动第一个子任务,它执行成功
	g.Go(func() error {
		time.Sleep(3 * time.Second)
		fmt.Println("exec #1")
		result[0] = nil // 保存成功或者失败的结果
		return nil
	})
	// 启动第二个子任务，它执行失败
	g.Go(func() error {
		time.Sleep(1 * time.Second)
		fmt.Println("exec #2")
		result[1] = errors.New("failed to exec #2") // 保存成功或者失败的结果
		return result[1]
	})
	// 启动第三个子任务，它执行成功
	g.Go(func() error {
		time.Sleep(2 * time.Second)
		fmt.Println("exec #3")
		result[2] = nil // 保存成功或者失败的结果
		return nil
	})
	if err := g.Wait(); err == nil {
		fmt.Printf("Successfully exec all. result: %v\n", result)
	} else {
		fmt.Printf("failed: %v\n", result)
	}
}

// ErrGroupDemo01 ==========ErrGroup Demo 01==========
func ErrGroupDemo01() {
	var g errgroup.Group
	g.Go(func() error { // 启动第一个子任务,它执行成功
		time.Sleep(3 * time.Second)
		fmt.Println("exec #1")
		return nil
	})
	g.Go(func() error { // 启动第二个子任务，它执行失败
		time.Sleep(1 * time.Second)
		fmt.Println("exec #2")
		return errors.New("failed to exec #2")
	})
	g.Go(func() error { // 启动第三个子任务，它执行成功
		time.Sleep(2 * time.Second)
		fmt.Println("exec #3")
		return nil
	})
	if err := g.Wait(); err == nil { // 等待三个任务都完成
		fmt.Println("Successfully exec all")
	} else {
		fmt.Println("failed:", err)
	}
}
