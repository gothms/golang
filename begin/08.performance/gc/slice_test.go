package gc

import "testing"

/*
arr := make([]int, 0, 0)
arr := make([]int, 0, N)
arr := make([]int, N, N)
arr := make([]int, N<<2, N<<2)

BenchmarkSliceAppend-8                     17626             68789 ns/op          357626 B/op         19 allocs/op
BenchmarkSliceAppendWithCap-8              67378             16335 ns/op           81920 B/op          1 allocs/op
BenchmarkSliceIndex-8                      87651             14941 ns/op           81920 B/op          1 allocs/op
BenchmarkSliceIndexOverSize-8              36188             32476 ns/op          327681 B/op          1 allocs/op
*/
const N = 10000

func sliceAppend() []int {
	arr := make([]int, 0, 0)
	for i := 0; i < N; i++ {
		arr = append(arr, i)
	}
	return arr
}
func sliceAppendWithCap() []int {
	arr := make([]int, 0, N)
	for i := 0; i < N; i++ {
		arr = append(arr, i)
	}
	return arr
}
func sliceIndex() []int {
	arr := make([]int, N, N)
	for i := 0; i < N; i++ {
		arr[i] = i
	}
	return arr
}
func sliceIndexOverSize() []int {
	M := N << 2
	arr := make([]int, M, M)
	for i := 0; i < N; i++ {
		arr[i] = i
	}
	return arr
}
func BenchmarkSliceAppend(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sliceAppend()
	}
	b.StopTimer()
}
func BenchmarkSliceAppendWithCap(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sliceAppendWithCap()
	}
	b.StopTimer()
}
func BenchmarkSliceIndex(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sliceIndex()
	}
	b.StopTimer()
}
func BenchmarkSliceIndexOverSize(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sliceIndexOverSize()
	}
	b.StopTimer()
}
