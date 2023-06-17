package performance

import (
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"testing"
	"time"
)

func TestProfile(t *testing.T) {
	// 创建输出文件
	f, err := os.Create("prof/cpu.prof")
	if err != nil {
		t.Fatal("could not create test profile:", err)
	}
	// 获取系统信息
	if err = pprof.StartCPUProfile(f); err != nil { // 监控 CPU
		t.Fatal("could not start CPU profile:", err)
	}
	defer pprof.StopCPUProfile()

	// 主逻辑区，进行一些简单的代码运算
	x := [r][c]int{}
	fillMatrix(&x)
	calculate(&x)

	f1, err := os.Create("prof/memory.prof")
	defer f1.Close()
	if err != nil {
		t.Fatal("could not create memory profile:", err)
	}
	runtime.GC()                                      // GC，获取最新的数据信息
	if err = pprof.WriteHeapProfile(f1); err != nil { // 写入内存信息
		t.Fatal("could not write memory profile:", err)
	}

	f2, err := os.Create("prof/goroutine.prof")
	defer f2.Close()
	if err != nil {
		t.Fatal("could not create goroutine profile:", err)
	}
	if gProf := pprof.Lookup("goroutine"); gProf == nil { // 协程：pprof.Lookup("goroutine")
		t.Log("could not write goroutine profile:")
	} else {
		gProf.WriteTo(f2, 0)
	}
}

const (
	r = 10000
	c = 10000
)

func fillMatrix(m *[r][c]int) {
	s := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			m[i][j] = s.Intn(100000)
			//rand.Intn(100000)
		}
	}
}
func calculate(m *[r][c]int) {
	for i := 0; i < r; i++ {
		sum := 0
		for j := 0; j < c; j++ {
			sum += m[i][j]
		}
	}
}
