package concurrent

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"
)

/*
仅需任意任务完成
	1.使用无缓冲 channel
		其他协程会阻塞，内存泄漏 以致造成 OOM
	2.使用有缓冲 channel
		缓冲区大小 "=" 协程数
所有任务完成
	1.可以使用 sync.WaitGroup
	2.这里使用 CSP
*/

// 所有任务完成 CSP
func AllResponse() string {
	numOfRunner := 10
	ch := make(chan string, numOfRunner)
	for i := 0; i < numOfRunner; i++ {
		go func(i int) {
			ret := runTask(i)
			ch <- ret
		}(i)
	}
	var sb strings.Builder
	for i := 0; i < numOfRunner; i++ {
		sb.WriteString(<-ch) // 完成一个，读一个
		//sb.WriteRune('\n')
	}
	return sb.String()
}
func TestAllResponse(t *testing.T) {
	t.Log("before:", runtime.NumGoroutine()) // 2 ?
	t.Log(AllResponse())
	time.Sleep(time.Second * 1)
	t.Log("after:", runtime.NumGoroutine()) // 11 / 2
}

// 仅需任意任务完成
func runTask(id int) string {
	time.Sleep(10 * time.Millisecond)
	return fmt.Sprintf("The result is from %d\n", id)
}
func FirstResponse() string {
	numOfRunner := 10
	ch := make(chan string) // 无缓冲
	//ch := make(chan string, numOfRunner) // 缓冲区大小 = 协程数
	for i := 0; i < numOfRunner; i++ {
		go func(i int) {
			ret := runTask(i)
			ch <- ret // 无缓冲时，阻塞
		}(i)
	}
	return <-ch
}
func TestFirstResponse(t *testing.T) {
	t.Log("before:", runtime.NumGoroutine()) // 2 ?
	t.Log(FirstResponse())
	time.Sleep(time.Second * 1)
	t.Log("after:", runtime.NumGoroutine()) // 11 / 2
}
