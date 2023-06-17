package gc

import (
	"os"
	"runtime"
	"runtime/trace"
	"testing"
)

/*
Go GC

	采用染色标记的 GC，和 Java 相当不同

避免内存分配和复制

	1.复杂对象尽量传递引用
		数组的传递
		结构体的传递

		测试：
		BenchmarkWithValue-8                  54          19943587 ns/op        80003123 B/op          1 allocs/op
		BenchmarkWithReference-8        1000000000               0.2553 ns/op          0 B/op          0 allocs/op
	2.Slice
		初始化至合适的大小：自动扩容是有代价的
		复用内存

打开 GC 日志

	windows：
		set GOGCTRACE=1
		set GODEBUG=gctrace=1
		日志输出到文件：go run main.go 2> gctrace.log
	Linux：
		在程序执行之前加上环境变量 GODEBUG=gctrace=1
		之后加上 2> gctrace.log：日志输出到文件，不加则输出到控制台
			如：
			GODEBUG=gctrace=1 go test -bench="." 2>gctrace.log
			GODEBUG=gctrace=1 go run main.go
		日志详细信息参考： https://godoc.org/runtime
			也可以查 YouTube/Google 视频

Go tool trace

	普通程序输出 trace 信息
		err := trace.Start(file)
		defer trace.Stop()
	测试程序输出 trace 信息
		go test -bench=BenchmarkWithValue --trace=trace_val.out
		go test -bench=BenchmarkWithReference --trace=trace_ref.out
	可视化 trace 信息
		go tool trace trace_val.out：跳转到浏览器->View trace
		go tool trace trace_ref.out

		我的 Chrome 报错：
		Trace Viewer is running with WebComponentsV0 polyfill, and some features may be broken. As a workaround, you may try running chrome with "--enable-blink-features=ShadowDOMV0,CustomElementsV0,HTMLImports" flag. See crbug.com/1036492.
*/
// 普通程序输出 trace 信息
func TestTrace(t *testing.T) {
	f, err := os.Create("trace_test.out")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	err = trace.Start(f)
	defer trace.Stop()
	if err != nil {
		panic(err)
	}
}

// GC 日志
const NumOfElems = 1000

type Content struct {
	Detail [10000]int
}

func withValue(arr [NumOfElems]Content) int {
	return 0
}
func withReference(arr *[NumOfElems]Content) int {
	return 0
}
func TestGC(t *testing.T) {
	var arr [NumOfElems]Content
	withValue(arr)
	withReference(&arr)
}
func BenchmarkWithValue(b *testing.B) {
	var arr [NumOfElems]Content
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		withValue(arr)
		runtime.GC()
	}
}
func BenchmarkWithReference(b *testing.B) {
	var arr [NumOfElems]Content
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		withReference(&arr)
	}
}
