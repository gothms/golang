package lock

import (
	"fmt"
	"sync"
	"testing"
)

/*
1.锁：别让性能被 “锁” 住

	很多程序的性能问题，都是由锁导致

	误区：sync.RWMutex
		互斥的写锁切换有一个性能消耗
		而读锁不互斥，没有很大的性能消耗，和没锁的性能差不多
	事实并非如此：没加锁 和 加读锁
		BenchmarkLockFree-8          285           4232609 ns/op             720 B/op         41 allocs/op
		BenchmarkLock-8                9         129882522 ns/op            1000 B/op         42 allocs/op
	通过 cpuprofile 分析
		RLock 和 RUnlock 都用了近 3s，但也是作为参考

2.sync.Map：Go 内置

	2.1.适合读多写少，且 key 相对稳定的环境
	2.2.采用了空间换时间的方案，并且采用指针的方式间接实现值的映射，所以存储空间会较 built-in map 大
		Read Only：R区域
		Dirty：RW区域
	2.3.参考：https://my.oschina.net/qiangmzsx/blog/1827059

3.ConcurrentMap：源自 Java

	partition 原理
	适用于读写都很频繁的情况

	https://github.com/easierway/concurrent_map

3.总结

	3.1.减少锁的影响范围
	3.2.减少发生锁冲突的概率
		sync.Map：读很多，写很少
		ConcurrentMap：读少，写多
	3.3.避免锁的使用
		LAMX Disruptor：https://martinfowler.com/articles/lmax.html
		lock free 非常高性能的数据交换的队列，可以在一台普通的 linepop？上实现百万 QPS
*/
const (
	NumOfReader = 40
	ReadTimes   = 100000
)

var cache map[string]string

func init() {
	cache = make(map[string]string)
	cache["a"] = "aa"
	cache["b"] = "bb"
}
func lockFreeAccess() {
	var wg sync.WaitGroup
	wg.Add(NumOfReader)
	for i := 0; i < NumOfReader; i++ {
		go func() {
			for j := 0; j < ReadTimes; j++ {
				if _, err := cache["a"]; !err {
					fmt.Println("Nothing")
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
func lockAccess() {
	var wg sync.WaitGroup
	wg.Add(NumOfReader)
	mut := new(sync.RWMutex)
	for i := 0; i < NumOfReader; i++ {
		go func() {
			for j := 0; j < ReadTimes; j++ {
				mut.RLock() // 测试读锁
				if _, err := cache["a"]; !err {
					fmt.Println("Nothing")
				}
				mut.RUnlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
func BenchmarkLockFree(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lockFreeAccess()
	}
}
func BenchmarkLock(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lockAccess()
	}
}
