package test

import (
	"bytes"
	"testing"
)

/*
$ go test -bench=.
$ go test -bench=方法名
	windows 下使用 go test 命令行时，-bench=. 应写为 -bench="."
	go test -bench="."
$ go test -bench="." -benchmem

*/
//func TestConcatStringByAdd(t *testing.T) {
//	strs := []string{"1", "2", "3", "4", "5"}
//	var ret string
//	for _, s := range strs {
//		ret += s
//	}
//	if ret != "12345" {
//		t.Error("string + error")
//	}
//}
//func TestConcatStringByBytesBuffer(t *testing.T) {
//	strs := []string{"1", "2", "3", "4", "5"}
//	var buf bytes.Buffer
//	for _, s := range strs {
//		buf.WriteString(s)
//	}
//	if buf.String() != "12345" {
//		t.Error("bytes buffer error")
//	}
//}
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
