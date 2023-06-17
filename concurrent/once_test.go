package concurrent

import (
	"fmt"
	"sync"
	"testing"
	"unsafe"
)

type Singleton struct {
}

var singleInstance *Singleton
var once sync.Once

// GetSingleton 单例：懒汉式
func GetSingleton() *Singleton {
	once.Do(func() { // 只运行了一次
		fmt.Println("Create Singleton")
		singleInstance = new(Singleton)
	})
	return singleInstance
}
func TestOnce(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			singleton := GetSingleton()
			//fmt.Printf("%d,%p\n", i, singleton)
			fmt.Println(i, unsafe.Pointer(singleton))
			wg.Done()
		}(i)
	}
	wg.Wait()
}
