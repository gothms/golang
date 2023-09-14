// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package semaphore provides a weighted semaphore implementation.

package source // import "golang.org/x/sync/semaphore"

import (
	"container/list"
	"context"
	"sync"
)

type waiter struct {
	n     int64
	ready chan<- struct{} // Closed when semaphore acquired.
}

// NewWeighted creates a new weighted semaphore with the given
// maximum combined weight for concurrent access.
func NewWeighted(n int64) *Weighted {
	w := &Weighted{size: n}
	return w
}

// Weighted provides a way to bound concurrent access to a resource.
// The callers can request access with a given weight.
type Weighted struct {
	size    int64      // 最大资源数
	cur     int64      // 当前已被使用的资源
	mu      sync.Mutex // 互斥锁，对字段的保护
	waiters list.List  // 等待队列
}

// Acquire acquires the semaphore with a weight of n, blocking until resources
// are available or ctx is done. On success, returns nil. On failure, returns
// ctx.Err() and leaves the semaphore unchanged.
//
// If ctx is already done, Acquire may still succeed without blocking.
// 它不仅仅要监控资源是否可用，而且还要检测 Context 的 Done 是否已关闭
func (s *Weighted) Acquire(ctx context.Context, n int64) error {
	s.mu.Lock()
	if s.size-s.cur >= n && s.waiters.Len() == 0 { // fast path, 如果有足够的资源，都不考虑 ctx.Done 的状态，将 cur 加上 n 就返回
		s.cur += n
		s.mu.Unlock()
		return nil // 新来的 go 有优先权，即使保证了一定的公平性
	}

	if n > s.size { // 如果是不可能完成的任务，请求的资源数大于能提供的最大的资源数
		// Don't make other Acquire calls block on one that's doomed to fail.
		s.mu.Unlock()
		<-ctx.Done() // 依赖 ctx 的状态返回，否则一直等待
		return ctx.Err()
	}
	// 否则就需要把调用者加入到等待队列中
	ready := make(chan struct{})    // 创建了一个 ready chan，以便被通知唤醒
	w := waiter{n: n, ready: ready} // 请求资源的实体对象
	elem := s.waiters.PushBack(w)   // 并加入等待队列中
	s.mu.Unlock()
	// 等待
	select {
	case <-ctx.Done(): // context 的 Done 被关闭
		err := ctx.Err()
		s.mu.Lock()
		select {
		case <-ready: // 如果被唤醒了，忽略 ctx 的状态
			// Acquired the semaphore after we were canceled.  Rather than trying to
			// fix up the queue, just pretend we didn't notice the cancelation.
			err = nil
		default: // 通知 waiter
			isFront := s.waiters.Front() == elem
			s.waiters.Remove(elem)
			// If we're at the front and there're extra tokens left, notify other waiters.
			if isFront && s.size > s.cur { // 通知其它的 waiters，检查是否有足够的资源
				s.notifyWaiters() // 逐个检查等待的调用者，尝试给它们分配资源
			}
		}
		s.mu.Unlock()
		return err

	case <-ready: // 被唤醒了，且资源已分配给它
		return nil
	}
}

// TryAcquire acquires the semaphore with a weight of n without blocking.
// On success, returns true. On failure, returns false and leaves the semaphore unchanged.
func (s *Weighted) TryAcquire(n int64) bool {
	s.mu.Lock()
	success := s.size-s.cur >= n && s.waiters.Len() == 0 // 有足够资源，且 waiter 队列为空，保证公平性
	if success {
		s.cur += n // 获取资源
	}
	s.mu.Unlock()
	return success
}

// Release releases the semaphore with a weight of n.
func (s *Weighted) Release(n int64) { // 将当前计数值减去释放的资源数 n，并唤醒等待队列中的调用者，看是否有足够的资源被获取
	s.mu.Lock()
	s.cur -= n // 将当前计数值减去释放的资源数 n
	if s.cur < 0 {
		s.mu.Unlock()
		panic("semaphore: released more than held")
	}
	s.notifyWaiters() // 逐个检查等待的调用者，尝试给它们分配资源
	s.mu.Unlock()
}

func (s *Weighted) notifyWaiters() { // 逐个检查等待的调用者，如果资源不够，或者是没有等待者了，就返回
	for { // 逐个检查等待的调用者，尝试给它们分配资源
		next := s.waiters.Front() // 依次检查队首的 waiter
		if next == nil {
			break // No more waiters blocked.
		}

		w := next.Value.(waiter)
		if s.size-s.cur < w.n { // 如果资源不够，则终止检查
			// Not enough tokens for the next waiter.  We could keep going (to try to
			// find a waiter with a smaller request), but under load that could cause
			// starvation for large requests; instead, we leave all remaining waiters
			// blocked.
			//
			// Consider a semaphore used as a read-write lock, with N tokens, N
			// readers, and one writer.  Each reader can Acquire(1) to obtain a read
			// lock.  The writer can Acquire(N) to obtain a write lock, excluding all
			// of the readers.  If we allow the readers to jump ahead in the queue,
			// the writer will starve — there is always one token available for every
			// reader.
			// 当释放 100 个资源的时候，如果第一个等待者需要 101 个资源，那么，队列中的所有等待者都会继续等待，即使有的等待者只需要 1 个资源
			// 这样做的目的是避免饥饿，否则的话，资源可能总是被那些请求资源数小的调用者获取
			// 这样一来，请求资源数巨大的调用者，就没有机会获得资源了
			break //避免饥饿，这里还是按照先入先出的方式处理
		}

		s.cur += w.n           // 分配资源
		s.waiters.Remove(next) // 并移除队首 waiter
		close(w.ready)         // 别忘记 close
	}
}
