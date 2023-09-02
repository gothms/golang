package test

import (
	"fmt"
	"golang/concurrent"
	"testing"
	"time"
)

func TestOrDone(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-concurrent.OrDone(
		sig(1*time.Second),
		sig(3*time.Second),
		sig(10*time.Second),
		sig(20*time.Second),
		sig(30*time.Second),
		sig(40*time.Second),
		sig(50*time.Second),
		sig(01*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}
func TestChannelMutexDemo(t *testing.T) {
	concurrent.ChannelMutexDemo()
}
func TestChannelReflectSelect(t *testing.T) {
	concurrent.ChannelReflectSelect()
}
