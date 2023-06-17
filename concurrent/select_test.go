package concurrent

import (
	"fmt"
	"testing"
	"time"
)

// select Test
func TestSelect(t *testing.T) {
	select {
	case ret := <-AsyncService():
		t.Log(ret)
	//case <-time.After(time.Millisecond * 100):
	case <-time.After(time.Millisecond * 30):
		t.Error("time out.")
	}
}

// channel Test
func TestChannel(t *testing.T) {
	ch := AsyncService()
	otherTask()
	t.Log(<-ch)
}
func otherTask() {
	fmt.Println("working...")
	time.Sleep(time.Millisecond * 100)
	fmt.Println("Task is done.")
}
func AsyncService() chan string {
	//ch := make(chan string)
	ch := make(chan string, 1)
	go func() {
		ret := service()
		fmt.Println("service return.")
		ch <- ret
		fmt.Println("sync service exited.")
	}()
	return ch
}
func service() string {
	time.Sleep(time.Millisecond * 50)
	return "service Done"
}
