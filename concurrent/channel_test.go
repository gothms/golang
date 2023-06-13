package advanced

import (
	"fmt"
	"sync"
	"testing"
)

/*
channel的关闭
	1.向关闭的channel发送数据，会导致panic
	2.v, ok <- ch; ok 为 bool 值，true正常接收，false通道已关闭
	3.所有channel接收者都会在channel关闭时，立刻从阻塞等待中返回，且 ok=false
		这个广播机制常被利用，进行向多个订阅者同时发送信号
			如 退出信号
			也避免了对订阅者数量的耦合
*/
func TestCloseChannel(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan int)
	wg.Add(1)
	dataProducer(ch, &wg)
	wg.Add(1)
	dataReceiver(ch, &wg)
	wg.Add(1)
	dataReceiver(ch, &wg) // 多个 receiver
	wg.Wait()
}
func dataProducer(ch chan int, wg *sync.WaitGroup) {
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
		wg.Done()
	}()
}
func dataReceiver(ch chan int, wg *sync.WaitGroup) {
	go func() {
		//for i := 0; i < 11; i++ { // channel 已 close 读取到默认值 0
		for {
			if data, ok := <-ch; ok {
				fmt.Println(data)
			} else {
				break
			}
		}
		wg.Done()
	}()
}
