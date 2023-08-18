package concurrent

import (
	"context"
	"fmt"
	"testing"
	"time"
)

/*
Context：go 1.9 后内置 Context
	1.根 Context：通过 context.Background() 创建
	2.子 Context：context.WithCancel(parentContext) 创建
		ctx, cancel := context.WithCancel(context.Background())
	3.当前 Context 被取消时，基于它的子 Context 都会被取消
	4.接收取消通知 <- ctx.Done()
*/
// 关联任务的取消
func ctxCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
func TestContextCancle(t *testing.T) {
	// ctx 组合了 context.Background()
	//ctx, cancel := context.WithCancel(context.Background())
	p := context.Background()            // context.emptyCtx
	ctx, cancel := context.WithCancel(p) // context.cancelCtx
	fmt.Printf("%[1]T,%[2]T\n%[1]p,%[2]p", p, ctx)
	for i := 0; i < 5; i++ {
		go func(i int, ctx context.Context) {
			for {
				if ctxCancelled(ctx) {
					break
				}
				time.Sleep(time.Millisecond * 50)
			}
			fmt.Println(i, "cancelled")
		}(i, ctx)
	}
	t.Log("cancel begin")
	cancel()
	time.Sleep(time.Second * 1)
}
