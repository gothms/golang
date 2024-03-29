// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package source

import (
	"sync/atomic"
)

// Map is like a Go map[interface{}]interface{} but is safe for concurrent use
// by multiple goroutines without additional locking or coordination.
// Loads, stores, and deletes run in amortized constant time.
//
// The Map type is specialized. Most code should use a plain Go map instead,
// with separate locking or coordination, for better type safety and to make it
// easier to maintain other invariants along with the map content.
//
// The Map type is optimized for two common use cases: (1) when the entry for a given
// key is only ever written once but read many times, as in caches that only grow,
// or (2) when multiple goroutines read, write, and overwrite entries for disjoint
// sets of keys. In these two cases, use of a Map may significantly reduce lock
// contention compared to a Go map paired with a separate Mutex or RWMutex.
//
// The zero Map is empty and ready for use. A Map must not be copied after first use.
//
// In the terminology of the Go memory model, Map arranges that a write operation
// “synchronizes before” any read operation that observes the effect of the write, where
// read and write operations are defined as follows.
// Load, LoadAndDelete, LoadOrStore, Swap, CompareAndSwap, and CompareAndDelete
// are read operations; Delete, LoadAndDelete, Store, and Swap are write operations;
// LoadOrStore is a write operation when it returns loaded set to false;
// CompareAndSwap is a write operation when it returns swapped set to true;
// and CompareAndDelete is a write operation when it returns deleted set to true.
type Map struct {
	mu Mutex

	// read contains the portion of the map's contents that are safe for
	// concurrent access (with or without mu held).
	//
	// The read field itself is always safe to load, but must only be stored with
	// mu held.
	//
	// Entries stored in read may be updated concurrently without mu, but updating
	// a previously-expunged entry requires that the entry be copied to the dirty
	// map and unexpunged with mu held.
	// 基本上你可以把它看成一个安全的只读的map
	read atomic.Pointer[readOnly] // 它包含的元素其实也是通过原子操作更新的，但是已删除的entry就需要加锁操作了

	// dirty contains the portion of the map's contents that require mu to be
	// held. To ensure that the dirty map can be promoted to the read map quickly,
	// it also includes all of the non-expunged entries in the read map.
	//
	// Expunged entries are not stored in the dirty map. An expunged entry in the
	// clean map must be unexpunged and added to the dirty map before a new value
	// can be stored to it.
	//
	// If the dirty map is nil, the next write to the map will initialize it by
	// making a shallow copy of the clean map, omitting stale entries.
	// 包含需要加锁才能访问的元素
	dirty map[any]*entry // 包括所有在read字段中但未被expunged（删除）的元素以及新加的元素

	// misses counts the number of loads since the read map was last updated that
	// needed to lock mu to determine whether the key was present.
	//
	// Once enough misses have occurred to cover the cost of copying the dirty
	// map, the dirty map will be promoted to the read map (in the unamended
	// state) and the next store to the map will make a new dirty copy.
	misses int // 记录从read中读取miss的次数，一旦miss数和dirty长度一样了，就会把dirty提升为read，并把dirty置为nil
}

// readOnly is an immutable struct stored atomically in the Map.read field.
type readOnly struct {
	m map[any]*entry
	// 当dirty中包含read没有的数据时为true，比如新增一条数据
	amended bool // true if the dirty map contains some key not in m.
}

// expunged is an arbitrary pointer that marks entries which have been deleted
// from the dirty map.
// expunged是用来标识此项已经删掉的指针
var expunged = new(any) // 当map中的一个项目被删除了，只是把它的值标记为expunged，以后才有机会真正删除此项

// An entry is a slot in the map corresponding to a particular key.
type entry struct { // entry代表一个值
	// p points to the interface{} value stored for the entry.
	//
	// If p == nil, the entry has been deleted, and either m.dirty == nil or
	// m.dirty[key] is e.
	//
	// If p == expunged, the entry has been deleted, m.dirty != nil, and the entry
	// is missing from m.dirty.
	//
	// Otherwise, the entry is valid and recorded in m.read.m[key] and, if m.dirty
	// != nil, in m.dirty[key].
	//
	// An entry can be deleted by atomic replacement with nil: when m.dirty is
	// next created, it will atomically replace nil with expunged and leave
	// m.dirty[key] unset.
	//
	// An entry's associated value can be updated by atomic replacement, provided
	// p != expunged. If p == expunged, an entry's associated value can be updated
	// only after first setting m.dirty[key] = e so that lookups using the dirty
	// map find the entry.
	p atomic.Pointer[any]
}

func newEntry(i any) *entry {
	e := &entry{}
	e.p.Store(&i)
	return e
}

func (m *Map) loadReadOnly() readOnly {
	if p := m.read.Load(); p != nil {
		return *p
	}
	return readOnly{}
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (m *Map) Load(key any) (value any, ok bool) {
	read := m.loadReadOnly() // 首先从read处理
	e, ok := read.m[key]
	if !ok && read.amended { // 如果不存在并且dirty不为nil(有新的元素)，即 dirty 中存在 read 中没有的 key
		m.mu.Lock() // dirty map 是 runtime.map，并发不安全
		// Avoid reporting a spurious miss if m.dirty got promoted while we were
		// blocked on m.mu. (If further loads of the same key will not miss, it's
		// not worth copying the dirty map for this key.)
		read = m.loadReadOnly() // double check 双检查（避免在上锁的过程中 dirty map 提升为 read map），看看read中现在是否存在此key
		e, ok = read.m[key]
		if !ok && read.amended { //依然不存在，并且dirty不为nil
			e, ok = m.dirty[key] // 从dirty中读取
			// Regardless of whether the entry was present, record a miss: this key
			// will take the slow path until the dirty map is promoted to the read
			// map.
			m.missLocked() // 不管dirty中存不存在，miss数都加1
		}
		m.mu.Unlock()
	}
	if !ok { // 没找到
		return nil, false
	}
	return e.load() //返回读取的对象，e既可能是从read中获得的，也可能是从dirty中获得的
}

func (e *entry) load() (value any, ok bool) {
	p := e.p.Load()
	if p == nil || p == expunged { // 对于 nil 和 expunged 状态的 entry，直接返回 ok=false
		return nil, false
	}
	return *p, true
}

// Store sets the value for a key.
func (m *Map) Store(key, value any) {
	_, _ = m.Swap(key, value)
}

// tryCompareAndSwap compare the entry with the given old value and swaps
// it with a new value if the entry is equal to the old value, and the entry
// has not been expunged.
//
// If the entry is expunged, tryCompareAndSwap returns false and leaves
// the entry unchanged.
func (e *entry) tryCompareAndSwap(old, new any) bool {
	p := e.p.Load()
	if p == nil || p == expunged || *p != old {
		return false
	}

	// Copy the interface after the first load to make this method more amenable
	// to escape analysis: if the comparison fails from the start, we shouldn't
	// bother heap-allocating an interface value to store.
	nc := new
	for {
		if e.p.CompareAndSwap(p, &nc) {
			return true
		}
		p = e.p.Load()
		if p == nil || p == expunged || *p != old {
			return false
		}
	}
}

// unexpungeLocked ensures that the entry is not marked as expunged.
//
// If the entry was previously expunged, it must be added to the dirty map
// before m.mu is unlocked.
func (e *entry) unexpungeLocked() (wasExpunged bool) { // 确保了 entry 没有被标记成已被清除
	// 如果 entry 先前被清除过了，那么在 mutex 解锁之前，它一定要被加入到 dirty map 中
	return e.p.CompareAndSwap(expunged, nil) // 将 p 的状态由 expunged  更改为 nil
}

// swapLocked unconditionally swaps a value into the entry.
//
// The entry must be known not to be expunged.
func (e *entry) swapLocked(i *any) *any {
	return e.p.Swap(i)
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *Map) LoadOrStore(key, value any) (actual any, loaded bool) { // 结合了 Load 和 Store 的功能
	// Avoid locking if it's a clean hit.
	read := m.loadReadOnly()
	if e, ok := read.m[key]; ok {
		actual, loaded, ok := e.tryLoadOrStore(value)
		if ok {
			return actual, loaded
		}
	}

	m.mu.Lock()
	read = m.loadReadOnly()
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
		actual, loaded, _ = e.tryLoadOrStore(value)
	} else if e, ok := m.dirty[key]; ok {
		actual, loaded, _ = e.tryLoadOrStore(value)
		m.missLocked()
	} else {
		if !read.amended {
			// We're adding the first new key to the dirty map.
			// Make sure it is allocated and mark the read-only map as incomplete.
			m.dirtyLocked()
			m.read.Store(&readOnly{m: read.m, amended: true})
		}
		m.dirty[key] = newEntry(value)
		actual, loaded = value, false
	}
	m.mu.Unlock()

	return actual, loaded
}

// tryLoadOrStore atomically loads or stores a value if the entry is not
// expunged.
//
// If the entry is expunged, tryLoadOrStore leaves the entry unchanged and
// returns with ok==false.
func (e *entry) tryLoadOrStore(i any) (actual any, loaded, ok bool) {
	p := e.p.Load()
	if p == expunged {
		return nil, false, false
	}
	if p != nil {
		return *p, true, true
	}

	// Copy the interface after the first load to make this method more amenable
	// to escape analysis: if we hit the "load" path or the entry is expunged, we
	// shouldn't bother heap-allocating.
	ic := i
	for {
		if e.p.CompareAndSwap(nil, &ic) {
			return i, false, true
		}
		p = e.p.Load()
		if p == expunged {
			return nil, false, false
		}
		if p != nil {
			return *p, true, true
		}
	}
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *Map) LoadAndDelete(key any) (value any, loaded bool) {
	read := m.loadReadOnly()
	e, ok := read.m[key]
	if !ok && read.amended { // 如果 read 中没有这个 key，且 dirty map 不为空
		m.mu.Lock()
		read = m.loadReadOnly() // 双检查
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			// 直接从 dirty 中删除这个 key
			delete(m.dirty, key) // 这一行长坤在1.15中实现的时候忘记加上了，导致在特殊的场景下有些key总是没有被回收...
			// Regardless of whether the entry was present, record a miss: this key
			// will take the slow path until the dirty map is promoted to the read
			// map.
			m.missLocked() // miss数加1
		}
		m.mu.Unlock()
	}
	if ok {
		return e.delete() // 如果在 read 中找到了这个 key，将 p 置为 nil
	}
	return nil, false
}

// Delete deletes the value for a key.
func (m *Map) Delete(key any) {
	m.LoadAndDelete(key)
}

func (e *entry) delete() (value any, ok bool) { // 将正常状态（指向一个 any）的 p 设置成 nil
	for {
		p := e.p.Load()
		if p == nil || p == expunged {
			return nil, false
		}
		// 没有设置成 expunged 的原因是，当 p 为 expunged 时，表示它已经不在 dirty 中了
		// 这是 p 的状态机决定的，而在 tryExpungeLocked 函数中，会将 nil 原子的设置成 expunged
		// tryExpungeLocked 是在新创建 dirty 时调用的，会将已被删除的 entry.p 从 nil 改成 expunged，这个 entry 就不会写入 dirty 了
		if e.p.CompareAndSwap(p, nil) { // 通过原子的（CAS 操作）设置 p 为 nil 被删除
			return *p, true
		}
	}
}

// trySwap swaps a value if the entry has not been expunged.
//
// If the entry is expunged, trySwap returns false and leaves the entry
// unchanged.
func (e *entry) trySwap(i *any) (*any, bool) {
	for { // 原子操作：for + CAS
		p := e.p.Load()
		if p == expunged { // 如果 p == expunged（entry 被删），那么仅当它初次被设置到 m.dirty 之后，才可以被更新
			return nil, false
		}
		if e.p.CompareAndSwap(p, i) { // CAS 更新
			return p, true
		}
	}
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (m *Map) Swap(key, value any) (previous any, loaded bool) {
	read := m.loadReadOnly()
	if e, ok := read.m[key]; ok { // 如果read字段包含这个项，说明是更新，cas更新项目的值即可
		if v, ok := e.trySwap(&value); ok { // 尝试更新，由于修改的是 entry 内部的 pointer，因此 dirty map 也可见
			if v == nil {
				return nil, false
			}
			return *v, true
		}
	}

	m.mu.Lock() // read中不存在，或者cas更新失败，就需要加锁访问dirty了
	read = m.loadReadOnly()
	if e, ok := read.m[key]; ok { // 双检查，看看read是否已经存在了
		if e.unexpungeLocked() {
			// read map 中存在该 key，但 p == expunged，说明 m.dirty != nil，且 m.dirty 中没有这个 key
			// 1.将 p 的状态由 expunged  更改为 nil
			// 2.dirty map 插入 key
			// The entry was previously expunged, which implies that there is a
			// non-nil dirty map and this entry is not in it.
			m.dirty[key] = e // 此项目先前已经被删除了，通过将它的值设置为nil，标记为unexpunged
		}
		if v := e.swapLocked(&value); v != nil { // 更新 entry.p = value (read map 和 dirty map 指向同一个 entry)
			loaded = true
			previous = *v
		}
	} else if e, ok := m.dirty[key]; ok { // 如果dirty中有此项
		if v := e.swapLocked(&value); v != nil { // 直接更新 entry(read map 中仍然没有这个 key)
			loaded = true
			previous = *v
		}
	} else { // 否则就是一个新的key
		// 1.如果 dirty map 为空，则需要创建 dirty map，并从 read map 中拷贝未删除的元素到新创建的 dirty map
		// 2.更新 amended 字段，标识 dirty map 中存在 read map 中没有的 key
		// 3.将 kv 写入 dirty map 中，read 不变
		if !read.amended { //如果dirty为nil
			// We're adding the first new key to the dirty map.
			// Make sure it is allocated and mark the read-only map as incomplete.
			m.dirtyLocked()                                   // 需要创建dirty对象，并且标记read的amended为true（说明有元素它不包含而dirty包含）
			m.read.Store(&readOnly{m: read.m, amended: true}) // 更新 amended
		}
		m.dirty[key] = newEntry(value) //将新值增加到dirty对象中
	}
	m.mu.Unlock()
	return previous, loaded
}

// CompareAndSwap swaps the old and new values for key
// if the value stored in the map is equal to old.
// The old value must be of a comparable type.
func (m *Map) CompareAndSwap(key, old, new any) bool {
	read := m.loadReadOnly()
	if e, ok := read.m[key]; ok {
		return e.tryCompareAndSwap(old, new)
	} else if !read.amended {
		return false // No existing value for key.
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	read = m.loadReadOnly()
	swapped := false
	if e, ok := read.m[key]; ok {
		swapped = e.tryCompareAndSwap(old, new)
	} else if e, ok := m.dirty[key]; ok {
		swapped = e.tryCompareAndSwap(old, new)
		// We needed to lock mu in order to load the entry for key,
		// and the operation didn't change the set of keys in the map
		// (so it would be made more efficient by promoting the dirty
		// map to read-only).
		// Count it as a miss so that we will eventually switch to the
		// more efficient steady state.
		m.missLocked()
	}
	return swapped
}

// CompareAndDelete deletes the entry for key if its value is equal to old.
// The old value must be of a comparable type.
//
// If there is no current value for key in the map, CompareAndDelete
// returns false (even if the old value is the nil interface value).
func (m *Map) CompareAndDelete(key, old any) (deleted bool) {
	read := m.loadReadOnly()
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read = m.loadReadOnly()
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			// Don't delete key from m.dirty: we still need to do the “compare” part
			// of the operation. The entry will eventually be expunged when the
			// dirty map is promoted to the read map.
			//
			// Regardless of whether the entry was present, record a miss: this key
			// will take the slow path until the dirty map is promoted to the read
			// map.
			m.missLocked()
		}
		m.mu.Unlock()
	}
	for ok {
		p := e.p.Load()
		if p == nil || p == expunged || *p != old {
			return false
		}
		if e.p.CompareAndSwap(p, nil) {
			return true
		}
	}
	return false
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently (including by f), Range may reflect any
// mapping for that key from any point during the Range call. Range does not
// block other methods on the receiver; even f itself may call any method on m.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (m *Map) Range(f func(key, value any) bool) { // Range 将遍历调用时刻 map 中的所有 k-v 对，将它们传给 f 函数，如果 f 返回 false，将停止遍历
	// We need to be able to iterate over all of the keys that were already
	// present at the start of the call to Range.
	// If read.amended is false, then read.m satisfies that property without
	// requiring us to hold m.mu for a long time.
	read := m.loadReadOnly()
	if read.amended { // dirty 中含有 read 中没有的 key
		// m.dirty contains keys not in read.m. Fortunately, Range is already O(N)
		// (assuming the caller does not break out early), so a call to Range
		// amortizes an entire copy of the map: we can promote the dirty copy
		// immediately!
		m.mu.Lock()
		read = m.loadReadOnly()
		if read.amended {
			read = readOnly{m: m.dirty} // 将 dirty 提升为 read，会将开销分摊开来，所以这里直接晋升
			m.read.Store(&read)
			m.dirty = nil
			m.misses = 0
		}
		m.mu.Unlock()
	}

	for k, e := range read.m { // 遍历 read，取出 entry 中的值，调用 f(k, v)
		v, ok := e.load()
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}

func (m *Map) missLocked() {
	m.misses++                   // misses计数加一
	if m.misses < len(m.dirty) { // 如果没达到阈值(dirty字段的长度),返回
		return
	}
	m.read.Store(&readOnly{m: m.dirty}) // 把dirty字段的内存提升为read字段
	m.dirty = nil                       // 清空dirty
	m.misses = 0                        // misses数重置为0
}

func (m *Map) dirtyLocked() {
	if m.dirty != nil { // 如果dirty字段已经存在，不需要创建了
		return
	}

	read := m.loadReadOnly() // 获取read字段
	m.dirty = make(map[any]*entry, len(read.m))
	for k, e := range read.m { // 遍历read字段
		if !e.tryExpungeLocked() { // 把非 expunged 的键值对复制到dirty中
			m.dirty[k] = e // 浅拷贝
		}
	}
}

func (e *entry) tryExpungeLocked() (isExpunged bool) { // 将 nil 原子的设置成 expunged
	p := e.p.Load()
	for p == nil {
		// 如果原来是 nil，说明原 key 已被删除，则将其转为 expunged
		if e.p.CompareAndSwap(nil, expunged) { // 如果之后创建 m.dirty，nil 又会被原子的设置为 expunged，且不会拷贝到 dirty 中
			return true
		}
		p = e.p.Load()
	}
	return p == expunged // 键值对已被删除，m.dirty != nil，且 m.dirty 中没有这个 key
}
