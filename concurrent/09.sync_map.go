package concurrent

import (
	"fmt"
	"sync"
)

/*
map：如何实现线程安全的map类型？

哈希表（Hash Table）
	它实现的就是 key-value 之间的映射关系，主要提供的方法包括 Add、Lookup、Delete 等
	因为这种数据结构是一个基础的数据结构，每个 key 都会有一个唯一的索引值，通过索引可以很快地找到对应的值，所以使用哈希表进行数据的插入和读取都是很快的
	Go 语言本身就内建了这样一个数据结构，也就是 map 数据类型

map 的基本使用方法
	类型
		key 类型的 K 必须是可比较的（comparable），也就是可以通过 == 和 != 操作符进行比较
		value 的值和类型无所谓，可以是任意的类型，或者为 nil
		通常情况下，我们会选择内建的基本类型，比如整数、字符串做 key 的类型
	comparable
		在 Go 语言中，bool、整数、浮点数、复数、字符串、指针、Channel、接口都是可比较的
		包含可比较元素的 struct 和数组，这俩也是可比较的
		而 slice、map、函数值都是不可比较的
	注意
		如果使用 struct 类型做 key 其实是有坑的
		因为如果 struct 的某个字段值修改了，查询 map 时无法获取它 add 进去的值
		如果要使用 struct 作为 key，我们要保证 struct 对象在逻辑上是不可变的，这样才会保证 map 的逻辑没有问题
	返回值
		在 Go 中，map[key]函数返回结果可以是一个值，也可以是两个值
		如果获取一个不存在的 key 对应的值时，会返回零值
	无序
		map 是无序的
		所以当遍历一个 map 对象的时候，迭代的元素的顺序是不确定的，无法保证两次遍历的顺序是一样的，也不能保证和插入的顺序一致
	有序
		人为排序
			如果我们想要按照 key 的顺序获取 map 的值
			需要先取出所有的 key 进行排序，然后按照这个排序的 key 依次获取对应的值
		辅助数据结构
			如果我们想要保证元素有序，比如按照元素插入的顺序进行遍历
			可以使用辅助的数据结构，比如 orderedmap，来记录插入顺序

使用 map 的 2 种常见错误
常见错误一：未初始化
	panic
		和 slice 或者 Mutex、RWmutex 等 struct 类型不同，map 对象必须在使用之前初始化
		如果不初始化就直接赋值的话，会出现 panic 异常
		目前还没有工具可以检查
	从一个 nil 的 map 对象中获取值不会 panic，而是会得到零值
常见错误二：并发读写
	panic
		如果没有注意到并发问题，程序在运行的时候就有可能出现并发读写导致的 panic
		Go 内建的 map 对象不是线程（goroutine）安全的，并发读写的时候运行时会有检查，遇到并发问题就会导致 panic
	示例：TestMapSyncPanic
		虽然读写 goroutine 各自操作不同的元素，貌似 map 也没有扩容的问题
		但是运行时检测到同时对 map 对象有并发访问，就会直接 panic
		报错
			fatal error: concurrent map read and map write
	Docker issue 40772
		在删除 map 对象的元素时忘记了加锁
		图示 09.map_docker_panic.jpg
	其他并发读写 map issue
		Docker issue 40772，Docker issue 35588、34540、39643 等
		Kubernetes 的 issue 84431、72464、68647、64484、48045、45593、37560 等，以及 TiDB 的 issue 14960 和 17494 等

如何实现线程安全的 map 类型？
	使用并发 map 的过程中，加锁和分片加锁这两种方案都比较常用
	如果是追求更高的性能，显然是分片加锁更好，因为它可以降低锁的粒度，进而提高访问此 map 对象的吞吐
	如果并发性能要求不是那么高的场景，简单加锁方式更简单
加读写锁：扩展 map，支持并发读写
	通过 interface{}来模拟泛型
		但还是要涉及接口和具体类型的转换，比较复杂，还不如将要发布的泛型方案更直接、性能更好
	示例：RWMap
		查询和遍历可以看做读操作，增加、修改和删除可以看做写操作
		可以通过读写锁对相应的操作进行保护
	锁是性能下降的万恶之源之一
		虽然使用读写锁可以提供线程安全的 map，但是在大量并发读写的情况下，锁的竞争会非常激烈
分片加锁：更高效的并发 map
	在并发编程中，我们的一条原则就是尽量减少锁的使用
		一些单线程单进程的应用（比如Redis 等），基本上不需要使用锁去解决并发线程访问的问题，所以可以取得很高的性能
		但是对于 Go 开发的应用程序来说，并发是常用的一个特性
		在这种情况下，我们能做的就是，尽量减少锁的粒度和锁的持有时间
	优化思路
		1.可以优化业务处理的代码，以此来减少锁的持有时间，比如将串行的操作变成并行的子任务执行
		2.对同步原语的优化，减少锁的粒度（这里使用的思路）
	分片（Shard）
		减少锁的粒度常用的方法就是分片（Shard），将一把锁分成几把锁，每个锁控制一个分片
		Go 比较知名的分片并发 map 的实现是 orcaman/concurrent-map
	orcaman/concurrent-map
		它默认采用 32 个分片，GetShard 是一个关键的方法，能够根据 key 计算出分片索引
			func (m ConcurrentMap[K, V]) GetShard(key K) *ConcurrentMapShared[K, V]
		增加或者查询的时候，首先根据分片索引得到分片对象，然后对分片对象加锁进行操作
			func (m ConcurrentMap[K, V]) Set(key K, value V)
			func (m ConcurrentMap[K, V]) Get(key K) (V, bool)
		除了 GetShard 方法，ConcurrentMap 还提供了很多其他的方法
			这些方法都是通过计算相应的分片实现的，目的是保证把锁的粒度限制在分片上
		https://github.com/orcaman/concurrent-map

应对特殊场景的 sync.Map
	sync.Map 简介
		Go 1.9 中增加了一个线程安全的 map，也就是 sync.Map
		Go 官方线程安全 map 的标准实现。虽然是官方标准，反而是不常用的，为什么呢？
		一句话来说就是 map 要解决的场景很难描述，很多时候在做抉择时根本就不知道该不该用它
		但是呢，确实有一些特定的场景，我们需要用到 sync.Map 来实现
	特殊场景
		官方文档：sync/map.go
			官方的文档中指出，在以下两个场景中使用 sync.Map，会比使用 map+RWMutex 的方式，性能要好得多
		1. 只会增长的缓存系统中，一个 key 只写入一次而被读很多次
		2. 多个 goroutine 为不相交的键集读、写和重写键值对
		这两个场景说得都比较笼统，而且，这些场景中还包含了一些特殊的情况
		所以，官方建议你针对自己的场景做性能评测，如果确实能够显著提高性能，再使用 sync.Map
	作者
		即使是 sync.Map 的作者 Bryan C.Mills，也很少使用 sync.Map
		即便是在使用 sync.Map 的时候，也是需要临时查询它的 API，才能清楚记住它的功能
		所以，我们可以把 sync.Map 看成一个生产环境中很少使用的同步原语
sync.Map 的实现
	sync.Map 的实现有几个优化点
		空间换时间
			通过冗余的两个数据结构（只读的 read 字段、可写的 dirty），来减少加锁对性能的影响
			对只读字段（read）的操作不需要加锁
		优先从 read 字段读取、更新、删除
			因为对 read 字段的读取不需要锁
		动态调整
			miss 次数多了之后，将 dirty 数据提升为 read，避免总是从 dirty 中加锁读取
		double-checking
			加锁之后先还要再检查 read 字段，确定真的不存在才操作 dirty 字段
		延迟删除
			删除一个键值只是打标记，只有在提升 dirty 字段为 read 字段的时候才清理删除的数据
	数据结构
		指向同一地址
			如果 dirty 字段非 nil 的话，map 的 read 字段和 dirty 字段会包含相同的非 expunged 的项
			所以如果通过 read 字段更改了这个项的值，从 dirty 字段中也会读取到这个项的新值
			因为本来它们指向的就是同一个地址
			即底层数据存储的是指针，指向的是同一份值
		dirty 包含重复项目的好处就是
			一旦 miss 数达到阈值需要将 dirty 提升为 read 的话，只需简单地把 dirty 设置为 read 对象即可
		不好的一点就是
			当创建新的 dirty 对象的时候，需要逐条遍历 read，把非 expunged 的项复制到 dirty 对象
	API
		Store、Load 和 Delete
		这三个核心函数的操作都是先从 read 字段中处理的，因为读取 read 字段的时候不用加锁
Store 方法
	Store 既可以是新增元素，也可以是更新元素
		如果运气好的话，更新的是已存在的未被删除的元素，直接更新即可，不会用到锁
		如果运气不好，需要更新（重用）删除的对象、更新还未提升的 dirty 中的对象，或者新增加元素的时候就会使用到了锁
		这个时候，性能就会下降
	dirtyLocked
		新加的元素需要放入到 dirty 中，如果 dirty 为 nil，那么需要从 read 字段中复制出来一个 dirty 对象
	从这一点来看，sync.Map 适合那些只会增长的缓存系统，可以进行更新，但是不要删除，并且不要频繁地增加新元素
Load 方法
	如果幸运的话，我们从 read 中读取到了这个 key 对应的值，那么就不需要加锁了，性能会非常好
		但是，如果请求的 key 不存在或者是新加的，就需要加锁从 dirty 中读取
		所以，读取不存在的 key 会因为加锁而导致性能下降，读取还没有提升的新值的情况下也会因为加锁性能下降
	missLocked
		missLocked 增加 miss 的时候，如果 miss 数等于 dirty 长度，会将 dirty 提升为 read，并将 dirty 置空
Delete 方法
	在 Go 1.15 中欧长坤提供了一个 LoadAndDelete 的实现（go#issue 33762）
		所以 Delete 方法的核心改在了对 LoadAndDelete 中实现了
	如果 read 中不存在，那么就需要从 dirty 中寻找这个项目
		最终，如果项目存在就删除（将它的值标记为 nil）
		如果项目不为 nil 或者没有被标记为 expunged，那么还可以把它的值返回
	nil vs expunged：read 中 key 被删除会有两个状态：nil 和 expunged
		nil和expunged都代表元素被删除了，只不过expunged比较特殊
		如果被删除的元素是expunged，代表它只存在于readonly之中，不存在于dirty中
		这样如果重新设置这个key的话，需要往dirty增加key
Len 方法
	sync.map 还有一些 LoadAndDelete、LoadOrStore、Range 等辅助方法
	但是没有 Len 这样查询 sync.Map 的包含项目数量的方法，并且官方也不准备提供
	如果你想得到 sync.Map 的项目数量的话，你可能不得不通过 Range 逐个计数
	为什么没有 Len 方法：https://link.zhihu.com/?target=https%3A//github.com/golang/go/issues/20680

总结
	Go map 三种实现
		Go 内置的 map 类型使用起来很方便，但是它有一个非常致命的缺陷，那就是它存在着并发问题，所以如果有多个 goroutine 同时并发访问这个 map，就会导致程序崩溃
		Go 官方 Blog 很早就提供了一种加锁的方法，还有后来提供了适用特定场景的线程安全的 sync.Map
		第三方实现的分片式的 map
	建议
		通过性能测试，看看某种线程安全的 map 实现是否满足你的需求
	其他功能的 map 实现
		带有过期功能的 timedmap、使用红黑树实现的 key 有序的 treemap 等
		但它们和并发问题没有关系

思考
	1.为什么 sync.Map 中的集合核心方法的实现中，如果 read 中项目不存在，加锁后还要双检查，再检查一次 read？
	2.你看到 sync.map 元素删除的时候只是把它的值设置为 nil，那么什么时候这个 key 才会真正从 map 对象中删除？

参考
	https://zhuanlan.zhihu.com/p/344834329
*/

func MapNilValue() {
	//m := make(map[bool]int)
	//a := make(map[float64]int)
	//b := make(map[*chan]int)
	//c := make(map[[3]int]int)

	m := make(map[bool][]int)
	m[true] = nil
	m[false] = []int{1}
	for k, v := range m {
		fmt.Println(k, v)
	}
}

// RWMap ==========支持并发读写==========

// RWMap ==========支持并发读写==========
type RWMap struct { // 一个读写锁保护的线程安全的map
	sync.RWMutex // 读写锁保护下面的map字段
	m            map[int]int
}

// NewRWMap 新建一个RWMap
func NewRWMap(n int) *RWMap {
	return &RWMap{
		m: make(map[int]int, n),
	}
}
func (m *RWMap) Get(k int) (int, bool) { //从map中读取一个值
	m.RLock()
	defer m.RUnlock()
	v, existed := m.m[k] // 在锁的保护下从map中读取
	return v, existed
}
func (m *RWMap) Set(k int, v int) { // 设置一个键值对
	m.Lock() // 锁保护
	defer m.Unlock()
	m.m[k] = v
}
func (m *RWMap) Delete(k int) { //删除一个键
	m.Lock() // 锁保护
	defer m.Unlock()
	delete(m.m, k)
}
func (m *RWMap) Len() int { // map的长度
	m.RLock() // 锁保护
	defer m.RUnlock()
	return len(m.m)
}
func (m *RWMap) Each(f func(k, v int) bool) { // 遍历map
	m.RLock() //遍历期间一直持有读锁
	defer m.RUnlock()
	for k, v := range m.m {
		if !f(k, v) {
			return
		}
	}
}
