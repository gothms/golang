package test

import (
	"golang/concurrent"
	"testing"
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
	i, ok = <-c // 已读完
	t.Log(i, ok)
}
