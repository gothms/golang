package concurrent

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

/*
sync.Pool 对象获取

	1.尝试从私有对象获取
		每个 Processor 有私有对象、共享池，私有对象协程安全，共享池协程不安全
	2.私有对象不存在，尝试从当前 Processor 的共享池获取
	3.如果当前 Processor 共享池也是空的，那么就尝试去其他 Processor 的共享池获取
	4.如果所有子池都是空的，最后就使用用户指定的 New 函数产生一个新的对象返回

Processor 请参考 goroutine_test.go，每个 Processor 都分为两个部分：

	私有对象：只能缓存一个对象，协程安全
	共享池：协程不安全，需要锁

sync.Pool 对象放回

	1.如果私有对象不存在，则保存为私有对象
	2.如果私有对象存在，放入当前 Processor 子池的共享池中

sync.Pool 对象的生命周期

	1.GC 会清除 sync.Pool 缓存的对象
	2.对象的缓存有效期为下一次 GC 之前

	正因如此，sync.Pool 对象的生命周期不可控，sync.Pool 不是 对象池

Get()

	对象已经从对象池里的取出，对象池里就没这个对象了

sync.Pool 总结

	1.适用于通过复用，降低复杂对象的创建和 GC 代价
	2.协程安全，会有锁的开销
	3.生命周期受 GC 影响，不适合于做连接池等，需要自己管理生命周期的资源的池化

使用 sync.Pool 的考虑：能否优化程序，取决于

	锁的开销
	创建对象的开销
	两个开销，哪个更大
*/
var i int

func TestSyncPool(t *testing.T) {
	pool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Create a new obj.")
			i++
			return 100 + i
		},
	}
	v := pool.Get().(int)
	fmt.Println(v)
	pool.Put(3)
	runtime.GC() // GC 会清除 sync.Pool 中缓存的对象
	time.Sleep(time.Second * 2)
	v1, _ := pool.Get().(int)
	fmt.Println(v1) // GC 后，v1=100
	v2, _ := pool.Get().(int)
	fmt.Println(v2) // v1 已经被取出，v2会被创建并取出
}

func TestSyncPoolInMultiGoroutine(t *testing.T) {
	pool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Create a new obj.")
			return 10
		},
	}
	pool.Put(101)
	pool.Put(102)
	pool.Put(103)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			//t.Log(pool.Get())
			fmt.Println(pool.Get())
			wg.Done()
		}(i)
	}
	wg.Wait()
}
