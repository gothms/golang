package concurrent

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// 原子操作 + lock-free 实现信号量的尝试
type waiter struct {
	n int64
	//next  unsafe.Pointer
	ready chan<- struct{}
}
type Semaphore struct {
	size int64
	left atomic.Int64
	//waiters list.List // api
	waiters SemLKQueue
}

func StructVal() {
	s := Semaphore{waiters: *NewSemLKQueue()}
	w1 := waiter{n: 1}
	w2 := waiter{n: 10}
	s.waiters.SemEnqueue(w1)
	s.waiters.SemEnqueue(w2)

	next := s.waiters.SemFront()
	w := next.val.(waiter)
	w.n++
	s.waiters.Range(nil)
	fmt.Println(next.val.(waiter).n) // 1

	fmt.Println()
	//next.val.(waiter).n = 3
	fmt.Println(next == s.waiters.SemFront()) // true

	// 修改第一个结点
	new := next
	w.n = 0
	new.val = w
	//new.val.(waiter).n = 0
	//if s.left.CompareAndSwap(cur, cur-w.n) {
	//	s.waiters.SemDequeueNode(next)
	//}
	s.waiters.SemDequeueNode(next, new)
	s.waiters.Range(nil) // 0
}

func (s *Semaphore) SemAcquire() {

}
func (s *Semaphore) SemTryAcquire() {

}
func (s *Semaphore) SemRelease() {

}
func (s *Semaphore) notifyWaiters() {
	for {
		next := s.waiters.SemFront()
		if next == nil {
			break
		}
		w := next.val.(waiter)
		if cur := s.left.Load(); cur < w.n {
			break
		} else {
			new := next
			w.n = 0
			new.val = w
			// 参考 Mutex、WaitGroup、RWMutex、Pool 等的实现
			if s.waiters.SemDequeueNode(next, new) { // 出队和获取资源都需要同步，一个条件变量不能保护两个资源
				s.waiters.SemDequeue()
			}
		}
	}
}

type SemLKQueue struct {
	head unsafe.Pointer // 辅助头指针，头指针不包含有意义的数据，只是一个辅助的节点
	tail unsafe.Pointer
}

type semNode struct { // 通过链表实现，这个数据结构代表链表中的节点
	val  any // 队列中实际的元素 elem
	next unsafe.Pointer
}

func NewSemLKQueue() *SemLKQueue {
	n := unsafe.Pointer(&semNode{})
	return &SemLKQueue{head: n, tail: n}
}

// SemEnqueue 入队
// 入队的时候，通过 CAS 操作将一个元素添加到队尾，并且移动尾指针
func (q *SemLKQueue) SemEnqueue(v any) {
	n := &semNode{val: v}
	for {
		tail := semLoad(&q.tail)
		next := semLoad(&tail.next)
		if tail == semLoad(&q.tail) { // 尾还是尾
			if next == nil { // 还没有新数据入队
				if semCas(&tail.next, next, n) { //保证了：增加到队尾
					semCas(&q.tail, tail, n) //入队成功，移动尾巴指针
					return
				}
			} else { // 已有新数据加到队列后面，需要移动尾指针
				semCas(&q.tail, tail, next)
			}
		}
	}
}

// SemFront 队列中的第一个节点
func (q *SemLKQueue) SemFront() *semNode {
	for {
		head := semLoad(&q.head)
		next := semLoad(&head.next)
		if head == semLoad(&q.head) {
			return next
		}
	}
}

// SemDequeue 出队，没有元素则返回nil
// 出队的时候移除一个节点，并通过 CAS 操作移动 head 指针，同时在必要的时候移动尾指针
func (q *SemLKQueue) SemDequeue() any {
	for {
		head := semLoad(&q.head)
		tail := semLoad(&q.tail)
		next := semLoad(&head.next)
		if head == semLoad(&q.head) { // head还是那个head
			if head == tail { // head和tail一样
				if next == nil { // 说明是空队列
					return nil
				}
				semCas(&q.tail, tail, next) // 只是尾指针还没有调整，尝试调整它指向下一个
			} else {
				v := next.val                    // 读取出队的数据
				if semCas(&q.head, head, next) { // 既然要出队了，头指针移动到下一个
					return v // Dequeue is done. return
				}
			}
		}
	}
}
func (q *SemLKQueue) SemDequeueNode(old, new *semNode) bool {
	head := semLoad(&q.head)
	nextUP := unsafe.Pointer(semLoad(&head.next))
	if nextUP == nil {
		return false
	}
	return atomic.CompareAndSwapPointer(
		&nextUP, unsafe.Pointer(old), unsafe.Pointer(new))

	//for {
	//	head := semLoad(&q.head)
	//	next := semLoad(&head.next)
	//	if head == semLoad(&q.head) {
	//		if next != n {
	//			return nil
	//		} else {
	//			v := next.val                    // 读取出队的数据
	//			if semCas(&q.head, head, next) { // 既然要出队了，头指针移动到下一个
	//				return v // Dequeue is done. return
	//			}
	//		}
	//	}
	//}
}
func (q *SemLKQueue) Range(n *semNode) {
	if n == nil {
		n = (*semNode)(load(&q.head).next)
	}
	//for cur := (*node)(load(&q.tail).next); cur != nil; cur = (*node)(cur.next) {
	for cur := n; cur != nil; cur = (*semNode)(cur.next) {
		fmt.Print(cur.val, " ")
	}
}
func semLoad(p *unsafe.Pointer) (n *semNode) { // 将unsafe.Pointer原子加载转换成node
	return (*semNode)(atomic.LoadPointer(p))
}
func semCas(p *unsafe.Pointer, old, new *semNode) (ok bool) { // 封装CAS,避免直接将*node转换成unsafe.Pointer
	return atomic.CompareAndSwapPointer(
		p, unsafe.Pointer(old), unsafe.Pointer(new))
}
