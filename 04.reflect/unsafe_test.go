package unsafe

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

/*
unsafe.Pointer 实际使用场景
	和外部 C 程序实现的高效的库，交互的时候
	其他使用场景比较少
	适用总结：
		性能极致的程序
		某些地方要用 C 程序来追求更高性能

“不安全”行为的危险性
	Go不支持强制类型转换
	unsafe.Pointer 可以把持有的指针，转换为任意类型的指针
		i := 10
		f := *(*float64)(unsafe.Pointer(&i))
示例
	1.合理的类型转换
	2.原子类型操作：atomic 包实现 共享buff 的安全并发读写
		提供了指针的原子操作，用于并发读写共享缓存，达到读写的安全性
		1.先写到另外一个新的地方
		2.写完后，用原子操作把读和写的地方，重新指向新的地方
		3.再读时，就会读到这块新写好的地方
*/

// 原子类型操作
func TestAtomic(t *testing.T) {
	var shareBufPtr unsafe.Pointer
	writeDataFn := func() {
		data := make([]int, 0) // 1 新的数据
		for i := 0; i < 10; i++ {
			data = append(data, i)
		}
		atomic.StorePointer(&shareBufPtr, unsafe.Pointer(&data)) // 2
	}
	readDataFn := func() {
		data := atomic.LoadPointer(&shareBufPtr) // 3
		t.Log(data, *(*[]int)(data))
	}
	var wg sync.WaitGroup
	writeDataFn()
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 5; j++ {
				writeDataFn()
				time.Sleep(time.Millisecond * 100)
			}
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			for j := 0; j < 5; j++ {
				readDataFn()
				time.Sleep(time.Millisecond * 100)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

type MyInt int

// 合理的类型转换
func TestConvert(t *testing.T) {
	a := []int{1, 2, 3, 4}
	b := *(*[]MyInt)(unsafe.Pointer(&a))
	t.Log(b)
}
func TestUnsafe(t *testing.T) {
	i := 10
	f := (*float64)(unsafe.Pointer(&i)) // 0xc000010390
	t.Log(unsafe.Pointer(&i))           // 0xc000010390
	t.Log(*f)                           // 5e-323
}
