package maps

import (
	"strconv"
	"sync"
	"testing"
)

const (
	NumOfReader = 2000
	NumOfWriter = 200
)

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
		cm := CreateConcurrentMapBenchmarkAdapter(199)
		benchmarkMap(b, cm)
	})
}
