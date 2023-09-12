package test

import (
	"golang/concurrent"
	"testing"
	"unsafe"
)

func TestChannelMapGoQueue(t *testing.T) {
	concurrent.ChannelMapGoQueue()
}
func TestChannelQueue(t *testing.T) {
	concurrent.ChannelQueue()
}
func TestAtomicQueue(t *testing.T) {
	concurrent.AtomicQueue()
}
func TestCondQueue(t *testing.T) {
	concurrent.CondQueue()
}
func TestChannel(t *testing.T) {
	c := make(chan int, 3)
	c <- 1
	c <- 2
	c <- 3
	i := <-c
	t.Log(i)
	close(c)
	i, ok := <-c // 2 true
	t.Log(i, ok)
	i, ok = <-c // 3 true
	t.Log(i, ok)
	i, ok = <-c // 已读完 0 false
	t.Log(i, ok)

	maxAlign := 8
	hchanSize := unsafe.Sizeof(hchanTest{}) + uintptr(-int(unsafe.Sizeof(hchanTest{}))&(maxAlign-1))
	t.Log(unsafe.Sizeof(hchanTest{}))                                 // 48
	t.Log(uintptr(-int(unsafe.Sizeof(hchanTest{})) & (maxAlign - 1))) // 0
	t.Log(-int(unsafe.Sizeof(hchanTest{})))                           // -48 = ^48 + 1
	t.Log(hchanSize)                                                  // 48

	t.Log(unsafe.Sizeof(uint(0)))   // 8
	t.Log(unsafe.Sizeof(uint16(0))) // 2
	t.Log(unsafe.Sizeof(uint32(0))) // 4
	v := 1
	t.Log(unsafe.Sizeof(unsafe.Pointer(&v))) // 8

	bug := EtcdBug()
	t.Log(bug)

	//tc := make(chan TChannel, 1)
	//tc = TChannelQuestion1()
	//tc = TChannelQuestion2()
}
func EtcdBug() int {
	n := make(chan chan int, 1)
	c := make(chan int)
	n <- c // fatal error: all goroutines are asleep - deadlock!
	go func() {
		c <- 56 // 解决方案
	}()
	return <-c // reader 一直阻塞
}
func TChannelQuestion1() <-chan TChannel {
	t := make(chan TChannel, 1)
	return t
}
func TChannelQuestion2() chan<- TChannel {
	t := make(chan TChannel, 1)
	return t
}

type TChannel struct{}

type hchanTest struct {
	qcount   uint           // total data in the queue：循环队列元素的数量
	dataqsiz uint           // size of the circular queue：循环队列的大小
	buf      unsafe.Pointer // points to an array of dataqsiz elements：循环队列的指针
	elemsize uint16         // chan 中元素的大小
	closed   uint32         // 是否已 close
	//elemtype *_type         // element type：chan 中元素类型
	sendx uint // send index：send 在 buf 中的索引
	recvx uint // receive index：recv 在 buf 中的索引
	//recvq    waitq          // list of recv waiters：receiver 的等待队列
	//sendq    waitq          // list of send waiters：sender 的等待队列
}
