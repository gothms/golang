package test

import (
	"fmt"
	"golang/concurrent"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestDeadlock 持有和等待
func TestDeadlock(t *testing.T) {
	// 派出所证明
	var psCertificate sync.Mutex
	// 物业证明
	var propertyCertificate sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2) // 需要派出所和物业都处理
	// 派出所处理goroutine
	go func() {
		defer wg.Done() // 派出所处理完成
		psCertificate.Lock()
		defer psCertificate.Unlock()
		// 检查材料
		time.Sleep(5 * time.Second)
		// 请求物业的证明
		propertyCertificate.Lock()
		propertyCertificate.Unlock()
	}()
	// 物业处理goroutine
	go func() {
		defer wg.Done() // 物业处理完成
		propertyCertificate.Lock()
		defer propertyCertificate.Unlock()
		// 检查材料
		time.Sleep(5 * time.Second)
		// 请求派出所的证明
		psCertificate.Lock()
		psCertificate.Unlock()
	}()
	wg.Wait()
	fmt.Println("成功完成")
}

// TokenRecursiveMutex 测试可重入锁 token
func TestTokenRecursiveMutex(t *testing.T) {
	token := rand.Intn(math.MaxInt32)
	l := &concurrent.TokenRecursiveMutex{}
	fooRecursiveMutex(l, int64(token))
}
func fooRecursiveMutex(l *concurrent.TokenRecursiveMutex, token int64) {
	fmt.Println("in foo")
	l.Lock(token)
	barRecursiveMutex(l, token)
	l.Unlock(token)
}
func barRecursiveMutex(l *concurrent.TokenRecursiveMutex, token int64) {
	l.Lock(token)
	fmt.Println("in bar")
	l.Unlock(token)
}

// TestRecursiveMutex 测试可重入锁 goroutine id
// hacker 方式
func TestRecursiveMutex(t *testing.T) {
	l := &concurrent.RecursiveMutex{}
	fooReentrantLock(l)
}

// TestGetGoId 通过 runtime.Stack 获取 goroutine id
func TestGetGoId(t *testing.T) {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// 得到id字符串
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))
	t.Log(idField, len(idField))
	id, err := strconv.Atoi(idField[0])
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	t.Log(id)
}

// TestReentrantLock 可重入锁
func TestReentrantLock(t *testing.T) {
	l := &sync.Mutex{}
	fooReentrantLock(l)
}

func fooReentrantLock(l sync.Locker) {
	fmt.Println("in foo")
	l.Lock()
	bar(l)
	l.Unlock()
}
func bar(l sync.Locker) {
	l.Lock()
	fmt.Println("in bar")
	l.Unlock()
}

// TestCheckdead Copy 已使用的 Mutex
func TestCheckdead(t *testing.T) {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("checkdead")
	foo(mu) // 复制锁
}

// 这里sync.Mutex的参数是通过复制的方式传入的
func foo(mu sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("in foo")
}
