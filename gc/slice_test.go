package slice

import "testing"

/*
arr := make([]int, 0, 0)
arr := make([]int, 0, N)
arr := make([]int, N, N)

BenchmarkSliceAppend-8                    177627              5811 ns/op           25208 B/op         12 allocs/op
BenchmarkSliceAppendWithCap-8            2418781               494.7 ns/op             0 B/op          0 allocs/op
BenchmarkSliceIndex-8                    4510444               264.4 ns/op             0 B/op          0 allocs/op
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
