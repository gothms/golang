package benchmark

import (
	"bytes"
	"strings"
	"testing"
)

/*
$ go test -bench=.
$ go test -bench=方法名
	windows 下使用 go test 命令行时，-bench=. 应写为 -bench="."
	go test -bench="."
		使用 . 前，命令行先 cd 到文件所在目录
	go test -bench=BenchmarkConcatStringByBuilder
		测试方法时，可以加引号，也可以不不加

$ go test -bench="." -benchmem
	内存大小
	allocs 次数，新的内存分配操作

卡住不动
	$ go test -bench="." -benchmem：如果某个目录卡住不动
	可以尝试新建个包，再测试
*/

func TestConcatStringByAdd(t *testing.T) {
	strs := []string{"1", "2", "3", "4", "5"}
	var ret string
	for _, s := range strs {
		ret += s
	}
	if ret != "12345" {
		t.Error("string + error")
	}
}
func TestConcatStringByBytesBuffer(t *testing.T) {
	strs := []string{"1", "2", "3", "4", "5"}
	var buf bytes.Buffer
	for _, s := range strs {
		buf.WriteString(s)
	}
	if buf.String() != "12345" {
		t.Error("bytes buffer error")
	}
}
func BenchmarkConcatStringByAdd(b *testing.B) {
	strs := []string{"1", "2", "3", "4", "5"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var ret string
		for _, s := range strs {
			ret += s
		}
	}
	b.StopTimer()
}
func BenchmarkConcatStringByBytesBuffer(b *testing.B) {
	strs := []string{"1", "2", "3", "4", "5"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		for _, s := range strs {
			buf.WriteString(s)
		}
	}
	b.StopTimer()
}
func BenchmarkConcatStringByBuilder(b *testing.B) {
	strs := []string{"1", "2", "3", "4", "5"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		//sb.Grow()
		for _, s := range strs {
			sb.WriteString(s)
		}
	}
}
