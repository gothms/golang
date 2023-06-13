package concurrent

import (
	"fmt"
	"sync"
	"testing"
)

type Singleton struct {
}

var singleInstance *Singleton
var once sync.Once

func GetSingleton() *Singleton {
	once.Do(func() {
		fmt.Println("Create Singleton")
		singleInstance = new(Singleton)
	})
	return singleInstance
}
func TestOnce(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			singleton := GetSingleton()
			fmt.Printf("%p\n", singleton)
			wg.Done()
		}()
	}
	wg.Wait()
}
