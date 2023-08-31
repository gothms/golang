package test

import (
	"golang/concurrent"
	"math/rand"
	"testing"
	"time"
)

func TestMutexMessage(t *testing.T) {
	var mu concurrent.Mutex
	for i := 0; i < 1000; i++ {
		go func() {
			mu.Lock()
			time.Sleep(time.Second)
			mu.Unlock()
		}()
	}
	time.Sleep(1 * time.Second)
	//go func() {
	//	mu.Lock()
	//	mu.Unlock()
	//}()
	//time.Sleep(300 * time.Millisecond)
	t.Logf("waitings:%d,isLocked:%t,woken:%t,starving:%t\n",
		mu.WaiterCount(), mu.IsLocked(), mu.IsWoken(), mu.IsStarving())
}

func TestWaiterCount(t *testing.T) {
	var mut concurrent.Mutex
	go func() {
		mut.Lock()
		defer mut.Unlock()
		time.Sleep(1 * time.Second)
	}()
	for i := 0; i < 5; i++ {
		go func() {
			mut.Lock()
			defer mut.Unlock()
			time.Sleep(60 * time.Millisecond)
		}()
	}
	time.Sleep(500 * time.Millisecond)
	cnt := mut.WaiterCount()
	mut.Lock()
	defer mut.Unlock()
	t.Log("cnt:", cnt)
}

func TestTryLock(t *testing.T) {
	var mu concurrent.Mutex
	go func() { // 启动一个g持有一段时间的锁
		mu.Lock()
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
		mu.Unlock()
	}()
	time.Sleep(time.Second)
	ok := mu.TryLock() // 尝试获取锁
	if ok {            // 获取锁成功
		t.Log("got the lock")
		// 开始你的业务
		mu.Unlock()
		return
	}
	// 没有获取到
	t.Log("can't get the lock")
}
