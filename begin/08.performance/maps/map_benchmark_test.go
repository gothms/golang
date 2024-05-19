package maps

import (
	"strconv"
	"sync"
	"testing"
)

const (
	NumOfReader = 200
	NumOfWriter = 2000
)

/*
$ go test -bench=BenchmarkSyncMAP golang/08.performance/maps -benchmem

NumOfReader = 2000
NumOfWriter = 200
	BenchmarkSyncMAP/map_with_RWLock-8                    44          26456686 ns/op         4530925 B/op     296285 allocs/op
	BenchmarkSyncMAP/sync.map-8                           40          28107385 ns/op         7758007 B/op     401346 allocs/op
	BenchmarkSyncMAP/Concurrent_Map-8                     66          17475029 ns/op        12898939 B/op     816797 allocs/op
	BenchmarkSyncMAP/OrcamanConcurrent_Map-8             127           9266371 ns/op         4561222 B/op     296612 allocs/op

NumOfReader = 20000
NumOfWriter = 2000
	BenchmarkSyncMAP/map_with_RWLock-8                     4         283859800 ns/op        46379198 B/op    2972238 allocs/op
	BenchmarkSyncMAP/sync.map-8                           10         201134830 ns/op        66360519 B/op    3700509 allocs/op
	BenchmarkSyncMAP/Concurrent_Map-8                      6         171040000 ns/op        129611954 B/op   8174462 allocs/op
	BenchmarkSyncMAP/OrcamanConcurrent_Map-8              12          93942700 ns/op        46124350 B/op    2971459 allocs/op

NumOfReader = 2000
NumOfWriter = 2000
	BenchmarkSyncMAP/map_with_RWLock-8                     7         152355314 ns/op        15797014 B/op    1142029 allocs/op
	BenchmarkSyncMAP/sync.map-8                            7         199874771 ns/op        40584574 B/op    2045450 allocs/op
	BenchmarkSyncMAP/Concurrent_Map-8                     19          58969505 ns/op        41483748 B/op    2742954 allocs/op
	BenchmarkSyncMAP/OrcamanConcurrent_Map-8              28          36860025 ns/op        15879554 B/op    1142910 allocs/op

NumOfReader = 200
NumOfWriter = 2000
	BenchmarkSyncMAP/map_with_RWLock-8                     8         141167000 ns/op        12782319 B/op     959274 allocs/op
	BenchmarkSyncMAP/sync.map-8                            7         181355014 ns/op        37396918 B/op    1866083 allocs/op
	BenchmarkSyncMAP/Concurrent_Map-8                     21          48180829 ns/op        32644279 B/op    2199526 allocs/op
	BenchmarkSyncMAP/OrcamanConcurrent_Map-8              38          31361611 ns/op        12797522 B/op     959455 allocs/op

总结：08.performance/lock/lock_test.go
*/

func benchmarkMap(b *testing.B, m Map) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for j := 0; j < NumOfWriter; j++ {
			wg.Add(1)
			go func() {
				for k := 0; k < 100; k++ {
					m.Set(strconv.Itoa(k), k*k)
					m.Set(strconv.Itoa(k), k*k)
					m.Del(strconv.Itoa(k))
				}
				wg.Done()
			}()
		}
		for j := 0; j < NumOfReader; j++ {
			wg.Add(1)
			go func() {
				for k := 0; k < 100; k++ {
					m.Get(strconv.Itoa(k))
				}
			}()
			wg.Done()
		}
		wg.Wait()
	}
}
func BenchmarkSyncMAP(b *testing.B) {
	b.Run("map with RWLock", func(b *testing.B) {
		m := CreateRWLockMap()
		benchmarkMap(b, m)
	})
	b.Run("sync.map", func(b *testing.B) {
		syncM := CreateSyncMapBenchmarkAdapter()
		benchmarkMap(b, syncM)
	})
	b.Run("Concurrent Map", func(b *testing.B) {
		cm := CreateConcurrentMapBenchmarkAdapter(32)
		benchmarkMap(b, cm)
	})
	b.Run("OrcamanConcurrent Map", func(b *testing.B) {
		cmap := CreateOrcamanConcurrentMapBenchmarkAdapter()
		benchmarkMap(b, cmap)
	})
}
