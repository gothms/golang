package concurrent

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

/*
Mutex：骇客编程，如何拓展额外功能？

性能下降的罪魁祸首之一
	锁是性能下降的“罪魁祸首”之一，所以，有效地降低锁的竞争，就能够很好地提高性能
需求分析
	不希望锁的 goroutine 继续等待
		如果互斥锁被某个 goroutine 获取了，而且还没有释放，那么，其他请求这把锁的 goroutine 就会阻塞等待，直到有机会获得这把锁
		有时候阻塞并不是一个很好的主意，比如你请求锁更新一个计数器，如果获取不到锁的话没必要等待，大不了这次不更新，我下次更新就好了
		如果阻塞的话会导致业务处理能力的下降
	监控锁
		如果我们要监控锁的竞争情况，一个监控指标就是，等待这把锁的 goroutine 数量
		我们可以把这个指标推送到时间序列数据库中，再通过一些监控系统（比如 Grafana）展示出
		监控关键互斥锁上等待的 goroutine 的数量，是我们分析锁竞争的激烈程度的一个重要指标
解决方案
	不论是不希望锁的 goroutine 继续等待，还是想监控锁
	我们都可以基于标准库中 Mutex 的实现，通过 Hacker 的方式，为 Mutex 增加一些额外的功能

TryLock
	可以为 Mutex 添加一个 TryLock 的方法，也就是尝试获取排外锁
	TryLock 逻辑
		当一个 goroutine 调用这个 TryLock 方法请求锁的时候，如果这把锁没有被其他 goroutine 所持有，那么这个 goroutine 就持有了这把锁，并返回 true
		如果这把锁已经被其他 goroutine 所持有，或者是正在准备交给某个被唤醒的 goroutine，那么这个请求锁的 goroutine 就直接返回 false
		不会阻塞在方法调用上
		图示 04.mutex_04_trylock.jpg
	应用场景举例
		在实际开发中，如果要更新配置数据，我们通常需要加锁，这样可以避免同时有多个 goroutine 并发修改数据
		有的时候，我们也会使用 TryLock
		这样一来，当某个 goroutine 想要更改配置数据时，如果发现已经有 goroutine 在更改了，则它的 goroutine 调用 TryLock，返回了 false
		这个 goroutine 就会放弃更改
	Go
		很多语言（比如 Java）都为锁提供了 TryLock 的方法
		但 Go 官方 issue 6123 有一个讨论（后来一些 issue 中也提到过），标准库的 Mutex 不会添加 TryLock 方法
		通过 Go 的 Channel 我们也可以实现 TryLock 的功能，此处基于 Mutex 去实现（而且传统的同步原语也不容易出错）
	示例
		TryLock & TestTryLock
获取等待者的数量等指标
	获取 Mutex 的未暴漏字段
		state 这个字段的第一位是用来标记锁是否被持有
		第二位用来标记是否已经唤醒了一个等待者，第三位标记锁是否处于饥饿状态
		通过分析这个 state 字段我们就可以得到这些状态信息
	unsafe.Pointer
		当前持有和等待这把锁的 goroutine 的总数
			示例：WaiterCount & TestWaiterCount
		锁是否被持有、是否有等待者被唤醒、锁是否处于饥饿状态
			示例：IsLocked、IsWoken、IsStarving & TestMutexMessage
	在获取 state 字段的时候，并没有通过 Lock 获取这把锁，所以获取的这个 state 的值是一个瞬态的值
		可能在你解析出这个字段之后，锁的状态已经发生了变化
		不过没关系，因为你查看的就是调用的那一时刻的锁的状态
使用 Mutex 实现一个线程安全的队列
	Locker 与 数据结构
		Mutex 经常会和其他非线程安全（对于 Go，指的是 goroutine 安全）的数据结构一起，组合成一个线程安全的数据结构
		新数据结构的业务逻辑由原来的数据结构提供，而 Mutex 提供了锁的机制，来保证线程安全
		比如队列，可以通过 Slice 实现队列，但是通过 Slice 实现的队列不是线程安全的
		出队（Dequeue）和入队（Enqueue）会有 data race 的问题
	标准库中没有线程安全的队列数据结构的实现
		示例：SliceQueue

总结
	sync 基石
		Mutex 是 package sync 的基石，其他的一些同步原语也是基于它实现
	骇客编程
		通过 Hacker 的方式，拓展 Mutex 额外功能
		基于 Mutex 实现 TryLock，通过 unsafe 的方式读取到 Mutex 内部的 state 字段
		解决了：不希望锁的 goroutine 继续等待，以及监控锁

思考
	为 Mutex 获取锁时加上 Timeout 机制吗？会有什么问题吗？
*/

// SliceQueue ==========并发安全的队列==========
type SliceQueue struct {
	data []interface{}
	mu   sync.Mutex
}

func NewSliceQueue(n int) (q *SliceQueue) {
	return &SliceQueue{data: make([]interface{}, 0, n)}
}

// Enqueue 把值放在队尾
func (q *SliceQueue) Enqueue(v interface{}) {
	q.mu.Lock()
	q.data = append(q.data, v)
	q.mu.Unlock()
}

// Dequeue 移去队头并返回
func (q *SliceQueue) Dequeue() interface{} {
	q.mu.Lock()
	//defer q.mu.Unlock()
	if len(q.data) == 0 {
		q.mu.Unlock()
		return nil
	}
	v := q.data[0]
	q.data = q.data[1:]
	q.mu.Unlock()
	return v
}

// ==========获取等待者的数量等指标==========
const (
	mutexLocked      = 1 << iota // 加锁标示位置
	mutexWoken                   // 唤醒标示位置
	mutexStarving                // 锁饥饿标示
	mutexWaiterShift = iota      // 标示Waiter的起始bit位置
)

type Mutex struct {
	sync.Mutex
}

// WaiterCount 获取等待者的数量等指标
func (m *Mutex) WaiterCount() int {
	state := unsafe.Pointer(&m.Mutex)
	sema := uintptr(state) + unsafe.Sizeof(state)
	stateVal := (*int32)(state)
	semaVal := (*uint32)(unsafe.Pointer(sema))
	fmt.Println(*stateVal, *stateVal>>mutexWaiterShift, *semaVal)

	// 获取state字段的值
	v := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	v = v >> mutexWaiterShift //得到等待者的数值
	v = v + (v & mutexLocked) //再加上锁持有者的数量，0或者1
	return int(v)             // 当前持有和等待这把锁的 goroutine 的总数
}

// IsLocked 锁是否被持有
func (m *Mutex) IsLocked() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexLocked == mutexLocked
}

// IsWoken 是否有等待者被唤醒
func (m *Mutex) IsWoken() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexWoken == mutexWoken
}

// IsStarving 锁是否处于饥饿状态
func (m *Mutex) IsStarving() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexStarving == mutexStarving
}

// TryLock 尝试获取锁
// ==========获取等待者的数量等指标==========
func (m *Mutex) TryLock() bool {
	fmt.Println(m.Mutex, *(*int32)(unsafe.Pointer(&m.Mutex)))

	// fast path 如果一开始就没有其他g争夺，那么直接获取锁
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked) {
		return true
	}
	// 如果处于唤醒，加锁或者饥饿状态，这次请求就不参与了竞争了，直接返回false
	old := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	if old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
		return false
	}
	// 尝试在竞争的状态下请求锁
	n := old | mutexLocked
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), old, n)
}
