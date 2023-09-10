// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package source

import (
	"internal/race"
	"sync/atomic"
	"unsafe"
)

// A WaitGroup waits for a collection of goroutines to finish.
// The main goroutine calls Add to set the number of
// goroutines to wait for. Then each of the goroutines
// runs and calls Done when finished. At the same time,
// Wait can be used to block until all goroutines have finished.
//
// A WaitGroup must not be copied after first use.
//
// In the terminology of the Go memory model, a call to Done
// “synchronizes before” the return of any Wait call that it unblocks.
type WaitGroup struct {
	noCopy noCopy // 避免复制使用的一个技巧，可以告诉vet工具违反了复制使用的规则
	// 64bit(8bytes) 的值分成两段，高 32bit 是计数值，低 32bit 是 waiter 的计数
	state atomic.Uint64 // high 32 bits are counter, low 32 bits are waiter count.
	sema  uint32        // 信号量
	// 因为64bit值的原子操作需要64bit对齐，但是32bit编译器不支持，所以数组中的元素在不同的架构...
	// 总之，会找到对齐的那64bit作为state，其余的32bit做信号量
	// state1 [3]uint32
}

// Add adds delta, which may be negative, to the WaitGroup counter.
// If the counter becomes zero, all goroutines blocked on Wait are released.
// If the counter goes negative, Add panics.
//
// Note that calls with a positive delta that occur when the counter is zero
// must happen before a Wait. Calls with a negative delta, or calls with a
// positive delta that start when the counter is greater than zero, may happen
// at any time.
// Typically this means the calls to Add should execute before the statement
// creating the goroutine or other event to be waited for.
// If a WaitGroup is reused to wait for several independent sets of events,
// new Add calls must happen after all previous Wait calls have returned.
// See the WaitGroup example.
func (wg *WaitGroup) Add(delta int) {
	if race.Enabled {
		if delta < 0 {
			// Synchronize decrements with Wait.
			race.ReleaseMerge(unsafe.Pointer(wg))
		}
		race.Disable()
		defer race.Enable()
	}
	state := wg.state.Add(uint64(delta) << 32) // 高32bit是计数值v，所以把delta左移32，增加到计数上
	v := int32(state >> 32)                    // 当前计数值
	w := uint32(state)                         // waiter count
	if race.Enabled && delta > 0 && v == int32(delta) {
		// The first increment must be synchronized with Wait.
		// Need to model this as a read, because there can be
		// several concurrent wg.counter transitions from 0.
		race.Read(unsafe.Pointer(&wg.sema))
	}
	if v < 0 { // 计数器设置为负值
		panic("sync: negative WaitGroup counter")
	}
	if w != 0 && delta > 0 && v == int32(delta) { // 先 Wait，后 Add，则 panic（但是 Wait 有 'if v == 0' 的判断）
		panic("sync: WaitGroup misuse: Add called concurrently with Wait") // 没测试到
	}
	if v > 0 || w == 0 { // 否则 v == 0 && w > 0
		return
	}
	// This goroutine has set counter to 0 when waiters > 0.
	// Now there can't be concurrent mutations of state:
	// - Adds must not happen concurrently with Wait,
	// - Wait does not increment waiters if it sees counter == 0.
	// Still do a cheap sanity check to detect WaitGroup misuse.
	if wg.state.Load() != state { // state 值已被修改
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	// Reset waiters count to 0.
	// 如果计数值v为0并且waiter的数量w不为0，那么state的值就是waiter的数量
	// 将waiter的数量设置为0，因为计数值v也是0,所以它们俩的组合*statep直接设置为0即可
	wg.state.Store(0)
	for ; w != 0; w-- {
		runtime_Semrelease(&wg.sema, false, 0) // 唤醒阻塞的 go
	}
}

// Done decrements the WaitGroup counter by one.
func (wg *WaitGroup) Done() { // Done方法实际就是计数器减1
	wg.Add(-1)
}

// Wait blocks until the WaitGroup counter is zero.
func (wg *WaitGroup) Wait() {
	if race.Enabled {
		race.Disable()
	}
	for {
		state := wg.state.Load()
		v := int32(state >> 32) // 当前计数值
		w := uint32(state)      // waiter的数量
		if v == 0 {             // 如果计数值为0, 调用这个方法的goroutine不必再等待，继续执行它后面的逻辑即可
			// Counter is 0, no need to wait.
			if race.Enabled {
				race.Enable()
				race.Acquire(unsafe.Pointer(wg))
			}
			return
		}
		// Increment waiters count.
		// 否则把waiter数量加1。期间可能有并发调用Wait的情况，所以最外层使用了一个for循环
		if wg.state.CompareAndSwap(state, state+1) { // 添加到低 32 位
			if race.Enabled && w == 0 {
				// Wait must be synchronized with the first Add.
				// Need to model this is as a write to race with the read in Add.
				// As a consequence, can do the write only for the first waiter,
				// otherwise concurrent Waits will race with each other.
				race.Write(unsafe.Pointer(&wg.sema))
			}
			runtime_Semacquire(&wg.sema) // 阻塞休眠等待
			if wg.state.Load() != 0 {    // state 值已被修改
				panic("sync: WaitGroup is reused before previous Wait has returned")
			}
			if race.Enabled {
				race.Enable()
				race.Acquire(unsafe.Pointer(wg))
			}
			return // 被唤醒，不再阻塞，返回
		}
	}
}

// state_ 示例
// WaitGroup 的数据结构定义以及得到 state 的地址和信号量的地址
func (wg *WaitGroup) state_() (statep *uint64, semap *uint32) {
	if uintptr(unsafe.Pointer(&wg.state1))%8 == 0 { // 64 位
		// 如果地址是64bit对齐的，数组前两个元素做state，后一个元素做信号量
		return (*uint64)(unsafe.Pointer(&wg.state1)), &wg.state1[2]
	} else { // 32 位
		// 如果地址是32bit对齐的，数组后两个元素用来做state，它可以用来做64bit的原子操作
		return (*uint64)(unsafe.Pointer(&wg.state1[1])), &wg.state1[0]
	}
}
