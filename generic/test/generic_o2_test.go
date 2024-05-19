package test

import (
	"fmt"
	"golang/generic"
	"testing"
)

/*
PS E:\gothmslee\golang> go test -bench="." generic\test\generic_o2_test.go -benchmem
goos: windows
goarch: amd64
cpu: Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz
BenchmarkMaxInt-8               435920424                2.773 ns/op           0 B/op          0 allocs/op
BenchmarkMaxAny-8               85000626                14.37 ns/op            0 B/op          0 allocs/op
BenchmarkMaxGenerics-8          393862182                3.175 ns/op           0 B/op          0 allocs/op
*/
func BenchmarkMaxInt(b *testing.B) {
	sl := []int{1, 2, 3, 4, 7, 8, 9, 0}
	for i := 0; i < b.N; i++ {
		generic.MaxInt(sl)
	}
}
func BenchmarkMaxAny(b *testing.B) {
	sl := []any{1, 2, 3, 4, 7, 8, 9, 0}
	for i := 0; i < b.N; i++ {
		generic.MaxAny(sl)
	}
}
func BenchmarkMaxGenerics(b *testing.B) {
	sl := []int{1, 2, 3, 4, 7, 8, 9, 0}
	for i := 0; i < b.N; i++ {
		generic.MaxGenerics[int](sl)
	}
}
func TestMaxGenerics(t *testing.T) {
	maxGenericsInt := generic.MaxGenerics[int] // 实例化后得到的新“机器”：maxGenericsInt
	fmt.Printf("%T\n", maxGenericsInt)         // func([]int) int
	genericsInt := maxGenericsInt([]int{1, 2, 3, 4, 7, 8, 9, 0})
	fmt.Println(genericsInt) // 输出：9
}
