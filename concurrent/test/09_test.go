package test

import (
	"golang/concurrent"
	"testing"
)

func TestMap(t *testing.T) {
	concurrent.MapNilValue()
}
func TestMapSyncPanic(t *testing.T) {
	var m = make(map[int]int, 10) // 初始化一个map
	go func() {
		for {
			m[1] = 1 //设置key
		}
	}()
	go func() {
		for {
			_ = m[2] //访问这个map
		}
	}()
	select {}
}
