package test

import (
	"golang/generic"
	"testing"
)

// benchmark_simple/add_test.go
func BenchmarkAddInt(b *testing.B) {
	b.ReportAllocs()
	var m, n int = 5, 6
	for i := 0; i < b.N; i++ {
		generic.AddInt(m, n)
	}
}
func BenchmarkAddIntGeneric(b *testing.B) {
	b.ReportAllocs()
	var m, n int = 5, 6
	for i := 0; i < b.N; i++ {
		generic.Add(m, n)
	}
}

func TestGenericSort(t *testing.T) {
	intSlice := []int{3, 6, 1, 8, 3}
	generic.SortGeneric[int](intSlice)
	t.Log(intSlice)
	floatSlice := []float64{3.14, 7.32, 0.15, 5.67, 1.11}
	generic.SortGeneric[float64](floatSlice)
	t.Log(floatSlice)
	strSlice := []string{"generic", "internet", "golang", "redis", "algorithms"}
	generic.SortGeneric[string](strSlice)
	t.Log(strSlice)
}
